package bitcask

import (
	"bytes"
	"fmt"
	"log/slog"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/Dev79844/bitcask/internal/datafile"
	"github.com/Dev79844/bitcask/internal/keydir"
)

var(
	HINTS_FILE = "bitcask.hint"
)

type Bitcask struct{
	sync.Mutex
	bufPool 	sync.Pool
	df      	*datafile.DataFile
	keydir		keydir.KeyDir
	staleFiles	map[int]*datafile.DataFile
	opts        *Options
	l			*slog.Logger
}

func Open(cfg ...Config) (*Bitcask, error) {
	options := DefaultOptions()
	for _, ops := range cfg{
		if err := ops(options); err!=nil{
			return nil, err
		}
	}

	files, err := getFiles(options.dir)
	if err!=nil{
		return nil, fmt.Errorf("error getting the existing files: %v", err)
	}

	var(
		index = 0
		staleFiles = map[int]*datafile.DataFile{}
		l = initLogger()
	)

	if len(files) > 0{
		ids, err := getIDs(files)
		if err!=nil{
			return nil, fmt.Errorf("error getting ids of files: %v", err)
		}

		index = ids[len(ids) - 1] + 1

		for _,idx := range ids{
			df, err := datafile.New(options.dir, idx)
			if err != nil {
				return nil, err
			}
			staleFiles[idx] = df
		}
	}

	df, err := datafile.New(options.dir, index)
	if err!=nil{
		return nil, fmt.Errorf("error creating a new datafile: %v", err)
	}
	
	kd := make(keydir.KeyDir, 0)

	hintPath := filepath.Join(options.dir, HINTS_FILE)
	if exists(hintPath){
		if err := kd.Decode(hintPath); err!=nil{
			return nil, fmt.Errorf("error populating keydir from hint file: %v", err)
		}
	}

	b := &Bitcask{
		df: df,
		bufPool: sync.Pool{New: func() any {
			return bytes.NewBuffer([]byte{})
		}},
		keydir: kd,
		staleFiles: staleFiles,
		opts: options,
		l: l,
	}

	// goroutine for compaction
	go b.RunCompactionWithInterval(options.compactInterval)

	// goroutine for running fsync periodically
	if b.opts.syncInterval != nil {
		go b.SyncFile(*options.syncInterval)
	}

	// goroutine for checking file size and rotating it
	go b.CheckFileSize(options.checkFileSizeInterval)

	return b, nil
}

func (b *Bitcask) Put(key string, value []byte) error {
	b.Lock()
	defer b.Unlock()

	if len(key) == 0 {
		return ErrEmptyKey
	}

	if b.opts.readOnly{
		return ErrReadOnly
	}

	b.l.Info("storing", "key", key, "value", value)
	return b.put(b.df, key, value)
}

func (b *Bitcask) Get(key string) ([]byte, error) {
	b.Lock()
	defer b.Unlock()

	if len(key) == 0{
		return nil, ErrEmptyKey
	}
	
	b.l.Info("getting value","key", key)
	row, err := b.get(key)
	if err!=nil{
		return nil, fmt.Errorf("error getting the record: %v", err)
	}

	if !row.isValidChecksum(){
		return nil, fmt.Errorf("error validating the checksum")
	}

	return row.Value, nil
}

func (b *Bitcask) Delete(key string) error {
	b.Lock()
	defer b.Unlock()

	if len(key) == 0 {
		return ErrEmptyKey
	}

	if b.opts.readOnly{
		return ErrReadOnly
	}

	b.l.Info("removing data", "key", key)
	return b.delete(key)
}

func (b *Bitcask) List_Keys() []string {
	b.Lock()
	defer b.Unlock()

	keys := make([]string, 0, len(b.keydir))

	for k := range b.keydir {
		keys = append(keys, k)
	}

	return keys
}

func (b *Bitcask) Fold(fn func(k string) error) error {
	b.Lock()
	defer b.Unlock()

	for k := range b.keydir{
		if err := fn(k); err!=nil{
			return err
		}
	}

	return nil
}

func (b *Bitcask) Sync() error {
	b.Lock()
	defer b.Unlock()

	return b.df.Sync()
}

func (b *Bitcask) Merge() error {
	b.Lock()
	defer b.Unlock()

	b.l.Info("merging all files")
	return b.RunCompaction()
}

func (b *Bitcask) Close() error {
	b.Lock()
	defer b.Unlock()

	if err := b.generateHintFiles(); err!=nil{
		b.l.Error("error generating hint file")
		return err
	}

	if err := b.df.Close(); err!=nil{
		b.l.Error("error closing the active data file")
		return err
	}

	for _, df := range b.staleFiles {
		if err := df.Close(); err!=nil{
			b.l.Error("error closing the stale file", slog.String("%d", strconv.Itoa(df.ID())))
			return err
		}
	}

	return nil
}
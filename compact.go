package bitcask

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/Dev79844/bitcask/internal/datafile"
)

func (b *Bitcask) RunCompactionWithInterval(interval time.Duration) error {
	ticker := time.NewTicker(interval)
	for range ticker.C{
		if err := b.RunCompaction(); err!=nil{
			return err
		}
	}
	return nil
}

func (b *Bitcask) RunCompaction() error {
	if err := b.merge(); err!=nil{
		return err
	}
	
	if err := b.generateHintFiles(); err!=nil{
		return err
	}

	return nil
}

func (b *Bitcask) merge() error {
	//create a temp dir to store meged file
	// create a merged file
	//get all the values from keydir and put them in merged file
	//close all stale files and delete them
	// reset the stale map
	// delete all .db files
	// rename the merged file and move it to the dir in config
	// set the active df as merged df

	var fsync bool

	if len(b.staleFiles) < 2 {
		return nil
	}

	if b.opts.alwaysFSync{
		fsync = true
		b.opts.alwaysFSync = false
	}

	tempDir, err := os.MkdirTemp(".", "merged")
	if err!=nil{
		return err
	}
	defer os.RemoveAll(tempDir)

	mergedDF, err := datafile.New(tempDir, 0)
	if err!=nil{
		return err
	}

	for k := range b.keydir {
		row, err := b.get(k)
		if err!=nil{
			return err
		}

		if err = b.put(mergedDF, k, row.Value); err!=nil{
			return err
		}
	}

	for _, df := range b.staleFiles {
		if err := df.Close(); err!=nil{
			return err
		}
	}

	b.staleFiles = make(map[int]*datafile.DataFile, 0)

	err = filepath.Walk(b.opts.dir, func(path string, info fs.FileInfo, err error) error{
		if err!=nil{
			return err
		}

		if info.IsDir(){
			return nil
		}

		if filepath.Ext(path) == ".db"{
			if err = os.Remove(path); err!=nil{
				return err
			}
		}

		return nil
	})
	if err !=nil{
		return err
	}

	os.Rename(filepath.Join(tempDir, fmt.Sprintf(datafile.ACTIVE_FILE, 0)), 
	filepath.Join(b.opts.dir, fmt.Sprintf(datafile.ACTIVE_FILE, 0)))

	b.df = mergedDF

	if fsync{
		b.opts.alwaysFSync = true
		b.df.Sync()
	}

	return nil
}

func (b *Bitcask) generateHintFiles() error {
	path := filepath.Join(b.opts.dir, HINTS_FILE)
	if err := b.keydir.Encode(path); err!=nil{
		return err
	}

	return nil
}

func(b *Bitcask) SyncFile(interval time.Duration) error {
	ticker := time.NewTicker(interval)

	for range ticker.C{
		if err := b.Sync(); err!=nil{
			return err
		}
	}
	
	return nil
}
package internal

import (
	"bytes"
	"fmt"
	// "log"
	// "os"
	"sync"

	"github.com/Dev79844/bitcask/internal/keydir"
	"github.com/Dev79844/bitcask/internal/datafile"
)

type Bitcask struct{
	sync.Mutex
	bufPool 	sync.Pool
	df      	*datafile.DataFile
	keydir		keydir.KeyDir
	staleFiles	map[int]*datafile.DataFile
}

func Open(dir string) (*Bitcask, error) {
	files, err := getFiles(dir)
	if err!=nil{
		return nil, fmt.Errorf("error getting the existing files: %v", err)
	}

	var(
		index = 0
		staleFiles = map[int]*datafile.DataFile{}
	)

	if len(files) > 0{
		ids, err := getIDs(files)
		if err!=nil{
			return nil, fmt.Errorf("error getting ids of files: %v", err)
		}

		index = ids[len(ids) - 1] + 1

		for _,idx := range ids{
			df, err := datafile.New(dir, idx)
			if err != nil {
				return nil, err
			}
			staleFiles[idx] = df
		}
	}

	df, err := datafile.New(dir, index)
	if err!=nil{
		return nil, fmt.Errorf("error creating a new datafile: %v", err)
	}
	
	kd := make(keydir.KeyDir, 0)

	return &Bitcask{
		df: df,
		bufPool: sync.Pool{New: func() any {
			return bytes.NewBuffer([]byte{})
		}},
		keydir: kd,
		staleFiles: staleFiles,
		}, nil
}

func (b *Bitcask) Put(key string, value []byte) error {
	b.Lock()
	defer b.Unlock()

	if len(key) == 0 {
		return ErrEmptyKey
	}

	return b.put(b.df, key, value)
}

func (b *Bitcask) Get(key string) ([]byte, error) {
	b.Lock()
	defer b.Unlock()

	if len(key) == 0{
		return nil, ErrEmptyKey
	}
	
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

	return b.delete(key)
}
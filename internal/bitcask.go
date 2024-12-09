package internal

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/Dev79844/bitcask/internal/crc"
	"github.com/Dev79844/bitcask/internal/keydir"
)

type Bitcask struct{
	bufPool sync.Pool
	file	*os.File
	keydir	keydir.KeyDir
}

func Open(filename string) (*Bitcask, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY | os.O_APPEND | os.O_CREATE | os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	
	kd := make(keydir.KeyDir, 0)

	return &Bitcask{
		file: file,
		bufPool: sync.Pool{New: func() any {
			return bytes.NewBuffer([]byte{})
		}},
		keydir: kd,
		}, nil
}

func (b *Bitcask) Put(key string, value []byte) error {
	header := Header{
		CRC: crc.CalcCRC(value),
		Tstamp: uint32(time.Now().Unix()),
		Ksz: uint32(len(key)),
		Valuesz: uint32(len(value)),
	}

	record := Row{
		Header: header,
		Key: key,
		Value: value,
	}

	buf := b.bufPool.Get().(*bytes.Buffer)
	defer b.bufPool.Put(&buf)
	defer buf.Reset()

	header.encode(buf)

	buf.WriteString(key)
	buf.Write(value)

	offset, err := b.file.Write(buf.Bytes())
	if err != nil {
		return fmt.Errorf("error writing data to file: %v", err)
	}

	b.keydir[key] = keydir.Meta{
		FileID: 1,
		RecordSz: len(buf.Bytes()),
		RecordPos: offset + len(buf.Bytes()),
		Tstamp: int(record.Header.Tstamp),
	}
	return nil
}

func (b *Bitcask) Get(key string) ([]byte, error) {
	val, ok := b.keydir[key]
	if !ok{
		return nil, ErrNoKey
	}

	reader := b.file

	record := make([]byte, 0)

	_, err := reader.ReadAt(record, int64(val.RecordPos))
	if err!=nil{
		log.Print(err)
	}
	return record, nil
}
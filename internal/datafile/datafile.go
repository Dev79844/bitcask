package datafile

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

const (
	ACTIVE_FILE = "bitcask_%d.db"
)

type DataFile struct {
	sync.RWMutex
	id		int
	writer	*os.File
	reader	*os.File
	offset	int
}

func New(dir string, index int) (*DataFile, error) {
	path := filepath.Join(dir, fmt.Sprintf(ACTIVE_FILE, index))
	writer, err := os.OpenFile(path, os.O_WRONLY | os.O_APPEND | os.O_CREATE | os.O_TRUNC, 0644)
	if err!=nil{
		return nil, fmt.Errorf("error opening file for writing to db: %v", err)
	}

	reader, err := os.Open(path)
	if err!=nil{
		return nil, fmt.Errorf("error opening file for reading db: %v", err)
	}

	file_data, err := os.Stat(path)
	if err!=nil{
		return nil, fmt.Errorf("error getting stats for active file: %v", err)
	}

	return &DataFile{
		id: index, 
		writer: writer,
		reader: reader,
		offset: int(file_data.Size()),
	}, nil
}

func (d *DataFile) ID() (int) {
	return d.id
}

func (d *DataFile) Write(data []byte) (int,error) {
	_, err := d.writer.Write(data)
	if err!=nil{
		return -1, err
	}

	offset := d.offset

	d.offset += len(data)

	return offset, nil
}

func (d *DataFile) Read(pos, size int) ([]byte, error) {
	start := int64(pos - size)
	record := make([]byte, size)

	n, err := d.reader.ReadAt(record, start)
	if err != nil {
		return nil, err
	}

	if n != int(size) {
		return nil, fmt.Errorf("error fetching record, invalid size")
	}

	return record, nil
}

func (d *DataFile) Sync() error {
	return d.writer.Sync()
}

func (d *DataFile) Close() error {
	if err := d.writer.Close(); err!=nil{
		return err
	}

	if err := d.reader.Close(); err!=nil{
		return err
	}

	return nil
}

func (d *DataFile) Size() (int64, error) {
	stat, err := d.writer.Stat()
	if err!=nil{
		return -1, fmt.Errorf("error fetching file stats: %v", err)
	}

	return stat.Size(), nil
}
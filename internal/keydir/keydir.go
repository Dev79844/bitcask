package keydir

import (
	"encoding/gob"
	"os"
)

type KeyDir map[string]Meta

type Meta struct {
	FileID		int
	RecordSz	int
	RecordPos	int
	Tstamp		int
}

func (k *KeyDir) Encode(path string) error {
	file, err := os.Create(path)
	if err!=nil{
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)

	err = encoder.Encode(k)
	if err!=nil{
		return err
	}

	return nil
}

func (k *KeyDir) Decode(path string) error {
	file, err := os.Open(path)
	if err!=nil{
		return err
	}

	decoder := gob.NewDecoder(file)

	err = decoder.Decode(k)
	if err!=nil{
		return err
	}

	return nil
}
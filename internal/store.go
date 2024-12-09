package internal

import (
	"bytes"
	"fmt"
	"time"

	"github.com/Dev79844/bitcask/internal/crc"
	"github.com/Dev79844/bitcask/internal/datafile"
	"github.com/Dev79844/bitcask/internal/keydir"
)


func (b *Bitcask) put(df *datafile.DataFile, k string, v []byte) error {
	header := Header{
		CRC: crc.CalcCRC(v),
		Tstamp: uint32(time.Now().Unix()),
		Ksz: uint32(len(k)),
		Valuesz: uint32(len(v)),
	}

	record := Row{
		Header: header,
		Key: k,
		Value: v,
	}

	buf := b.bufPool.Get().(*bytes.Buffer)
	defer b.bufPool.Put(buf)
	defer buf.Reset()

	header.encode(buf)

	buf.WriteString(k)
	buf.Write(v)

	offset, err := df.Write(buf.Bytes())
	if err != nil {
		return fmt.Errorf("error writing data to file: %v", err)
	}

	b.keydir[k] = keydir.Meta{
		FileID: 1,
		RecordSz: len(buf.Bytes()),
		RecordPos: offset + len(buf.Bytes()),
		Tstamp: int(record.Header.Tstamp),
	}

	return nil
}

func (b *Bitcask) get(k string) (Row, error) {
	val, ok := b.keydir[k]
	if !ok {
		return Row{}, ErrNoKey
	}
	var header Header
	
	reader := b.df

	row, err := reader.Read(val.RecordPos, val.RecordSz)
	if err!=nil{
		return Row{}, fmt.Errorf("error reading the value from file: %v", err)
	}

	if err := header.decode(row); err!=nil{
		return Row{}, fmt.Errorf("error decoding the header")
	}

	valPos := val.RecordSz - int(header.Valuesz)
	value := row[valPos:]
	
	return Row{
		Header: header,
		Key: k,
		Value: value,
	}, nil
}
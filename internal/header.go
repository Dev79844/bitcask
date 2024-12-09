package internal

import (
	"bytes"
	"encoding/binary"

	"github.com/Dev79844/bitcask/internal/crc"
)

type Header struct {
	CRC			uint32
	Tstamp		uint32
	Ksz			uint32
	Valuesz		uint32
}

type Row struct{
	Header		Header	
	Key 		string
	Value		[]byte
}

func (h *Header) encode(buf *bytes.Buffer) error{
	return binary.Write(buf, binary.LittleEndian, h)
}

func (h *Header) decode(row []byte) error {
	return binary.Read(bytes.NewReader(row), binary.LittleEndian, h)
}

func (r *Row) isValidChecksum() bool {
	return crc.CalcCRC(r.Value) == r.Header.CRC
}
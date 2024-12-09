package internal

import (
	"bytes"
	"encoding/binary"
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

func (h *Header)encode(buf *bytes.Buffer) error{
	return binary.Write(buf, binary.LittleEndian, h)
}
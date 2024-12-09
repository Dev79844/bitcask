package crc

import (
	"hash/crc32"
)

const IEEE = 0xedb88320

func CalcCRC(val []byte) uint32 {
	return crc32.Checksum(val, crc32.MakeTable(IEEE))
}


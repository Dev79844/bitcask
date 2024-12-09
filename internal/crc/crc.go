package crc

import (
	"hash/crc32"
)

func CalcCRC(val []byte) uint32 {
	return crc32.ChecksumIEEE(val)
}


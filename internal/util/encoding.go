package util

import (
	"encoding/binary"
)

func Uint64ToByteArray(
	num uint64,
) []byte {
	buf := make(
		[]byte,
		binary.MaxVarintLen64,
	)
	n := binary.PutUvarint(buf, num)
	return buf[:n]
}

func ByteArrayToUint64(
	buf []byte,
) uint64 {
	res, _ := binary.Varint(buf)
	return uint64(res)
}

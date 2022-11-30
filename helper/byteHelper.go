package helper

import "encoding/binary"

func uint64ToByteArrayU8(u64 uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(u64))
	return b
}

func byteArrayToString(byteArray []byte) string {
	str1 := string(byteArray[:])
	return str1
}

func byteArrayToUint64(byteArray []byte) uint64 {
	return binary.LittleEndian.Uint64(byteArray)
}

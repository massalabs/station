package helper

import "encoding/binary"

func Uint64ToByteArrayU8(u64 uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(u64))
	return b
}

func StringtoByteArray(str string) []byte {
	return []byte(str)
}

func ByteArrayToString(byteArray []byte) string {
	str1 := string(byteArray[:])
	return str1
}

func ByteArrayToUint64(byteArray []byte) uint64 {
	return binary.LittleEndian.Uint64(byteArray)
}

func ByteArrayToUint32(byteArray []byte) uint32 {
	return binary.LittleEndian.Uint32(byteArray)
}

package helper

import "encoding/binary"

func Uint64ToByteArrayU8(u64 uint64) []byte {
	encode := 8
	b := make([]byte, encode)
	binary.LittleEndian.PutUint64(b, u64)

	return b
}

func StringToByteArray(str string) []byte {
	return []byte(str)
}

func ByteArrayToUint64(byteArray []byte) uint64 {
	return binary.LittleEndian.Uint64(byteArray)
}

func ByteArrayToUint32(byteArray []byte) uint32 {
	return binary.LittleEndian.Uint32(byteArray)
}

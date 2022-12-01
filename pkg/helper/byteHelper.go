package helper

import (
	"encoding/binary"
	"unicode/utf16"
)

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

func EncodeUint8ToUTF16String(numberToEncode uint32) string {
	//nolint:gomnd
	buffer := make([]byte, 4)
	binary.LittleEndian.PutUint32(buffer, numberToEncode)
	//nolint:gomnd
	runesBuffer := make([]rune, 4)

	for i := 0; i < len(buffer); i++ {
		runesBuffer[i] = utf16.Decode([]uint16{uint16(buffer[i])})[0]
	}

	encodedString := string(runesBuffer)

	return encodedString
}

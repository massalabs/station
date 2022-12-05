package convert

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

func EncodeUint32ToUTF8String(numberToEncode uint32) string {
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

func ByteArrayWithSize(arr []byte) []byte {
	return append(EncodeIntToByteArray(len(arr)), arr...)
}

func EncodeIntToByteArray(numberToEncode int) []byte {
	//nolint:gomnd
	buffer := make([]byte, 4)
	binary.LittleEndian.PutUint32(buffer, uint32(numberToEncode))

	return buffer
}

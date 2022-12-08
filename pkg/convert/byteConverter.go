package convert

import (
	"encoding/binary"
	"unicode/utf16"
)

const bytesPerUint64 = 8

// Encode uint64 to byte array.
func U64ToBytes(nb int) (bytes []byte) {
	u64 := uint64(nb)
	bytes = make([]byte, bytesPerUint64)
	binary.LittleEndian.PutUint64(bytes, u64)

	return
}

const bytesPerUint32 = 4

// Encode uint32 to byte array.
func U32ToBytes(nb int) (bytes []byte) {
	u32 := uint32(nb)
	bytes = make([]byte, bytesPerUint32)
	binary.LittleEndian.PutUint32(bytes, u32)

	return
}

func StringToBytes(str string) []byte {
	// let's start by encoding the string length.
	lenBytes := U32ToBytes(len(str))

	runesBuffer := make([]rune, bytesPerUint32)

	for i := 0; i < len(lenBytes); i++ {
		runesBuffer[i] = utf16.Decode([]uint16{uint16(lenBytes[i])})[0]
	}

	encodedLength := string(runesBuffer)
	encodedLengthBytes := []byte(encodedLength)

	return append(encodedLengthBytes, str...)
}

func BytesToString(entry []byte) string {
	content := entry[bytesPerUint32:] // content is always prefixed by its size encoded using a u32.

	return string(content)
}

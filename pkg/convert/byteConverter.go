package convert

import (
	"encoding/binary"
	"unicode/utf16"
)

const bytesPerUint64 = 8

// Encode uint64 to byte array.
func FromU64(u64 uint64) (bytes []byte) {
	bytes = make([]byte, bytesPerUint64)
	binary.LittleEndian.PutUint64(bytes, u64)

	return
}

const bytesPerUint32 = 4

// Encode uint32 to byte array.
func FromU32(u32 uint32) (bytes []byte) {
	bytes = make([]byte, bytesPerUint32)
	binary.LittleEndian.PutUint32(bytes, u32)

	return
}

func EncodeStringToByteArray(str string) []byte {
        //let's start by encoding the string length
	lenBytes := FromU32(len(str))
	//nolint:gomnd
	runesBuffer := make([]rune, 4)

	for i := 0; i < len(buffer); i++ {
		runesBuffer[i] = utf16.Decode([]uint16{uint16(buffer[i])})[0]
	}

	encodedString := string(runesBuffer)
	slice := []byte(encodedString)

	return append(slice, str...)
}

func RemoveStringEncodingPrefix(entry []byte) string {
	prefix := 4 // 4 first bytes representing the Length of the string
	entryWithoutPrefix := entry[prefix:]

	return string(entryWithoutPrefix)
}

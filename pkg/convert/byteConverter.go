package convert

import (
	"encoding/binary"
	"unicode/utf16"
)

// Encode uint64 to byte array.
func U64NbBytes(u64 uint64) []byte {
	const encode = 8
	b := make([]byte, encode)
	binary.LittleEndian.PutUint64(b, u64)

	return b
}

// Decode uint32 to string.
func U32ToString(numberToEncode uint32) string {
	//nolint:gomnd
	buffer := make([]byte, 4)
	binary.LittleEndian.PutUint32(buffer, numberToEncode)
	//nolint:gomnd
	runesBuffer := make([]rune, 4)

	for i := 0; i < len(buffer); i++ {
		runesBuffer[i] = utf16.Decode([]uint16{uint16(buffer[i])})[0]
	}

	return string(runesBuffer)
}

func EncodeStringUint32ToUTF8(str string) []byte {
	numberToEncode := len(str)
	//nolint:gomnd
	buffer := make([]byte, 4)
	binary.LittleEndian.PutUint32(buffer, uint32(numberToEncode))
	//nolint:gomnd
	runesBuffer := make([]rune, 4)

	for i := 0; i < len(buffer); i++ {
		runesBuffer[i] = utf16.Decode([]uint16{uint16(buffer[i])})[0]
	}

	encodedString := string(runesBuffer)
	slice := []byte(encodedString)

	return append(slice, str...)
}

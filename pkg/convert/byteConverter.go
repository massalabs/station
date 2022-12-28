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

func ByteToStringArray(entry []byte) []string {
	var result []string

	content := entry
	// with args you will have at least 5 bytes for a string (4 for the size and 1 for the value)
	minimumNbBytes := bytesPerUint32 + 1
	// we parse the content until there is no more string left inside
	for len(content) >= minimumNbBytes {
		// we check the string length and we update the offset
		stringLength := binary.LittleEndian.Uint32(content[:bytesPerUint32])
		offset := bytesPerUint32 + int(stringLength)

		str := string(content[bytesPerUint32:offset])

		result = append(result, str)

		// we remove the string and its length header to the content
		content = content[bytesPerUint32+int(stringLength):]
	}

	return result
}

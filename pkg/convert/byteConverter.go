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
	var stringContent string
	content := entry
	// with args you will have at least 5 bytes for a string (4 for the size and 1 for the value)
	minBytesForString := 5
	// we parse the content until there is no more string left inside
	for len(content) >= minBytesForString {
		//we check the string length and we update the offset
		stringLength := binary.LittleEndian.Uint32(content[:bytesPerUint32])
		offsetDown := len(content) - bytesPerUint32 - int(stringLength)

		//we check offset because if = 0 will throw an error and we decode the string
		if offsetDown > 0 {
			stringContent = string(content[bytesPerUint32:offsetDown])
		} else {
			stringContent = string(content[bytesPerUint32:])
		}

		result = append(result, stringContent)

		// we remove the string and its length header to the content
		content = content[bytesPerUint32+int(stringLength):]

	}
	return result
}

package convert

import (
	"bytes"
	"encoding/binary"
	"math/big"
	"unicode/utf16"

	"github.com/massalabs/station/pkg/logger"
)

const BytesPerUint64 = 8

// Encode uint64 to byte array.
func U64ToBytes(nb int) (bytes []byte) {
	u64 := uint64(nb)
	bytes = make([]byte, BytesPerUint64)
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

func ToBytes(str string) []byte {
	return []byte(str)
}

func ToBytesWithPrefixLength(str string) []byte {
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

func ToString(entry []byte) string {
	content := entry[bytesPerUint32:] // content is always prefixed by its size encoded using a u32.

	return string(content)
}

func ToStringArray(entry []byte) []string {
	var result []string

	content := entry
	// with args you will have at least 5 bytes for a string (4 for the size and 1 for the value).
	minimumNbBytes := bytesPerUint32 + 1
	// we parse the content until there is no more string left inside.
	for len(content) >= minimumNbBytes {
		// we check the string length and we update the offset.
		stringLength := binary.LittleEndian.Uint32(content[:bytesPerUint32])
		offset := bytesPerUint32 + int(stringLength)

		str := string(content[bytesPerUint32:offset])

		result = append(result, str)

		// we remove the string and its length header from the content.
		content = content[bytesPerUint32+int(stringLength):]
	}

	return result
}

// this function encodes a string array to an array of byte arrays.
func StringArrayToArrayOfByteArray(stringArray []string) [][]byte {
	stringArrayLength := len(stringArray)

	var result [][]byte

	for i := 0; i < stringArrayLength; i++ {
		result = append(result, ToBytesWithPrefixLength(stringArray[i]))
	}

	return result
}

func BytesToU64(byteArray []byte) uint64 {
	var u64 uint64

	err := binary.Read(bytes.NewReader(byteArray), binary.LittleEndian, &u64)
	if err != nil {
		logger.Errorf("error converting bytesToU64 :%v\n", err)
	}

	return u64
}

// ReverseBytes creates and returns a new byte slice with reversed order.
func ReverseBytes(bytes []byte) []byte {
	reversedBytes := make([]byte, len(bytes))
	for i, j := 0, len(bytes)-1; i < len(bytes); i, j = i+1, j-1 {
		reversedBytes[j] = bytes[i]
	}
	return reversedBytes
}

// BytesToU256 decodes the given bytes, representing a 256-bit unsigned integer in big-endian format,
// into a big.Int. It constructs a new big.Int with the bytes interpreted in little-endian order.
// The function returns the resulting big.Int representing the 256-bit integer.
func BytesToU256(bytes []byte) (*big.Int, error) {
	// Reverse the bytes to convert from big-endian to little-endian
	reversedBytes := ReverseBytes(bytes)

	// Create a big.Int and set its bytes representation
	u256Value := new(big.Int).SetBytes(reversedBytes)

	return u256Value, nil
}

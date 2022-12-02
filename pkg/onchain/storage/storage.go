package storage

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strconv"

	"github.com/massalabs/thyra/pkg/node"
)

func readZipFile(z *zip.File) ([]byte, error) {
	file, err := z.Open()
	if err != nil {
		return nil, fmt.Errorf("opening zip content: %w", err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("reading zip content: %w", err)
	}

	return content, nil
}

//nolint:nolintlint,ireturn,funlen
func Get(client *node.Client, address string, key string) (map[string][]byte, error) {
	chunkNumberKey := "total_chunks"
	fmt.Println(" ~ file: storage.go:33 ~ funcGet ~ chunkNumberKey", chunkNumberKey)

	keyNumber, err := node.DatastoreEntry(client, address, []byte(chunkNumberKey))
	if err != nil {
		return nil, fmt.Errorf("reading datastore entry '%s' at '%s': %w", address, chunkNumberKey, err)
	}

	fmt.Println("~ file: storage.go:36 ~ funcGet ~ keyNumber", keyNumber)

	chunkNumber := binary.LittleEndian.Uint64(keyNumber.CandidateValue)
	fmt.Println(" ~ file: storage.go:44 ~ funcGet ~ chunkNumber", chunkNumber)

	entries := []node.DatastoreEntriesKeysAsString{}

	for i := uint64(0); i < chunkNumber; i++ {
		entry := node.DatastoreEntriesKeysAsString{
			Address: address,
			Key:     []byte("massa_web_" + strconv.Itoa(int(i))),
		}
		entries = append(entries, entry)
	}

	response, err := node.DatastoreEntries(client, entries)
	if err != nil {
		return nil, fmt.Errorf("calling get_datastore_entries '%+v': %w", entries, err)
	}

	var dataStore []byte
	for i := uint64(0); i < chunkNumber; i++ {
		dataStore = append(dataStore, response[i].CandidateValue...)
	}

	zipReader, err := zip.NewReader(bytes.NewReader(dataStore), int64(len(dataStore)))
	if err != nil {
		return nil, fmt.Errorf("instantiating zip reader from decoded datastore entry '%s' at '%s': %w", address, key, err)
	}

	content := make(map[string][]byte)

	// Read all the files from zip archive
	for _, zipFile := range zipReader.File {
		rsc, err := readZipFile(zipFile)
		if err != nil {
			return nil, err
		}

		content[zipFile.Name] = rsc
	}
	fmt.Println(" ~ file: storage.go:85 ~ funcGet ~ content", content)
	return content, nil
}

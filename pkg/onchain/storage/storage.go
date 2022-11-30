package storage

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"strconv"

	"github.com/massalabs/thyra/pkg/helper"
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

	keyNumber, err := node.DatastoreEntry(client, address, helper.StringToByteArray(chunkNumberKey))
	if err != nil {
		return nil, fmt.Errorf("reading datastore entry '%s' at '%s': %w", address, chunkNumberKey, err)
	}

	chunkNumber, err := strconv.Atoi(string(keyNumber.CandidateValue))
	if err != nil {
		return nil, fmt.Errorf("error converting String to integer")
	}

	entries := []node.DatastoreEntriesKeysAsString{}

	for i := 0; i < chunkNumber; i++ {
		entry := node.DatastoreEntriesKeysAsString{
			Address: address,
			Key:     helper.StringToByteArray("massa_web_" + strconv.Itoa(i)),
		}
		entries = append(entries, entry)
	}

	response, err := node.DatastoreEntries(client, entries)
	if err != nil {
		return nil, fmt.Errorf("calling get_datastore_entries '%+v': %w", entries, err)
	}

	dataStore := ""
	for i := 0; i < chunkNumber; i++ {
		dataStore += string(response[i].CandidateValue)
	}

	b64, err := base64.StdEncoding.DecodeString(dataStore)
	if err != nil {
		return nil, fmt.Errorf("base64 decoding datastore entry '%s' at '%s': %w", address, key, err)
	}

	zipReader, err := zip.NewReader(bytes.NewReader(b64), int64(len(b64)))
	if err != nil {
		return nil, fmt.Errorf("instanciating zip reader from decoded datastore entry '%s' at '%s': %w", address, key, err)
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

	return content, nil
}

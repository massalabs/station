package storage

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strconv"

	"github.com/massalabs/thyra/pkg/convert"
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

/*
	This function fetch the datastore entries required to display
	a website in the browser from the website storer contract, unzip them and
	return the full unzipped website content.
	Datastore entries fetched :
	- total_chunks : Total number of chunks that are to be fetched.
	- massa_web_XXX : Keys containing the website data, with XXX being the chunk ID.
*/
//nolint:nolintlint,ireturn,funlen
func Get(client *node.Client, websiteStorerAddress string) (map[string][]byte, error) {
	chunkNumberKey := "total_chunks"

	keyNumber, err := node.DatastoreEntry(client, websiteStorerAddress, convert.StringToBytes(chunkNumberKey))
	if err != nil {
		return nil, fmt.Errorf("reading datastore entry '%s' at '%s': %w", websiteStorerAddress, chunkNumberKey, err)
	}

	chunkNumber := int(binary.LittleEndian.Uint64(keyNumber.CandidateValue))

	entries := []node.DatastoreEntriesKeysAsString{}

	for i := 0; i < chunkNumber; i++ {
		entry := node.DatastoreEntriesKeysAsString{
			Address: websiteStorerAddress,
			Key:     convert.StringToBytes("massa_web_" + strconv.Itoa(i)),
		}
		entries = append(entries, entry)
	}

	response, err := node.DatastoreEntries(client, entries)
	if err != nil {
		return nil, fmt.Errorf("calling get_datastore_entries '%+v': %w", entries, err)
	}

	var dataStore []byte
	for i := 0; i < chunkNumber; i++ {
		// content is prefixed with it's length encoded using a u32 (4 bytes).
		dataStore = append(dataStore, response[i].CandidateValue[4:]...)
	}

	zipReader, err := zip.NewReader(bytes.NewReader(dataStore), int64(len(dataStore)))
	if err != nil {
		return nil, fmt.Errorf("instantiating zip reader from decoded datastore entries '%s' at  '%w'", dataStore, err)
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

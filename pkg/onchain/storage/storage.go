package storage

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"

	"github.com/massalabs/station/pkg/cache"
	"github.com/massalabs/station/pkg/convert"
	"github.com/massalabs/station/pkg/node"
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

// Fetch website chunks.
func Fetch(client *node.Client, websiteStorerAddress string) ([]byte, error) {
	chunkNumberKey := "NB_CHUNKS"

	key := convert.ToBytes(chunkNumberKey)

	chunkNbByte, err := node.FetchDatastoreEntry(client, websiteStorerAddress, key)
	if err != nil {
		return nil, fmt.Errorf("reading datastore entry '%s' at '%s': %w", websiteStorerAddress, chunkNumberKey, err)
	}

	chunkNumber, err := convert.BytesToI32(chunkNbByte.CandidateValue)
	if err != nil {
		return nil, fmt.Errorf("converting fetched data for key '%s': %w ", chunkNumberKey, err)
	}

	keys := make([][]byte, chunkNumber)

	for i := 0; i < int(chunkNumber); i++ {
		keys[i] = convert.I32ToBytes(i)
	}

	response, err := node.ContractDatastoreEntries(client, websiteStorerAddress, keys)
	if err != nil {
		return nil, fmt.Errorf("calling get_datastore_entries '%+v': %w", keys, err)
	}

	if len(response) != int(chunkNumber) {
		return nil, fmt.Errorf("expected %d entries, got %d", chunkNumber, len(response))
	}

	var dataStore []byte

	for index := 0; index < int(chunkNumber); index++ {
		dataStore = append(dataStore, response[index].CandidateValue...)
	}

	return dataStore, nil
}

// Get tries to get the file from the cache and fallback to Fetch from the datastore.
// New files are automatically added to the cache.
// New website version at the same address are handled thanks to the LastUpdateTimestamp.
func Get(client *node.Client, websiteStorerAddress, configDir string) (map[string][]byte, error) {
	content := make(map[string][]byte)

	fileName := websiteStorerAddress

	var fileContent []byte

	var err error

	cacheManager := cache.Cache{ConfigDir: configDir}

	// we check if the website is in cache, if not we fetch it from the blockchain
	if cacheManager.IsPresent(fileName) {
		fileContent, err = cacheManager.Read(fileName)
		if err != nil {
			return nil, fmt.Errorf("reading file '%s': %w", fileName, err)
		}
	} else {
		fileContent, err = Fetch(client, websiteStorerAddress)
		if err != nil {
			return nil, fmt.Errorf("fetching website content at %s from blockchain: %w", websiteStorerAddress, err)
		}

		err = cacheManager.Save(fileName, fileContent)
		if err != nil {
			return nil, fmt.Errorf("caching %s: %w", fileName, err)
		}
	}

	zipReader, err := zip.NewReader(bytes.NewReader(fileContent), int64(len(fileContent)))
	if err != nil {
		return nil, fmt.Errorf("instantiating zip reader from '%s': %w", fileName, err)
	}

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

package storage

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strconv"

	"github.com/massalabs/thyra/pkg/cache"
	"github.com/massalabs/thyra/pkg/convert"
	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/onchain/dns"
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

// Fetch retrieves data from the blockchain ledger at the given address.
func Fetch(client *node.Client, websiteStorerAddress string) ([]byte, error) {
	chunkNumberKey := "total_chunks"

	keyNumber, err := node.DatastoreEntry(client, websiteStorerAddress, convert.StringToBytes(chunkNumberKey))
	if err != nil {
		return nil, fmt.Errorf("reading datastore entry '%s' at '%s': %w", websiteStorerAddress, chunkNumberKey, err)
	}

	chunkNumber := int(binary.LittleEndian.Uint64(keyNumber.CandidateValue))

	entries := []node.DatastoreEntriesKeys{}

	for i := 0; i < chunkNumber; i++ {
		entry := node.DatastoreEntriesKeys{
			Address: websiteStorerAddress,
			Key:     convert.StringToBytes("massa_web_" + strconv.Itoa(i)),
		}
		entries = append(entries, entry)
	}
	fmt.Println("🚀 ~ file: storage.go:54 ~ funcFetch ~ entries:", entries)

	response, err := node.DatastoreEntries(client, entries)
	if err != nil {
		return nil, fmt.Errorf("calling get_datastore_entries '%+v': %w", entries, err)
	}

	var dataStore []byte
	for i := 0; i < chunkNumber; i++ {
		// content is prefixed with it's length encoded using a u32 (4 bytes).
		dataStore = append(dataStore, response[i].CandidateValue[4:]...)
	}

	return dataStore, nil
}

// Get tries to get the file from the cache and fallback to Fetch from the datastore.
// New files are automatically added to the cache.
// New website version at the same address are handled thanks to the LastUpdateTimestamp.
func Get(client *node.Client, websiteStorerAddress string) (map[string][]byte, error) {
	content := make(map[string][]byte)

	metaData, err := dns.FetchRecordMetaData(client, websiteStorerAddress)
	if err != nil {
		return nil, fmt.Errorf("getting metadata for '%s' : %w", websiteStorerAddress, err)
	}

	lastTimestamp := metaData.LastUpdateTimestamp

	if lastTimestamp == 0 {
		lastTimestamp = metaData.CreationTimeStamp
	}

	fileName := fmt.Sprintf("%s-%d", websiteStorerAddress, lastTimestamp)

	var fileContent []byte

	// we check if the website is in cache, if not we fetch it from the blockchain
	if cache.IsPresent(fileName) {
		fileContent, err = cache.Read(fileName)
		if err != nil {
			return nil, fmt.Errorf("reading file '%s': %w", fileName, err)
		}
	} else {
		fileContent, err = Fetch(client, websiteStorerAddress)
		if err != nil {
			return nil, fmt.Errorf("fetching website content at %s from blockchain: %w", websiteStorerAddress, err)
		}

		err = cache.Save(fileName, fileContent)
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

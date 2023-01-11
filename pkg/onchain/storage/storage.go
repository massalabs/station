package storage

import (
	"archive/zip"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strconv"

	"github.com/massalabs/thyra/pkg/config"
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
	This function check the cache, fetch the datastore entries required to display
	a website in the browser from the website storer contract and save it in cache
	if not present yet, read the cache zipfile,
	unzip them and	return the full unzipped website content.
	Datastore entries fetched :
	- total_chunks : Total number of chunks that are to be fetched.
	- massa_web_XXX : Keys containing the website data, with XXX being the chunk ID.
*/
//nolint:nolintlint,ireturn,funlen
func Get(client *node.Client, websiteStorerAddress string) (map[string][]byte, error) {
	content := make(map[string][]byte)
	filepath, err := getFilePath(client, websiteStorerAddress)
	if err != nil {
		return nil, fmt.Errorf("getting file path '%s' : %w", websiteStorerAddress, err)
	}

	// we check if the website is in cache, if not we fetch it from the blockchain
	if !IsInCache(filepath) {
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

		response, err := node.DatastoreEntries(client, entries)
		if err != nil {
			return nil, fmt.Errorf("calling get_datastore_entries '%+v': %w", entries, err)
		}

		var dataStore []byte
		for i := 0; i < chunkNumber; i++ {
			// content is prefixed with it's length encoded using a u32 (4 bytes).
			dataStore = append(dataStore, response[i].CandidateValue[4:]...)
		}

		// Once we get the zip content, we save it in the cache
		saveInCache(filepath, dataStore)
	}

	// We read the website from the cache
	zipReader, err := zip.OpenReader(filepath)
	if err != nil {
		return nil, fmt.Errorf("reading zipfile from cache '%s' at  '%w'", filepath, err)
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

// Checks if the website is in cache and return true if yes.
func IsInCache(filepath string) bool {
	_, err := os.Stat(filepath)

	return !os.IsNotExist(err)
}

// Checks the timestamp ("META") key and return the uint64 UNIX Timestamp.
func getTimeStamp(client *node.Client, websiteStorerAddress string) uint64 {
	timestamp, _ := node.DatastoreEntry(client, websiteStorerAddress, convert.StringToBytes("META"))

	return convert.BytesToU64(timestamp.CandidateValue)
}

// Returns the file path in the cache, cache folder is inside your .config/thyra
// file name is the websiteStorer SC address+_+the UNIX timestamp.zip

func getFilePath(client *node.Client, websiteStorerAddress string) (string, error) {
	timestamp := strconv.Itoa(int(getTimeStamp(client, websiteStorerAddress)))
	configDir, _ := config.GetConfigDir()
	cachePath := path.Join(configDir, "websitesCache")
	_, err := os.Stat(cachePath)

	if os.IsNotExist(err) {
		err := os.MkdirAll(cachePath, os.ModePerm)
		if err != nil {
			return "", fmt.Errorf("error creating folder: %w", err)
		}
	}

	return path.Join(cachePath, websiteStorerAddress+"_"+timestamp+".zip"), nil
}

// Save the raw content as a zip file.
func saveInCache(filePath string, content []byte) {
	file, err := os.Create(filePath)
	if err != nil {
		log.Printf("error while creating file: %v\n", err)

		return
	}
	defer file.Close()

	err = os.WriteFile(filePath, content, 0600)
	if err != nil {
		log.Printf("error saving zipfile: %v\n", err)

		return
	}
}

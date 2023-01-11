package storage

import (
	"archive/zip"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strconv"

	"github.com/massalabs/thyra/pkg/config"
	"github.com/massalabs/thyra/pkg/convert"
	"github.com/massalabs/thyra/pkg/node"
)

const websitePath = "websitesCache/"

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

	content := make(map[string][]byte)
	filepath := getFilePath(client, websiteStorerAddress)
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

		saveInCache(filepath, dataStore)

	}

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

func IsInCache(filepath string) bool {
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func getTimeStamp(client *node.Client, websiteStorerAddress string) uint64 {
	timestamp, _ := node.DatastoreEntry(client, websiteStorerAddress, convert.StringToBytes("META"))
	return convert.BytesToU64(timestamp.CandidateValue)
}

func getFilePath(client *node.Client, websiteStorerAddress string) string {
	timestamp := strconv.Itoa(int(getTimeStamp(client, websiteStorerAddress)))
	configDir, _ := config.GetConfigDir()
	cachePath := path.Join(configDir, "websitesCache")
	_, err := os.Stat(cachePath)
	if os.IsNotExist(err) {
		// folder does not exist
		// create the folder
		err := os.MkdirAll(cachePath, os.ModePerm)
		if err != nil {
			fmt.Printf("Error while creating folder: %v\n", err)
		}
	} else if err != nil {
		fmt.Printf("Error while checking if folder exists: %v\n", err)
	}

	return path.Join(cachePath, websiteStorerAddress+"_"+timestamp+".zip")
}

func saveInCache(filePath string, content []byte) {
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Error while creating file: %v\n", err)
		return
	}
	defer file.Close()

	err = ioutil.WriteFile(filePath, content, 0644)
	if err != nil {
		fmt.Printf("error saving zipfile: %w", err)
		return
	}
}

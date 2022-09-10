package storage

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/massalabs/thyra/pkg/node"
)

func readZipFile(z *zip.File) ([]byte, error) {
	file, err := z.Open()
	if err != nil {
		return nil, fmt.Errorf("opening zip content: %w", err)
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("reading zip content: %w", err)
	}

	return content, nil
}

func Get(client *node.Client, address string, key string) (map[string][]byte, error) {
	entry, err := node.DatastoreEntry(client, address, key)
	if err != nil {
		return nil, fmt.Errorf("reading datastore entry '%s' at '%s': %w", address, key, err)
	}

	if len(entry.CandidateValue) == 0 {
		return nil, errors.New("no data in candidate value key")
	}

	b64, err := base64.StdEncoding.DecodeString(string(entry.CandidateValue))
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

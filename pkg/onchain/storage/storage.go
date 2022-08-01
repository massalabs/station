package storage

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"errors"
	"io/ioutil"

	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/node/getters"
)

func readZipFile(z *zip.File) ([]byte, error) {
	file, err := z.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return ioutil.ReadAll(file)
}

func Get(client *node.Client, address string, key string) (map[string][]byte, error) {
	entry, err := getters.DatastoreEntry(client, address, key)
	if err != nil {
		return nil, err
	}

	if len(entry.CandidateValue) == 0 {
		return nil, errors.New("no data in candidate value key")
	}

	b64, err := base64.StdEncoding.DecodeString(string(entry.CandidateValue))
	if err != nil {
		return nil, err
	}

	zipReader, err := zip.NewReader(bytes.NewReader(b64), int64(len(b64)))
	if err != nil {
		return nil, err
	}

	m := make(map[string][]byte)

	// Read all the files from zip archive
	for _, zipFile := range zipReader.File {
		uzb, err := readZipFile(zipFile)
		if err != nil {
			return nil, err
		}

		m[zipFile.Name] = uzb
	}

	return m, nil
}

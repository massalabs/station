package storage

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type ledgerInfo struct {
	Datastore map[string][]byte
}

type ledger struct {
	Info ledgerInfo `json:"candidate_sce_ledger_info"`
}

type jsonRPCResponse struct {
	Result *[]ledger
	Error  *jsonRPCError
}

type jsonRPCError struct {
	Code    int64
	Message string
}

func readZipFile(z *zip.File) ([]byte, error) {
	f, err := z.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

func Get(address string, key string) (map[string][]byte, error) {
	body := []byte(`{
                "jsonrpc": "2.0",
                "method": "get_addresses",
                "id": 111,
                "params": [["`)
	tail := []byte(`"]]}`)

	body = append(body, address...)
	body = append(body, tail...)

	req, err := http.NewRequest("POST", "https://test.massa.net/api/v2", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var j jsonRPCResponse
	err = json.Unmarshal(b, &j)

	if err != nil {
		return nil, err
	}

	if j.Result == nil || len(*j.Result) == 0 {
		return nil, errors.New("no ledger")
	}

	s, ok := (*j.Result)[0].Info.Datastore[key]

	if !ok {
		return nil, errors.New("no web site in datastore")
	}

	b64, err := base64.StdEncoding.DecodeString(string(s))
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

package website

import (
	"encoding/json"
	"strconv"

	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/node/base58"
	"github.com/massalabs/thyra/pkg/onchain"
	"github.com/massalabs/thyra/pkg/onchain/dns"
	"github.com/massalabs/thyra/pkg/sc"
	"github.com/massalabs/thyra/pkg/wallet"
)

func PrepareForUpload(url string, wallet *wallet.Wallet) (string, error) {
	client := node.NewDefaultClient()

	// Prepare address to webstorage.
	scAddress, err := onchain.DeploySC(client, *wallet, []byte(sc.WebsiteStorer))
	if err != nil {
		return "", err
	}

	// Set DNS.
	_, err = dns.SetRecord(client, *wallet, url, scAddress)
	if err != nil {
		return "", err
	}

	return scAddress, nil
}

type UploadWebsiteParam struct {
	Data        string `json:"data"`
	TotalChunks string `json:"totalChunks"`
}

func Upload(atAddress string, content string, wallet *wallet.Wallet) (string, error) {
	const blockLength = 260000
	client := node.NewDefaultClient()

	addr, _, err := base58.VersionedCheckDecode(atAddress[1:])
	if err != nil {
		return "", err
	}
	blocks := chunk(content, blockLength)

	if len(blocks) == 1 {
		_, err = uploadLight(client, addr, content, wallet)
	} else {
		_, err = uploadHeavy(client, addr, blocks, wallet)
	}

	if err != nil {
		return "", err
	}
	return "1", nil
}

func Status(websiteCreator string) (*models.UploadState, error) {
	client := node.NewDefaultClient()

	actualChunkEntry, err := node.DatastoreEntry(client, websiteCreator, "actual_chunk")
	if err != nil {
		return nil, err
	}

	totalChunksEntry, err := node.DatastoreEntry(client, websiteCreator, "total_chunks")
	if err != nil {
		return nil, err
	}

	actualChunk, err := strconv.Atoi(string(actualChunkEntry.CandidateValue))
	if err != nil {
		actualChunk = -1
	}
	totalChunks, err := strconv.Atoi(string(totalChunksEntry.CandidateValue))
	if err != nil {
		totalChunks = -1
	}

	var status string
	if actualChunk == -1 && totalChunks == -1 {
		status = "NOT_STARTED"
	} else if actualChunk == totalChunks {
		status = "COMPLETED"
	} else {
		status = "IN_PROGRESS"
	}

	result := models.UploadState{LastChunk: int64(actualChunk), TotalChunk: int64(totalChunks), Status: status}

	return &result, nil
}

func uploadLight(client *node.Client, addr []byte, content string, wallet *wallet.Wallet) (string, error) {
	param, err := json.Marshal(UploadWebsiteParam{
		Data: content,
	})
	if err != nil {
		return "", err
	}
	op, err := onchain.CallFunction(client, *wallet, addr, "initializeWebsite", param)
	if err != nil {
		return "", err
	}
	return op, nil
}

func uploadHeavy(client *node.Client, addr []byte, chunks []string, wallet *wallet.Wallet) (string, error) {
	op := ""
	param, err := json.Marshal(UploadWebsiteParam{
		Data:        chunks[0],
		TotalChunks: strconv.FormatInt(int64(len(chunks)), 10),
	})
	if err != nil {
		return "", err
	}
	_, err = onchain.CallFunction(client, *wallet, addr, "initializeWebsite", param)
	if err != nil {
		return "", err
	}
	for i := 1; i < len(chunks); i++ {

		param, err = json.Marshal(UploadWebsiteParam{
			Data: chunks[i],
		})
		if err != nil {
			return "", err
		}
		op, err = onchain.CallFunction(client, *wallet, addr, "appendBytesToWebsite", param)
		if err != nil {
			return "", err
		}
	}
	return op, nil
}

func chunk(s string, chunkSize int) []string {

	chunks := []string{}
	counter := 0

	chunkNumber := len(s)/chunkSize + 1
	for i := 1; i < chunkNumber; i++ {
		counter += chunkSize
		chunks = append(chunks, s[(i-1)*chunkSize:(i)*chunkSize])
	}
	chunks = append(chunks, s[(chunkNumber-1)*chunkSize:])

	return chunks

}

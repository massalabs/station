package website

import (
	"embed"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/node/base58"
	"github.com/massalabs/thyra/pkg/onchain"
	"github.com/massalabs/thyra/pkg/onchain/dns"
	"github.com/massalabs/thyra/pkg/wallet"
)

//go:embed sc
var content embed.FS

func PrepareForUpload(url string, wallet *wallet.Wallet) (string, error) {
	client := node.NewDefaultClient()

	basePath := "sc/"

	websiteStorer, err := content.ReadFile(basePath + "websiteStorer.wasm")
	if err != nil {
		return "", fmt.Errorf("SC file not retrieved: %w", err)
	}

	// Prepare address to webstorage.
	scAddress, err := onchain.DeploySC(client, *wallet, websiteStorer)
	if err != nil {
		return "", fmt.Errorf("deploying webstorage SC: %w", err)
	}

	// Set DNS.
	_, err = dns.SetRecord(client, *wallet, url, scAddress)
	if err != nil {
		return "", fmt.Errorf("adding DNS record '%s' => '%s': %w", url, scAddress, err)
	}

	return scAddress, nil
}

type UploadWebsiteParam struct {
	Data        string `json:"data"`
	TotalChunks string `json:"total_chunks"`
}

func Upload(atAddress string, content string, wallet *wallet.Wallet) (string, error) {
	const blockLength = 260000

	client := node.NewDefaultClient()

	addr, _, err := base58.VersionedCheckDecode(atAddress[1:])
	if err != nil {
		return "", fmt.Errorf("decoding address '%s': %w", atAddress[1:], err)
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

func uploadLight(client *node.Client, addr []byte, content string, wallet *wallet.Wallet) (string, error) {
	param, err := json.Marshal(UploadWebsiteParam{
		Data:        content,
		TotalChunks: "1",
	})
	if err != nil {
		return "", fmt.Errorf("marshaling '%s': %w", param, err)
	}

	op, err := onchain.CallFunction(client, *wallet, addr, "initializeWebsite", param)
	if err != nil {
		return "", fmt.Errorf("calling initializeWebsite at '%s': %w", addr, err)
	}

	return op, nil
}

func uploadHeavy(client *node.Client, addr []byte, chunks []string, wallet *wallet.Wallet) (string, error) {
	const baseFormatInt = 10

	param, err := json.Marshal(UploadWebsiteParam{
		Data:        chunks[0],
		TotalChunks: strconv.FormatInt(int64(len(chunks)), baseFormatInt),
	})
	if err != nil {
		return "", fmt.Errorf("marshaling '%s': %w", param, err)
	}

	_, err = onchain.CallFunction(client, *wallet, addr, "initializeWebsite", param)
	if err != nil {
		return "", fmt.Errorf("calling initializeWebsite at '%s': %w", addr, err)
	}

	var opID string

	for index := 1; index < len(chunks); index++ {
		//nolint:exhaustruct
		param, err = json.Marshal(UploadWebsiteParam{
			Data: chunks[index],
		})
		if err != nil {
			//nolint:exhaustruct
			return "", fmt.Errorf("marshaling '%s': %w", UploadWebsiteParam{Data: chunks[index]}, err)
		}

		opID, err = onchain.CallFunction(client, *wallet, addr, "appendBytesToWebsite", param)
		if err != nil {
			return "", fmt.Errorf("calling initializeWebsite at '%s': %w", addr, err)
		}
	}

	return opID, nil
}

func chunk(data string, chunkSize int) []string {
	counter := 0

	chunkNumber := len(data)/chunkSize + 1

	var chunks []string

	for i := 1; i < chunkNumber; i++ {
		counter += chunkSize
		chunks = append(chunks, data[(i-1)*chunkSize:(i)*chunkSize])
	}

	chunks = append(chunks, data[(chunkNumber-1)*chunkSize:])

	return chunks
}

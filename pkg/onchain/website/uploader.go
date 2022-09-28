package website

import (
	"encoding/json"
	"strconv"

	"fmt"

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
	Data    string `json:"data"`
	ChunkID string `json:"chunkID"`
}

type WebsiteInitialisationParams struct {
	TotalChunks string `json:"totalChunks"`
}

func Upload(atAddress string, content string, wallet *wallet.Wallet) (string, error) {
	const blockLength = 260000

	client := node.NewDefaultClient()

	addr, _, err := base58.VersionedCheckDecode(atAddress[1:])
	if err != nil {
		return "", fmt.Errorf("decoding address '%s': %w", atAddress[1:], err)
	}

	blocks := chunk(content, blockLength)
	_, err = uploadHeavy(client, addr, blocks, wallet)
	if err != nil {
		return "", err
	}

	return "Website deployed", nil
}

func uploadHeavy(client *node.Client, addr []byte, chunks []string, wallet *wallet.Wallet) (string, error) {
	paramInit, err := json.Marshal(WebsiteInitialisationParams{
		TotalChunks: strconv.Itoa(len(chunks)),
	})

	if err != nil {
		return "", fmt.Errorf("marshaling '%s': %w", UploadWebsiteParam{Data: chunks[0]}, err)
	}

	_, err = onchain.CallFunction(client, *wallet, addr, "initializeWebsite", paramInit)
	if err != nil {
		return "", fmt.Errorf("calling initializeWebsite at '%s': %w", addr, err)
	}

	var opID string

	for i := 0; i < len(chunks); i++ {

		param, err := json.Marshal(UploadWebsiteParam{
			Data:    chunks[i],
			ChunkID: strconv.Itoa(i),
		})

		if err != nil {
			return "", fmt.Errorf("marshaling '%s': %w", UploadWebsiteParam{Data: chunks[i]}, err)
		}

		opID, err = onchain.CallFunctionUnwaited(client, *wallet, addr, "appendBytesToWebsite", param)

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

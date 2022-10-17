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

const baseOffset = 5

const multiplicator = 2

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

type InitialisationParams struct {
	TotalChunks string `json:"total_chunks"`
}
type AppendParams struct {
	Data    string `json:"data"`
	ChunkID string `json:"chunk_id"`
}

func Upload(atAddress string, content string, wallet *wallet.Wallet) (*[]string, error) {
	const blockLength = 260000

	client := node.NewDefaultClient()

	addr, _, err := base58.VersionedCheckDecode(atAddress[1:])
	if err != nil {
		return nil, fmt.Errorf("decoding address '%s': %w", atAddress[1:], err)
	}

	blocks := chunk(content, blockLength)

	operations, err := upload(client, addr, blocks, wallet)
	if err != nil {
		return nil, err
	}

	return operations, nil
}

func upload(client *node.Client, addr []byte, chunks []string, wallet *wallet.Wallet) (*[]string, error) {
	operations := make([]string, len(chunks)+1)

	paramInit, err := json.Marshal(InitialisationParams{
		TotalChunks: strconv.Itoa(len(chunks)),
	})
	if err != nil {
		return nil, fmt.Errorf("marshaling '%s': %w", InitialisationParams{TotalChunks: strconv.Itoa(len(chunks))}, err)
	}

	opID, err := onchain.CallFunction(client, *wallet, addr, "initializeWebsite", paramInit)
	if err != nil {
		return nil, fmt.Errorf("calling initializeWebsite at '%s': %w", addr, err)
	}

	operations[0] = opID

	for index := 0; index < len(chunks); index++ {
		param, err := json.Marshal(AppendParams{
			Data:    chunks[index],
			ChunkID: strconv.Itoa(index),
		})
		if err != nil {
			return nil,
				fmt.Errorf("marshaling '%s': %w", AppendParams{Data: chunks[index], ChunkID: strconv.Itoa(index)}, err)
		}

		//nolint:lll
		opID, err = onchain.CallFunctionUnwaited(client, *wallet, baseOffset+uint64(index)*multiplicator, addr, "appendBytesToWebsite", param)
		if err != nil {
			return nil, fmt.Errorf("calling initializeWebsite at '%s': %w", addr, err)
		}

		operations[index+1] = opID
	}

	return &operations, nil
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

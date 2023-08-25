package website

import (
	"embed"
	"fmt"
	"strconv"
	"strings"

	"github.com/massalabs/station/int/config"
	"github.com/massalabs/station/pkg/convert"
	"github.com/massalabs/station/pkg/node"
	sendOperation "github.com/massalabs/station/pkg/node/sendoperation"
	"github.com/massalabs/station/pkg/node/sendoperation/signer"
	"github.com/massalabs/station/pkg/onchain"
	"github.com/massalabs/station/pkg/onchain/dns"
)

//go:embed sc
var content embed.FS

const baseOffset = 5

const blockLength = 260000

// function calculating the max expiry period, this calculation is empiric

func maxExpiryPeriod(index int) uint64 {
	return baseOffset + uint64(index)*2
}

func PrepareForUpload(
	config config.AppConfig,
	url string,
	description string,
	nickname string,
) (string, string, error) {
	client := node.NewClient(config.NodeURL)

	basePath := "sc/"

	websiteStorer, err := content.ReadFile(basePath + "websiteStorer.wasm")
	if err != nil {
		return "", "", fmt.Errorf("SC file not retrieved: %w", err)
	}
	// Prepare address to webstorage.
	operationWithEventResponse, err := onchain.DeploySC(
		client,
		nickname,
		sendOperation.DefaultGasLimit,
		sendOperation.AccountCreationStorageCost,
		sendOperation.DefaultFee,
		sendOperation.DefaultSlotsDuration,
		websiteStorer,
		nil,
		sendOperation.OperationBatch{NewBatch: true, CorrelationID: ""},
		&signer.WalletPlugin{},
	)
	if err != nil {
		return "", "", fmt.Errorf("deploying webstorage SC: %w", err)
	}

	scAddress := strings.Split(operationWithEventResponse.Event, ":")[1]

	// Set DNS.
	_, err = dns.SetRecord(
		config,
		client,
		nickname,
		url,
		description,
		scAddress,
		sendOperation.OperationBatch{
			NewBatch:      false,
			CorrelationID: operationWithEventResponse.OperationResponse.CorrelationID,
		},
	)
	if err != nil {
		return "", "", fmt.Errorf("adding DNS record '%s' => '%s': %w", url, scAddress, err)
	}

	return scAddress, operationWithEventResponse.OperationResponse.CorrelationID, nil
}

func Upload(
	config config.AppConfig,
	atAddress string,
	content []byte,
	nickname string,
	operationBatch sendOperation.OperationBatch,
) ([]string, error) {
	client := node.NewClient(config.NodeURL)

	blocks := chunk(content, blockLength)

	operations, err := upload(client, atAddress, blocks, nickname, operationBatch)
	if err != nil {
		return nil, err
	}

	return operations, nil
}

//nolint:funlen
func upload(
	client *node.Client,
	addr string,
	chunks [][]byte,
	nickname string,
	operationBatch sendOperation.OperationBatch,
) ([]string, error) {
	nbChunks := len(chunks)
	operations := make([]string, nbChunks+1)

	operationWithEventResponse, err := onchain.CallFunction(
		client,
		nickname,
		addr,
		"initializeWebsite",
		convert.U64ToBytes(nbChunks),
		sendOperation.OneMassa, // To be updated when readonly is working to estimate storage cost
		sendOperation.DefaultSlotsDuration,
		false,
		operationBatch,
		&signer.WalletPlugin{},
	)
	if err != nil {
		return nil, fmt.Errorf("calling initializeWebsite at '%s': %w", addr, err)
	}

	operations[0] = operationWithEventResponse.OperationResponse.OperationID

	for index := 0; index < nbChunks; index++ {
		chunkSize := len(chunks[index])
		// Chunk ID encoding
		params := convert.U64ToBytes(index)

		// Chunk data length encoding
		params = append(params, convert.U32ToBytes(chunkSize)...)

		// Chunk data encoding
		params = append(params, chunks[index]...)

		operationResponse, err := onchain.CallFunction(
			client,
			nickname,
			addr,
			"appendBytesToWebsite",
			params,
			sendOperation.StorageCostPerByte*uint64(chunkSize),
			maxExpiryPeriod(index),
			false,
			sendOperation.OperationBatch{
				NewBatch:      false,
				CorrelationID: operationWithEventResponse.OperationResponse.CorrelationID,
			},
			&signer.WalletPlugin{},
		)
		if err != nil {
			return nil, fmt.Errorf("calling appendBytesToWebsite at '%s': %w", addr, err)
		}

		operations[index+1] = operationResponse.OperationResponse.OperationID
	}

	return operations, nil
}

func UploadMissedChunks(
	config config.AppConfig,
	atAddress string,
	content []byte,
	nickname string,
	missedChunks string,
	operationBatch sendOperation.OperationBatch,
) ([]string, error) {
	client := node.NewClient(config.NodeURL)

	blocks := chunk(content, blockLength)

	operations, err := uploadMissedChunks(
		client,
		atAddress,
		blocks,
		missedChunks,
		nickname,
		operationBatch,
	)
	if err != nil {
		return nil, err
	}

	return operations, nil
}

func uploadMissedChunks(
	client *node.Client,
	addr string,
	chunks [][]byte,
	missedChunks string,
	nickname string,
	operationBatch sendOperation.OperationBatch,
) ([]string, error) {
	operations := make([]string, len(chunks)+1)
	arrMissedChunks := strings.Split(missedChunks, ",")

	for index := 0; index < len(arrMissedChunks); index++ {
		chunkID, err := strconv.Atoi(arrMissedChunks[index])
		if err != nil {
			return nil, fmt.Errorf("error while converting chunk ID")
		}

		chunkSize := len(chunks[chunkID])

		params := convert.U64ToBytes(chunkID)

		// Chunk data length encoding
		//nolint:ineffassign,nolintlint
		params = append(params, convert.U32ToBytes(chunkSize)...)
		// Chunk data encoding
		//nolint:ineffassign,nolintlint
		params = append(params, chunks[chunkID]...)

		operationResponse, err := onchain.CallFunction(
			client,
			nickname,
			addr,
			"appendBytesToWebsite",
			params,
			sendOperation.StorageCostPerByte*uint64(chunkSize),
			maxExpiryPeriod(index),
			false,
			operationBatch,
			&signer.WalletPlugin{},
		)
		if err != nil {
			return nil, fmt.Errorf("calling appendBytesToWebsite at '%s': %w", addr, err)
		}

		operations[index] = operationResponse.OperationResponse.OperationID
		operationBatch.NewBatch = false
		operationBatch.CorrelationID = operationResponse.OperationResponse.CorrelationID
	}

	return operations, nil
}

func chunk(data []byte, chunkSize int) [][]byte {
	counter := 0

	chunkNumber := len(data)/chunkSize + 1

	var chunks [][]byte

	for i := 1; i < chunkNumber; i++ {
		counter += chunkSize
		chunks = append(chunks, data[(i-1)*chunkSize:(i)*chunkSize])
	}

	chunks = append(chunks, data[(chunkNumber-1)*chunkSize:])

	return chunks
}

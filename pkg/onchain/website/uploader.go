package website

import (
	"embed"
	"fmt"
	"strconv"
	"strings"

	"github.com/massalabs/thyra/pkg/convert"
	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/node/base58"
	sendOperation "github.com/massalabs/thyra/pkg/node/sendoperation"
	"github.com/massalabs/thyra/pkg/onchain"
	"github.com/massalabs/thyra/pkg/onchain/dns"
)

//go:embed sc
var content embed.FS

const baseOffset = 5

const blockLength = 260000

// function calculating the max expiry period, this calculation is empiric

func maxExpiryPeriod(index int) uint64 {
	return baseOffset + uint64(index)*2
}

func PrepareForUpload(url string, nickname string) (string, string, error) {
	client := node.NewDefaultClient()

	basePath := "sc/"

	websiteStorer, err := content.ReadFile(basePath + "websiteStorer.wasm")
	if err != nil {
		return "", "", fmt.Errorf("SC file not retrieved: %w", err)
	}

	// Prepare address to webstorage.
	operationWithEventResponse, err := onchain.DeploySC(
		client,
		nickname,
		sendOperation.DefaultGazLimit,
		sendOperation.NoCoin,
		sendOperation.NoFee,
		sendOperation.DefaultSlotsDuration,
		websiteStorer,
		nil,
		sendOperation.OperationBatch{NewBatch: true, CorrelationID: ""},
	)
	if err != nil {
		return "", "", fmt.Errorf("deploying webstorage SC: %w", err)
	}

	scAddress := strings.Split(operationWithEventResponse.Event, ":")[1]

	// Set DNS.
	_, err = dns.SetRecord(
		client,
		nickname,
		url,
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
	atAddress string,
	content []byte,
	nickname string,
	operationBatch sendOperation.OperationBatch,
) ([]string, error) {
	client := node.NewDefaultClient()

	addr, _, err := base58.VersionedCheckDecode(atAddress[2:])
	if err != nil {
		return nil, fmt.Errorf("decoding address '%s': %w", atAddress[2:], err)
	}

	blocks := chunk(content, blockLength)

	operations, err := upload(client, addr, blocks, nickname, operationBatch)
	if err != nil {
		return nil, err
	}

	return operations, nil
}

func upload(
	client *node.Client,
	addr []byte,
	chunks [][]byte,
	nickname string,
	operationBatch sendOperation.OperationBatch,
) ([]string, error) {
	operations := make([]string, len(chunks)+1)

	operationWithEventResponse, err := onchain.CallFunction(
		client,
		nickname,
		addr,
		"initializeWebsite",
		convert.U64ToBytes(len(chunks)),
		sendOperation.OneMassa,
		operationBatch,
	)
	if err != nil {
		return nil, fmt.Errorf("calling initializeWebsite at '%s': %w", addr, err)
	}

	operations[0] = operationWithEventResponse.OperationResponse.OperationID

	for index := 0; index < len(chunks); index++ {
		// Chunk ID encoding
		params := convert.U64ToBytes(index)

		// Chunk data length encoding
		params = append(params, convert.U32ToBytes(len(chunks[index]))...)

		// Chunk data encoding
		params = append(params, chunks[index]...)

		operationResponse, err := onchain.CallFunctionUnwaited(
			client,
			nickname,
			maxExpiryPeriod(index),
			addr,
			"appendBytesToWebsite",
			params,
			operationBatch,
		)
		if err != nil {
			return nil, fmt.Errorf("calling appendBytesToWebsite at '%s': %w", addr, err)
		}

		operations[index+1] = operationResponse.OperationID
	}

	return operations, nil
}

func UploadMissedChunks(atAddress string, content []byte, nickname string, missedChunks string) ([]string, error) {
	client := node.NewDefaultClient()

	addr, _, err := base58.VersionedCheckDecode(atAddress[2:])
	if err != nil {
		return nil, fmt.Errorf("decoding address '%s': %w", atAddress[2:], err)
	}

	blocks := chunk(content, blockLength)

	operations, err := uploadMissedChunks(
		client,
		addr,
		blocks,
		missedChunks,
		nickname,
		sendOperation.OperationBatch{NewBatch: false, CorrelationID: ""},
	)
	if err != nil {
		return nil, err
	}

	return operations, nil
}

func uploadMissedChunks(
	client *node.Client,
	addr []byte,
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

		params := convert.U64ToBytes(chunkID)

		// Chunk data length encoding
		//nolint:ineffassign,nolintlint
		params = append(params, convert.U32ToBytes(len(chunks[chunkID]))...)
		// Chunk data encoding
		//nolint:ineffassign,nolintlint
		params = append(params, chunks[chunkID]...)

		operationResponse, err := onchain.CallFunctionUnwaited(
			client,
			nickname,
			maxExpiryPeriod(index),
			addr,
			"appendBytesToWebsite",
			params,
			operationBatch,
		)
		if err != nil {
			return nil, fmt.Errorf("calling appendBytesToWebsite at '%s': %w", addr, err)
		}

		operations[index] = operationResponse.OperationID
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

package website

import (
	"embed"
	"fmt"
	"strconv"
	"strings"

	"github.com/massalabs/station/int/config"
	"github.com/massalabs/station/pkg/convert"
	"github.com/massalabs/station/pkg/logger"
	"github.com/massalabs/station/pkg/node"
	sendOperation "github.com/massalabs/station/pkg/node/sendoperation"
	"github.com/massalabs/station/pkg/node/sendoperation/signer"
	"github.com/massalabs/station/pkg/onchain"
	"github.com/massalabs/station/pkg/onchain/dns"
)

//go:embed sc
var content embed.FS

const blockLength = 260000

//nolint:funlen
func PrepareForUpload(
	config config.AppConfig,
	url string,
	description string,
	nickname string,
) (string, string, error) {
	client := node.NewClient(config.NodeURL)

	basePath := "sc/"

	websiteDeployer, err := content.ReadFile(basePath + "websiteDeployer.wasm")
	if err != nil {
		return "", "", fmt.Errorf("websiteDeployer contract not found: %w", err)
	}

	websiteStorer, err := content.ReadFile(basePath + "websiteStorer.wasm")
	if err != nil {
		return "", "", fmt.Errorf("websiteStorer contract not found: %w", err)
	}

	// webSiteInitCost correspond to the cost of owner initialization
	//nolint:lll,gomnd
	webSiteInitCost := sendOperation.StorageKeyCreationCost + 53 /*owner Addr max byteLenght*/ *sendOperation.StorageCostPerByte
	deployCost := sendOperation.AccountCreationStorageCost + (len(websiteStorer) * sendOperation.StorageCostPerByte)
	totalStorageCost := webSiteInitCost + deployCost

	logger.Debug("Website deployment cost estimation: ", totalStorageCost)

	// Prepare address to webstorage.
	operationResponse, events, err := onchain.DeploySC(
		client,
		nickname,
		sendOperation.DefaultGasLimit,
		uint64(totalStorageCost),
		sendOperation.DefaultFee,
		sendOperation.DefaultExpiryInSlot,
		websiteDeployer,
		nil,
		sendOperation.OperationBatch{NewBatch: true, CorrelationID: ""},
		&signer.WalletPlugin{},
	)
	if err != nil {
		return "", "", fmt.Errorf("deploying webstorage SC: %w", err)
	}

	scAddress, found := onchain.FindDeployedAddress(events)
	if !found {
		return "", "",
			fmt.Errorf("unable to retrieve SC address from deployment event of Opid: %s", operationResponse.OperationID)
	}

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
			CorrelationID: operationResponse.CorrelationID,
		},
	)
	if err != nil {
		return "", "", fmt.Errorf("adding DNS record '%s' => '%s': %w", url, scAddress, err)
	}

	return scAddress, operationResponse.CorrelationID, nil
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

func upload(
	client *node.Client,
	addr string,
	chunks [][]byte,
	nickname string,
	operationBatch sendOperation.OperationBatch,
) ([]string, error) {
	nbChunks := len(chunks)
	operations := make([]string, nbChunks)

	for chunkIndex := 0; chunkIndex < nbChunks; chunkIndex++ {
		chunkSize := len(chunks[chunkIndex])
		// Chunk ID encoding
		params := convert.I32ToBytes(chunkIndex)

		// Chunk data length encoding
		params = append(params, convert.U32ToBytes(chunkSize)...)

		// Chunk data encoding
		params = append(params, chunks[chunkIndex]...)

		// Total cost is the cost for storage bytes plus a fixed cost for key creation
		uploadCost := sendOperation.StorageKeyCreationCost + sendOperation.StorageCostPerByte*chunkSize

		if chunkIndex == 0 {
			// if chunkID == 0, we need to add the cost of the key creation for the NB_CHUNKS key
			chunkKeyCost := sendOperation.StorageKeyCreationCost + sendOperation.StorageCostPerByte*convert.BytesPerUint32
			uploadCost += chunkKeyCost
		}

		logger.Debug("Website chunk upload cost estimation: ", uploadCost)

		operationResponse, err := onchain.CallFunction(
			client,
			nickname,
			addr,
			"appendBytesToWebsite",
			params,
			sendOperation.DefaultFee,
			sendOperation.DefaultGasLimit,
			uint64(uploadCost),
			sendOperation.DefaultExpiryInSlot,
			false,
			operationBatch,
			&signer.WalletPlugin{},
		)
		if err != nil {
			return nil, fmt.Errorf("calling appendBytesToWebsite at '%s': %w", addr, err)
		}

		operations[chunkIndex] = operationResponse.OperationResponse.OperationID
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

		params := convert.I32ToBytes(chunkID)

		// Chunk data length encoding
		//nolint:ineffassign,nolintlint
		params = append(params, convert.U32ToBytes(chunkSize)...)
		// Chunk data encoding
		//nolint:ineffassign,nolintlint
		params = append(params, chunks[chunkID]...)

		// Total cost is the cost for storage bytes plus a fixed cost for key creation
		uploadCost := sendOperation.StorageKeyCreationCost + sendOperation.StorageCostPerByte*chunkSize

		if chunkID == 0 {
			// if chunkID == 0, we may need to add the cost of the key creation for the NB_CHUNKS key
			chunkKeyCost := sendOperation.StorageKeyCreationCost + sendOperation.StorageCostPerByte*convert.BytesPerUint32
			uploadCost += chunkKeyCost
		}

		operationResponse, err := onchain.CallFunction(
			client,
			nickname,
			addr,
			"appendBytesToWebsite",
			params,
			sendOperation.DefaultFee,
			sendOperation.DefaultGasLimit,
			uint64(uploadCost),
			sendOperation.DefaultExpiryInSlot,
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

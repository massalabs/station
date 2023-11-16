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

const (
	blockLength = 260000
	nbChunkKey  = "NB_CHUNKS"
	ownerKey    = "OWNER"
)

//nolint:funlen,cyclop
func PrepareForUpload(
	network config.NetworkInfos,
	url string,
	description string,
	nickname string,
) (string, string, error) {
	client := node.NewClient(network.NodeURL)

	versionFloat, err := strconv.ParseFloat(network.Version, 64)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse nodeversion %s: %w", network.Version, err)
	}

	scSuffix := ".wasm"

	//nolint:gomnd
	if versionFloat >= 26 {
		// current version is higher than 26.0
		scSuffix = "_v26.wasm"
	}

	basePath := "sc/"

	websiteDeployer, err := content.ReadFile(basePath + "websiteDeployer" + scSuffix)
	if err != nil {
		return "", "", fmt.Errorf("websiteDeployer contract not found: %w", err)
	}

	websiteStorer, err := content.ReadFile(basePath + "websiteStorer" + scSuffix)
	if err != nil {
		return "", "", fmt.Errorf("websiteStorer contract not found: %w", err)
	}

	// webSiteInitCost correspond to the cost of owner initialization
	//nolint:lll,gomnd
	webSiteInitCost, err := sendOperation.StorageCostForEntry(network.Version, len([]byte(ownerKey)), 53 /*owner Addr max byteLenght*/)
	if err != nil {
		return "", "", fmt.Errorf("unable to compute storage cost for website init: %w", err)
	}

	deployCost, err := sendOperation.StorageCostForEntry(network.Version, 0, len(websiteStorer))
	if err != nil {
		return "", "", fmt.Errorf("unable to compute storage cost for website deployment: %w", err)
	}

	accountCreationCost, err := sendOperation.AccountCreationStorageCost(network.Version)
	if err != nil {
		return "", "", fmt.Errorf("unable to compute storage cost for account creation: %w", err)
	}

	totalStorageCost := webSiteInitCost + deployCost + accountCreationCost

	logger.Debug("Website deployment cost estimation: ", totalStorageCost)

	// Prepare address to webstorage.
	operationResponse, events, err := onchain.DeploySC(
		client,
		nickname,
		sendOperation.DefaultGasLimitExecuteSC,
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
		network,
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
	network config.NetworkInfos,
	atAddress string,
	content []byte,
	nickname string,
	operationBatch sendOperation.OperationBatch,
) ([]string, error) {
	blocks := chunk(content, blockLength)

	operations, err := upload(network, atAddress, blocks, nickname, operationBatch)
	if err != nil {
		return nil, err
	}

	return operations, nil
}

//nolint:funlen
func upload(
	network config.NetworkInfos,
	addr string,
	chunks [][]byte,
	nickname string,
	operationBatch sendOperation.OperationBatch,
) ([]string, error) {
	nbChunks := len(chunks)
	operations := make([]string, nbChunks)

	client := node.NewClient(network.NodeURL)

	for chunkIndex := 0; chunkIndex < nbChunks; chunkIndex++ {
		chunkSize := len(chunks[chunkIndex])
		// Chunk ID encoding
		params := convert.I32ToBytes(chunkIndex)

		// Chunk data length encoding
		params = append(params, convert.U32ToBytes(chunkSize)...)

		// Chunk data encoding
		params = append(params, chunks[chunkIndex]...)

		uploadCost, err := sendOperation.StorageCostForEntry(network.Version, convert.BytesPerUint32, chunkSize)
		if err != nil {
			return nil, fmt.Errorf("unable to compute storage cost chunk upload: %w", err)
		}

		if chunkIndex == 0 {
			// if chunkID == 0, we need to add the cost of the key creation for the NB_CHUNKS key
			chunkKeyCost, err := sendOperation.StorageCostForEntry(
				network.Version,
				len([]byte(nbChunkKey)),
				convert.BytesPerUint32)
			if err != nil {
				return nil, fmt.Errorf("unable to compute storage cost for chunk key creation: %w", err)
			}

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
			sendOperation.DefaultGasLimitCallSC,
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
	config config.NetworkInfos,
	atAddress string,
	content []byte,
	nickname string,
	missedChunks string,
	operationBatch sendOperation.OperationBatch,
) ([]string, error) {
	blocks := chunk(content, blockLength)

	operations, err := uploadMissedChunks(
		config,
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

//nolint:funlen
func uploadMissedChunks(
	network config.NetworkInfos,
	addr string,
	chunks [][]byte,
	missedChunks string,
	nickname string,
	operationBatch sendOperation.OperationBatch,
) ([]string, error) {
	client := node.NewClient(network.NodeURL)

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

		uploadCost, err := sendOperation.StorageCostForEntry(network.Version, convert.BytesPerUint32, chunkSize)
		if err != nil {
			return nil, fmt.Errorf("unable to compute storage cost for chunk upload: %w", err)
		}

		if chunkID == 0 {
			// if chunkID == 0, we may need to add the cost of the key creation for the NB_CHUNKS key
			chunkKeyCost, err := sendOperation.StorageCostForEntry(
				network.Version,
				len([]byte(nbChunkKey)),
				convert.BytesPerUint32)
			if err != nil {
				return nil, fmt.Errorf("unable to compute storage cost for chunk key creation: %w", err)
			}

			uploadCost += chunkKeyCost
		}

		operationResponse, err := onchain.CallFunction(
			client,
			nickname,
			addr,
			"appendBytesToWebsite",
			params,
			sendOperation.DefaultFee,
			sendOperation.DefaultGasLimitCallSC,
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

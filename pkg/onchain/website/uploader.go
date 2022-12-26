package website

import (
	"embed"
	"fmt"
	"strconv"
	"strings"

	"github.com/massalabs/thyra/pkg/convert"
	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/node/base58"
	"github.com/massalabs/thyra/pkg/node/sendoperation"
	"github.com/massalabs/thyra/pkg/onchain"
	"github.com/massalabs/thyra/pkg/onchain/dns"
	"github.com/massalabs/thyra/pkg/wallet"
)

//go:embed sc
var content embed.FS

const baseOffset = 5

const blockLength = 260000

// function calculating the max expiry period, this calculation is empiric

func maxExpiryPeriod(index int) uint64 {
	return baseOffset + uint64(index)*2
}

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

func Upload(atAddress string, content []byte, wallet wallet.Wallet, url string) ([]string, error) {
	client := node.NewDefaultClient()

	addr, _, err := base58.VersionedCheckDecode(atAddress[1:])
	if err != nil {
		return nil, fmt.Errorf("decoding address '%s': %w", atAddress[1:], err)
	}

	blocks := chunk(content, blockLength)

	operations, err := upload(client, addr, blocks, wallet, url)
	if err != nil {
		return nil, err
	}

	return operations, nil
}

func upload(client *node.Client, addr []byte, chunks [][]byte, wallet wallet.Wallet, url string) ([]string, error) {
	operations := make([]string, len(chunks)+1)

	// add encoded nbChunks to arguments
	argsInitializeWebsite := convert.U64ToBytes(len(chunks))

	websiteName := convert.StringToBytes(url)

	// add encoded websiteName to arguments
	argsInitializeWebsite = append(argsInitializeWebsite, websiteName...)

	opID, err := onchain.CallFunction(client, wallet, addr, "initializeWebsite", argsInitializeWebsite,
		sendoperation.OneMassa)
	if err != nil {
		return nil, fmt.Errorf("calling initializeWebsite at '%s': %w", addr, err)
	}

	operations[0] = opID

	for index := 0; index < len(chunks); index++ {
		// Chunk ID encoding
		params := convert.U64ToBytes(index)

		// add websiteName to params
		params = append(params, websiteName...)

		// Chunk data length encoding
		params = append(params, convert.U32ToBytes(len(chunks[index]))...)

		// Chunk data encoding
		params = append(params, chunks[index]...)

		opID, err = onchain.CallFunctionUnwaited(client, wallet, maxExpiryPeriod(index), addr, "appendBytesToWebsite", params)
		if err != nil {
			return nil, fmt.Errorf("calling appendBytesToWebsite at '%s': %w", addr, err)
		}

		operations[index+1] = opID
	}

	return operations, nil
}

//nolint:lll
func UploadMissedChunks(atAddress string, content []byte, wallet *wallet.Wallet, missedChunks string, url string) ([]string, error) {
	client := node.NewDefaultClient()

	addr, _, err := base58.VersionedCheckDecode(atAddress[1:])
	if err != nil {
		return nil, fmt.Errorf("decoding address '%s': %w", atAddress[1:], err)
	}

	blocks := chunk(content, blockLength)

	operations, err := uploadMissedChunks(client, addr, blocks, missedChunks, wallet, url)
	if err != nil {
		return nil, err
	}

	return operations, nil
}

//nolint:lll
func uploadMissedChunks(client *node.Client, addr []byte, chunks [][]byte, missedChunks string, wallet *wallet.Wallet, url string) ([]string, error) {
	operations := make([]string, len(chunks)+1)
	arrMissedChunks := strings.Split(missedChunks, ",")

	for index := 0; index < len(arrMissedChunks); index++ {
		chunkID, err := strconv.Atoi(arrMissedChunks[index])
		if err != nil {
			return nil, fmt.Errorf("error while converting chunk ID")
		}

		params := convert.U64ToBytes(chunkID)

		// add websiteName to params
		params = append(params, convert.StringToBytes(url)...)

		// Chunk data length encoding
		//nolint:ineffassign,nolintlint
		params = append(params, convert.U32ToBytes(len(chunks[chunkID]))...)
		// Chunk data encoding
		//nolint:ineffassign,nolintlint
		params = append(params, chunks[chunkID]...)

		//nolint:lll
		opID, err := onchain.CallFunctionUnwaited(client, *wallet, maxExpiryPeriod(index), addr, "appendBytesToWebsite", params)
		if err != nil {
			return nil, fmt.Errorf("calling appendBytesToWebsite at '%s': %w", addr, err)
		}

		operations[index] = opID
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

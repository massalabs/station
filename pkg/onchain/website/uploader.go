package website

import (
	"embed"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf16"

	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/node/base58"
	"github.com/massalabs/thyra/pkg/onchain"
	"github.com/massalabs/thyra/pkg/onchain/dns"
	"github.com/massalabs/thyra/pkg/wallet"
	"golang.org/x/text/encoding/unicode"
)

//go:embed sc
var content embed.FS

const baseOffset = 5

const multiplicator = 2

const blockLength = 260000

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

func Upload(atAddress string, content string, wallet *wallet.Wallet) ([]string, error) {
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

func upload(client *node.Client, addr []byte, chunks []string, wallet *wallet.Wallet) ([]string, error) {
	operations := make([]string, len(chunks)+1)
	totalChunks := make([]byte, 8)
	binary.LittleEndian.PutUint64(totalChunks, uint64(len(chunks)))
	//binary.LittleEndian.PutUint64(totalChunks, 18_446_744_073_709_551_615)

	totalChunksRunes := make([]rune, 8)

	for i := 0; i < len(totalChunks); i++ {
		totalChunksRunes[i] = utf16.Decode([]uint16{uint16(totalChunks[i])})[0]
	}

	totalChunksUTF8 := string(totalChunksRunes)

	totalChunksUTF16, err := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder().String(totalChunksUTF8)

	if err != nil {
		panic(err)
	}

	fmt.Println("%#v", totalChunks, totalChunksRunes, totalChunksUTF8, totalChunksUTF16)
	fmt.Println([]byte(totalChunksUTF16))

	opID, err := onchain.CallFunction(client, *wallet, addr, "initializeWebsite", []byte(totalChunksUTF16), 1000000000)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("calling initializeWebsite at '%s': %w", addr, err)
	}

	operations[0] = opID

	for index := 0; index < len(chunks); index++ {

		params := make([]byte, 8+8+len(chunks[index]))
		binary.LittleEndian.PutUint64(params, uint64(index))
		binary.LittleEndian.PutUint64(params, uint64(len(chunks[index])))
		params = append(params, []byte(chunks[index])...)

		//nolint:lll
		opID, err = onchain.CallFunctionUnwaited(client, *wallet, baseOffset+uint64(index)*multiplicator, addr, "appendBytesToWebsite", params)
		if err != nil {
			return nil, fmt.Errorf("calling appendBytesToWebsite at '%s': %w", addr, err)
		}

		operations[index+1] = opID
	}

	return operations, nil
}

//nolint:lll
func UploadMissedChunks(atAddress string, content string, wallet *wallet.Wallet, missedChunks string) ([]string, error) {
	client := node.NewDefaultClient()

	addr, _, err := base58.VersionedCheckDecode(atAddress[1:])
	if err != nil {
		return nil, fmt.Errorf("decoding address '%s': %w", atAddress[1:], err)
	}

	blocks := chunk(content, blockLength)

	operations, err := uploadMissedChunks(client, addr, blocks, missedChunks, wallet)
	if err != nil {
		return nil, err
	}

	return operations, nil
}

//nolint:lll
func uploadMissedChunks(client *node.Client, addr []byte, chunks []string, missedChunks string, wallet *wallet.Wallet) ([]string, error) {
	operations := make([]string, len(chunks)+1)
	arrMissedChunks := strings.Split(missedChunks, "")

	for index := 0; index < len(arrMissedChunks); index++ {
		rawParams := appendParams(index, chunks)

		param, err := json.Marshal(rawParams)
		if err != nil {
			return nil,
				fmt.Errorf("marshaling '%s': %w", rawParams, err)
		}

		//nolint:lll
		opID, err := onchain.CallFunctionUnwaited(client, *wallet, baseOffset+uint64(index)*multiplicator, addr, "appendBytesToWebsite", param)
		if err != nil {
			return nil, fmt.Errorf("calling initializeWebsite at '%s': %w", addr, err)
		}

		operations[index] = opID
	}

	return operations, nil
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

func appendParams(index int, chunks []string) AppendParams {
	return AppendParams{Data: chunks[index], ChunkID: strconv.Itoa(index)}
}

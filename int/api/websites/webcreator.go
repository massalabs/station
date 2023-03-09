package websites

import (
	"archive/zip"
	"bytes"
	"io"
	"net/http"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/gui"
	"github.com/massalabs/thyra/pkg/onchain/website"
	"github.com/massalabs/thyra/pkg/wallet"
	"golang.org/x/exp/slices"
)

const UploadMaxSize = "UPLOAD_MAX_SIZE"

const defaultMaxArchiveSize = 1500000

func CreatePrepareForWebsiteHandler(
	app *fyne.App,
) func(params operations.WebsiteCreatorPrepareParams) middleware.Responder {
	return func(params operations.WebsiteCreatorPrepareParams) middleware.Responder {
		return prepareForWebsiteHandler(params, app)
	}
}

func listFileName(zipReader *zip.Reader) []string {
	FilesInArchive := []string{}
	for _, zipFile := range zipReader.File {
		FilesInArchive = append(FilesInArchive, zipFile.Name)
	}

	return FilesInArchive
}

func prepareForWebsiteHandler(params operations.WebsiteCreatorPrepareParams, app *fyne.App) middleware.Responder {
	wallet, _ := LoadAndUnprotectWallet(params.Nickname, app)
	archive, _ := readAndCheckArchive(params.Zipfile, app)

	address, err := website.PrepareForUpload(params.URL, wallet)
	if err != nil {
		return createInternalServerError(errorCodeWebCreatorPrepare, err.Error())
	}

	_, err = website.Upload(address, archive, *wallet)
	if err != nil {
		return createInternalServerError(errorCodeWebCreatorUpload, err.Error())
	}

	return operations.NewWebsiteCreatorPrepareOK().
		WithPayload(
			&models.Websites{
				Name:         params.URL,
				Address:      address,
				BrokenChunks: nil,
			})
}

func GetMaxArchiveSize() int {
	uploadMaxSizeStr := os.Getenv(UploadMaxSize)

	if uploadMaxSizeStr == "" {
		return defaultMaxArchiveSize
	}

	uploadMaxSizeInt, err := strconv.Atoi(uploadMaxSizeStr)
	if err != nil {
		return defaultMaxArchiveSize
	}

	return uploadMaxSizeInt
}

func createInternalServerError(errorCode string, errorMessage string) middleware.Responder {
	return operations.NewWebsiteCreatorPrepareInternalServerError().
		WithPayload(
			&models.Error{
				Code:    errorCode,
				Message: errorMessage,
			})
}

func CreateUploadWebsiteHandler(app *fyne.App) func(params operations.WebsiteCreatorUploadParams) middleware.Responder {
	return func(params operations.WebsiteCreatorUploadParams) middleware.Responder {
		return uploadWebsiteHandler(params, app)
	}
}

func uploadWebsiteHandler(params operations.WebsiteCreatorUploadParams, app *fyne.App) middleware.Responder {
	wallet, _ := LoadAndUnprotectWallet(params.Nickname, app)
	archive, _ := readAndCheckArchive(params.Zipfile, app)

	_, err := website.Upload(params.Address, archive, *wallet)
	if err != nil {
		return createInternalServerError(errorCodeWebCreatorUpload, err.Error())
	}

	return operations.NewWebsiteCreatorUploadOK().
		WithPayload(&models.Websites{
			Name:         "",
			Address:      params.Address,
			BrokenChunks: nil,
		})
}

func checkContentType(archive []byte, fileType string) bool {
	contentType := http.DetectContentType(archive)

	return contentType == fileType
}

//nolint:lll
func CreateUploadMissingChunksHandler(app *fyne.App) func(params operations.WebsiteUploadMissingChunksParams) middleware.Responder {
	return func(params operations.WebsiteUploadMissingChunksParams) middleware.Responder {
		return websiteUploadMissingChunksHandler(params, app)
	}
}

//nolint:lll
func websiteUploadMissingChunksHandler(params operations.WebsiteUploadMissingChunksParams, app *fyne.App) middleware.Responder {
	wallet, _ := LoadAndUnprotectWallet(params.Nickname, app)
	archive, _ := readAndCheckArchive(params.Zipfile, app)

	_, err := website.UploadMissedChunks(params.Address, archive, wallet, params.MissedChunks)
	if err != nil {
		return createInternalServerError(errorCodeWebCreatorUpload, err.Error())
	}

	return operations.NewWebsiteUploadMissingChunksOK().
		WithPayload(&models.Websites{
			Name:         "",
			Address:      params.Address,
			BrokenChunks: nil,
		})
}

func LoadAndUnprotectWallet(nickname string, app *fyne.App) (*wallet.Wallet, middleware.Responder) {
	wallet, err := wallet.Load(nickname)
	if err != nil {
		return nil, createInternalServerError(errorCodeGetWallet, err.Error())
	}

	clearPassword, err := gui.AskPassword(wallet.Nickname, app)
	if err != nil {
		return nil, createInternalServerError(ErrorCodeWalletCanceledAction, err.Error())
	}

	if len(clearPassword) == 0 {
		return nil, createInternalServerError(ErrorCodeWalletPasswordEmptyWebCreator, ErrorCodeWalletPasswordEmptyWebCreator)
	}

	err = wallet.Unprotect(clearPassword, 0)
	if err != nil {
		return nil, createInternalServerError(errorCodeWalletWrongPassword, err.Error())
	}

	return wallet, nil
}

//nolint:unparam
func readAndCheckArchive(zipFile io.ReadCloser, app *fyne.App) ([]byte, middleware.Responder) {
	archive, err := io.ReadAll(zipFile)
	if err != nil {
		return nil, createInternalServerError(errorCodeWebCreatorReadArchive, err.Error())
	}

	maxArchiveSize := GetMaxArchiveSize()

	if len(archive) > maxArchiveSize {
		return nil, createInternalServerError(errorCodeWebCreatorArchiveSize, errorCodeWebCreatorArchiveSize)
	}

	zipReader, _ := zip.NewReader(bytes.NewReader(archive), int64(len(archive)))
	FilesOfArchive := listFileName(zipReader)

	if slices.Index(FilesOfArchive, "index.html") == -1 {
		return nil, createInternalServerError(errorCodeWebCreatorHTMLNotInSource, errorCodeWebCreatorHTMLNotInSource)
	}

	if !checkContentType(archive, "application/zip") {
		return nil, createInternalServerError(errorCodeWebCreatorFileType, errorCodeWebCreatorFileType)
	}

	return archive, nil
}

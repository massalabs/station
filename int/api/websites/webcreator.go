package websites

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"fmt"
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

func readZipFile(z *zip.File) ([]byte, error) {
	file, err := z.Open()
	if err != nil {
		return nil, fmt.Errorf("opening zip content: %w", err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("reading zip content: %w", err)
	}

	return content, nil
}
func contains(elems []string, v string) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}

	return false
}

//nolint:nolintlint,ireturn,funlen
func prepareForWebsiteHandler(params operations.WebsiteCreatorPrepareParams, app *fyne.App) middleware.Responder {
	list_of_files := []string{}
	wallet, err := wallet.Load(params.Nickname)
	if err != nil {
		return createInternalServerError(errorCodeGetWallet, err.Error())
	}

	clearPassword, err := gui.AskPassword(wallet.Nickname, app)
	if err != nil {
		return createInternalServerError(ErrorCodeWalletCanceledAction, err.Error())
	}

	if len(clearPassword) == 0 {
		return createInternalServerError(ErrorCodeWalletPasswordEmptyWebCreator, ErrorCodeWalletPasswordEmptyWebCreator)
	}

	err = wallet.Unprotect(clearPassword, 0)

	if err != nil {
		return createInternalServerError(errorCodeWalletWrongPassword, err.Error())
	}

	archive, err := io.ReadAll(params.Zipfile)
	if err != nil {
		return createInternalServerError(errorCodeWebCreatorReadArchive, err.Error())
	}

	zipReader, err := zip.NewReader(bytes.NewReader(archive), int64(len(archive)))
	// Read all the files from zip archive
	for _, zipFile := range zipReader.File {
		if err != nil {
			return nil
		}
		list_of_files = append(list_of_files, zipFile.Name)
	}
	// check if zip archive exist
	if !contains(list_of_files, "index.html") {
		return createInternalServerError(errorCodeWebCreatorHTMLNotInSource, err.Error())
	}
	maxArchiveSize := GetMaxArchiveSize()

	if len(archive) > maxArchiveSize {
		return createInternalServerError(errorCodeWebCreatorArchiveSize, errorCodeWebCreatorArchiveSize)
	}

	if !checkContentType(archive, "application/zip") {
		return createInternalServerError(errorCodeWebCreatorFileType, errorCodeWebCreatorFileType)
	}

	b64 := base64.StdEncoding.EncodeToString(archive)

	address, err := website.PrepareForUpload(params.URL, wallet)
	if err != nil {
		return createInternalServerError(errorCodeWebCreatorPrepare, err.Error())
	}

	_, err = website.Upload(address, b64, wallet)
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

//nolint:nolintlint,ireturn
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

//nolint:nolintlint,ireturn
func uploadWebsiteHandler(params operations.WebsiteCreatorUploadParams, app *fyne.App) middleware.Responder {
	wallet, err := wallet.Load(params.Nickname)
	if err != nil {
		return operations.NewWebsiteCreatorUploadInternalServerError().
			WithPayload(
				&models.Error{
					Code:    errorCodeGetWallet,
					Message: err.Error(),
				})
	}

	clearPassword, err := gui.AskPassword(wallet.Nickname, app)
	if err != nil {
		return operations.NewWebsiteCreatorUploadInternalServerError().
			WithPayload(
				&models.Error{
					Code:    ErrorCodeWalletCanceledAction,
					Message: ErrorCodeWalletCanceledAction,
				})
	}

	err = wallet.Unprotect(clearPassword, 0)
	if err != nil {
		return operations.NewWebsiteCreatorUploadInternalServerError().
			WithPayload(
				&models.Error{
					Code:    errorCodeWalletWrongPassword,
					Message: err.Error(),
				})
	}

	archive, err := io.ReadAll(params.Zipfile)
	if err != nil {
		return operations.NewWebsiteCreatorUploadInternalServerError().
			WithPayload(&models.Error{
				Code:    errorCodeWebCreatorReadArchive,
				Message: err.Error(),
			})
	}

	if !checkContentType(archive, "application/zip") {
		return createInternalServerError(errorCodeWebCreatorFileType, errorCodeWebCreatorFileType)
	}

	b64 := base64.StdEncoding.EncodeToString(archive)

	_, err = website.Upload(params.Address, b64, wallet)
	if err != nil {
		return createInternalServerError(errorCodeWebCreatorUpload, err.Error())
	}

	return operations.NewWebsiteCreatorUploadOK().
		WithPayload(&models.Websites{
			Name:         "Name",
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

//nolint:nolintlint,ireturn,lll
func websiteUploadMissingChunksHandler(params operations.WebsiteUploadMissingChunksParams, app *fyne.App) middleware.Responder {
	wallet, err := wallet.Load(params.Nickname)
	if err != nil {
		return operations.NewWebsiteCreatorUploadInternalServerError().
			WithPayload(
				&models.Error{
					Code:    errorCodeGetWallet,
					Message: err.Error(),
				})
	}

	clearPassword, err := gui.AskPassword(wallet.Nickname, app)
	if err != nil {
		return createInternalServerError(ErrorCodeWalletCanceledAction, err.Error())
	}

	err = wallet.Unprotect(clearPassword, 0)
	if err != nil {
		return operations.NewWebsiteCreatorUploadInternalServerError().
			WithPayload(
				&models.Error{
					Code:    errorCodeWalletWrongPassword,
					Message: err.Error(),
				})
	}

	archive, err := io.ReadAll(params.Zipfile)
	if err != nil {
		return operations.NewWebsiteCreatorUploadInternalServerError().
			WithPayload(&models.Error{
				Code:    errorCodeWebCreatorReadArchive,
				Message: err.Error(),
			})
	}

	if !checkContentType(archive, "application/zip") {
		return createInternalServerError(errorCodeWebCreatorFileType, errorCodeWebCreatorFileType)
	}

	b64 := base64.StdEncoding.EncodeToString(archive)

	_, err = website.UploadMissedChunks(params.Address, b64, wallet, params.MissedChunks)
	if err != nil {
		return createInternalServerError(errorCodeWebCreatorUpload, err.Error())
	}

	return operations.NewWebsiteUploadMissingChunksOK().
		WithPayload(&models.Websites{
			Name:         "Name",
			Address:      params.Address,
			BrokenChunks: nil,
		})
}

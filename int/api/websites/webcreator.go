package websites

import (
	"encoding/base64"
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

//nolint:nolintlint,ireturn,funlen
func prepareForWebsiteHandler(params operations.WebsiteCreatorPrepareParams, app *fyne.App) middleware.Responder {
	wallet, err := wallet.Load(params.Nickname)
	if err != nil {
		return createInternalServerError(errorCodeGetWallet, err.Error())
	}

	password := gui.AskPassword(wallet.Nickname, app)

	err = wallet.Unprotect(password, 0)
	if err != nil {
		return createInternalServerError(errorCodeWalletWrongPassword, err.Error())
	}

	archive, err := io.ReadAll(params.Zipfile)
	if err != nil {
		return createInternalServerError(errorCodeWebCreatorReadArchive, err.Error())
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
				Name:           params.URL,
				Address:        address,
				HasBrokenChunk: false,
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

	password := gui.AskPassword(wallet.Nickname, app)

	err = wallet.Unprotect(password, 0)
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
			Name:           "Name",
			Address:        params.Address,
			HasBrokenChunk: false,
		})
}

func checkContentType(archive []byte, fileType string) bool {
	contentType := http.DetectContentType(archive)

	return contentType == fileType
}

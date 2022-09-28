package api

import (
	"encoding/base64"
	"io"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/onchain/website"
	"github.com/massalabs/thyra/pkg/wallet"
)

const maxArchiveSize = 1500000

//nolint:nolintlint,ireturn,funlen
func PrepareForWebsiteHandler(params operations.WebsiteCreatorPrepareParams) middleware.Responder {
	wallet, err := wallet.Load(params.Nickname)
	if err != nil {
		return createInternalServerError(errorCodeGetWallet, err.Error())
	}

	err = wallet.Unprotect(params.HTTPRequest.Header.Get("Authorization"), 0)
	if err != nil {
		return createInternalServerError(errorCodeWalletWrongPassword, err.Error())
	}

	archive, err := io.ReadAll(params.Zipfile)
	if err != nil {
		return createInternalServerError(errorCodeWebCreatorReadArchive, err.Error())
	}

	if len(archive) > maxArchiveSize {
		return createInternalServerError(errorCodeWebCreatorArchiveSize, errorCodeWebCreatorArchiveSize)
	}

	if checkContentType(archive, "application/zip") == false {
		{
			return operations.NewWebsiteCreatorPrepareInternalServerError().
				WithPayload(&models.Error{
					Code:    errorCodeWebCreatorFileType,
					Message: err.Error(),
				})
		}
	}

	b64 := base64.StdEncoding.EncodeToString(archive)

	address, err := website.PrepareForUpload(params.URL, wallet)
	if err != nil {
		return createInternalServerError(errorCodeWebCreatorPrepare, err.Error())
	}

	_, err = website.Upload(address, b64, wallet)
	if err != nil {
		return operations.NewWebsiteCreatorPrepareInternalServerError().
			WithPayload(&models.Error{
				Code:    errorCodeWebCreatorUpload,
				Message: err.Error(),
			})
	}

	return operations.NewWebsiteCreatorPrepareOK().
		WithPayload(
			&models.Websites{
				Name:    params.URL,
				Address: address,
			})
}

//nolint:nolintlint,ireturn
func UploadWebsiteHandler(params operations.WebsiteCreatorUploadParams) middleware.Responder {
	wallet, err := wallet.Load(params.Nickname)
	if err != nil {
		return operations.NewWebsiteCreatorUploadInternalServerError().
			WithPayload(
				&models.Error{
					Code:    errorCodeGetWallet,
					Message: err.Error(),
				})
	}

	err = wallet.Unprotect(params.HTTPRequest.Header.Get("Authorization"), 0)
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

	if checkContentType(archive, "application/zip") == false {
		{
			return operations.NewWebsiteCreatorPrepareInternalServerError().
				WithPayload(&models.Error{
					Code:    errorCodeWebCreatorFileType,
					Message: err.Error(),
				})
		}
	}

	b64 := base64.StdEncoding.EncodeToString(archive)

	_, err = website.Upload(params.Address, b64, wallet)
	if err != nil {
		return createInternalServerError(errorCodeWebCreatorUpload, err.Error())
	}

	return operations.NewWebsiteCreatorUploadOK().
		WithPayload(&models.Websites{
			Name:    "Name",
			Address: params.Address,
		})
}

func checkContentType(archive []byte, fileType string) bool {
	contentType := http.DetectContentType(archive)

	if contentType != fileType {
		{
			return false
		}
	}
	return true
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

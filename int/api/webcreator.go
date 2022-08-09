package api

import (
	"encoding/base64"
	"io/ioutil"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/onchain/website"
)

func PrepareForWebsiteHandler(params operations.WebsiteCreatorPrepareParams) middleware.Responder {
	scAddress, err := website.PrepareForUpload(*params.Body.URL)

	if err != nil {
		return operations.NewWebsiteCreatorPrepareInternalServerError().
			WithPayload(
				&models.Error{
					Code:    "",
					Message: err.Error(),
				})
	}

	return operations.NewWebsiteCreatorPrepareOK().
		WithPayload(
			&models.Websites{
				Name:    *params.Body.URL,
				Address: scAddress,
			})
}

func UploadWebsiteHandler(params operations.WebsiteCreatorUploadParams) middleware.Responder {
	archive, err := ioutil.ReadAll(params.Zipfile)
	if err != nil {
		return operations.NewFillWebPostInternalServerError().
			WithPayload(&models.Error{
				Code:    "",
				Message: err.Error(),
			})
	}

	b64 := base64.StdEncoding.EncodeToString(archive)

	_, err = website.Upload(b64, params.Address)
	if err != nil {
		return operations.NewFillWebPostInternalServerError().
			WithPayload(&models.Error{
				Code:    "",
				Message: err.Error(),
			})
	}

	return operations.NewFillWebPostOK().
		WithPayload(&models.Websites{
			Name:    "Name",
			Address: params.Address})
}

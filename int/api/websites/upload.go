package websites

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/config"
	sendOperation "github.com/massalabs/thyra/pkg/node/sendoperation"
	"github.com/massalabs/thyra/pkg/onchain/website"
)

func NewWebsiteUploadHandler(config *config.AppConfig) operations.WebsiteCreatorUploadHandler {
	return &websiteUploadHandler{config: config}
}

type websiteUploadHandler struct {
	config *config.AppConfig
}

func (h *websiteUploadHandler) Handle(params operations.WebsiteCreatorUploadParams) middleware.Responder {
	archive, errorResponse := readAndCheckArchive(params.Zipfile)
	if errorResponse != nil {
		return errorResponse
	}
	_, err := website.Upload(
		*h.config,
		params.Address,
		archive,
		params.Nickname,
		sendOperation.OperationBatch{
			NewBatch:      true,
			CorrelationID: "",
		},
	)
	if err != nil {
		return createInternalServerError(errorCodeWebCreatorUpload, err.Error())
	}

	return operations.NewWebsiteCreatorUploadOK().
		WithPayload(&models.Websites{
			Name:         "",
			Description:  "",
			Address:      params.Address,
			BrokenChunks: nil,
		})
}

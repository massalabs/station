package websites

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station/api/swagger/server/models"
	"github.com/massalabs/station/api/swagger/server/restapi/operations"
	"github.com/massalabs/station/pkg/config"
	sendOperation "github.com/massalabs/station/pkg/node/sendoperation"
	"github.com/massalabs/station/pkg/onchain/website"
)

func NewWebsiteUploadHandler(config *config.AppConfig) operations.WebsiteUploaderUploadHandler {
	return &websiteUploadHandler{config: config}
}

type websiteUploadHandler struct {
	config *config.AppConfig
}

func (h *websiteUploadHandler) Handle(params operations.WebsiteUploaderUploadParams) middleware.Responder {
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

	return operations.NewWebsiteUploaderUploadOK().
		WithPayload(&models.Websites{
			Name:         "",
			Description:  "",
			Address:      params.Address,
			BrokenChunks: nil,
		})
}

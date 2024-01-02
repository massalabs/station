package websites

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station/api/swagger/server/models"
	"github.com/massalabs/station/api/swagger/server/restapi/operations"
	"github.com/massalabs/station/int/api/utils"
	"github.com/massalabs/station/int/config"
	sendOperation "github.com/massalabs/station/pkg/node/sendoperation"
	"github.com/massalabs/station/pkg/onchain/website"
)

func NewWebsiteUploadHandler(config *config.NetworkInfos) operations.WebsiteUploaderUploadHandler {
	return &websiteUploadHandler{networkInfos: config}
}

type websiteUploadHandler struct {
	networkInfos *config.NetworkInfos
}

func (h *websiteUploadHandler) Handle(params operations.WebsiteUploaderUploadParams) middleware.Responder {
	//nolint:revive
	return utils.NewGoneResponder()

	archive, errorResponse := readAndCheckArchive(params.Zipfile)
	if errorResponse != nil {
		return errorResponse
	}

	_, err := website.Upload(
		h.networkInfos,
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

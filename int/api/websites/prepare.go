package websites

import (
	"archive/zip"
	"bytes"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station/api/swagger/server/models"
	"github.com/massalabs/station/api/swagger/server/restapi/operations"
	"github.com/massalabs/station/int/api/utils"
	"github.com/massalabs/station/int/config"
	sendOperation "github.com/massalabs/station/pkg/node/sendoperation"
	"github.com/massalabs/station/pkg/onchain/website"
	"golang.org/x/exp/slices"
)

const UploadMaxSize = "UPLOAD_MAX_SIZE"

const defaultMaxArchiveSize = 1500000

func NewWebsitePrepareHandler(config *config.NetworkInfos) operations.WebsiteUploaderPrepareHandler {
	return &websitePrepare{networkInfos: config}
}

type websitePrepare struct {
	networkInfos *config.NetworkInfos
}

func listFileName(zipReader *zip.Reader) []string {
	FilesInArchive := []string{}
	for _, zipFile := range zipReader.File {
		FilesInArchive = append(FilesInArchive, zipFile.Name)
	}

	return FilesInArchive
}

func (h *websitePrepare) Handle(params operations.WebsiteUploaderPrepareParams) middleware.Responder {
	//nolint:revive
	return utils.NewGoneResponder()

	//nolint:govet
	archive, errorResponse := readAndCheckArchive(params.Zipfile)
	if errorResponse != nil {
		return errorResponse
	}

	address, correlationID, err := website.PrepareForUpload(
		h.networkInfos,
		params.URL,
		params.Description,
		params.Nickname,
	)
	if err != nil {
		return createInternalServerError(errorCodeWebCreatorPrepare, err.Error())
	}

	_, err = website.Upload(
		h.networkInfos,
		address,
		archive,
		params.Nickname,
		sendOperation.OperationBatch{
			NewBatch:      false,
			CorrelationID: correlationID,
		},
	)
	if err != nil {
		return createInternalServerError(errorCodeWebCreatorUpload, err.Error())
	}

	return operations.NewWebsiteUploaderPrepareOK().
		WithPayload(
			&models.Websites{
				Name:         params.URL,
				Description:  params.Description,
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
	return operations.NewWebsiteUploaderPrepareInternalServerError().
		WithPayload(
			&models.Error{
				Code:    errorCode,
				Message: errorMessage,
			})
}

func checkContentType(archive []byte, fileType string) bool {
	contentType := http.DetectContentType(archive)

	return contentType == fileType
}

func readAndCheckArchive(zipFile io.ReadCloser) ([]byte, middleware.Responder) {
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

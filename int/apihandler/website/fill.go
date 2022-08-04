package websites

import (
	"encoding/base64"
	"io"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/onchain/website"
)

func NewFillWebsitePost() operations.FillWebPostHandler {
	return &fillWeb{todelete: "todelete"}
}

type fillWeb struct {
	todelete string
}

func (c *fillWeb) Handle(params operations.FillWebPostParams) middleware.Responder {

	buf := new(strings.Builder)
	_, err := io.Copy(buf, params.Zipfile)
	if err != nil {
		return operations.NewFillWebPostInternalServerError()
	}
	sEnc := base64.StdEncoding.EncodeToString([]byte(buf.String()))
	_, err = website.UploadWebsite(sEnc, params.Website)
	if err != nil {
		return operations.NewFillWebPostInternalServerError()
	}

	website := &models.Websites{
		Name:    "Name",
		Address: params.Website}

	return operations.NewFillWebPostOK().WithPayload(website)
}

package websites

import (
	b64 "encoding/base64"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/onchain/website"
)

func NewFillWebsitePost() operations.FillWebPostHandler {
	return &fillWeb{as: "a"}
}

type fillWeb struct {
	as string
}

func (c *fillWeb) Handle(params operations.FillWebPostParams) middleware.Responder {

	data := params.HTTPRequest.MultipartForm.Value["zipfile"]
	sEnc := b64.StdEncoding.EncodeToString([]byte(data[0]))
	_, err := website.UploadWebsite([]byte(sEnc), params.Website)
	if err != nil {
		return operations.NewFillWebPostInternalServerError()
	}

	return operations.NewFillWebPostOK()
}

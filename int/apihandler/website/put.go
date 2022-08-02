package website

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
)

func NewWebsitePut() operations.UploadWebPutHandler {
	return &newWebsitePut{as: "a"}
}

type newWebsitePut struct {
	as string
}

func (c *newWebsitePut) Handle(params operations.UploadWebPutParams) middleware.Responder {

	// _, err := website.UploadWebsite(params.Website)
	// if err != nil {
	// 	return operations.NewUploadWebPutInternalServerError()
	// }
	return operations.NewUploadWebPutOK()
}

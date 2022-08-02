package website

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/onchain/website"
)

func NewWebsiteGet() operations.UploadWebGetHandler {
	return &newWebsiteGet{as: "a"}
}

type newWebsiteGet struct {
	as string
}

func (c *newWebsiteGet) Handle(params operations.UploadWebGetParams) middleware.Responder {
	_, err := website.RefreshDeployers()
	if err != nil {
		return operations.NewUploadWebGetInternalServerError()
	}
	return operations.NewWebsiteGetOK()
}

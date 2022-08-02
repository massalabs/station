package website

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/onchain/website"
)

func NewWebsitePost() operations.UploadWebPostHandler {
	return &newWebsitePost{as: "a"}
}

type newWebsitePost struct {
	as string
}

func (c *newWebsitePost) Handle(params operations.UploadWebPostParams) middleware.Responder {
	_, err := website.CreateWebsiteDeployer()
	if err != nil {
		return operations.NewUploadWebPostInternalServerError()
	}
	return operations.NewUploadWebPostOK()
}

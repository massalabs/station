package websites

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/onchain/website"
)

func NewWebsitePost() operations.UploadWebPostHandler {
	return &newWebsitePost{todelete: "todelete"}
}

type newWebsitePost struct {
	todelete string
}

func (c *newWebsitePost) Handle(params operations.UploadWebPostParams) middleware.Responder {
	smartContract, err := website.PostWebsite(params.Dnsname)

	if err != nil {
		return operations.NewUploadWebPostInternalServerError()
	}
	newWebsite := &models.Websites{
		Name:    params.Dnsname,
		Address: *smartContract}

	return operations.NewUploadWebPostOK().WithPayload(newWebsite)
}

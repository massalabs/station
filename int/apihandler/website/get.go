package websites

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/onchain/website"
)

func NewWebsiteGet() operations.UploadWebGetHandler {
	return &newWebsiteGet{todelete: "todelete"}
}

type newWebsiteGet struct {
	todelete string
}

func (c *newWebsiteGet) Handle(params operations.UploadWebGetParams) middleware.Responder {
	deployers, err := website.GetDeployers()
	if err != nil {
		return operations.NewUploadWebGetInternalServerError()
	}

	var websites []*models.Websites

	for i := 0; i < len(deployers); i++ {
		newWebsite := &models.Websites{
			Name:    *deployers[i].DnsName,
			Address: *deployers[i].Address}
		websites = append(websites, newWebsite)
	}
	return operations.NewUploadWebGetOK().WithPayload(websites)
}

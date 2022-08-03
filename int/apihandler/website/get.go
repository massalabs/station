package websites

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
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
	deployers, err := website.RefreshDeployers()
	if err != nil {
		return operations.NewUploadWebGetInternalServerError()
	}

	var websites []*models.Websites

	for i := 0; i < len(deployers); i++ {
		newWebsite := &models.Websites{
			Name:    "Name",
			Address: deployers[i]}
		websites = append(websites, newWebsite)
	}
	return operations.NewUploadWebGetOK().WithPayload(websites)
}

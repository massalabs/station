package api

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/my"
	"github.com/massalabs/thyra/pkg/node"
)

func DomainsHandler(params operations.MyDomainsParams) middleware.Responder {
	client := node.NewDefaultClient()

	myDomainNames, err := my.GetDomains(client, params.Nickname)
	if err != nil {
		return operations.NewUploadWebGetInternalServerError().
			WithPayload(
				&models.Error{
					Code:    errorCodeGetDomainNames,
					Message: err.Error(),
				})
	}

	myDomains, err := my.GetOwnedWebsites(client, myDomainNames)
	if err != nil {
		return operations.NewUploadWebGetInternalServerError().
			WithPayload(
				&models.Error{
					Code:    errorCodeGetDomainAddresses,
					Message: err.Error(),
				})
	}
	return operations.NewMyDomainsOK().WithPayload(myDomains)
}

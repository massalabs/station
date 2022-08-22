package api

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/onchain/dns"
)

func DomainsHandler(params operations.MyDomainsParams) middleware.Responder {
	client := node.NewDefaultClient()

	myDomainNames, err := dns.GetMyDomainNames(client, params.Nickname)
	if err != nil {
		return operations.NewUploadWebGetInternalServerError().
			WithPayload(
				&models.Error{
					Code:    errorCodeGetDomainNames,
					Message: err.Error(),
				})
	}

	myDomains, err := dns.GetOwnedDomains(client, myDomainNames)
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

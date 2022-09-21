package api

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/my"
	"github.com/massalabs/thyra/pkg/node"
)

//nolint:nolintlint,ireturn
func DomainsHandler(params operations.MyDomainsGetterParams) middleware.Responder {
	client := node.NewDefaultClient()

	myDomainNames, err := my.Domains(client, params.Nickname)
	if err != nil {
		return operations.NewMyDomainsGetterInternalServerError().
			WithPayload(
				&models.Error{
					Code:    errorCodeGetDomainNames,
					Message: err.Error(),
				})
	}

	myDomains, err := my.Websites(client, myDomainNames)
	if err != nil {
		return operations.NewMyDomainsGetterInternalServerError().
			WithPayload(
				&models.Error{
					Code:    errorCodeGetDomainAddresses,
					Message: err.Error(),
				})
	}

	return operations.NewMyDomainsGetterOK().WithPayload(myDomains)
}

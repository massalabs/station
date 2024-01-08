package websites

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station/api/swagger/server/models"
	"github.com/massalabs/station/api/swagger/server/restapi/operations"
	"github.com/massalabs/station/int/api/utils"
	"github.com/massalabs/station/int/config"
	"github.com/massalabs/station/pkg/my"
	"github.com/massalabs/station/pkg/node"
)

func NewDomainsHandler(config *config.NetworkInfos) operations.MyDomainsGetterHandler {
	return &domainsHandler{config: config}
}

type domainsHandler struct {
	config *config.NetworkInfos
}

func (h *domainsHandler) Handle(params operations.MyDomainsGetterParams) middleware.Responder {
	//nolint:revive
	return utils.NewGoneResponder()

	//nolint:govet
	client := node.NewClient(h.config.NodeURL)

	myDomainNames, err := my.Domains(*h.config, client, params.Nickname)
	if err != nil {
		return operations.NewMyDomainsGetterInternalServerError().
			WithPayload(
				&models.Error{
					Code:    errorCodeGetDomainNames,
					Message: err.Error(),
				})
	}

	myDomains, err := my.GetWebsites(*h.config, client, myDomainNames)
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

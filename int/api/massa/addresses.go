package massa

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/station/api/swagger/server/models"
	"github.com/massalabs/station/api/swagger/server/restapi/operations"
	"github.com/massalabs/station/int/config"
	"github.com/massalabs/station/pkg/node"
)

func NewGetAddressHandler(configManager *config.MSConfigManager) operations.MassaGetAddressesHandler {
	return &getAddress{configManager: configManager}
}

type getAddress struct {
	configManager *config.MSConfigManager
}

func (g *getAddress) Handle(params operations.MassaGetAddressesParams) middleware.Responder {
	currentNetwork := g.configManager.CurrentNetwork()
	client := node.NewClient(currentNetwork.NodeURL)

	addressesDetails, err := node.Addresses(client, params.Addresses)
	if err != nil {
		return operations.NewMassaGetAddressesInternalServerError().
			WithPayload(
				&models.Error{
					Code:    errorCodeMassaAddresses,
					Message: fmt.Sprintf("while getting details of addresses %v: %s\n", params.Addresses, err),
				},
			)
	}
	//nolint: lll
	addressesAttributes := make(map[string]operations.MassaGetAddressesOKBodyAddressesAttributesAnon, len(addressesDetails))

	for _, details := range addressesDetails {
		//nolint: exhaustruct
		attributes := operations.MassaGetAddressesOKBodyAddressesAttributesAnon{}

		if hasAttribute(params.Attributes, "balance") || hasNotAttribute(params.Attributes) {
			attributes.Balance = &operations.MassaGetAddressesOKBodyAddressesAttributesAnonBalance{
				Pending: details.CandidateBalance,
				Final:   details.FinalBalance,
			}
		}

		addressesAttributes[details.Address] = attributes
	}

	return operations.NewMassaGetAddressesOK().
		WithPayload(&operations.MassaGetAddressesOKBody{
			AddressesAttributes: addressesAttributes,
		})
}

func hasAttribute(request []string, attribute string) bool {
	for _, v := range request {
		if v == attribute {
			return true
		}
	}

	return false
}

func hasNotAttribute(request []string) bool {
	return len(request) == 0
}

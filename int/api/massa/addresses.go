package massa

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/node"
)

func AddressesHandler(params operations.MassaGetAddressesParams) middleware.Responder {
	client := node.NewDefaultClient()

	addressesDetails, err := node.Addresses(client, params.Query)
	if err != nil {
		return operations.NewMassaGetAddressesInternalServerError().
			WithPayload(
				&models.Error{
					Code:    errorCodeMassaAddresses,
					Message: fmt.Sprintf("while getting details of addresses %v: %s\n", params.Query, err),
				},
			)
	}
	//nolint: lll
	addressesAttributes := make(map[string]operations.MassaGetAddressesOKBodyAddressesAttributesAnon, len(addressesDetails))

	for _, details := range addressesDetails {
		//nolint: exhaustruct
		attributes := operations.MassaGetAddressesOKBodyAddressesAttributesAnon{}

		if hasAttribute(params, "balance") || hasNotAttribute(params) {
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

func hasAttribute(request operations.MassaGetAddressesParams, attribute string) bool {
	for _, v := range request.Query {
		if v == attribute {
			return true
		}
	}

	return false
}

func hasNotAttribute(request operations.MassaGetAddressesParams) bool {
	return request.Query[0] == ""
}

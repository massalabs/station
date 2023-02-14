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

	addressMap := make(map[string]operations.MassaGetAddressesOKBodyAddressesAttributesAnon, len(addressesDetails))

	for _, details := range addressesDetails {
		//nolint: exhaustruct
		attribute := operations.MassaGetAddressesOKBodyAddressesAttributesAnon{}

		if requestedAttributesContains(params, "balance") || requestedAttributeIsEmpty(params) {
			attribute.Balance = &operations.MassaGetAddressesOKBodyAddressesAttributesAnonBalance{
				Pending: details.CandidateBalance,
				Final:   details.FinalBalance,
			}
		}

		addressMap[details.Address] = attribute
	}

	return operations.NewMassaGetAddressesOK().
		//nolint: govet,nolintlint
		WithPayload(&operations.MassaGetAddressesOKBody{addressMap})
}

func requestedAttributesContains(requestedAttributes operations.MassaGetAddressesParams, valueToCheck string) bool {
	for _, v := range requestedAttributes.Query {
		if v == valueToCheck {
			return true
		}
	}

	return false
}

func requestedAttributeIsEmpty(requestedAttributes operations.MassaGetAddressesParams) bool {
	return requestedAttributes.Query[0] == ""
}

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

	//nolint: prealloc
	var pendingBalances []string
	for _, addressDetails := range addressesDetails {
		//nolint: staticcheck, nolintlint
		pendingBalances = append(pendingBalances, addressDetails.CandidateBalance)
	}

	//nolint: prealloc
	var finalBalances []string
	for _, addressDetails := range addressesDetails {
		//nolint: staticcheck, nolintlint
		finalBalances = append(finalBalances, addressDetails.CandidateBalance)
	}

	addressMap := make(models.AddressesAttributes, len(addressesDetails))

	for _, details := range addressesDetails {
		balance := &models.AddressesAttributesAnonBalance{Pending: details.CandidateBalance, Final: details.FinalBalance}

		addressMap[details.Address] = models.AddressesAttributesAnon{Balance: balance}
	}

	return operations.NewMassaGetAddressesOK().
		WithPayload(addressMap)
}

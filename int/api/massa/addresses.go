package massa

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/node"
)

func AddressesHandler(params operations.MassaGetAddressesParams) middleware.Responder {
	client := node.NewDefaultClient()

	addressesDetails, err := node.Addresses(client, params.Body.Addresses)
	if err != nil {
		return operations.NewMassaGetAddressesInternalServerError().
			WithPayload(
				&models.Error{
					Code:    "get_Addresses error",
					Message: "Error : Cannot get result from Address: " + err.Error(),
				},
			)
	}
	//nolint: prealloc
	var PendingBalances []string
	for _, addressDetails := range addressesDetails {
		//nolint: staticcheck
		PendingBalances = append(PendingBalances, addressDetails.CandidateBalance)
	}
	//nolint: prealloc
	var FinalBalances []string
	for _, addressDetails := range addressesDetails {
		//nolint: staticcheck
		FinalBalances = append(FinalBalances, addressDetails.CandidateBalance)
	}

	if checkIfPresentInStringArray(params.Body.Options, "balances") {
		return operations.NewMassaGetAddressesOK().
			WithPayload(
				&models.GetAddresses{
					FinalBalances:   FinalBalances,
					PendingBalances: FinalBalances,
				})
	}

	return operations.NewMassaGetAddressesBadRequest().
		WithPayload(
			&models.Error{
				Code:    "MassaAddress Error",
				Message: "Options missing or invalid",
			})
}

func checkIfPresentInStringArray(arr []string, toCheck string) bool {
	for _, ar := range arr {
		if ar == toCheck {
			return true
		}
	}

	return false
}

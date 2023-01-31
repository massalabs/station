package node

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type Balance struct {
	Candidate decimal.Decimal
	Final     decimal.Decimal
}

// FetchBalance returns as a decimal the candidate balance and final balance of an address.
func FetchBalance(client *Client, address string) (*Balance, error) {
	addressDetails, err := Addresses(client, []string{address})
	if err != nil {
		return nil, err
	}

	candidate, err := decimal.NewFromString(addressDetails[0].CandidateBalance)
	if err != nil {
		return nil, fmt.Errorf("converting candidateBalance %s to decimal :%w", addressDetails[0].CandidateBalance, err)
	}

	final, err := decimal.NewFromString(addressDetails[0].FinalBalance)
	if err != nil {
		return nil, fmt.Errorf("converting FinalBalance %s to decimal :%w", addressDetails[0].FinalBalance, err)
	}

	return &Balance{Candidate: candidate, Final: final}, nil
}

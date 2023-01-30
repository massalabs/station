package node

import (
	"fmt"
	"strconv"
)

type Balance struct {
	Candidate float64
	Final     float64
}

// FetchBalance returns as a float64 the candidate balance and final balance of an address.
func FetchBalance(client *Client, address string) (*Balance, error) {
	addressDetails, err := Addresses(client, []string{address})
	if err != nil {
		return nil, err
	}

	candidate, err := strconv.ParseFloat(addressDetails[0].CandidateBalance, 64)
	if err != nil {
		return nil, fmt.Errorf("converting candidateBalance %s f64 :%w", addressDetails[0].CandidateBalance, err)
	}

	final, err := strconv.ParseFloat(addressDetails[0].FinalBalance, 64)
	if err != nil {
		return nil, fmt.Errorf("converting candidateBalance %s f64 :%w", addressDetails[0].FinalBalance, err)
	}

	return &Balance{Candidate: candidate, Final: final}, nil
}

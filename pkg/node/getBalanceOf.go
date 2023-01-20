package node

import (
	"fmt"
	"strconv"
)

type GetBalanceOfResponse struct {
	CandidateBalance float64
	FinalBalance     float64
}

func GetBalanceOf(client *Client, walletAddress string) (*GetBalanceOfResponse, error) {
	addressDetails, err := Addresses(client, []string{walletAddress})
	if err != nil {
		return nil, fmt.Errorf("calling get_addresses :%w", err)
	}

	candidateBalance, err := strconv.ParseFloat(addressDetails[0].CandidateBalance, 64)
	if err != nil {
		return nil, fmt.Errorf("converting string to float :%w", err)
	}

	finalBalance, err := strconv.ParseFloat(addressDetails[0].FinalBalance, 64)
	if err != nil {
		return nil, fmt.Errorf("converting string to float :%w", err)
	}

	return &GetBalanceOfResponse{CandidateBalance: candidateBalance, FinalBalance: finalBalance}, nil
}

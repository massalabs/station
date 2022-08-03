package dns

import (
	"encoding/json"
	"time"

	"github.com/massalabs/thyra/pkg/node"
	"github.com/massalabs/thyra/pkg/node/base58"
	"github.com/massalabs/thyra/pkg/node/getters"
	sendOperation "github.com/massalabs/thyra/pkg/node/sendoperation"
	"github.com/massalabs/thyra/pkg/node/sendoperation/callsc"
	"github.com/massalabs/thyra/pkg/wallet"
)

type setApprovalForAll struct {
	Operator string `json:"operator"`
	Approved bool   `json:"approved"`
}
type setResolvers struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

func SetDnsResolver(c *node.Client, wallet wallet.Wallet, dnsName string, smartContract string, expire uint64) (*string, error) {
	dnsAddress, _, err := base58.VersionedCheckDecode("A12jkDPTcdhkqGg9VoKsTwvkBwZeSHQw7wJqQYKrNesKnjnGejuR"[1:])
	if err != nil {
		return nil, err
	}
	setResolvers := &setResolvers{
		Name:    dnsName,
		Address: smartContract}

	b, err := json.Marshal(setResolvers)
	if err != nil {
		return nil, err
	}
	callSC := callsc.New(dnsAddress, "setResolver", b, 0, 700000000, 0, 0)
	operationId, err := sendOperation.Call(c, expire, 0, callSC, wallet.KeyPairs[0].PublicKey, wallet.KeyPairs[0].PrivateKey)
	if err != nil {
		return nil, err
	}

	// Wait DNS is set through events
	dnsSet := false
	n := 0
	for n < 3 && dnsSet {

		time.Sleep(15 * time.Second)
		events, err := getters.GetEvents(c, nil, nil, nil, nil, &operationId)
		if err != nil {
			return nil, err
		}

		eventsValue := *events
		if len(eventsValue) > 0 {
			dnsSet = true
		}
		n++
	}
	return &operationId, nil

}

func SetDnsApproval(c *node.Client, wallet wallet.Wallet, a bool, expire uint64) (*string, error) {
	dnsAddress, _, err := base58.VersionedCheckDecode("A1Q65NojVV5YPyZruVkeU1CGeS3tjLNwGSzAmZfAJPE5vuvus4C"[1:])
	if err != nil {
		return nil, err
	}

	// Set Resolver prepare data
	setApprovalForAll := &setApprovalForAll{
		Operator: wallet.Address,
		Approved: true}

	b, err := json.Marshal(setApprovalForAll)
	if err != nil {
		return nil, err
	}
	callSC := callsc.New(dnsAddress, "setApprovalForAll", b, 0, 700000000, 0, 0)
	operationId, err := sendOperation.Call(c, expire, 0, callSC, wallet.KeyPairs[0].PublicKey, wallet.KeyPairs[0].PrivateKey)
	if err != nil {
		return nil, err
	}
	return &operationId, nil

}

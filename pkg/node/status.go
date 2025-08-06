package node

import (
	"context"
	"fmt"
	"regexp"
)

type State struct {
	Config         *Config         `json:"config"`
	ConsensusStats *ConsensusStats `json:"consensus_stats"`
	CurrentCycle   *uint           `json:"current_cycle"`
	CurrentTime    *uint           `json:"current_time"`
	LastSlot       *Slot           `json:"last_slot"`
	NetworkStats   *NetworkStats   `json:"network_stats"`
	NextSlot       *Slot           `json:"next_slot"`
	NodeID         *string         `json:"node_id"`
	NodeIP         *string         `json:"node_ip"`
	PoolStats      *[]uint         `json:"pool_stats"`
	Version        *string         `json:"version"`
	ChainID        *uint           `json:"chain_id"`
	MinimalFees    *string         `json:"minimal_fees"`
}

//nolint:tagliatelle
type Config struct {
	BlockReward             *string `json:"block_reward"`
	DeltaF0                 *uint   `json:"delta_f0"`
	EndTimeStamp            *uint   `json:"end_timestamp"`
	GenesisTimestamp        *uint   `json:"genesis_timestamp"`
	OperationValidityParios *uint   `json:"operation_validity_periods"`
	PeriodsPerCycle         *uint   `json:"periods_per_cycle"`
	PosLockCycles           *uint   `json:"pos_lock_cycles"`
	PosLookbackCycle        *uint   `json:"pos_lookback_cycles"`
	RollPrice               *string `json:"roll_price"`
	T0                      *uint   `json:"t0"`
	ThreadCount             *uint   `json:"thread_count"`
}

type ConsensusStats struct {
	CliqueCount         *uint `json:"clique_count"`
	EndTimespan         *uint `json:"end_timespan"`
	FinalBlockCount     *uint `json:"final_block_count"`
	FinalOperationCount *uint `json:"final_operation_count"`
	StakerCount         *uint `json:"staker_count"`
	StaleBlockCount     *uint `json:"stale_block_count"`
	StartTimespan       *uint `json:"start_timespan"`
}

//nolint:tagliatelle
type NetworkStats struct {
	ActiveNodeCount    *uint `json:"active_node_count"`
	BannedPeerCount    *uint `json:"banned_peer_count"`
	InConnectionCount  *uint `json:"in_connection_count"`
	KnowPeerCount      *uint `json:"known_peer_count"`
	OutConnectionCount *uint `json:"out_connection_count"`
}

func Status(client *Client) (*State, error) {
	rawResponse, err := client.RPCClient.Call(
		context.Background(),
		"get_status",
	)
	if err != nil {
		return nil, fmt.Errorf("calling get_status: %w", err)
	}

	if rawResponse.Error != nil {
		return nil, rawResponse.Error
	}

	var resp State

	err = rawResponse.GetObject(&resp)
	if err != nil {
		return nil, fmt.Errorf("parsing get_status jsonrpc response '%+v': %w", rawResponse, err)
	}

	return &resp, nil
}

// GetVersionDigits extracts the version digits from the node version string.
// Example: "DEVN.1.2" -> "1.2"
// Example: "MAIN.6.6" -> "6.6"
func GetVersionDigits(status *State) (string, error) {
	pattern := `.+\.(\d+\.\d+)`

	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(*status.Version)

	//nolint:gomnd
	if len(matches) != 2 {
		return "", fmt.Errorf("failed to parse node version from: %s", *status.Version)
	}

	return matches[1], nil
}

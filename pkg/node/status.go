package node

import (
	"context"
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
	PoolStats      *PoolStats      `json:"pool_stats"`
	Version        *string         `json:"version"`
}
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

type NetworkStats struct {
	ActiveNodeCount    *uint `json:"active_node_count"`
	BannedPeerCount    *uint `json:"banned_peer_count"`
	InConnectionCount  *uint `json:"in_connection_count"`
	KnowPeerCount      *uint `json:"known_peer_count"`
	OutConnectionCount *uint `json:"out_connection_count"`
}

type PoolStats struct {
	EndorsementCount *uint `json:"endorsement_count"`
	OperationCount   *uint `json:"operation_count"`
}

func Status(client *Client) (*State, error) {
	rawResponse, err := client.RPCClient.Call(
		context.Background(),
		"get_status",
	)
	if err != nil {
		return nil, err
	}

	if rawResponse.Error != nil {
		return nil, rawResponse.Error
	}

	var resp State

	err = rawResponse.GetObject(&resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

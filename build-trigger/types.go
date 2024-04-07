package main

import (
	"fmt"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

type BlockID struct {
	Hash   common.Hash `json:"hash"`
	Number uint64      `json:"number"`
}

type L2BlockRef struct {
	Hash           common.Hash `json:"hash"`
	Number         uint64      `json:"number"`
	ParentHash     common.Hash `json:"parentHash"`
	Time           uint64      `json:"timestamp"`
	L1Origin       BlockID     `json:"l1origin"`
	SequenceNumber uint64      `json:"sequenceNumber"` // distance to first block of epoch
}

type Uint64Quantity = hexutil.Uint64
type Bytes32 [32]byte

func (b *Bytes32) UnmarshalJSON(text []byte) error {
	return hexutil.UnmarshalFixedJSON(reflect.TypeOf(b), text, b[:])
}

func (b *Bytes32) UnmarshalText(text []byte) error {
	return hexutil.UnmarshalFixedText("Bytes32", text, b[:])
}

func (b Bytes32) MarshalText() ([]byte, error) {
	return hexutil.Bytes(b[:]).MarshalText()
}

func (b Bytes32) String() string {
	return hexutil.Encode(b[:])
}

// TerminalString implements log.TerminalStringer, formatting a string for console
// output during logging.
func (b Bytes32) TerminalString() string {
	return fmt.Sprintf("%x..%x", b[:3], b[29:])
}

type Data = hexutil.Bytes

type PayloadAttributes struct {
	// value for the timestamp field of the new payload
	Timestamp Uint64Quantity `json:"timestamp"`
	// value for the random field of the new payload
	PrevRandao Bytes32 `json:"prevRandao"`
	// suggested value for the coinbase field of the new payload
	SuggestedFeeRecipient common.Address `json:"suggestedFeeRecipient"`
	// Withdrawals to include into the block -- should be nil or empty depending on Shanghai enablement
	Withdrawals *types.Withdrawals `json:"withdrawals,omitempty"`
	// parentBeaconBlockRoot optional extension in Dencun
	ParentBeaconBlockRoot *common.Hash `json:"parentBeaconBlockRoot,omitempty"`

	// Optimism additions

	// Transactions to force into the block (always at the start of the transactions list).
	Transactions []Data `json:"transactions,omitempty"`
	// NoTxPool to disable adding any transactions from the transaction-pool.
	NoTxPool bool `json:"noTxPool,omitempty"`
	// GasLimit override
	GasLimit *Uint64Quantity `json:"gasLimit,omitempty"`
}

type AttributesWithParent struct {
	Attributes   PayloadAttributes `json:"attributes"`
	Parent       L2BlockRef        `json:"parent"`
	IsLastInSpan bool              `json:"isLastInSpan"`
}

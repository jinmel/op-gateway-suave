package utils

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/suave/sdk"
)

func NewPrivKeyFromHex(hex string) *PrivKey {
	key, err := crypto.HexToECDSA(hex)
	if err != nil {
		panic(fmt.Sprintf("failed to parse private key: %v", err))
	}
	return &PrivKey{priv: key}
}

type PrivKey struct {
	priv  *ecdsa.PrivateKey
	nonce uint64
}

func (p *PrivKey) Address() common.Address {
	return crypto.PubkeyToAddress(p.priv.PublicKey)
}

func (p *PrivKey) MarshalPrivKey() []byte {
	return crypto.FromECDSA(p.priv)
}

func (p *PrivKey) Nonce() uint64 {
	return p.nonce
}

func (p *PrivKey) StepNonce() {
	p.nonce = p.nonce + 1
}

type Client struct {
	*sdk.Client
	Key *PrivKey
}

func NewClient(rpc *rpc.Client, key *PrivKey, addr common.Address) *Client {
	return &Client{Client: sdk.NewClient(rpc, key.priv, addr), Key: key}
}

func DeployContract(artifact *Artifact, clt *Client) (*sdk.Contract, error) {
	txnResult, err := sdk.DeployContract(artifact.Code, clt.Client)

	if err != nil {
		return nil, err
	}

	receipt, err := txnResult.Wait()
	if err != nil {
		return nil, err
	}

	if receipt.Status == 0 {
		return nil, fmt.Errorf("failed to deploy contract")
	}

	return sdk.GetContract(receipt.ContractAddress, artifact.Abi, clt.Client), nil
}

type Artifact struct {
	Abi          *abi.ABI
	DeployedCode []byte
	Code         []byte
}

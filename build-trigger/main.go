package main

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/r3labs/sse/v2"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/suave/sdk"
)

var (
	kettleAddress = common.HexToAddress("b5feafbdd752ad52afb7e1bd2e40432a485bbb7f")
	exNodeNetAddr = "http://localhost:8545"

	// This account is funded in both devnev networks
	// address: 0xb5feafbdd752ad52afb7e1bd2e40432a485bbb7f
	fundedAccount = newPrivKeyFromHex("6c45335a22461ccdb978b78ab61b238bad2fae4544fb55c14eb096c875ccfc52")
)

func newPrivKeyFromHex(hex string) *privKey {
	key, err := crypto.HexToECDSA(hex)
	if err != nil {
		panic(fmt.Sprintf("failed to parse private key: %v", err))
	}
	return &privKey{priv: key}
}

type privKey struct {
	priv  *ecdsa.PrivateKey
	nonce uint64
}

func (p *privKey) Address() common.Address {
	return crypto.PubkeyToAddress(p.priv.PublicKey)
}

func (p *privKey) MarshalPrivKey() []byte {
	return crypto.FromECDSA(p.priv)
}

func (p *privKey) Nonce() uint64 {
	return p.nonce
}

func (p *privKey) StepNonce() {
	p.nonce = p.nonce + 1
}

type Client struct {
	*sdk.Client
	Key *privKey
}

func NewClient(rpc *rpc.Client, key *privKey, addr common.Address) *Client {
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

func main() {
	mevmRpc, _ := rpc.Dial(exNodeNetAddr)
	mevmClt := NewClient(mevmRpc, fundedAccount, kettleAddress)

	contract, err := DeployContract(GatewayContract, mevmClt)

	if err != nil {
		log.Error("failed to deploy contract: %v", err)
	}

	log.Info("contract deployed at: %v", contract.Address())

	events := make(chan *sse.Event)
	client := sse.NewClient("http://localhost:8080/events")

	client.SubscribeChan("payload_attributes", events)

	for {
		select {
		case event := <-events:
			fmt.Println("Event: ", event)
		}
	}
}

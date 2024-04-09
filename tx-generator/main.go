package main

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jinmel/op-gateway-suave/utils"
)

var (
	address       = common.HexToAddress("f39Fd6e51aad88F6F4ce6aB8827279cffFb92266")
	fundedAccount = utils.NewPrivKeyFromHex("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
)

func sendEthToAddress(value uint64, address common.Address) error {

	return nil
}

func main() {
	client, err := ethclient.Dial("http://localhost:9545")

	if err != nil {
		panic(err)
	}

	balance, err := client.BalanceAt(context.Background(), address, nil)

	if err != nil {
		panic(err)
	}
	fmt.Println("balance", balance)

	nonce, err := client.NonceAt(context.Background(), address, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("nonce", nonce)
}

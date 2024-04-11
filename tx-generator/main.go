package main

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jinmel/op-gateway-suave/utils"
)

var (
	address       = common.HexToAddress("f39Fd6e51aad88F6F4ce6aB8827279cffFb92266")
	fundedAccount = utils.NewPrivKeyFromHex("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
)

func sendEthToAddress(client *ethclient.Client, value *big.Int, address common.Address, gasPrice *big.Int) error {
	nonce, err := client.PendingNonceAt(context.Background(), fundedAccount.Address())
	if err != nil {
		return err
	}

	gasLimit := uint64(21000)
	tx := types.NewTransaction(nonce, address, value, gasLimit, gasPrice, nil)
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fundedAccount.PrivateKey())
	if err != nil {
		return err
	}
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return err
	}

	fmt.Printf("tx sent: %s\n", signedTx.Hash().Hex())
	return nil
}

// sends n random recipients
func sendEthRandomRecipients(client *ethclient.Client, n int) error {
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get gas price: %w", err)
	}

	for i := 0; i < n; i++ {
		address, err := utils.RandomAddress()
		if err != nil {
			return fmt.Errorf("failed to generate random address: %w", err)
		}

		value := big.NewInt(rand.Int63n(100000000))
		fmt.Printf("sending %s wei to %s\n", value.String(), address.Hex())
		err = sendEthToAddress(client, value, address, gasPrice)
		if err != nil {
			return fmt.Errorf("failed to send eth to address: %w", err)
		}
		fmt.Printf("Sent %s eth to address %s\n", value.String(), address.Hex())
	}
	return nil
}

func main() {
	client, err := ethclient.Dial("http://localhost:5545")

	if err != nil {
		panic(err)
	}

	balance, err := client.BalanceAt(context.Background(), address, nil)

	if err != nil {
		panic(err)
	}
	fmt.Println("balance", balance)

	fmt.Println("Starting loop")

	// sends batch of txs every block time.
	ticker := time.NewTicker(2 * time.Second)
	batchSize := 20
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			fmt.Println("Sending batch")
			sendEthRandomRecipients(client, batchSize)
			fmt.Println("done sending batch")
		}
	}
}

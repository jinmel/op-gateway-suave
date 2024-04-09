package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/jinmel/op-gateway-suave/utils"
)

var (
	kettleAddress = common.HexToAddress("b5feafbdd752ad52afb7e1bd2e40432a485bbb7f")
	suaveAddr     = "http://localhost:11545"

	// This account is funded in both devnev networks
	// address: 0xb5feafbdd752ad52afb7e1bd2e40432a485bbb7f
	fundedAccount = utils.NewPrivKeyFromHex("6c45335a22461ccdb978b78ab61b238bad2fae4544fb55c14eb096c875ccfc52")
)

func main() {
	mevmRpc, _ := rpc.Dial(suaveAddr)
	mevmClt := utils.NewClient(mevmRpc, fundedAccount, kettleAddress)

	balance, err := mevmClt.Client.RPC().BalanceAt(context.Background(), fundedAccount.Address(), nil)

	fmt.Println("Address: ", fundedAccount.Address())
	fmt.Println("balance: ", balance)

	contract, err := DeployContract(GatewayContract, mevmClt)

	if err != nil {
		fmt.Printf("failed to deploy contract: %v\n", err)
		return
	}

	fmt.Printf("contract deployed at: %v\n", contract.Address())

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Println("Triggering build")
			// peekers := []common.Address{common.HexToAddress("0x0000000000000000000000000000000042100001")}
			result, err := contract.SendTransaction("build", []interface{}{
				uint64(0), // not used.
				"local",   // submits to local relay.
				[]common.Address{common.HexToAddress("0xC8df3686b4Afb2BB53e60EAe97EF043FE03Fb829")},
				[]common.Address{},
			}, []byte{0x20, 0x30})

			if err != nil {
				fmt.Println("failed to trigger build: ", err)
				continue
			}
			fmt.Printf("build triggered: %+v\n", result)
		}
	}
}

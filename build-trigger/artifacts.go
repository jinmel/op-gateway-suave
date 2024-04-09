package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/jinmel/op-gateway-suave/utils"
)

var (
	GatewayContract = newArtifact("Builder.sol/Builder.json")
)

func newArtifact(name string) *utils.Artifact {
	// Get the caller's file path.
	_, filename, _, _ := runtime.Caller(1)

	// Resolve the directory of the caller's file.
	callerDir := filepath.Dir(filename)

	// Construct the absolute path to the target file.
	targetFilePath := filepath.Join(callerDir, "../out", name)

	data, err := os.ReadFile(targetFilePath)
	if err != nil {
		panic(fmt.Sprintf("failed to read artifact %s: %v. Maybe you forgot to generate the artifacts? `cd suave && forge build`", name, err))
	}

	var artifactObj struct {
		Abi              *abi.ABI `json:"abi"`
		DeployedBytecode struct {
			Object string
		} `json:"deployedBytecode"`
		Bytecode struct {
			Object string
		} `json:"bytecode"`
	}
	if err := json.Unmarshal(data, &artifactObj); err != nil {
		panic(fmt.Sprintf("failed to unmarshal artifact %s: %v", name, err))
	}

	return &utils.Artifact{
		Abi:          artifactObj.Abi,
		Code:         hexutil.MustDecode(artifactObj.Bytecode.Object),
		DeployedCode: hexutil.MustDecode(artifactObj.DeployedBytecode.Object),
	}
}

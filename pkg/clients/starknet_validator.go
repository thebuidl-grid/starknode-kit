package clients

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/thebuidl-grid/starknode-kit/pkg/constants"
	"github.com/thebuidl-grid/starknode-kit/pkg/process"
)

type stakingValidator struct {
	provider stakingValidatorProviderConfig
	wallet   stakingValidatorWalletConfig
}

type stakingValidatorProviderConfig struct {
	starkentHttp string
	starkentWS   string
}

type stakingValidatorWalletConfig struct {
	address    string
	privatekey string
}

func (_ stakingValidator) getCommand() string {
	return filepath.Join(constants.InstallStarknetDir, "starknet-staking-v2", "validator")
}

func (c stakingValidator) buildArgs() []string {
	args := []string{
		"--provider-http" + c.provider.starkentHttp,
		"--provider-ws" + c.provider.starkentWS,
		"--signer-op-address" + c.wallet.address,
		"--signer-priv-key" + c.wallet.privatekey,
	}
	return args
}

func (c *stakingValidator) Start() error {
	args := c.buildArgs()
	command := c.getCommand()
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filepath.Join(constants.InstallStarknetDir, "starknet-staking-v2", "validator")

	logFilePath := filepath.Join(constants.InstallStarknetDir, "starknet-staking-v2",
		"logs",
		fmt.Sprintf("starknet-staking-v2_%s.log", timestamp))
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}

	return process.StartClient("staking-validator", command, logFile, args...)
}

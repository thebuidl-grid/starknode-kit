package clients

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/thebuidl-grid/starknode-kit/pkg/constants"
	"github.com/thebuidl-grid/starknode-kit/pkg/process"
)

type StakingValidator struct {
	Provider stakingValidatorProviderConfig
	Wallet   stakingValidatorWalletConfig
}

type stakingValidatorProviderConfig struct {
	starknetHttp string
	starkentWS   string
}

type stakingValidatorWalletConfig struct {
	address    string
	privatekey string
}

func (_ StakingValidator) getCommand() string {
	return filepath.Join(constants.InstallStarknetDir, "starknet-staking-v2", "validator")
}

func (c StakingValidator) buildArgs() []string {
	args := []string{
		"--provider-http", c.Provider.starknetHttp,
		"--provider-ws", c.Provider.starkentWS,
		"--signer-op-address", c.Wallet.address,
		"--signer-priv-key", c.Wallet.privatekey,
	}
	return args
}

func (c *StakingValidator) Start() error {
	args := c.buildArgs()
	command := c.getCommand()
	timestamp := time.Now().Format("2006-01-02_15-04-05")

	logFilePath := filepath.Join(constants.InstallStarknetDir, "starknet-staking-v2",
		"logs",
		fmt.Sprintf("starknet-staking-v2_%s.log", timestamp))
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	return process.StartClient("staking-validator", command, multiWriter, args...)
}

<<<<<<< HEAD
package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
)

const (
	predeployedClassHash = "0x61dac032f228abef9c6626f995015233097ae253a7f72d68552db02f2971b8f"
	rpcURL               = "https://starknet-sepolia.public.blastapi.io/rpc/v0_8"
	strkTokenAddress     = "0x04718f5a0fc34cc1af16a1cdee98ffb20c31f5cd61d6ab07201858f4287c938d"
)

// checkBalance queries the STRK balance of the given address
func checkBalance(client *rpc.Provider, address *felt.Felt) (*felt.Felt, error) {
	strkAddr, err := utils.HexToFelt(strkTokenAddress)
	if err != nil {
		return nil, err
	}

	// Call balanceOf function on STRK token contract
	callReq := rpc.FunctionCall{
		ContractAddress:    strkAddr,
		EntryPointSelector: utils.GetSelectorFromNameFelt("balanceOf"),
		Calldata:           []*felt.Felt{address},
	}

	result, err := client.Call(context.Background(), callReq, rpc.BlockID{Tag: "latest"})
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return &felt.Zero, nil
	}

	balance := result[0]
	return balance, nil
}

// monitorFunding monitors the account balance and notifies when funded
func monitorFunding(client *rpc.Provider, address *felt.Felt, requiredAmount *felt.Felt, funded chan<- bool) {
	fmt.Printf("🔍 Monitoring funding for address: %s\n", FormatStarknetAddress(address))
	fmt.Printf("📊 Required amount: %.6f STRK\n", utils.FRIToSTRK(requiredAmount))
	fmt.Println("⏳ Checking balance every 10 seconds...")

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			balance, err := checkBalance(client, address)
			if err != nil {
				fmt.Printf("❌ Error checking balance: %v\n", err)
				continue
			}

			if balance.Cmp(requiredAmount) >= 0 {
				fmt.Println("✅ Account has been funded! Proceeding with deployment...")
				funded <- true
				return
			}
		}
	}
}

func DeployAccount() error {
	client, err := createRPCProvider()
	if err != nil {
		return fmt.Errorf("failed to create RPC provider: %w", err)
	}

	ks, pub, priv := generateKeys()

	accnt, err := createAccount(client, pub, ks)
	if err != nil {
		return fmt.Errorf("failed to create account: %w", err)
	}

	classHash, err := getClassHash()
	if err != nil {
		return fmt.Errorf("failed to get class hash: %w", err)
	}

	deployTxn, precomputedAddr, err := buildDeployTransaction(accnt, pub, classHash)
	if err != nil {
		return fmt.Errorf("failed to build deploy transaction: %w", err)
	}

	requiredAmount, err := displayFundingInfoAndStartMonitoring(deployTxn, precomputedAddr)
	if err != nil {
		return fmt.Errorf("failed to start funding monitoring: %w", err)
	}

	waitForFundingWithMonitoring(client, precomputedAddr, requiredAmount)

	resp, err := executeDeployment(accnt, deployTxn)
	if err != nil {
		return fmt.Errorf("failed to execute deployment: %w", err)
	}

	fmt.Println("✅ Account deployment transaction successfully submitted!")
	fmt.Printf("🔗 Transaction hash: %v\n", FormatTransactionHash(resp.Hash))
	fmt.Printf("📍 Contract address: %v\n", FormatStarknetAddress(resp.ContractAddress))
	walletKS := map[string]string{
		"STARKNET_WALLET":      FormatStarknetAddress(resp.ContractAddress),
		"STARKNET_PRIVATE_KEY": FormatStarknetAddress(priv),
	}
	err = writeToENV(walletKS)
	if err != nil {
		fmt.Println("Error writing to env file")
		return err
	}
	fmt.Println("⏰ Wait a few minutes to see it in Voyager.")

	return err
}

// createRPCProvider initializes and returns an RPC provider
func createRPCProvider() (*rpc.Provider, error) {
	client, err := rpc.NewProvider(rpcURL)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// generateKeys generates a new keypair for the account
func generateKeys() (*account.MemKeystore, *felt.Felt, *felt.Felt) {
	return account.GetRandomKeys()
}

// createAccount creates a new account instance
func createAccount(client *rpc.Provider, pub *felt.Felt, ks *account.MemKeystore) (*account.Account, error) {
	accnt, err := account.NewAccount(client, pub, pub.String(), ks, account.CairoV2)
	if err != nil {
		return nil, err
	}
	return accnt, nil
}

// getClassHash converts the predefined class hash string to felt
func getClassHash() (*felt.Felt, error) {
	classHash, err := utils.HexToFelt(predeployedClassHash)
	if err != nil {
		return nil, err
	}
	return classHash, nil
}

// buildDeployTransaction builds and estimates the deploy account transaction
func buildDeployTransaction(accnt *account.Account, pub *felt.Felt, classHash *felt.Felt) (*rpc.BroadcastDeployAccountTxnV3, *felt.Felt, error) {
	deployAccountTxn, precomputedAddress, err := accnt.BuildAndEstimateDeployAccountTxn(
		context.Background(),
		pub,
		classHash,
		[]*felt.Felt{pub},
		nil,
	)
	if err != nil {
		return nil, nil, err
	}
	return deployAccountTxn, precomputedAddress, nil
}

// displayFundingInfoAndStartMonitoring shows the precomputed address, required funding amount, and returns the required amount
func displayFundingInfoAndStartMonitoring(deployTxn *rpc.BroadcastDeployAccountTxnV3, precomputedAddr *felt.Felt) (*felt.Felt, error) {
	fmt.Printf("🏠 PrecomputedAddress: %s\n", FormatStarknetAddress(precomputedAddr))

	overallFee, err := utils.ResBoundsMapToOverallFee(deployTxn.ResourceBounds, 1)
	if err != nil {
		return nil, err
	}

	feeInSTRK := utils.FRIToSTRK(overallFee)

	fmt.Println("\n📋 Funding Instructions:")
	fmt.Printf("💸 The precomputed address needs approximately %.6f STRK to perform the transaction.\n", feeInSTRK)
	fmt.Printf("🚰 You can use the Starknet faucet or send STRK to: %s\n", FormatStarknetAddress(precomputedAddr))

	return overallFee, nil
}

// waitForFundingWithMonitoring uses goroutine to monitor funding automatically
func waitForFundingWithMonitoring(client *rpc.Provider, precomputedAddr *felt.Felt, requiredAmount *felt.Felt) {
	funded := make(chan bool)

	go monitorFunding(client, precomputedAddr, requiredAmount, funded)

	// Wait for funding notification
	<-funded
}

func executeDeployment(accnt *account.Account, deployTxn *rpc.BroadcastDeployAccountTxnV3) (*rpc.TransactionResponse, error) {
	fmt.Println("🚀 Deploying account...")
	resp, err := accnt.SendTransaction(context.Background(), deployTxn)
	if err != nil {
		fmt.Println("❌ Error returned from SendTransaction:")
		return nil, err
	}
	return resp, nil
}

// NOTE utility func not using
// CheckAccountBalance is a utility function to check STRK balance of any address
func CheckAccountBalance(address string) (float64, error) {
	client, err := createRPCProvider()
	if err != nil {
		return 0, fmt.Errorf("failed to create RPC provider: %w", err)
	}

	addr, err := utils.HexToFelt(address)
	if err != nil {
		return 0, fmt.Errorf("invalid address format: %w", err)
	}

	balance, err := checkBalance(client, addr)
	if err != nil {
		return 0, fmt.Errorf("failed to check balance: %w", err)
	}

	return utils.FRIToSTRK(balance), nil
}

// NOTE currently not using
// MonitorAddressFunding is a utility function to monitor any address for funding
func MonitorAddressFunding(address string, requiredAmount float64, callback func()) error {
	client, err := createRPCProvider()
	if err != nil {
		return fmt.Errorf("failed to create RPC provider: %w", err)
	}

	addr, err := utils.HexToFelt(address)
	if err != nil {
		return fmt.Errorf("invalid address format: %w", err)
	}

	requiredFRI := utils.STRKToFRI(requiredAmount)

	funded := make(chan bool)
	go monitorFunding(client, addr, requiredFRI, funded)

	<-funded
	if callback != nil {
		callback()
	}

	return nil
}

package utils

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/thebuidl-grid/starknode-kit/pkg/constants"
	"github.com/thebuidl-grid/starknode-kit/pkg/types"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/rpc"
	starkutils "github.com/NethermindEth/starknet.go/utils"
)

// checkBalance queries the STRK balance of the given address
func CheckBalance(client *rpc.Provider, address *felt.Felt) (*felt.Felt, error) {
	strkAddr, err := starkutils.HexToFelt(constants.StrkTokenAddress)
	if err != nil {
		return nil, err
	}

	// Call balanceOf function on STRK token contract
	callReq := rpc.FunctionCall{
		ContractAddress:    strkAddr,
		EntryPointSelector: starkutils.GetSelectorFromNameFelt("balanceOf"),
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

// CreateRPCProvider initializes and returns an RPC provider
func CreateRPCProvider(network string) (*rpc.Provider, error) {
	url, ok := constants.RPCURL[network]
	if !ok {
		return nil, fmt.Errorf("Invalid network: %s", network)
	}
	client, err := rpc.NewProvider(url)
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
	classHash, err := starkutils.HexToFelt(constants.PredeployedClassHash)
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

	overallFee, err := starkutils.ResBoundsMapToOverallFee(deployTxn.ResourceBounds, 1)
	if err != nil {
		return nil, err
	}

	feeInSTRK := starkutils.FRIToSTRK(overallFee)

	fmt.Println("\n📋 Funding Instructions:")
	fmt.Printf("💸 The precomputed address needs approximately %.6f STRK to perform the transaction.\n", feeInSTRK)
	fmt.Printf("🚰 You can use the Starknet faucet or send STRK to: %s\n", FormatStarknetAddress(precomputedAddr))

	return overallFee, nil
}

// waitForFundingWithMonitoring waits for the user to fund the account and press Enter.
func waitForFundingWithMonitoring(client *rpc.Provider, precomputedAddr *felt.Felt, requiredAmount *felt.Felt) {
	fmt.Println("Press Enter to check balance and proceed with deployment...")
	fmt.Scanln()

	for {
		balance, err := CheckBalance(client, precomputedAddr)
		if err != nil {
			fmt.Printf("❌ Error checking balance: %v\n", err)
			fmt.Println("Press Enter to try again...")
			fmt.Scanln()
			continue
		}

		if balance.Cmp(requiredAmount) >= 0 {
			fmt.Printf("✅ Sufficient balance found: %.6f STRK. Proceeding with deployment...\n", starkutils.FRIToSTRK(balance))
			return
		}

		fmt.Printf("❌ Insufficient balance: %.6f STRK. Required: %.6f STRK\n", starkutils.FRIToSTRK(balance), starkutils.FRIToSTRK(requiredAmount))
		fmt.Println("Please fund the account and press Enter to re-check balance...")
		fmt.Scanln()
	}
}

func executeDeployment(accnt *account.Account, deployTxn *rpc.BroadcastDeployAccountTxnV3) (rpc.TransactionResponse, error) {
	fmt.Println("🚀 Deploying account...")
	resp, err := accnt.SendTransaction(context.Background(), deployTxn)
	if err != nil {
		fmt.Println("❌ Error returned from SendTransaction:")
		return rpc.TransactionResponse{}, err
	}
	return resp, nil
}

func DeployAccount(netowork string) (*types.Wallet, error) {
	client, err := CreateRPCProvider(netowork)
	if err != nil {
		return nil, fmt.Errorf("failed to create RPC provider: %w", err)
	}

	ks, pub, priv := generateKeys()

	accnt, err := createAccount(client, pub, ks)
	if err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	classHash, err := getClassHash()
	if err != nil {
		return nil, fmt.Errorf("failed to get class hash: %w", err)
	}

	deployTxn, precomputedAddr, err := buildDeployTransaction(accnt, pub, classHash)
	if err != nil {
		return nil, fmt.Errorf("failed to build deploy transaction: %w", err)
	}
	requiredAmount, err := displayFundingInfoAndStartMonitoring(deployTxn, precomputedAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to start funding monitoring: %w", err)
	}

	waitForFundingWithMonitoring(client, precomputedAddr, requiredAmount)

	resp, err := executeDeployment(accnt, deployTxn)
	if err != nil {
		return nil, fmt.Errorf("failed to execute deployment: %w", err)
	}

	fmt.Println("✅ Account deployment transaction successfully submitted!")
	fmt.Printf("🔗 Transaction hash: %v\n", FormatTransactionHash(resp.Hash))
	fmt.Printf("📍 Contract address: %v\n", FormatStarknetAddress(resp.ContractAddress))

	// Set all wallet-related environment variables for validator configuration
	// These variables will be used in the config YAML with ${VAR_NAME} syntax
	walletKS := map[string]string{
		"STARKNET_WALLET":      FormatStarknetAddress(resp.ContractAddress), // Wallet contract address
		"STARKNET_CLASS_HASH":  FormatStarknetAddress(classHash),            // Account contract class hash
		"STARKNET_PRIVATE_KEY": FormatStarknetAddress(priv),                 // Private key for signing
		"STARKNET_PUBLIC_KEY":  FormatStarknetAddress(pub),                  // Public key derived from private key
		"STARKNET_SALT":        FormatStarknetAddress(pub),                  // Salt used for deployment (using pub as salt)
		"STARKNET_DEPLOYED":    "true",                                      // Account deployment status
		"STARKNET_LEGACY":      "false",                                     // Account type (Cairo v2)
	}
	err = writeToENV(walletKS)
	if err != nil {
		fmt.Println("Error writing to env file")
		return nil, err
	}
	fmt.Println("⏰ Wait a few minutes to see it in Voyager.")

	// Create and return the Wallet struct
	wallet := &types.Wallet{
		Address:    FormatStarknetAddress(resp.ContractAddress),
		ClassHash:  FormatStarknetAddress(classHash),
		Deployed:   true,
		Legacy:     false,
		PrivateKey: FormatStarknetAddress(priv),
		PublicKey:  FormatStarknetAddress(pub),
		Salt:       FormatStarknetAddress(pub), // Using pub as salt
	}
	transactionUrl := fmt.Sprintf("https://sepolia.voyager.online/tx/%s", FormatTransactionHash(resp.Hash))
	fmt.Println("Transaction successfull, view here: ", transactionUrl)

	return wallet, nil
}

func StakeStark(network string, wallet types.WalletConfig) error {

	stackingAddr, err := starkutils.HexToFelt(constants.StakingContract)
	if err != nil {
		return fmt.Errorf("failed to convert staking contrct addr to felt: %w", err)
	}

	commisionIntConv, err := strconv.Atoi(wallet.StakeCommision)
	if err != nil {
		return err
	}
	commisionInt := new(big.Int).SetUint64(uint64(commisionIntConv * 100))

	commisionFelt := starkutils.BigIntToFelt(commisionInt)
	starkTokenAdress, err := starkutils.HexToFelt(constants.StrkTokenAddress)
	if err != nil {
		return fmt.Errorf("failed to convert starktoken address to felt: %w", err)
	}
	rewardAddress, err := starkutils.HexToFelt(wallet.RewardAddress)
	if err != nil {
		return err
	}

	userWalletAddress, err := starkutils.HexToFelt(wallet.Wallet.Address)
	if err != nil {
		return err
	}

	client, err := CreateRPCProvider(network)
	if err != nil {
		return err
	}

	ks := account.NewMemKeystore()
	privKeyBI, ok := new(big.Int).SetString(wallet.Wallet.PrivateKey, 0)
	if !ok {
		return fmt.Errorf("Fail to convert privKey to bitInt")
	}
	ks.Put(wallet.Wallet.PublicKey, privKeyBI)

	userAccount, _ := account.NewAccount(
		client,
		userWalletAddress,
		wallet.Wallet.PublicKey,
		ks,
		account.CairoV2,
	)

	txn1 := rpc.InvokeFunctionCall{
		ContractAddress: starkTokenAdress,
		FunctionName:    "approve",
		CallData:        []*felt.Felt{stackingAddr, constants.Stakes[network][0], constants.Stakes[network][1]},
	}
	txn2 := rpc.InvokeFunctionCall{
		ContractAddress: stackingAddr,
		FunctionName:    "stake",
		CallData:        []*felt.Felt{rewardAddress, userAccount.Address, constants.Stakes[network][0]},
	}
	txn3 := rpc.InvokeFunctionCall{
		ContractAddress: stackingAddr,
		FunctionName:    "set_commission",
		CallData:        []*felt.Felt{commisionFelt},
	}

	resp, err := userAccount.BuildAndSendInvokeTxn(context.Background(), []rpc.InvokeFunctionCall{txn1, txn2, txn3}, nil)
	if err != nil {
		return err
	}

	txnReceipt, err := userAccount.WaitForTransactionReceipt(context.Background(), resp.Hash, 10*time.Second)
	if err != nil {
		transactionUrl := fmt.Sprintf("https://sepolia.voyager.online/tx/%s", FormatTransactionHash(resp.Hash))
		fmt.Println("Transaction error, view here: ", transactionUrl)
		return err
	}
	transactionUrl := fmt.Sprintf("https://sepolia.voyager.online/tx/%s", FormatTransactionHash(txnReceipt.Hash))
	fmt.Println("Transaction successfull, view here: ", transactionUrl)
	return nil
}

// TODO delegation pool
// NOTE check balance may return early before the transaction is complete

package clients

import (
	"path/filepath"
	"testing"

	"github.com/thebuidl-grid/starknode-kit/pkg/constants"
	"github.com/thebuidl-grid/starknode-kit/pkg/types"
)

func TestGethClient(t *testing.T) {
	config := &gethConfig{
		port:          30303,
		executionType: "full",
		network:       "mainnet",
	}

	args := config.buildArgs()

	expectedArgs := []string{
		"--mainnet",
		"--port=30303",
		"--discovery.port=30303",
		"--http",
		"--http.api=eth,net,engine,admin",
		"--http.corsdomain=*",
		"--http.addr=0.0.0.0",
		"--http.port=8545",
		"--authrpc.jwtsecret=" + constants.JWTPath,
		"--authrpc.addr=0.0.0.0",
		"--authrpc.port=8551",
		"--authrpc.vhosts=*",
		"--metrics",
		"--metrics.addr=0.0.0.0",
		"--metrics.port=7878",
		"--syncmode=snap",
		"--datadir=" + filepath.Join(constants.InstallClientsDir, "geth", "database"),
	}

	if len(args) != len(expectedArgs) {
		t.Errorf("Expected %d arguments, got %d", len(expectedArgs), len(args))
	}

	for i, expected := range expectedArgs {
		if args[i] != expected {
			t.Errorf("Expected argument %d to be '%s', got '%s'", i, expected, args[i])
		}
	}

	
}

func TestRethClient(t *testing.T) {
	config := &rethConfig{
		port:          30303,
		executionType: "full",
		network:       "mainnet",
	}

	args := config.buildArgs()

	expectedArgs := []string{
		"node",
		"--chain", "mainnet",
		"--http",
		"--http.addr", "0.0.0.0",
		"--http.port", "8545",
		"--http.api", "eth,net,admin",
		"--http.corsdomain", "*",
		"--authrpc.addr", "0.0.0.0",
		"--authrpc.port", "8551",
		"--authrpc.jwtsecret", constants.JWTPath,
		"--port", "30303",
		"--metrics", "0.0.0.0:7878",
		"--datadir", filepath.Join(constants.InstallClientsDir, "reth", "database"),
	}

	if len(args) != len(expectedArgs) {
		t.Errorf("Expected %d arguments, got %d", len(expectedArgs), len(args))
	}

	for i, expected := range expectedArgs {
		if args[i] != expected {
			t.Errorf("Expected argument %d to be '%s', got '%s'", i, expected, args[i])
		}
	}

	
}

func TestLighthouseClient(t *testing.T) {
	config := &lightHouseConfig{
		port:                []int{9000, 9001},
		consensusCheckpoint: "https://checkpoint.sync",
		network:             "mainnet",
	}

	args := config.buildArgs()

	expectedArgs := []string{
		"bn",
		"--network",
		"mainnet",
		"--port=9000",
		"--quic-port=9001",
		"--execution-endpoint",
		"http://localhost:8551",
		"--checkpoint-sync-url",
		"https://checkpoint.sync",
		"--checkpoint-sync-url-timeout",
		"1200",
		"--disable-deposit-contract-sync",
		"--execution-jwt",
		constants.JWTPath,
		"--metrics",
		"--metrics-address",
		"127.0.0.1",
		"--metrics-port",
		"5054",
		"--http",
		"--disable-upnp",
		"--datadir=" + filepath.Join(constants.InstallClientsDir, "lighthouse", "database"),
	}

	if len(args) != len(expectedArgs) {
		t.Errorf("Expected %d arguments, got %d", len(expectedArgs), len(args))
	}

	for i, expected := range expectedArgs {
		if args[i] != expected {
			t.Errorf("Expected argument %d to be '%s', got '%s'", i, expected, args[i])
		}
	}

	
}

func TestPrysmClient(t *testing.T) {
	config := &prysmConfig{
		port:                []int{9000, 9001},
		consensusCheckpoint: "https://checkpoint.sync",
		network:             "mainnet",
	}

	args := config.buildArgs()

	expectedArgs := []string{
		"beacon-chain",
		"--mainnet",
		"--p2p-udp-port=9001",
		"--p2p-quic-port=9000",
		"--p2p-tcp-port=9000",
		"--execution-endpoint",
		"http://localhost:8551",
		"--grpc-gateway-host=0.0.0.0",
		"--grpc-gateway-port=5052",
		"--checkpoint-sync-url=https://checkpoint.sync",
		"--genesis-beacon-api-url=https://checkpoint.sync",
		"--accept-terms-of-use=true",
		"--jwt-secret",
		constants.JWTPath,
		"--monitoring-host",
		"127.0.0.1",
		"--monitoring-port",
		"5054",
		"--datadir=" + filepath.Join(constants.InstallClientsDir, "prsym", "database"),
	}

	if len(args) != len(expectedArgs) {
		t.Errorf("Expected %d arguments, got %d", len(expectedArgs), len(args))
	}

	for i, expected := range expectedArgs {
		if args[i] != expected {
			t.Errorf("Expected argument %d to be '%s', got '%s'", i, expected, args[i])
		}
	}

	
}

func TestJunoClient(t *testing.T) {
	config := types.JunoConfig{
		Port:    6060,
		EthNode: "ws://localhost:8546",
	}

	client := &JunoClient{config: config, network: "mainnet"}
	args := client.buildJunoArgs()

	expectedArgs := []string{
		"--http",
		"--http-port=6060",
		"--http-host=0.0.0.0",
		"--db-path=" + filepath.Join(constants.InstallStarknetDir, "juno", "database"),
		"--eth-node=ws://localhost:8546",
		"--ws=false",
		"--ws-port=6061",
		"--ws-host=0.0.0.0",
		"--network=mainnet",
	}

	if len(args) != len(expectedArgs) {
		t.Errorf("Expected %d arguments, got %d", len(expectedArgs), len(args))
	}

	for i, expected := range expectedArgs {
		if args[i] != expected {
			t.Errorf("Expected argument %d to be '%s', got '%s'", i, expected, args[i])
		}
	}

	
}

func TestStarknetValidatorClient(t *testing.T) {
	config := &StakingValidator{
		Provider: stakingValidatorProviderConfig{
			starknetHttp: "http://localhost:6060",
			starkentWS:   "ws://localhost:6061",
		},
		Wallet: stakingValidatorWalletConfig{
			address:    "0x123",
			privatekey: "0x456",
		},
	}

	args := config.buildArgs()

	expectedArgs := []string{
		"--provider-http", "http://localhost:6060",
		"--provider-ws", "ws://localhost:6061",
		"--signer-op-address", "0x123",
		"--signer-priv-key", "0x456",
	}

	if len(args) != len(expectedArgs) {
		t.Errorf("Expected %d arguments, got %d", len(expectedArgs), len(args))
	}

	for i, expected := range expectedArgs {
		if args[i] != expected {
			t.Errorf("Expected argument %d to be '%s', got '%s'", i, expected, args[i])
		}
	}

	
}

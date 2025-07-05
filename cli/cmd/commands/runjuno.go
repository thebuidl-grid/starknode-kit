package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"starknode-kit/pkg/clients"

	"github.com/spf13/cobra"
)

var (
	junoNetwork string
	junoPort    string
	junoDataDir string
	junoEthNode string
	useSnapshot bool
)

// RunJunoCmd represents the run juno command
var RunJunoCmd = &cobra.Command{
	Use:   "run-juno",
	Short: "Run a local Juno Starknet node",
	Long: `Run a local Juno Starknet node with configurable options.
Juno is a Go-based Starknet node implementation by Nethermind that provides
full JSON-RPC support for Starknet networks.

Juno requires an Ethereum node connection to verify L1 state. You can specify
an Ethereum node URL using the --eth-node flag.

Example:
  starknode-kit run-juno --network mainnet --port 6060 --data-dir ./juno-data --eth-node ws://localhost:8546`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create data directory if it doesn't exist
		if err := os.MkdirAll(junoDataDir, 0755); err != nil {
			return fmt.Errorf("failed to create data directory: %w", err)
		}

		// Create absolute path for data directory
		absDataDir, err := filepath.Abs(junoDataDir)
		if err != nil {
			return fmt.Errorf("failed to get absolute path for data directory: %w", err)
		}

		// Create Juno configuration
		config := &clients.JunoConfig{
			Network:     junoNetwork,
			Port:        junoPort,
			UseSnapshot: useSnapshot,
			DataDir:     absDataDir,
			EthNode:     junoEthNode,
			Environment: []string{
				fmt.Sprintf("JUNO_NETWORK=%s", junoNetwork),
				fmt.Sprintf("JUNO_HTTP_PORT=%s", junoPort),
				"JUNO_HTTP_HOST=0.0.0.0",
			},
		}

		// Create Juno client
		junoClient, err := clients.NewJunoClient(config)
		if err != nil {
			return fmt.Errorf("failed to create juno client: %w", err)
		}
		defer junoClient.Close()

		// Create context with cancellation
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Start the node
		fmt.Println("Starting Juno Starknet node...")
		if err := junoClient.StartNode(ctx); err != nil {
			return fmt.Errorf("failed to start juno node: %w", err)
		}

		// Get node status
		status, err := junoClient.GetNodeStatus(ctx)
		if err != nil {
			return fmt.Errorf("failed to get node status: %w", err)
		}

		fmt.Printf("Juno Starknet node is running with status: %s\n", status)
		fmt.Printf("Node data directory: %s\n", absDataDir)
		fmt.Printf("HTTP endpoint: http://localhost:%s\n", junoPort)
		fmt.Printf("Metrics endpoint: http://localhost:6060\n")

		// Wait for user interrupt
		fmt.Println("\nPress Ctrl+C to stop the node...")
		<-ctx.Done()

		// Stop the node
		fmt.Println("\nStopping Juno Starknet node...")
		if err := junoClient.StopNode(ctx); err != nil {
			return fmt.Errorf("failed to stop juno node: %w", err)
		}

		return nil
	},
}

func init() {
	// Add flags
	RunJunoCmd.Flags().StringVar(&junoNetwork, "network", "mainnet", "Network to connect to (mainnet, sepolia, sepolia-integration)")
	RunJunoCmd.Flags().StringVar(&junoPort, "port", "6060", "Port to run the node on")
	RunJunoCmd.Flags().StringVar(&junoDataDir, "data-dir", "./juno-data", "Directory to store node data")
	RunJunoCmd.Flags().StringVar(&junoEthNode, "eth-node", "ws://localhost:8546", "Ethereum node WebSocket URL (required for L1 verification)")
	RunJunoCmd.Flags().BoolVar(&useSnapshot, "use-snapshot", true, "Whether to use snapshots for faster sync")
}

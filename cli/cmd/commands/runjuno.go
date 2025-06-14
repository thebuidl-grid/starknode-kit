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
	useSnapshot bool
)

// RunJunoCmd represents the run juno command
var RunJunoCmd = &cobra.Command{
	Use:   "run-juno",
	Short: "Run a local Nethermind Juno node",
	Long: `Run a local Nethermind Juno node with configurable options.
Example:
  starknode-kit run-juno --network mainnet --port 5050 --data-dir ./juno-data`,
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
		fmt.Println("Starting Juno node...")
		if err := junoClient.StartNode(ctx); err != nil {
			return fmt.Errorf("failed to start juno node: %w", err)
		}

		// Get node status
		status, err := junoClient.GetNodeStatus(ctx)
		if err != nil {
			return fmt.Errorf("failed to get node status: %w", err)
		}

		fmt.Printf("Juno node is running with status: %s\n", status)
		fmt.Printf("Node data directory: %s\n", absDataDir)
		fmt.Printf("HTTP endpoint: http://localhost:%s\n", junoPort)

		// Wait for user interrupt
		fmt.Println("\nPress Ctrl+C to stop the node...")
		<-ctx.Done()

		// Stop the node
		fmt.Println("\nStopping Juno node...")
		if err := junoClient.StopNode(ctx); err != nil {
			return fmt.Errorf("failed to stop juno node: %w", err)
		}

		return nil
	},
}

func init() {
	// Add flags
	RunJunoCmd.Flags().StringVar(&junoNetwork, "network", "mainnet", "Network to connect to (mainnet, testnet)")
	RunJunoCmd.Flags().StringVar(&junoPort, "port", "5050", "Port to run the node on")
	RunJunoCmd.Flags().StringVar(&junoDataDir, "data-dir", "./juno-data", "Directory to store node data")
	RunJunoCmd.Flags().BoolVar(&useSnapshot, "use-snapshot", true, "Whether to use snapshots for faster sync")
}

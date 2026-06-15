// Package commands exposes cobra sub-commands for managing a stellar-core node.
// Wire these into the root CLI with: root.AddCommand(stellar.Commands())
package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thebuidl-grid/starknode-kit/stellar"
)

// Commands returns the top-level 'stellar' cobra command with all sub-commands attached
func Commands() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stellar",
		Short: "Manage a Stellar Core (stellar-core) validator node",
	}

	cmd.AddCommand(
		installCmd(),
		startCmd(),
		statusCmd(),
		httpCmd(),
		initDBCmd(),
		removeCmd(),
		versionCmd(),
	)

	return cmd
}

// --- flags shared across sub-commands ---

var (
	flagNetwork    string
	flagConfigPath string
	flagHTTPPort   int
)

func addConfigFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&flagNetwork, "network", stellar.NetworkPubnet, `Network to join: "pubnet" or "testnet"`)
	cmd.Flags().StringVar(&flagConfigPath, "conf", "", "Path to stellar-core.cfg (defaults to managed path)")
	cmd.Flags().IntVar(&flagHTTPPort, "http-port", stellar.DefaultHTTPPort, "stellar-core local HTTP admin port")
}

func buildConfig() stellar.Config {
	return stellar.Config{
		Network:    flagNetwork,
		ConfigPath: flagConfigPath,
		HTTPPort:   flagHTTPPort,
	}
}

// --- sub-commands ---

func installCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "install",
		Short: "Install stellar-core via the platform package manager",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Installing stellar-core...")
			return stellar.Install()
		},
	}
}

func startCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start stellar-core in the background",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := buildConfig()
			client, err := stellar.New(cfg)
			if err != nil {
				return err
			}
			return client.Start()
		},
	}
	addConfigFlags(cmd)
	return cmd
}

func statusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show stellar-core node status (calls the info HTTP command)",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := buildConfig()
			client, err := stellar.New(cfg)
			if err != nil {
				return err
			}
			info, err := client.Info()
			if err != nil {
				return err
			}
			fmt.Println(info)
			return nil
		},
	}
	addConfigFlags(cmd)
	return cmd
}

func httpCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "http-command <command>",
		Short: "Send an HTTP command to the running stellar-core instance",
		Long: `Sends a command to the local stellar-core HTTP admin endpoint.
Common commands: info, peers, quorum, scp, metrics, catchup, self-check`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := buildConfig()
			client, err := stellar.New(cfg)
			if err != nil {
				return err
			}
			resp, err := client.HTTPCommand(args[0])
			if err != nil {
				return err
			}
			fmt.Println(resp)
			return nil
		},
	}
	addConfigFlags(cmd)
	return cmd
}

func initDBCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init-db",
		Short: "Initialise the stellar-core database (run once before first start)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return stellar.InitDB(buildConfig())
		},
	}
	addConfigFlags(cmd)
	return cmd
}

func removeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove",
		Short: "Remove the managed stellar-core installation",
		RunE: func(cmd *cobra.Command, args []string) error {
			return stellar.Remove()
		},
	}
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the installed stellar-core version",
		Run: func(cmd *cobra.Command, args []string) {
			v := stellar.InstalledVersion()
			if v == "" {
				fmt.Println("stellar-core is not installed")
				return
			}
			fmt.Println(v)
		},
	}
}

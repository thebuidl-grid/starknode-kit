package commands

import (
	"buidlguidl-go/cli/cmd/options"
	"buidlguidl-go/pkg"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var SetCommand = &cobra.Command{
	Use:   "set",
	Short: "Set config values",
	Long:  "Set allows you to modify configuration values used by the application. You can set individual keys or multiple values at once.",
	Run:   setCommand,
}

func setCommand(cmd *cobra.Command, args []string) {
	cfg, err := pkg.LoadConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	if options.ConsensusClient != "" {
		client, err := pkg.GetConsensusClient(options.ConsensusClient)
		if err != nil {
			fmt.Println(err)
			return
		}
		cfg.ConsensusCientSettings.Name = pkg.ClientType(client)
	}
	if options.ExecutionClient != "" {
		client, err := pkg.GetExecutionClient(options.ExecutionClient)
		if err != nil {
			fmt.Println(err)
			return
		}
		cfg.ExecutionCientSettings.Name = pkg.ClientType(client)
	}
	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 {
			fmt.Printf("Invalid argument (must be key=value): %s\n", arg)
			continue
		}
		key := strings.ToLower(parts[0])
		value := parts[1]

		if err := applyConfigUpdate(&cfg, key, value); err != nil {
			fmt.Printf("Failed to set %s: %v\n", key, err)
		}
	}
	err = pkg.UpdateStackNodeConfig(cfg)
	if err != nil {
		fmt.Println(err)
		return

	}
}
func applyConfigUpdate(cfg *pkg.StarkNodeKitConfig, key, value string) error {
	switch key {
	case "network":
		cfg.ExecutionCientSettings.Network = value
	case "port":
		cfg.ExecutionCientSettings.Port = parsePorts(value)

	default:
		return fmt.Errorf("unknown config key: %s", key)
	}
	return nil
}

func parsePorts(value string) []string {
	parts := strings.Split(value, ",")
	var ports []string
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			ports = append(ports, trimmed)
		}
	}
	return ports
}

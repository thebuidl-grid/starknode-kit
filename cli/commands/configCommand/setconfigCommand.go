package configcommand

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/thebuidl-grid/starknode-kit/cli/options"
	t "github.com/thebuidl-grid/starknode-kit/pkg/types"
	"github.com/thebuidl-grid/starknode-kit/pkg/utils"

	"github.com/spf13/cobra"
)

var setELCmd = &cobra.Command{
	Use:   "el key=value [key=value...]",
	Short: "Set execution layer (EL) client configuration",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runSetCommand("execution", args)
	},
}
var setCLCmd = &cobra.Command{
	Use:   "cl key=value [key=value...]",
	Short: "Set consensus layer (CL) client configuration",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runSetCommand("consensus", args)
	},
}

func runSetCommand(target string, args []string) {
	if !options.LoadedConfig {
		fmt.Println(utils.Red("❌ No config found."))
		fmt.Println(utils.Yellow("💡 Run `starknode-kit config new` to create a config file."))
		return
	}

	if err := processConfigArgs(&options.Config, args, target); err != nil {
		fmt.Println(utils.Red(fmt.Sprintf("❌ Error processing config arguments: %v", err)))
		return
	}

	if err := utils.UpdateStarkNodeConfig(options.Config); err != nil {
		fmt.Println(utils.Red(fmt.Sprintf("❌ Failed to save config: %v", err)))
	}

	fmt.Println(utils.Green("✅ Configuration updated successfully!"))
}

func processConfigArgs(cfg *t.StarkNodeKitConfig, args []string, target string) error {
	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 {
			fmt.Println(utils.Yellow(fmt.Sprintf("⚠️ Invalid argument (must be key=value): %s", arg)))
			continue
		}
		key := strings.ToLower(parts[0])
		value := parts[1]

		if err := applyConfigUpdate(cfg, key, value, target); err != nil {
			return err
		}
	}
	return nil
}

func applyConfigUpdate(cfg *t.StarkNodeKitConfig, key, value, target string) error {
	var (
		updated t.ClientConfig
		err     error
	)

	switch target {
	case "execution":
		updated, err = setClientConfigValue(cfg.ExecutionCientSettings, key, value, target)
		if err == nil {
			cfg.ExecutionCientSettings = updated
		}
	case "consensus":
		updated, err = setClientConfigValue(cfg.ConsensusCientSettings, key, value, target)
		if err == nil {
			cfg.ConsensusCientSettings = updated
		}
	default:
		return fmt.Errorf("invalid config target: %s (must be 'execution' or 'consensus')", target)
	}

	if err != nil {
		return err
	}

	return nil
}

func setClientConfigValue(clientCfg t.ClientConfig, key, value, target string) (t.ClientConfig, error) {
	switch key {
	case "client":
		var client t.ClientType
		var err error

		switch target {
		case "execution":
			client, err = utils.GetExecutionClient(value)
			if err != nil {
				return clientCfg, fmt.Errorf(`%w\nSupported execution clients are:\n  - geth\n  - reth`, err)
			}
			clientCfg.Name = t.ClientType(client)
		case "consensus":
			client, err = utils.GetConsensusClient(value)
			if err != nil {
				return clientCfg, fmt.Errorf(`%w\nSupported consensus clients are:\n  - lighthouse 
  - prysm`, err)
			}
			clientCfg.Name = t.ClientType(client)
		}
	case "port":
		ports, err := parsePorts(value)
		if err != nil {
			return clientCfg, err
		}
		clientCfg.Port = ports
	case "type":
		clientCfg.ExecutionType = value

	default:
		return clientCfg, fmt.Errorf(`
"unknown config key: %s", key
Available keys you can set:
  - client           (client name)
  - port             (client ports, comma-separated)`, key)
	}
	return clientCfg, nil
}

func parsePorts(value string) ([]int, error) {
	parts := strings.Split(value, ",")
	var ports []int
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed == "" {
			continue
		}
		port, err := strconv.Atoi(trimmed)
		if err != nil {
			return nil, err // you might want to wrap this with more context
		}
		ports = append(ports, port)
	}
	return ports, nil
}


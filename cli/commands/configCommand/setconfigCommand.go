package configcommand

import (
	"fmt"
	"net/url"
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

var setStarknetCmd = &cobra.Command{
	Use:   "starknet key=value [key=value...]",
	Short: "Set starknet client configuration",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runSetCommand("starknet", args)
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
		updated any
		err     error
	)

	switch target {
	case "execution":
		_, err := utils.GetExecutionClient(value)
		if err != nil {
			return fmt.Errorf(`%w\nSupported execution clients are:\n  - geth\n  - reth`, err)
		}
		updated, err = setClientConfigValue(cfg.ExecutionCientSettings, key, value)
		if err == nil {
			cfg.ExecutionCientSettings = updated.(t.ClientConfig)
		}
	case "consensus":
		_, err := utils.GetConsensusClient(value)
		if err != nil {
			return fmt.Errorf(`%w\nSupported execution clients are:\n  - geth\n  - reth`, err)
		}
		updated, err = setClientConfigValue(cfg.ConsensusCientSettings, key, value)
		if err == nil {
			cfg.ConsensusCientSettings = updated.(t.ClientConfig)
		}
	case "starknet":
		if options.Config.JunoConfig.EthNode == "" {
			return fmt.Errorf("This is not a starknet node")
		}
		updated, err = setClientConfigValue(cfg.JunoConfig, key, value)
		if err == nil {
			cfg.JunoConfig = updated.(t.JunoConfig)
		}
	default:
		return fmt.Errorf("invalid config target: %s (must be 'execution', 'consensus' or 'starknet')", target)
	}

	if err != nil {
		return err
	}

	return nil
}

func setClientConfigValue[T t.ClientConfig | t.JunoConfig](clientCfg T, key, value string) (T, error) {
	switch c := any(clientCfg).(type) {
	case t.JunoConfig:
		if key != "eth_node" {
			return clientCfg, fmt.Errorf("invalid key '%s' for starknet config: only 'eth_node' is accepted", key)
		}
		if _, err := url.ParseRequestURI(value); err != nil {
			return clientCfg, fmt.Errorf("invalid URL format for eth_node: '%s'", value)
		}
		c.EthNode = value
		return any(c).(T), nil
	case t.ClientConfig:
		switch key {
		case "client":
			c.Name = t.ClientType(value)
		case "port":
			ports, err := parsePorts(value)
			if err != nil {
				return clientCfg, err
			}
			c.Port = ports
		case "type":
			c.ExecutionType = value
		default:
			return clientCfg, fmt.Errorf(`
"unknown config key: %s", key
Available keys you can set:
  - client           (client name)
  - port             (client ports, comma-separated)`, key)
		}
		return any(c).(T), nil
	default:
		return clientCfg, fmt.Errorf("unsupported config type")
	}
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
			return nil, fmt.Errorf("invalid port '%s': must be a number", trimmed)
		}
		ports = append(ports, port)
	}
	return ports, nil
}

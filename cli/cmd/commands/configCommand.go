package commands

import (
	"fmt"
	t "starknode-kit/pkg/types"
	"starknode-kit/pkg/utils"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var ConfigCommand = &cobra.Command{
	Use:   "config",
	Short: "show the configured Ethereum clients",
	Long: `The show command shows the Ethereum clients (e.g., Prysm, Lighthouse, Geth, etc.)
that have been added to your local configuration.`,
	Run: configCommand,
}

func configCommand(cmd *cobra.Command, args []string) {
	config, err := utils.LoadConfig()
	if err != nil {
		fmt.Println("No config found")
		fmt.Println("Run `starknode init` to create config file")
		return
	}

	var configBytes []byte
	all, _ := cmd.Flags().GetBool("all")
	el, _ := cmd.Flags().GetBool("el")
	cl, _ := cmd.Flags().GetBool("cl")

	if all {

		configBytes, err = yaml.Marshal(config)

	} else if el {

		configBytes, err = yaml.Marshal(config.ExecutionCientSettings)

	} else if cl {
		configBytes, err = yaml.Marshal(config.ExecutionCientSettings)

	}
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("=== Configuration ===")
	fmt.Println()
	fmt.Println(string(configBytes))
	fmt.Println("=== === === === === ===")
	return
}

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
	cfg, err := utils.LoadConfig()
	if err != nil {
		fmt.Println("No config found")
		fmt.Println("Run `starknode init` to create config file")
		return
	}

	if err := processConfigArgs(&cfg, args, target); err != nil {
		fmt.Println(err)
		return
	}

	if err := utils.UpdateStarkNodeConfig(cfg); err != nil {
		fmt.Println("Failed to save config:", err)
	}
}

func processConfigArgs(cfg *t.StarkNodeKitConfig, args []string, target string) error {
	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 {
			fmt.Printf("Invalid argument (must be key=value): %s\n", arg)
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
				return clientCfg, fmt.Errorf(`%w
Supported execution clients are:
  - geth
  - reth`, err)
			}
			clientCfg.Name = t.ClientType(client)
		case "consensus":
			client, err = utils.GetConsensusClient(value)
			if err != nil {
				return clientCfg, fmt.Errorf(`%w
Supported consensus clients are:
  - lighthouse 
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
  - network          (client network)
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

func init() {
	ConfigCommand.Flags().Bool("all", false, "Show all client settings")
	ConfigCommand.Flags().Bool("el", false, "Show execution client settings")
	ConfigCommand.Flags().Bool("cl", false, "Show consensus client settings")
	ConfigCommand.AddCommand(setCLCmd)
	ConfigCommand.AddCommand(setELCmd)
}

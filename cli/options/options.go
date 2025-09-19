package options

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/thebuidl-grid/starknode-kit/pkg"
	"github.com/thebuidl-grid/starknode-kit/pkg/types"
	"github.com/thebuidl-grid/starknode-kit/pkg/utils"
)

var (
	ConsensusClient string
	ExecutionClient string
	StarknetClient  string
	Installer       = pkg.NewInstaller()
	Config          types.StarkNodeKitConfig
	LoadedConfig    = false
)

func InitGlobalOptions(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&ConsensusClient, "consensus-client", "c", "", "Specify the consensus client")
	cmd.PersistentFlags().StringVarP(&ExecutionClient, "execution-client", "e", "", "Specify the execution client")
	cmd.PersistentFlags().StringVarP(&StarknetClient, "starknet-client", "s", "", "Specify the Starknet client")
}

func LoadLogs(clients []string) {

	time.Sleep(3 * time.Second)

	logs := []string{"-f"}

	for _, i := range clients {

		ilog, err := getLatestLogFile(i)
		if err != nil {
			log.Fatalf(utils.Red("❌ Could not find client log file: %v for client %s"), err, i)
			continue
		}
		logs = append(logs, ilog)
	}

	tailCmd := exec.Command("tail", logs...)
	tailCmd.Stdout = os.Stdout
	tailCmd.Stderr = os.Stderr

	if err := tailCmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if strings.Contains(exitErr.Error(), "signal: interrupt") {
				fmt.Println(utils.Green("\n✅ Stopped tailing logs. Clients are still running in the background."))
				fmt.Println(utils.Yellow("💡 Use `starknode-kit stop --all` to stop them."))
				return
			}
		}
		log.Printf(utils.Red("❌ Error tailing logs: %%v"), err)
	}

}

package initcommand

import (
	"fmt"

	"starknode-kit/pkg/styles"
	"starknode-kit/pkg/utils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

const (
	stepSelectNodeType = iota
	stepFullNodeSetup
	stepClientInstallation
)

type InitFlowModel struct {
	isQuitting               bool
	currentStep              int
	selectNodeScreen         tea.Model
	fullNodeSetupScreen      tea.Model
	clientInstallationScreen tea.Model
}

func (m InitFlowModel) Init() tea.Cmd {
	return nil
}

func (m InitFlowModel) View() string {
	if m.isQuitting {
		return ""
	}
	switch m.currentStep {
	case stepSelectNodeType:
		return m.selectNodeScreen.View()
	case stepFullNodeSetup:
		return m.fullNodeSetupScreen.View()
	case stepClientInstallation:
		return m.clientInstallationScreen.View()
	default:
		return "Unknown step"
	}
}

func (m InitFlowModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Global key handling
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "q", "esc", "ctrl+c":
			m.isQuitting = true
			return m, tea.Quit

		}
	}
	// Screen-specific update
	switch m.currentStep {

	case stepSelectNodeType:
		updated, cmd := m.selectNodeScreen.Update(msg)
		if model, ok := updated.(*Screen); ok {
			m.selectNodeScreen = model
			if model.done {
				m.currentStep = stepFullNodeSetup
				return m, m.fullNodeSetupScreen.Init()
			}
		}

		return m, cmd

	case stepFullNodeSetup:
		updated, cmd := m.fullNodeSetupScreen.Update(msg)
		if model, ok := updated.(*Screen); ok {
			m.fullNodeSetupScreen = model
			if model.done {
				m.currentStep = stepClientInstallation // TODO we need to use a slise to store screens
				return m, m.clientInstallationScreen.Init()
			}
		}
		return m, cmd

	case stepClientInstallation:
		updated, cmd := m.clientInstallationScreen.Update(msg)
		if model, ok := updated.(*Screen); ok {
			m.clientInstallationScreen = model
			if model.done {
				return m, tea.Quit
			}
		}
		return m, cmd
	}

	return m, nil
}

var InitCommand = &cobra.Command{
	Use:   "init",
	Short: "Create a default configuration file",
	Run:   runInitFlow,
}

func runInitFlow(cmd *cobra.Command, args []string) {
	utils.ClearScreen()

	// Init full node screen
	fullNodeSetupModel := NewFullNodeConfigScreen()

	// Init select screen
	selectNodeModel := NewNodeSelectionScreen()

	// Init client installation screen
	clientInstallationModel := NewInstallationScreen()

	// Create root program model
	programModel := InitFlowModel{
		isQuitting:               false,
		fullNodeSetupScreen:      fullNodeSetupModel,
		selectNodeScreen:         selectNodeModel,
		currentStep:              stepSelectNodeType,
		clientInstallationScreen: clientInstallationModel,
	}

	p := tea.NewProgram(programModel)
	if _, err := p.Run(); err != nil {
		fmt.Println(styles.Danger.Render("could not start program:"), err)
	}
}

package initcommand

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"starknode-kit/cli/cmd/commands"
	"starknode-kit/pkg"
	"starknode-kit/pkg/styles"
	"starknode-kit/pkg/types"
	"starknode-kit/pkg/utils"

	tea "github.com/charmbracelet/bubbletea"
	"gopkg.in/yaml.v3"
)

// Constants and Types
type Step int

// Select network Steps
const (
	stepSelectNetwork Step = iota
	stepSelectElClient
	stepSelectClClient
	stepAcceptConfig
)

// Install Client Steps
const (
	stepInstallClClient = iota
	stepInstallElClient
)

const (
	numMainChoices = 3
)

// Scene interface defines common behavior for all screens
type Scene interface {
	View() string
	Enter()
}

// Base screen type that delegates to active Scene
type Screen struct {
	choice    int
	done      bool
	numChoice int
	step      Step
	current   Scene
}

// Base screen methods
func (m *Screen) Init() tea.Cmd { return nil }

func (m *Screen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {

		switch key.String() {
		case "j", "down":
			m.choice = (m.choice + 1) % (m.numChoice + 1)
		case "k", "up":
			m.choice = (m.choice - 1 + m.numChoice + 1) % (m.numChoice + 1)
		case "-":
			if m.step > 0 {
				m.step--
			}
		case "enter":
			if m.current != nil {
				m.current.Enter()
			}
		}
		return m, nil
	}
	return m, nil
}

func (m *Screen) View() string {
	if m.current != nil {
		return m.current.View()
	}
	return ""
}

func (m *Screen) SetScene(scene Scene) {
	m.current = scene
}

func (m *Screen) Done() bool {
	return m.done
}

// Node Selection Screen
type NodeSelectionScreen struct {
	*Screen
}

func NewNodeSelectionScreen() *Screen {
	s := &Screen{numChoice: numMainChoices - 1, choice: 0}
	sel := &NodeSelectionScreen{Screen: s}
	s.SetScene(sel)
	return s
}

func (m *NodeSelectionScreen) View() string {
	header := pkg.Banner.String() + "\nWhat type of node do you want to run?\n\n%s\n\n"
	instructions := fmt.Sprintf(
		"Press %s to select, %s to confirm",
		styles.Primary.Render("↑/↓ or j/k"),
		styles.Primary.Render("Enter"),
	)

	choices := fmt.Sprintf(
		"%s\n%s\n%s",
		styles.Checkbox("Full node", m.choice == 0),
		styles.Checkbox("Full Starknet node", m.choice == 1),
		styles.Checkbox("Validator node", m.choice == 2),
	)

	return fmt.Sprintf(header, choices) + instructions
}

func (m *NodeSelectionScreen) Enter() {
	m.done = true
}

// Full Node Configuration Screen
type FullNodeConfigScreen struct {
	network  string
	elClient types.ClientType
	clClient types.ClientType
	*Screen
}

func NewFullNodeConfigScreen() *Screen {
	s := &Screen{
		numChoice: len(supportedNetorks) - 1,
		step:      stepSelectNetwork,
		choice:    0,
	}
	full := &FullNodeConfigScreen{Screen: s}
	s.SetScene(full)
	return s
}

func (m *FullNodeConfigScreen) View() string {
	switch m.step {
	case stepSelectNetwork:
		return m.renderSelectionScreen("Which network do you want to use?", supportedNetorks)
	case stepSelectElClient:
		return m.renderSelectionScreen("Which execution client do you want to use?", clientTypesToStrings(elClientOptions))
	case stepSelectClClient:
		return m.renderSelectionScreen("Which consensus client do you want to use?", clientTypesToStrings(clClientOptions))
	case stepAcceptConfig:
		return m.renderConfigurationScreen()
	default:
		return ""
	}
}

func (m *FullNodeConfigScreen) renderSelectionScreen(prompt string, options []string) string {
	var choices string
	for i, option := range options {
		choices += fmt.Sprintf("%s\n", styles.Checkbox(option, i == m.choice))
	}
	return fmt.Sprintf("%s\n\n%s", prompt, choices)
}

func (m *FullNodeConfigScreen) renderConfigurationScreen() string {
	config := types.StarkNodeKitConfig{
		Network: m.network,
		ExecutionCientSettings: types.ClientConfig{
			Name:          m.elClient,
			Port:          elClientPort,
			ExecutionType: "full",
		},
		ConsensusCientSettings: types.ClientConfig{
			Name:                m.clClient,
			Port:                clClientPort,
			ConsensusCheckpoint: fmt.Sprintf("https://%s-checkpoint-sync.stakely.io/", m.network),
		},
	}

	if err := utils.UpdateStarkNodeConfig(config); err != nil {
		return fmt.Sprintf("Error saving configuration: %v", err)
	}

	configBytes, _ := yaml.Marshal(config)
	return fmt.Sprintf("Configuration generated:\n\n%s\n\nPress Enter to confirm and continue...", string(configBytes))
}

func (m *FullNodeConfigScreen) Enter() {
	switch m.step {
	case stepSelectNetwork:
		m.network = supportedNetorks[m.choice]
		m.step = stepSelectElClient
		m.choice = 0
		m.numChoice = len(elClientOptions) - 1

	case stepSelectElClient:
		m.elClient = elClientOptions[m.choice]
		m.step = stepSelectClClient
		m.choice = 0
		m.numChoice = len(clClientOptions) - 1

	case stepSelectClClient:
		m.clClient = clClientOptions[m.choice]
		m.step = stepAcceptConfig
		m.numChoice = 0

	case stepAcceptConfig:
		m.done = true
	}
}

type InstallationScreen struct {
	config     types.StarkNodeKitConfig
	step       int
	output     string
	outputChan <-chan string
}

func NewInstallationScreen(config types.StarkNodeKitConfig) *InstallationScreen {
	installScreen := &InstallationScreen{config: config, step: stepInstallClClient}
	return installScreen

}

type installMsg struct {
	output     string
	done       bool
	err        error
	outputChan <-chan string
}

func (m *InstallationScreen) Init() tea.Cmd {
	switch m.step {
	case stepInstallClClient:
		return m.installClientCmd(m.config.ConsensusCientSettings.Name)
	case stepInstallElClient:
		return m.installClientCmd(m.config.ExecutionCientSettings.Name)
	}
	return nil
}

func (m *InstallationScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case installMsg:
		if msg.outputChan != nil {
			m.outputChan = msg.outputChan
			return m, m.readFromChannelCmd()
		}

		if msg.output != "" {
			m.output += msg.output
		}

		if msg.done {
			m.step++
			m.outputChan = nil
			return m, m.Init() // Start next installation
		}
	}
	return m, nil
}

func (m *InstallationScreen) readFromChannelCmd() tea.Cmd {
	return func() tea.Msg {
		if m.outputChan == nil {
			return installMsg{done: true}
		}

		select {
		case line, ok := <-m.outputChan:
			if !ok {
				return installMsg{done: true}
			}
			return installMsg{output: line}
		default:
			// No output available, try again shortly
			return tea.Tick(time.Millisecond*50, func(t time.Time) tea.Msg {
				return m.readFromChannelCmd()()
			})
		}
	}
}

func (m InstallationScreen) View() string {
	return m.output
}

func (m *InstallationScreen) installClientCmd(clientName types.ClientType) tea.Cmd {
	return func() tea.Msg {
		// Create pipe for capturing output
		r, w, err := os.Pipe()
		if err != nil {
			return installMsg{err: err}
		}

		// Save and redirect stdout
		originalStdout := os.Stdout
		os.Stdout = w

		// Channel to collect output
		outputChan := make(chan string, 100)

		// Start reader goroutine
		go func() {
			defer close(outputChan)
			scanner := bufio.NewScanner(r)
			for scanner.Scan() {
				outputChan <- scanner.Text() + "\n"
			}
		}()

		// Start installation goroutine
		go func() {
			defer w.Close()
			defer func() { os.Stdout = originalStdout }()
			_ = commands.Installer.InstallClient(clientName)
		}()

		// Return the channel immediately
		return installMsg{
			outputChan: outputChan,
			done:       false,
		}
	}
}

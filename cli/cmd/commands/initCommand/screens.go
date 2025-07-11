package initcommand

import (
	"fmt"
	"starknode-kit/pkg"
	"starknode-kit/pkg/styles"
	"starknode-kit/pkg/types"

	tea "github.com/charmbracelet/bubbletea"
)

// Step represents the flow step
type Step int

const (
	stepSelectElClient Step = iota
	stepSelectClClient
)

// Scene interface to allow sub-screens to define their logic
type Scene interface {
	View() string
	Enter()
}

// screen is the base type that delegates to an active Scene
type screen struct {
	choice    int
	done      bool
	numChoice int
	step      Step

	current Scene
}

func (m *screen) Init() tea.Cmd { return nil }

func (m *screen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "j", "down":
			if m.choice < m.numChoice {
				m.choice++
			} else {
				m.choice = 0
			}
		case "k", "up":
			if m.choice > 0 {
				m.choice--
			} else {
				m.choice = m.numChoice
			}
		case "-":
			if m.step > 0 {
				m.step--
			}
		case "enter":
			if m.current != nil {
				m.current.Enter()
			}
			return m, nil
		}
	}
	return m, nil
}

func (m *screen) View() string {
	if m.current != nil {
		return m.current.View()
	}
	return ""
}

func (m *screen) SetScene(scene Scene) {
	m.current = scene
}

func (m *screen) Done() bool {
	return m.done
}

type selectModel struct {
	*screen
}

func (m *selectModel) View() string {
	c := m.screen.choice
	header := pkg.Banner.String()
	header += "\nWhat type of node do you want to run?\n\n"
	header += "%s\n\n"
	header += "Press %s to select, %s to confirm\n"

	choices := fmt.Sprintf(
		"%s\n%s\n%s",
		styles.Checkbox("Full node", c == 0),
		styles.Checkbox("Full Starknet node", c == 1),
		styles.Checkbox("Validator node", c == 2),
	)

	return fmt.Sprintf(
		header,
		choices,
		styles.Primary.Render("↑/↓ or j/k"),
		styles.Primary.Render("Enter"),
	)
}

func (m *selectModel) Enter() {
	// handle selection logic here
	m.screen.done = true
}

type fullNodeModel struct {
	elClient types.ClientType
	clClient types.ClientType
	*screen
}

func (m *fullNodeModel) View() string {
	var prompt string
	var options string
	choice := m.choice

	switch m.step {
	case stepSelectElClient:
		prompt = "Which execution client do you want to use?"
		for i, c := range elClientOptions {
			options += fmt.Sprintf("%s\n", styles.Checkbox(c.String(), i == choice))
		}
	case stepSelectClClient:
		prompt = "Which consensus client do you want to use?"
		for i, c := range clClientOptions {
			options += fmt.Sprintf("%s\n", styles.Checkbox(c.String(), i == choice))
		}
	default:
		return "Configuration complete."
	}

	return fmt.Sprintf("%s\n\n%s", prompt, options)
}

func (m *fullNodeModel) Enter() {
	switch m.step {
	case stepSelectElClient:
		index := m.choice
		m.elClient = elClientOptions[index]
		m.step = stepSelectClClient
		m.choice = 0
		m.done = false

	case stepSelectClClient:
		index := m.choice
		m.clClient = clClientOptions[index]
		m.choice = 0
		m.done = true
	}
}

func NewSelectScreen() *screen {
	s := &screen{numChoice: 3, choice: 0}
	sel := &selectModel{screen: s}
	s.SetScene(sel)
	return s
}

func NewFullNodeScreen() *screen {
	s := &screen{numChoice: len(elClientOptions), step: stepSelectElClient, choice: 0}
	full := &fullNodeModel{screen: s}
	s.SetScene(full)
	return s
}

// NOTE going back does not work 

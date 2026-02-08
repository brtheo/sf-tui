package newMetadata

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MetadataType string
const (
	LWC MetadataType = "LWC"
	ApexClass MetadataType = "ApexClass"
	ApexTrigger MetadataType = "ApexTrigger"
)

type viewState string
const (
	IDLE viewState = "IDLE"
	ERR  viewState = "ERR"
)
type triggerState string
const (
	TRIGGER_NAME triggerState = "IDLE"
	TRIGGER_SOBJECT triggerState = "NEED_SOBJECT"
)

type Model struct {
	input textinput.Model
	metadataType MetadataType
	output, name, sobject string
	state viewState
	triggerState triggerState
}

func (m Model) sfCmdAndLeave(args []string) (Model, tea.Cmd) {
	return m, tea.Sequence(
		func() tea.Cmd {
			return func() tea.Msg {
				exec.Command("sf", args...).Run()
				return nil
			}
		}(), tea.Quit,
	)
}

func generateLWC(name, path string) []string {
	raw := fmt.Sprintf("lightning generate component --name %s --type lwc --output-dir %s/lwc", name, path)
	return strings.Split(raw, " ")
}
func generateApexClass(name, path string) []string {
	raw := fmt.Sprintf("apex generate class --name %s --output-dir %s/classes", name, path)
	return strings.Split(raw, " ")
}
func generateApexTrigger(name,sobject, path string) []string {
	raw := fmt.Sprintf("apex generate trigger --name %s --sobject %s --output-dir %s/triggers", name,sobject, path)
	return strings.Split(raw, " ")
}

func New(metadataType MetadataType, output string) Model {
	input := textinput.New()
	input.Width = 50
	input.Placeholder = fmt.Sprintf("Enter %s name", metadataType)
	input.Focus()

	return Model{
		input: input,
		metadataType: metadataType,
		output: output,
		state: IDLE,
		triggerState: TRIGGER_NAME,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
				case tea.KeyEnter:
				  value := strings.TrimSpace(m.input.Value())
			 		if value == "" {
						m.state = ERR
						return m, nil
					}

					if(m.metadataType == ApexTrigger && m.triggerState == TRIGGER_NAME){
						m.sobject = value
					}	else {
						m.name = value
					}

					switch m.metadataType {
						case LWC:
							return m.sfCmdAndLeave(generateLWC(m.name, m.output))
						case ApexClass:
							return m.sfCmdAndLeave(generateApexClass(m.name, m.output))
						case ApexTrigger:
							switch m.triggerState {
								case TRIGGER_NAME:
									m.state = IDLE
									m.triggerState = TRIGGER_SOBJECT
									m.input.Placeholder = "Enter SObject Name"
									m.input.Reset()
									return m, nil
								case TRIGGER_SOBJECT:
									return m.sfCmdAndLeave(generateApexTrigger(m.name, m.sobject, m.output))
							}
					}
			}
	}
	var inputCmd tea.Cmd
	m.input, inputCmd = m.input.Update(msg)
	return m, inputCmd
}

func (m Model) View() string {
	input := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#1AB9FF")).
		Render(m.input.View())
	errMsg := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#D83A00")).
		PaddingLeft(1).
		Render(" Name cannot be empty")

	switch m.state {
	case IDLE:
		return input
	case ERR:
		return lipgloss.JoinVertical(lipgloss.Left, input, errMsg)
	}
	return m.input.View()
}

package mdRetriever

import (
	"fmt"
	"math"

	"github.com/brtheo/sf-tui/models/mdRetriever/mdTypePicker"
	"github.com/brtheo/sf-tui/models/mdRetriever/mdTable"
	tea "github.com/charmbracelet/bubbletea"
)

type WizardStep int
const (
	PickMetadataType WizardStep = iota
	PickMetadataRecord
)

type Model struct {
	mdTable         mdTable.Model
	mdTypePicker     mdTypePicker.Model

	wizardStep WizardStep
	frameSize []int
	selectedMetadataType string
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			m.mdTypePicker.List.SetSize(msg.Width - 1, msg.Height - 2)
			m.mdTable.Table.SetHeight(msg.Height - 7)
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyLeft:
				m.wizardStep = WizardStep(math.Abs(float64((m.wizardStep - 1) % 2)))
			case tea.KeyRight:
				m.wizardStep = (m.wizardStep + 1) % 2
			}
		case mdTypePicker.HasPickedTypeMsg:
			m.selectedMetadataType = string(msg)
			m.wizardStep = PickMetadataRecord
			var mdTableCmd tea.Cmd
			mdTableModel, mdTableCmd := m.mdTable.Update(mdTable.HasSelectedMdTypeMsg(string(msg)))
			m.mdTable = mdTableModel.(mdTable.Model)
			return m, mdTableCmd
		case mdTypePicker.HasPickedTypeRequiringFolderMsg:
			m.selectedMetadataType = string(msg)
			m.wizardStep = PickMetadataRecord
			var mdTableCmd tea.Cmd
			mdTableModel, mdTableCmd := m.mdTable.Update(mdTable.HasSelectedMdTypeRequiringFolderMsg(string(msg)))
			m.mdTable = mdTableModel.(mdTable.Model)
			return m, mdTableCmd
	}

	switch m.wizardStep {
		case PickMetadataType:
			var mdTypePickerCmd tea.Cmd
			mdTypePickerModel, mdTypePickerCmd := m.mdTypePicker.Update(msg)
			m.mdTypePicker = mdTypePickerModel.(mdTypePicker.Model)
			return m, mdTypePickerCmd

		case PickMetadataRecord:
			var mdTableCmd tea.Cmd
			mdTableModel, mdTableCmd := m.mdTable.Update(msg)
			m.mdTable = mdTableModel.(mdTable.Model)
			return m, mdTableCmd
	}

	return m, nil
}

func (m Model) View() string {
	switch m.wizardStep {
		case PickMetadataType:
			return m.mdTypePicker.View()

		case PickMetadataRecord:
			if m.selectedMetadataType == "" {
				return "Select metadata type first"
			}
			return m.mdTable.View()
	}
	return fmt.Sprintf("%d", m.wizardStep)
}

func New() Model {
	mdTypePicker := mdTypePicker.New()
	mdTable := mdTable.New()

	return Model{
		mdTypePicker:   mdTypePicker,
		mdTable:        mdTable,
		wizardStep:   PickMetadataType,
		selectedMetadataType: "",
	}
}

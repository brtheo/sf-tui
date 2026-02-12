package metadataRetriever

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var METADATA_LIST = []string{"org", "list","metadata","--json","--metadata-type"}

type FetchState int
const (
	Fetch_Idle FetchState = iota
	Fetch_Fetching
	Fetch_Error
	Fetch_Done
)

type WizardStep int
const (
	PickMetadataType WizardStep = iota
	PickMetadataRecord
)

type ColumnID int
const (
	Col_Checkbox ColumnID = iota
	Col_FullName
	Col_CreatedBy
	Col_CreatedAt
	Col_UpdatedBy
	Col_UpdatedAt
)

func (c ColumnID) String() string {
	return [...]string{"Metadata name", "Created by", "Created at", "Updated by", "Updated at"}[c]
}

var (
	baseStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))

	highlightStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		Bold(true)
)

type Model struct {
	textInput     textinput.Model
	table         table.Model
	list list.Model
	originalRows  []table.Row
	filterColumn  ColumnID
	selectedRows  map[int]bool
	fetchState FetchState
	wizardStep WizardStep
	frameSize []int
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			m.table.SetHeight(msg.Height - 7)
			h, v := m.frameSize[0], m.frameSize[1]
			m.list.SetSize(msg.Width-h, msg.Height-v)
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyCtrlC, tea.KeyEsc:
				return m, tea.Quit
			case tea.KeyTab:
				m.filterColumn = (m.filterColumn + 2) % 5
			case tea.KeyLeft:
				m.wizardStep = (m.wizardStep + 1) % 2
			case tea.KeyRight:
				m.wizardStep = (m.wizardStep - 1) % 2
			case tea.KeyEnter:
				val, ok := m.selectedRows[m.table.Cursor()]
				if !ok {
					m.selectedRows[m.table.Cursor()] = true
				}
				m.selectedRows[m.table.Cursor()] = !val
			}
	}

	m.textInput, cmd = m.textInput.Update(msg)

	searchTerm := strings.ToLower(m.textInput.Value())
	var filteredRows []table.Row

	for i, originalRow := range m.originalRows {
		targetValue := strings.ToLower(originalRow[int(m.filterColumn)])
		if strings.Contains(targetValue, searchTerm) {
			checkbox := "🞅"
			if m.selectedRows[i] {
				checkbox = "🞊"
			}
			newRow := table.Row{checkbox}
			newRow = append(newRow, originalRow[1:]...)
			filteredRows = append(filteredRows, newRow)
		}
	}

	m.table.SetRows(filteredRows)

	var tCmd tea.Cmd
	m.table, tCmd = m.table.Update(msg)

	return m, tea.Batch(cmd, tCmd)
}

func (m Model) View() string {
	switch m.wizardStep {
		case PickMetadataType:
			return fmt.Sprintf(
				"Choose metadata type\n%s\n%s",
				m.textInput.View(),
				m.list.View(),
			)
		case PickMetadataRecord:
			return fmt.Sprintf(
				" Filtering by: %s (Press [Tab] to switch)\n Input: %s\n\n%s\n",
				highlightStyle.Render(m.filterColumn.String()),
				m.textInput.View(),
				baseStyle.Render(m.table.View()),
			)
	}
	return fmt.Sprintf(
		" Filtering by: %s (Press [Tab] to switch)\n Input: %s\n\n%s\n",
		highlightStyle.Render(m.filterColumn.String()),
		m.textInput.View(),
		baseStyle.Render(m.table.View()),
	)
}

func New(frameSize... int) Model {
	ti := textinput.New()
	ti.Placeholder = "Search..."
	ti.Focus()


	columns := []table.Column{
		{Title: "", Width: 3},
		{Title: "Metadata name", Width: 40},
		{Title: "Created by", Width: 15},
		{Title: "Created at", Width: 30},
		{Title: "Updated by", Width: 15},
		{Title: "Updated at", Width: 30},
	}
	METADATA_LIST = append(METADATA_LIST, "ApexClass")
	raw, err := exec.Command("sf",METADATA_LIST...).Output()
	if err != nil {
		fmt.Println(err)
	}
	metadata, err := UnmarshalMetadata(raw)
	if err != nil {
		fmt.Println(err)
	}

	var rows = []table.Row{}
	for _, field := range metadata.Result {
		rows = append(rows,
			table.Row {
				"🞅",
				field.FullName,
				field.CreatedByName,
				field.CreatedDate.String(),
				field.LastModifiedByName,
				field.LastModifiedDate.String(),
			},
		)
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	list := list.New(MetadataTypes, list.NewDefaultDelegate(), 0, 0)

	return Model{
		textInput:    ti,
		list: list,
		table:        t,
		originalRows: rows,
		filterColumn: Col_FullName, // Default filter
		selectedRows: make(map[int]bool),
		fetchState:   Fetch_Idle,
		wizardStep:   PickMetadataType,
		frameSize: frameSize,
	}
}

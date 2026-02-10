package metadataRetriever

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var METADATA_LIST = []string{"org", "list","metadata","--json","--metadata-type"}

type ColumnID int

const (
	Col_FullName ColumnID = iota
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
	originalRows  []table.Row
	filterColumn  ColumnID
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			m.table.SetHeight(msg.Height - 7)
		case tea.KeyMsg:
			switch msg.String() {
			case "esc", "ctrl+c":
				return m, tea.Quit
			case "tab":
				m.filterColumn = (m.filterColumn + 1) % 5
			}
	}

	m.textInput, cmd = m.textInput.Update(msg)

	searchTerm := strings.ToLower(m.textInput.Value())
	var filteredRows []table.Row

	for _, row := range m.originalRows {
		targetValue := strings.ToLower(row[int(m.filterColumn)])
		if strings.Contains(targetValue, searchTerm) {
			filteredRows = append(filteredRows, row)
		}
	}

	m.table.SetRows(filteredRows)

	var tCmd tea.Cmd
	m.table, tCmd = m.table.Update(msg)

	return m, tea.Batch(cmd, tCmd)
}

func (m Model) View() string {
	return fmt.Sprintf(
		" Filtering by: %s (Press [Tab] to switch)\n Input: %s\n\n%s\n",
		highlightStyle.Render(m.filterColumn.String()),
		m.textInput.View(),
		baseStyle.Render(m.table.View()),
	)
}

func New() Model {
	ti := textinput.New()
	ti.Placeholder = "Search..."
	ti.Focus()

	columns := []table.Column{
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

	return Model{
		textInput:    ti,
		table:        t,
		originalRows: rows,
		filterColumn: Col_FullName, // Default filter
	}
}

package mdTable

import (
	"fmt"
	"strings"

	"github.com/brtheo/sf-tui/models/mdRetriever/shared"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type HasFetchedRowsMsg []table.Row
type HasSelectedMdTypeMsg string
type HasSelectedMdTypeRequiringFolderMsg string

var columns = []table.Column{
	{Title: "", Width: 1},
	{Title: "Metadata name", Width: 40},
	{Title: "Created by", Width: 15},
	{Title: "Created at", Width: 20},
	{Title: "Updated by", Width: 15},
	{Title: "Updated at", Width: 20},
}

var (
	baseStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))

	highlightStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		Bold(true)
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
	return [...]string{
		"Metadata name",
		"Created by",
		"Created at",
		"Updated by",
		"Updated at"}[c]
}

type Model struct {
	filterInput textinput.Model
	Table         table.Model
	spinner       spinner.Model
	help          help.Model
	originalRows  []table.Row
	filterColumn  ColumnID
	selectedRows  map[string]map[string]bool
	SelectedMdType string
	isFetching bool
}

func New() Model {
	filterInput := textinput.New()
	filterInput.Placeholder = "Search ..."
	filterInput.Focus()

	loader := spinner.New()
    loader.Spinner = spinner.MiniDot
    loader.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	mdTable := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(3),
	)

	return Model{
		filterInput   : filterInput,
		Table         : mdTable,
		spinner       : loader,
		help          : help.New(),
		originalRows  : []table.Row{},
		filterColumn  : Col_FullName,
		selectedRows  : make(map[string]map[string]bool),
		SelectedMdType: "",
		isFetching: true,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var tableCmd, filterInputCmd, spinnerCmd tea.Cmd

	switch msg := msg.(type) {
		case HasSelectedMdTypeMsg:
			return m.handleSelectedMdType(msg)

		case HasSelectedMdTypeRequiringFolderMsg:
			return m.handleSelectedMdTypeRequiringFolder(msg)

		case HasFetchedRowsMsg:
			m.Table.SetRows(msg)
			m.originalRows = msg
			m.isFetching = false
		case tea.KeyMsg:
			switch msg.Type {
				case tea.KeyCtrlG:
					m.handleWritePackageXML()
				case tea.KeyTab:
					m.filterColumn = (m.filterColumn + 2) % 5
				case tea.KeyEnter:
					m.handleCheckboxClick()
			}
	}
	searchTerm := strings.ToLower(m.filterInput.Value())

	m.Table.SetRows(m.getRowsWithCheckboxes(searchTerm))

	m.Table, tableCmd = m.Table.Update(msg)
	m.filterInput, filterInputCmd = m.filterInput.Update(msg)
	m.spinner, spinnerCmd = m.spinner.Update(msg)

	return m, tea.Batch(
		tableCmd,
		filterInputCmd,
		spinnerCmd,
	)
}

func (m Model) View() string {
	if m.isFetching {
		return fmt.Sprintf(
			"Fetching %s %s",
			m.SelectedMdType,
			m.spinner.View(),
		)
	}
	return fmt.Sprintf(
		"Filtering %s by: %s\n Input: %s\n\n%s\n%s",
		m.SelectedMdType,
		highlightStyle.Render(m.filterColumn.String()),
		m.filterInput.View(),
		baseStyle.Render(m.Table.View()),
		m.help.View(shared.CtrlKeys),
	)
}

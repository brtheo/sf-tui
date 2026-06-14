package mdTable

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)
func (m Model) handleSelectedMdTypeRequiringFolder(msg HasSelectedMdTypeRequiringFolderMsg) (tea.Model, tea.Cmd) {
	m.SelectedMdType = string(msg)
	m.isFetching = true
	return m, tea.Batch(
		m.fetchMdListWithFolder,
		m.spinner.Tick,
	)
}

func (m Model) handleSelectedMdType(msg HasSelectedMdTypeMsg) (tea.Model, tea.Cmd) {
	m.SelectedMdType = string(msg)
	m.isFetching = true
	return m, tea.Batch(
		m.fetchMdList,
		m.spinner.Tick,
	)
}

func (m Model) handleCheckboxClick() {
	name := m.Table.SelectedRow()[1]
	m.setKeyIfNil(m.SelectedMdType)
	val, ok := m.selectedRows[m.SelectedMdType][name]
	if !ok {
		m.selectedRows[m.SelectedMdType][name] = true
	}
	m.selectedRows[m.SelectedMdType][name] = !val
}


func (m Model) handleWritePackageXML() {
	res, err := m.generatePackageXML()
	if err != nil {
		fmt.Println("Error generating package XML:", err)
	}

	if _, err := os.Stat("./manifest"); os.IsNotExist(err) {
		if err := os.Mkdir("./manifest", os.ModePerm); err != nil {
			panic(err)
		}
	}
	os.WriteFile("./manifest/package.xml", []byte(res), os.ModePerm)
}

package mdTable

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"text/template"
	"encoding/xml"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) getRowsWithCheckboxes (searchTerm string) (filteredRows []table.Row) {
	m.parsePackageXML()

	for _, originalRow := range m.originalRows {
		targetValue := strings.ToLower(originalRow[int(m.filterColumn)])

		if strings.Contains(targetValue, searchTerm) {
			checkbox := "🞅"
			if m.selectedRows[m.SelectedMdType][originalRow[1]] { checkbox = "🞊" }

			filteredRow := append(
				table.Row{checkbox},
				originalRow[1:]...
			)

			filteredRows = append(filteredRows, filteredRow)
		}
	}
	return filteredRows
}


func (m Model) fetchMdList() tea.Msg {
	return HasFetchedRowsMsg(m.getMetadataRows(""))
}

func (m Model) fetchFolders() []string {
	raw, err := exec.Command(
		"sf","org","list","metadata","--json","--metadata-type", "EmailFolder",
	).Output()
	if err != nil {
		fmt.Println(err)
	}
	folders, err := UnmarshalMetadata(raw)
	if err != nil {
		fmt.Println(err)
	}

	var folderNames = []string{}
	for _, folder := range folders.Result {
		folderNames = append(folderNames, folder.FullName)
	}
	return folderNames
}

func (m Model) getMetadataRows(folderName string) (rows []table.Row) {
	args := []string{"org", "list", "metadata", "--json", "--metadata-type", m.SelectedMdType}
	// m.setKeyIfNil(m.SelectedMdType)
	if folderName != "" {
		args = append(args, "--folder", folderName)
	}
	raw, err := exec.Command("sf", args...).Output()
	if err != nil {
		fmt.Println(err)
	}
	metadata, err := UnmarshalMetadata(raw)
	if err != nil {
		fmt.Println(err)
	}
	sort.Slice(metadata.Result, func(i, j int) bool {
		return metadata.Result[i].LastModifiedDate.After(metadata.Result[j].LastModifiedDate)
	})
	for _, field := range metadata.Result {
		rows = append(rows,
			table.Row {
				"",
				field.FullName,
				field.CreatedByName,
				field.CreatedDate.String(),
				field.LastModifiedByName,
				field.LastModifiedDate.String(),
			},
		)
	}
	return rows
}

func (m Model) fetchMdListWithFolder() tea.Msg {
	var rows = []table.Row{}
	folderNames := m.fetchFolders()
	for _, folderName := range folderNames {
		rows = append(rows, m.getMetadataRows(folderName)...)
	}
	return HasFetchedRowsMsg(rows)
}

func (m Model) setKeyIfNil(key string) {
	var val, ok = m.selectedRows[key]
	if !ok {
		val = make(map[string]bool)
		m.selectedRows[key] = val
	}
}

func (m Model) generatePackageXML() (string, error) {
	const packageTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<Package xmlns="http://soap.sforce.com/2006/04/metadata">{{range $mdType, $mdNames := .values}}
    <types>{{range $index, $mdName := $mdNames}}
      	<members>{{$mdName}}</members>{{end}}
        <name>{{$mdType}}</name>
    </types>{{end}}
    <version>{{.version}}</version>
</Package>`

	tmpl, err := template.New("package").Parse(packageTemplate)
	if err != nil {
		return "", err
	}

	var result strings.Builder
	values := map[string][]string{}
	for mdType, row := range m.selectedRows {

		for mdName, selected := range row {
			if selected {
				val, ok := values[mdType]
				if !ok {
					val = []string{}
				}
				values[mdType] = append(val, mdName)
			}
		}
	}
	apiVersion := "66.0"

	data, err := os.ReadFile("./sfdx-project.json")
	if err == nil {
		sfdxProject := map[string]any{}
		if err := json.Unmarshal(data, &sfdxProject); err == nil {
			if version, ok := sfdxProject["sourceApiVersion"]; ok {
				apiVersion = version.(string)
			}
		}
	}

	args := map[string]any{
		"version": apiVersion,
		"values":  values,
	}

	err = tmpl.Execute(&result, args)
	if err != nil {
		return "", err
	}

	return result.String(), nil
}

func (m Model) parsePackageXML() error  {
	data, err := os.ReadFile("./manifest/package.xml")
	if err != nil {
		return err
	}

	var pkg Package
	err = xml.Unmarshal(data, &pkg)
	if err != nil {
		return err
	}

	for _, types := range pkg.Types {
		mdType := types.Name
		m.setKeyIfNil(mdType)
		for _, member := range types.Members {
			m.selectedRows[mdType][member] = true
		}
	}
	return nil
}

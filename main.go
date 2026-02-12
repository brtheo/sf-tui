package main

import (
	"github.com/alexflint/go-arg"
	genMetadata "github.com/brtheo/sf-tui/models/genMetadata"
	metadataRetriever "github.com/brtheo/sf-tui/models/metadataRetriever"
	orgPicker "github.com/brtheo/sf-tui/models/orgPicker"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ModelType string
const (
	OrgPicker ModelType = "org-picker"
	GenMetadata ModelType = "gen"
	MetadataRetriever ModelType = "metadata-retriever"
)

var args struct {
	Model ModelType `arg:"positional" default:"org-picker" help:"Model type (newMetadata, orgPicker)"`
	MetadataType genMetadata.MetadataType `arg:"-t, --type" help:"Metadata type (LWC, ApexClass, ApexTrigger)"`
	Output string `arg:"-o, --output" help:"Path to force-app"`
}

var globalStyle = lipgloss.NewStyle().Margin(0, 1)

type Model struct {
	subModel tea.Model
}

func New(modelType ModelType, p *arg.Parser) Model {
	switch modelType {
		case GenMetadata:
			if args.Output == "" || args.MetadataType == "" {
				p.Fail("Missing required arguments.\nMust provide output path and metadata type.\nSee --help for more information.")
			}
			return Model{subModel: genMetadata.New(args.MetadataType, args.Output)}
		case OrgPicker:
			return Model{subModel: orgPicker.New(globalStyle.GetFrameSize())}
		case MetadataRetriever:
			return Model{subModel: metadataRetriever.New(globalStyle.GetFrameSize())}
	}
	return Model{}
}

func (m Model) Init() tea.Cmd {
	return m.subModel.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
			case tea.KeyCtrlC, tea.KeyEsc:
				return m, tea.Quit
		}
	}
	var subModelCmd tea.Cmd
	m.subModel, subModelCmd = m.subModel.Update(msg)
	return m, subModelCmd
}

func (m Model) View() string {
	return globalStyle.Render(m.subModel.View())
}

func main() {
	p := arg.MustParse(&args)
	// var model tea.Model
	// switch args.Model {
	// 	case GenMetadata:
	// 	if args.Output == "" || args.MetadataType == "" {
	// 		p.Fail("Missing required arguments.\nMust provide output path and metadata type.\nSee --help for more information.")
	// 	}
	// 	model = genMetadata.New(args.MetadataType, args.Output)
	// 	case OrgPicker:
	// 		model = orgPicker.New()
	// }
	model := New(args.Model, p)
	tea.NewProgram(model, tea.WithAltScreen()).Run()
}

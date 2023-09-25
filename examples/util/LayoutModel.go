package util

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/winder/layout"
)

type LayoutModel interface {
}

type layoutModel struct {
	models []tea.Model
	layout layout.BubbleLayout
	view   func([]tea.Model) string
}

func NewLayoutModel(models []tea.Model, layout layout.BubbleLayout, view func([]tea.Model) string) tea.Model {
	return layoutModel{
		models: models,
		layout: layout,
		view:   view,
	}
}

func (m layoutModel) Init() tea.Cmd {
	return func() tea.Msg {
		return m.layout.Resize(100, 15)
	}
}

func (m layoutModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		// Convert WindowSizeMsg to LayoutMsg.
		return m, func() tea.Msg {
			return m.layout.Resize(msg.Width, msg.Height)
		}
	}

	// Dispatch to all models.
	var commands []tea.Cmd
	for idx, model := range m.models {
		var cmd tea.Cmd
		m.models[idx], cmd = model.Update(msg)
		commands = append(commands, cmd)
	}

	return m, tea.Batch(commands...)
}

func (m layoutModel) View() string {
	return m.view(m.models)
}

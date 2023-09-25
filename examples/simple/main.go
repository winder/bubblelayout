package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	bl "github.com/winder/layout"
)

type layoutModel struct {
	layout bl.BubbleLayout

	leftID  bl.ID
	rightID bl.ID

	leftSize  bl.Size
	rightSize bl.Size
}

func New() tea.Model {
	layoutModel := layoutModel{
		layout: bl.New(),
	}
	layoutModel.leftID = layoutModel.layout.Add(bl.Layout{MaxWidth: 10})
	layoutModel.rightID = layoutModel.layout.Add(bl.Layout{})
	return layoutModel
}

func (m layoutModel) Init() tea.Cmd {
	return func() tea.Msg {
		return m.layout.Resize(80, 40)
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
		// Convert WindowSizeMsg to BubbleLayoutMsg.
		return m, func() tea.Msg {
			return m.layout.Resize(msg.Width, msg.Height)
		}
	case bl.BubbleLayoutMsg:
		m.leftSize, _ = msg.Size(m.leftID)
		m.rightSize, _ = msg.Size(m.rightID)
	}

	return m, nil
}

func boxStyle(size bl.Size, bg lipgloss.Color) lipgloss.Style {
	return lipgloss.NewStyle().
		Background(bg).
		Width(size.Width).
		Height(size.Height).
		Align(lipgloss.Center)
}

func (m layoutModel) View() string {
	return lipgloss.JoinHorizontal(0,
		boxStyle(m.leftSize, "9").Render("left"),
		boxStyle(m.rightSize, "13").Render("right"))
}

func main() {
	p := tea.NewProgram(New())
	p.Run()
}

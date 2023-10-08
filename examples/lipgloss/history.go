package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	bl "github.com/winder/bubblelayout"
)

var historyStyle = lipgloss.NewStyle().
	Align(lipgloss.Left).
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(highlight).
	Padding(1, 2, 1, 2)

type historyModel struct {
	ID         bl.ID
	size       bl.Size
	position   lipgloss.Position
	history    string
	leftMargin int
}

func MakeHistory(id bl.ID, history string, alignment lipgloss.Position, extraMarginLeft int) historyModel {
	lipgloss.NewStyle().Align()
	return historyModel{
		ID:         id,
		history:    history,
		position:   alignment,
		leftMargin: extraMarginLeft,
	}
}

func (m historyModel) Init() tea.Cmd {
	return nil
}

func (m historyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case bl.BubbleLayoutMsg:
		var err error
		m.size, err = msg.Size(m.ID)
		if err != nil {
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m historyModel) View() string {
	history := historyStyle.
		Align(m.position).
		Width(m.size.Width - 3 - m.leftMargin).
		Height(m.size.Height - 1).
		MaxHeight(m.size.Height - 1).
		Render(m.history)
	return lipgloss.NewStyle().
		MarginLeft(1 + m.leftMargin).
		MarginRight(2).
		MaxHeight(m.size.Height).
		Height(m.size.Height).
		MaxWidth(m.size.Width).
		Render(history)
}

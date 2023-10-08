package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	bl "github.com/winder/bubblelayout"
)

var (
	divider = lipgloss.NewStyle().
		SetString("•").
		Padding(0, 1).
		Foreground(subtle).
		String()

	url = lipgloss.NewStyle().Foreground(special).Render

	descStyle = lipgloss.NewStyle().MarginTop(1)

	infoStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderTop(true).
			BorderForeground(subtle)
)

type descModel struct {
	ID   bl.ID
	size bl.Size
}

func (m descModel) Init() tea.Cmd {
	return nil
}

func (m descModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m descModel) View() string {
	desc := lipgloss.JoinVertical(lipgloss.Left,
		descStyle.Render("Style Definitions for Nice Terminal Layouts"),
		infoStyle.Render("From Charm"+divider+url("https://github.com/charmbracelet/lipgloss")),
	)

	return lipgloss.NewStyle().MaxWidth(m.size.Width).Render(desc)
}

package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	bl "github.com/winder/bubblelayout"
)

var dialogBoxStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("#874BFD")).
	Padding(1, 0).
	BorderTop(true).
	BorderLeft(true).
	BorderRight(true).
	BorderBottom(true)

var buttonStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FFF7DB")).
	Background(lipgloss.Color("#888B7E")).
	Padding(0, 3).
	MarginTop(1).
	MarginRight(2)

var activeButtonStyle = buttonStyle.Copy().
	Foreground(lipgloss.Color("#FFF7DB")).
	Background(lipgloss.Color("#F25D94")).
	Underline(true)

type dialogModel struct {
	ID     bl.ID
	size   bl.Size
	active int
}

func (m dialogModel) Init() tea.Cmd {
	return nil
}

func (m dialogModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "right", "l", "n", "tab":
			fallthrough
		case "left", "h", "p", "shift+tab":
			m.active = (m.active + 1) % 2
			return m, nil
		}
	case bl.BubbleLayoutMsg:
		var err error
		m.size, err = msg.Size(m.ID)
		if err != nil {
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m dialogModel) View() string {
	var okButton, cancelButton string
	if m.active == 0 {
		okButton = activeButtonStyle.Render("Yes")
		cancelButton = buttonStyle.Render("Maybe")
	} else {
		okButton = buttonStyle.Render("Yes")
		cancelButton = activeButtonStyle.Render("Maybe")
	}

	question := lipgloss.NewStyle().Width(50).Align(lipgloss.Center).Render("Are you sure you want to eat marmalade?")
	buttons := lipgloss.JoinHorizontal(lipgloss.Top, okButton, cancelButton)
	ui := lipgloss.JoinVertical(lipgloss.Center, question, buttons)

	//These whitespace characters are weird, had to fuss around with  the pad amount.
	dialog := lipgloss.Place(m.size.Width-6, m.size.Height,
		lipgloss.Center, lipgloss.Center,
		dialogBoxStyle.Render(ui),
		lipgloss.WithWhitespaceChars("猫咪"),
		lipgloss.WithWhitespaceForeground(subtle),
	)

	return lipgloss.NewStyle().PaddingLeft(3).MaxWidth(m.size.Width).MaxHeight(m.size.Height).Render(dialog)
}

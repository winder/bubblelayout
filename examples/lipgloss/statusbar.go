package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	bl "github.com/winder/bubblelayout"
)

var (
	statusNugget = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Padding(0, 1)

	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
			Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"})

	statusStyle = lipgloss.NewStyle().
			Inherit(statusBarStyle).
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#FF5F87")).
			Padding(0, 1).
			MarginRight(1)

	encodingStyle = statusNugget.Copy().
			Background(lipgloss.Color("#A550DF")).
			Align(lipgloss.Right)

	statusText = lipgloss.NewStyle().Inherit(statusBarStyle)

	fishCakeStyle = statusNugget.Copy().Background(lipgloss.Color("#6124DF"))
)

type statusbarModel struct {
	ID      bl.ID
	size    bl.Size
	lastKey string
	repeat  int
}

func (m statusbarModel) Init() tea.Cmd {
	return nil
}

func (m statusbarModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.lastKey == msg.String() {
			m.repeat++
		} else {
			m.lastKey = msg.String()
			m.repeat = 1
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

func (m statusbarModel) View() string {
	w := lipgloss.Width

	statusKey := statusStyle.Render("STATUS")
	encoding := encodingStyle.Render("UTF-8")
	fishCake := fishCakeStyle.Render("🍥 Fish Cake")
	var statusString string
	if m.repeat == 1 {
		statusString = fmt.Sprintf("Ravishing (last key: %s)", m.lastKey)
	} else if m.repeat > 1 {
		statusString = fmt.Sprintf("Ravishing (last key: %s x%d)", m.lastKey, m.repeat)
	} else {
		statusString = "Ravishing"
	}
	statusVal := statusText.Copy().
		Width(m.size.Width - w(statusKey) - w(encoding) - w(fishCake)).
		Render(statusString)

	bar := lipgloss.JoinHorizontal(lipgloss.Top,
		statusKey,
		statusVal,
		encoding,
		fishCake,
	)

	return statusBarStyle.Width(m.size.Width).MaxWidth(m.size.Width).Render(bar)
}

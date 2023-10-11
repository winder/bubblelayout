package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	bl "github.com/winder/bubblelayout"
)

var (
	titleStyle = lipgloss.NewStyle().
			MarginLeft(1).
			MarginRight(5).
			Padding(0, 1).
			Italic(true).
			Foreground(lipgloss.Color("#FFF7DB")).
			SetString("Lip Gloss")

	colors = colorGrid(1, 5)
)

type titleModel struct {
	ID          bl.ID
	size        bl.Size
	colorOffset int
}

func (m titleModel) Init() tea.Cmd {
	return nil
}

func (m titleModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "right", "l", "n", "tab":
			m.colorOffset--
			if m.colorOffset < 0 {
				m.colorOffset = len(colors) - 1
			}
			return m, nil
		case "left", "h", "p", "shift+tab":
			m.colorOffset++
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

func (m titleModel) View() string {
	var buf strings.Builder
	for i := 0; i < m.size.Height-2; i++ {
		const offset = 2
		c := lipgloss.Color(colors[(i+m.colorOffset)%len(colors)][0])
		fmt.Fprint(&buf, titleStyle.Copy().MarginLeft(2+i*offset).Background(c))
		if i < m.size.Height {
			buf.WriteRune('\n')
		}
	}
	return lipgloss.NewStyle().Margin(1, 0).MaxHeight(m.size.Height).MaxWidth(m.size.Width).Render(buf.String())
}

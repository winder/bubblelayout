package main

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	bl "github.com/winder/bubblelayout"
)

type gridModel struct {
	ID          bl.ID
	size        bl.Size
	colorOffset int
	cached      string
}

func MakeGridModel(id bl.ID) gridModel {
	return gridModel{
		ID:          id,
		colorOffset: 0,
		cached:      grid(0, 15, 8),
	}
}

func (m gridModel) Init() tea.Cmd {
	return nil
}

func grid(offset, width, height int) string {
	colors := colorGrid(width/2, height)

	b := strings.Builder{}
	for xi := range colors {
		x := (xi + offset) % len(colors)
		for yi := range colors[x] {
			y := colors[x][(yi+offset)%len(colors[x])]
			s := lipgloss.NewStyle().SetString("  ").Background(lipgloss.Color(y))
			b.WriteString(s.String())
		}
		b.WriteRune('\n')
	}

	return b.String()
}

func (m gridModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "right", "l", "n", "tab":
			m.colorOffset--
			if m.colorOffset < 0 {
				m.colorOffset = m.size.Width/2 - 3
			}
			m.cached = grid(m.colorOffset, m.size.Width-5, m.size.Height-2)
			return m, nil
		case "left", "h", "p", "shift+tab":
			m.colorOffset++
			m.cached = grid(m.colorOffset, m.size.Width-5, m.size.Height-2)
			return m, nil
		}
	case bl.BubbleLayoutMsg:
		var err error
		m.size, err = msg.Size(m.ID)
		m.cached = grid(m.colorOffset, m.size.Width-5, m.size.Height-2)
		if err != nil {
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m gridModel) View() string {
	return lipgloss.NewStyle().
		MarginLeft(2).
		MarginRight(3).
		MarginTop(1).
		MaxHeight(m.size.Height).
		MaxWidth(m.size.Width).
		Render(m.cached)
}

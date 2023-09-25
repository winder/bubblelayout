package util

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/winder/layout"
)

// simpleModel listens to BubbleLayoutMsg events and displays a colored box and its ID.
type simpleModel struct {
	id         layout.ID
	background lipgloss.Color
	w, h       int
	message    string
}

func NewSimpleModel(bg lipgloss.Color, id layout.ID) tea.Model {
	return simpleModel{
		background: bg,
		id:         id,
	}
}

func NewSimpleModelWithMessage(bg lipgloss.Color, id layout.ID, message string) tea.Model {
	return simpleModel{
		background: bg,
		id:         id,
		message:    message,
	}
}

func (m simpleModel) Init() tea.Cmd {
	return nil
}

func (m simpleModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(layout.BubbleLayoutMsg); ok {
		size, err := msg.Size(m.id)
		if err != nil {
			panic(err)
		}
		m.w = size.Width
		m.h = size.Height
	}
	return m, nil
}

func (m simpleModel) View() string {
	st := lipgloss.NewStyle().
		Background(m.background).
		Width(m.w).
		Height(m.h).
		Align(lipgloss.Center)
	if m.message != "" {
		return st.Render(m.message)
	}
	return st.Render(fmt.Sprintf("%d", m.id))
}

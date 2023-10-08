package main

import (
	"math/rand"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	bl "github.com/winder/bubblelayout"
)

var list = lipgloss.NewStyle().
	PaddingLeft(4).
	PaddingRight(2)

var listHeader = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderBottom(true).
	BorderForeground(subtle).
	PaddingRight(2)

var listItem = lipgloss.NewStyle().PaddingLeft(2).Render

var checkMark = lipgloss.NewStyle().SetString("✓").
	Foreground(special).
	PaddingRight(1).
	String()

var listDone = func(s string) string {
	return checkMark + lipgloss.NewStyle().
		Strikethrough(true).
		Foreground(lipgloss.AdaptiveColor{Light: "#969B86", Dark: "#696969"}).
		Render(s)
}

var finalWithSeparator = func(width, height int, view string) string {
	return lipgloss.NewStyle().Align(lipgloss.Left).
		Border(lipgloss.NormalBorder(), false, true, false, false).
		BorderForeground(subtle).
		MarginTop(1).
		Width(width - 1). // space for border
		Height(height).
		MaxWidth(width).
		MaxHeight(height - 1). // the right border  was going too low...
		Render(view)
}

type listItem2 struct {
	done bool
	text string
}

type listModel struct {
	ID    bl.ID
	size  bl.Size
	title string
	items []listItem2
}

func MakeListModel(id bl.ID, title string, items []listItem2) listModel {
	m := listModel{
		ID:    id,
		title: title,
		items: items,
	}
	m.randomize()
	return m
}

func (m listModel) randomize() {
	src := rand.NewSource(time.Now().UnixNano())
	for i := range m.items {
		m.items[i].done = src.Int63()%2 == 0
	}
}

func (m listModel) Init() tea.Cmd {
	m.randomize()
	return nil
}

func (m listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "right", "l", "n", "tab":
			fallthrough
		case "left", "h", "p", "shift+tab":
			m.randomize()
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

func (m listModel) View() string {
	var items []string
	items = append(items, listHeader.MaxWidth(m.size.Width).Render(m.title))
	for _, item := range m.items {
		if item.done {
			items = append(items, listDone(item.text))
		} else {
			items = append(items, listItem(item.text))
		}
	}
	view := list.Copy().MaxWidth(m.size.Width - 2).Render(
		lipgloss.JoinVertical(lipgloss.Left, items...),
	)

	return finalWithSeparator(m.size.Width, m.size.Height, view)
}

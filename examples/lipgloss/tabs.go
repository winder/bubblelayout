// Mostly copied from charmbracelet/bubbletea/examples/tabs/main.go
package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	bl "github.com/winder/bubblelayout"
)

var (
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	docStyle          = lipgloss.NewStyle().PaddingLeft(2).PaddingRight(2)
	highlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Copy().Border(activeTabBorder, true)
	remainderStyle    = inactiveTabStyle.Copy().Border(inactiveTabBorder, false, false, true, false)
)

type tabModel struct {
	ID        bl.ID
	Tabs      []string
	activeTab int
	size      bl.Size
}

func (m tabModel) Init() tea.Cmd {
	return nil
}

func (m tabModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "right", "l", "n", "tab":
			m.activeTab = (m.activeTab + 1) % len(m.Tabs)
			return m, nil
		case "left", "h", "p", "shift+tab":
			m.activeTab--
			if m.activeTab < 0 {
				m.activeTab = len(m.Tabs) - 1
			}
			return m, nil
		}
	case bl.BubbleLayoutMsg:
		var err error
		m.size, err = msg.Size(m.ID)
		if err != nil {
			return m, tea.Quit
		}
		if m.size.Height > 3 {
			m.size.Height = m.size.Height
		}
	}

	return m, nil
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

func (m tabModel) View() string {
	var renderedTabs []string
	w := 0

	for i, t := range m.Tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(m.Tabs)-1, i == m.activeTab
		if isActive {
			style = activeTabStyle.Copy()
		} else {
			style = inactiveTabStyle.Copy()
		}
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "┘"
		} else if isFirst && !isActive {
			border.BottomLeft = "┴"
		} else if isLast && isActive {
			border.BottomRight = "└"
		} else if isLast && !isActive {
			border.BottomRight = "┴"
		}
		style = style.Border(border)
		tab := style.Render(t)
		w += lipgloss.Width(tab)
		renderedTabs = append(renderedTabs, tab)
	}

	remainder := m.size.Width - w - 4
	if remainder > 0 {
		// Height = 2 + 1 for the border
		rs := remainderStyle.Copy().Width(remainder).Height(m.size.Height - 1)
		renderedTabs = append(renderedTabs, rs.Render(""))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)

	return docStyle.Copy().MaxWidth(m.size.Width).MaxHeight(m.size.Height).Render(row)
}

package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"examples/util"
	bl "github.com/winder/bubblelayout"
)

func New() tea.Model {
	layout := bl.New()

	// Each model is initialized along with a layout ID.
	title := titleModel{ID: layout.Add("height 7")}
	description := descModel{ID: layout.Add("span 2, height 3,  wrap")}
	dialog := dialogModel{ID: layout.Add("span 3, wrap")}
	list1 := MakeListModel(layout.Add("height 9"), list1Title, list1)
	list2 := MakeListModel(layout.Add("height 9"), list2Title, list2)
	grid := MakeGridModel(layout.Add("wrap"))
	history1 := MakeHistory(layout.Add("spanh 2"), historyA, lipgloss.Right, 2)
	history2 := MakeHistory(layout.Add("spanh 2"), historyB, lipgloss.Center, 0)
	history3 := MakeHistory(layout.Add("spanh 2"), historyC, lipgloss.Left, 0)

	// Tab header and status bar are initialized as usual and docked north and south.
	tabs := tabModel{
		Tabs: []string{"Lip Gloss", "Blush", "Eye Shadow", "Mascara", "Foundation"},
		ID:   layout.Add("north 3!"),
	}
	statusbar := statusbarModel{
		ID: layout.Add("south 1!"),
	}

	// The models are collected into a slice for the layout model.
	models := []tea.Model{
		title, description,
		dialog,
		list1, list2, grid,
		history1, history2, history3,
		tabs,
		statusbar,
	}

	// The view function glues everything together. It is called by the LayoutModel utility.
	view := func(models []tea.Model) string {
		return lipgloss.JoinVertical(0,
			models[tabs.ID-1].View(),
			lipgloss.JoinHorizontal(0,
				models[title.ID-1].View(),
				models[description.ID-1].View(),
			),
			models[dialog.ID-1].View(),
			lipgloss.JoinHorizontal(0,
				models[list1.ID-1].View(),
				models[list2.ID-1].View(),
				models[grid.ID-1].View(),
			),
			lipgloss.JoinHorizontal(0,
				models[history1.ID-1].View(),
				models[history2.ID-1].View(),
				models[history3.ID-1].View(),
			),
			models[statusbar.ID-1].View(),
		)
	}

	// This is an example utility that calls update and converts tea.WindowSizeMsg to bl.BubbleLayoutMsg.
	// It may be useful for real programs but is not part of the bubblelayout library.
	return util.NewLayoutModel(models, layout, view)
}

func main() {
	p := tea.NewProgram(New())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

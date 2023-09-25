package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"examples/util"
	"github.com/winder/layout"
)

func New() tea.Model {
	// -------------------------
	// |   |     NORTH     |   |
	// |   |----------------   |
	// | W |   0   |   1   | E |
	// | E |---------------- A |
	// | S |       -       | S |
	// | T | -  -  2  -  - | T |
	// |   |       -       |   |
	// |   |----------------   |
	// |   |     SOUTH     |   |
	// -------------------------

	bl := layout.New()
	var models []tea.Model
	models = append(models, util.NewSimpleModel("9", bl.Add(layout.Layout{})))
	models = append(models, util.NewSimpleModel("10", bl.Add(layout.Layout{})))
	bl.Wrap()
	models = append(models, util.NewSimpleModel("11", bl.Add(layout.Layout{SpanWidth: 2, SpanHeight: 2})))

	models = append(models, util.NewSimpleModel("12", bl.Dock(layout.NORTH, 1, 1, 1)))
	models = append(models, util.NewSimpleModel("13", bl.Dock(layout.SOUTH, 1, 1, 1)))
	models = append(models, util.NewSimpleModel("14", bl.Dock(layout.WEST, 1, 10, 10)))
	models = append(models, util.NewSimpleModel("15", bl.Dock(layout.EAST, 1, 10, 10)))

	view := func(models []tea.Model) string {
		// Note: docks should be joined in the order they are defined.
		center := lipgloss.JoinVertical(0,
			models[3].View(), // north
			lipgloss.JoinHorizontal(0, models[0].View(), models[1].View()),
			models[2].View(),
			models[4].View()) // south
		return lipgloss.JoinHorizontal(0,
			models[5].View(), // west
			center,
			models[6].View()) // east
	}
	return util.NewLayoutModel(models, bl, view)
}

func main() {
	p := tea.NewProgram(New())
	p.Run()
}

package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"examples/util"
	bl "github.com/winder/bubblelayout"
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

	layout := bl.New()
	var models []tea.Model
	models = append(models, util.NewSimpleModel("9", layout.Add(bl.Cell{})))
	models = append(models, util.NewSimpleModel("10", layout.Add(bl.Cell{})))
	layout.Wrap()
	models = append(models, util.NewSimpleModel("11", layout.Add(bl.Cell{SpanWidth: 2, SpanHeight: 2})))

	models = append(models, util.NewSimpleModel("12", layout.Dock(bl.Dock{Cardinal: bl.NORTH, Min: 1, Preferred: 1, Max: 1})))
	models = append(models, util.NewSimpleModel("13", layout.Dock(bl.Dock{Cardinal: bl.SOUTH, Min: 1, Preferred: 1, Max: 1})))
	models = append(models, util.NewSimpleModel("14", layout.Dock(bl.Dock{Cardinal: bl.WEST, Min: 1, Preferred: 10, Max: 10})))
	models = append(models, util.NewSimpleModel("15", layout.Dock(bl.Dock{Cardinal: bl.EAST, Min: 1, Preferred: 10, Max: 10})))

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
	return util.NewLayoutModel(models, layout, view)
}

func main() {
	p := tea.NewProgram(New())
	p.Run()
}

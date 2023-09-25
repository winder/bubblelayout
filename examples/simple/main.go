package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"examples/util"
	"github.com/winder/layout"
)

func New() tea.Model {
	bl := layout.New()
	var models []tea.Model
	models = append(models, util.NewSimpleModel("9", bl.Add(layout.Layout{})))
	models = append(models, util.NewSimpleModel("10", bl.Add(layout.Layout{SpanWidth: 2, SpanHeight: 2})))
	models = append(models, util.NewSimpleModel("11", bl.Add(layout.Layout{})))
	bl.Wrap()
	models = append(models, util.NewSimpleModel("12", bl.Add(layout.Layout{SpanHeight: 2})))
	models = append(models, util.NewSimpleModel("13", bl.Add(layout.Layout{})))
	bl.Wrap()
	models = append(models, util.NewSimpleModel("14", bl.Add(layout.Layout{})))
	models = append(models, util.NewSimpleModel("15", bl.Add(layout.Layout{SpanWidth: 2})))

	view := func(models []tea.Model) string {
		// Glue the views together.
		// ---------------------------------
		// |   0   |       -       |   2   |
		// --------- -  -  1  -  - |--------
		// |   -   |       -       |   4   |
		// | - 3 - -------------------------
		// |   -   |   5   |       6       |
		// ---------------------------------
		left := lipgloss.JoinVertical(0, models[0].View(), models[3].View())
		right := lipgloss.JoinVertical(0, models[2].View(), models[4].View())
		bottom := lipgloss.JoinHorizontal(0, models[5].View(), models[6].View())
		center := lipgloss.JoinHorizontal(0, models[1].View(), right)
		right = lipgloss.JoinVertical(0, center, bottom)
		return lipgloss.JoinHorizontal(0, left, right)
	}
	return util.NewLayoutModel(models, bl, view)
}

func main() {
	p := tea.NewProgram(New())
	p.Run()
}

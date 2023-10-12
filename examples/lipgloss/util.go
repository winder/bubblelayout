package main

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
)

// Random UI constants
var (
	list1Title = "Citrus Fruits to Try"
	list1      = []listItem2{
		{done: true, text: "Grapefruit"},
		{done: true, text: "Yuzu"},
		{done: true, text: "Citron"},
		{done: true, text: "Kumquat"},
		{done: true, text: "Pomelo"},
	}

	list2Title = "Actual Lip Gloss Vendors"
	list2      = []listItem2{
		{done: true, text: "Glossier"},
		{done: true, text: "Claire‘s Boutique"},
		{done: true, text: "Nyx"},
		{done: true, text: "Mac"},
		{done: true, text: "Milk"},
	}

	historyA = "The Romans learned from the Greeks that quinces slowly cooked with honey would “set” when cool. The Apicius gives a recipe for preserving whole quinces, stems and leaves attached, in a bath of honey diluted with defrutum: Roman marmalade. Preserves of quince and lemon appear (along with rose, apple, plum and pear) in the Book of ceremonies of the Byzantine Emperor Constantine VII Porphyrogennetos."
	historyB = "Medieval quince preserves, which went by the French name cotignac, produced in a clear version and a fruit pulp version, began to lose their medieval seasoning of spices in the 16th century. In the 17th century, La Varenne provided recipes for both thick and clear cotignac."
	historyC = "In 1524, Henry VIII, King of England, received a “box of marmalade” from Mr. Hull of Exeter. This was probably marmelada, a solid quince paste from Portugal, still made and sold in southern Europe today. It became a favourite treat of Anne Boleyn and her ladies in waiting."
)

var (
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
)

func colorGrid(xSteps, ySteps int) [][]string {
	if xSteps <= 0 || ySteps <= 0 {
		return nil
	}
	x0y0, _ := colorful.Hex("#F25D94")
	x1y0, _ := colorful.Hex("#EDFF82")
	x0y1, _ := colorful.Hex("#643AFF")
	x1y1, _ := colorful.Hex("#14F9D5")

	x0 := make([]colorful.Color, ySteps)
	for i := range x0 {
		x0[i] = x0y0.BlendLuv(x0y1, float64(i)/float64(ySteps))
	}

	x1 := make([]colorful.Color, ySteps)
	for i := range x1 {
		x1[i] = x1y0.BlendLuv(x1y1, float64(i)/float64(ySteps))
	}

	result := make([][]string, ySteps)
	for x := 0; x < ySteps; x++ {
		y0 := x0[x]
		result[x] = make([]string, xSteps)
		for y := 0; y < xSteps; y++ {
			result[x][y] = y0.BlendLuv(x1[x], float64(y)/float64(xSteps)).Hex()
		}
	}

	return result
}

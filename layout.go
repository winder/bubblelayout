package bubblelayout

import (
	"fmt"
)

type ID uint64
type Cardinal uint

type Size struct {
	Width  int
	Height int
}

type BubbleLayoutMsg struct {
	size map[ID]*Size
}

// Size returns the size allocated for a view.
func (l BubbleLayoutMsg) Size(id ID) (Size, error) {
	s, ok := l.size[id]
	if !ok {
		return Size{}, fmt.Errorf("view not registered")
	}
	return *s, nil
}

const (
	NORTH Cardinal = iota + 1
	SOUTH
	EAST
	WEST
)

type preferenceGroup []sizePreference

func (pg preferenceGroup) computeDims(allocated int) []int {
	if len(pg) == 0 {
		return nil
	}

	largestPref := func(maximum int, p []int) int {
		if len(p) == 0 {
			return 0
		}
		ret := 0
		for _, v := range p {
			ret = max(ret, v)
		}
		if maximum != 0 {
			return min(maximum, ret)
		}
		return ret
	}

	dims := make([]int, len(pg))
	sum := 0
	// start at min or preference and grow to the allocated size.
	growers := make(map[int]struct{})
	// no preference, start at min and maybe grow to max.
	empty := make(map[int]struct{})
	evenMax := allocated / len(pg)
	for idx, p := range pg {
		pref := largestPref(p.max, p.preferred)
		pg[idx].largestPref = pref
		if pref != 0 {
			dims[idx] = min(pref, evenMax)
			sum += min(pref, evenMax)
		} else {
			dims[idx] = p.min
			sum += p.min
			// don't count as empty and a grower
			if !p.grow && p.max != 0 {
				empty[idx] = struct{}{}
			}
		}
		if p.grow || p.max == 0 {
			growers[idx] = struct{}{}
		}
	}

	// if all preferences are fullfilled and nothing is growing, return the preferred sizes
	if len(growers) == 0 && len(empty) == 0 {
		return dims
	}

	// offer an even split to all empty and growing components.
	// a second pass is made for growers in case the empty components have a max.
	remainder := allocated - sum

	// keep loping until space runs out or the maximizable components are maximized.
	for remainder > (len(empty)+len(growers)) && len(empty) > 0 {
		split := remainder / (len(growers) + len(empty))
		for idx, p := range pg {
			if _, ok := empty[idx]; ok {
				// otherwise grow up to the max or split
				diff := p.max - dims[idx]
				if p.max != 0 && diff < split {
					// grow to max
					dims[idx] = dims[idx] + diff
					remainder -= diff
					delete(empty, idx)
				} else {
					// grow to split
					dims[idx] = dims[idx] + split
					remainder -= split
				}
			}
		}
	}

	// if there is a remainder, loop through again but this time only add to the growers.
	last := -1
	if remainder != 0 && len(growers) != 0 {
		split := remainder / len(growers)
		for idx := range growers {
			dims[idx] = dims[idx] + split
			remainder -= split
			last = idx
		}
	}

	// if there is still a remainder, it is a rounding error. Add it to the last grower
	if last != -1 && remainder != 0 {
		dims[last] = dims[last] + remainder
	}

	return dims
}

type sizePreference struct {
	min         int
	preferred   []int
	largestPref int
	max         int
	grow        bool
}

type Grid [][]layout

func (g Grid) makeMessage(wDims, hDims []int) BubbleLayoutMsg {
	msg := BubbleLayoutMsg{
		size: make(map[ID]*Size),
	}

	// to avoid double counting spanning cells, keep track of which rows and column was used to process a layout size.
	idRow := make(map[ID]int)
	idCol := make(map[ID]int)
	for rowIdx, row := range g {
		for colIdx, l := range row {
			if _, ok := msg.size[l.id]; !ok {
				msg.size[l.id] = &Size{}
			}
			if _, ok := idRow[l.id]; !ok {
				idRow[l.id] = rowIdx
			}
			if _, ok := idCol[l.id]; !ok {
				idCol[l.id] = colIdx
			}
			if idRow[l.id] == rowIdx {
				msg.size[l.id].Width += wDims[colIdx]
			}
			if idCol[l.id] == colIdx {
				msg.size[l.id].Height += hDims[rowIdx]
			}
		}
	}
	return msg
}

// TODO:
//   print function?
//   compare function?

// layout holds the Cell or Dock information in addition to the ID.
type layout struct {
	id ID

	// wrap indicates that the grid should wrap to the next row after this Layout.
	wrap bool

	Cell
	Dock
}

// Cell defines the size and position that should be allocated for a view.
type Cell struct {
	// SpanWidth defines the number of columns that the view should span. Defaults to 1.
	SpanWidth int
	// SpanHeight defines the number of rows that the view should span. Defaults to 1.
	SpanHeight int

	// MinWidth overrides the minimum width that should be allocated for the view.
	MinWidth int
	// PreferredWidth overrides the preferred width that should be allocated for the view.
	PreferredWidth int
	// MaxWidth overrides the maximum width that should be allocated for the view.
	MaxWidth int

	// MinHeight overrides the minimum height that should be allocated for the view.
	MinHeight int
	// PreferredHeight overrides the preferred height that should be allocated for the view.
	PreferredHeight int
	// MaxHeight overrides the maximum height that should be allocated for the view.
	MaxHeight int

	// GrowWidth indicates that the horizontal size should be maximized.
	GrowWidth bool
	// GrowHeight indicates that the vertical size should be maximized.
	GrowHeight bool

	// wDuplicate is used as part of horizontal spanning for calculating dimensions.
	wDuplicate bool
	// hDuplicate is used as part of vertical spanning for calculating dimensions.
	hDuplicate bool
}

// Dock defines a component that should span an entire side of the layout.
type Dock struct {
	// Cardinal indicates which side of the layout the view should be docked to.
	Cardinal Cardinal

	// Min overrides the minimum width or height that should be allocated for the view.
	Min int

	// Preferred overrides the preferred width or height that should be allocated for the view.
	Preferred int

	// Max overrides the maximum width or height that should be allocated for the view.
	Max int
}

type BubbleLayout interface {
	AddStr(string) ID
	Add(Cell) ID
	Dock(Dock) ID
	Wrap()
	Resize(width, height int) BubbleLayoutMsg
	Validate() error
}

func New() BubbleLayout {
	return &bubbleLayout{
		layouts: [][]layout{{}},
	}
}

type bubbleLayout struct {
	idCounter ID
	layouts   Grid
	docks     []layout

	// resizeCache is the layouts after being merged with the docks.
	resizeCache Grid
	hPref       preferenceGroup
	wPref       preferenceGroup
}

// AddStr uses the string notation to define the layout. This is often shorter and easier to read than using the Layout struct.
func (bl *bubbleLayout) AddStr(str string) ID {
	// see http://www.miglayout.com/QuickStart.pdf
	panic("not implemented")
}

// Add a tea.Model to the layout. The model will be placed in the next available cell.
func (bl *bubbleLayout) Add(c Cell) ID {
	bl.idCounter++
	l := layout{
		id:   bl.idCounter,
		Cell: c,
	}
	idx := len(bl.layouts) - 1
	l.id = bl.idCounter
	bl.layouts[idx] = append(bl.layouts[idx], l)

	if l.wrap {
		bl.layouts = append(bl.layouts, []layout{})
	}

	// TODO: Debug mode which panics here as soon as a constraint violation is detected.
	return bl.idCounter
}

// Wrap inserts a new row into the layout, subsequent calls to Add will place models in the new row.
func (bl *bubbleLayout) Wrap() {
	bl.layouts = append(bl.layouts, []layout{})
}

// Dock places a model on the edge of the layout, spanning the entire width or height.
// For NORTH and SOUTH components, the width is fixed and the height is defined by min, preferred and max.
// For EAST and WEST components, the height is fixed and the width is defined by min, preferred and max.
func (bl *bubbleLayout) Dock(dock Dock) ID {
	bl.idCounter++
	bl.docks = append(bl.docks, layout{id: bl.idCounter, Dock: dock})
	return bl.idCounter
}

type preferenceConstraintError struct {
	row           bool
	idx, min, max int
}

func (p preferenceConstraintError) Error() string {
	var dir string
	var dim string
	if p.row {
		dir = "row"
		dim = "width"
	} else {
		dir = "col"
		dim = "height"
	}
	return fmt.Sprintf("constraint violation: %s %d: min %s (%d), max %s (%d)", dir, p.idx, dim, p.min, dim, p.max)
}

func makeRowViolation(idx, min, max int) error {
	return preferenceConstraintError{row: true, idx: idx, min: min, max: max}
}

func makeColViolation(idx, min, max int) error {
	return preferenceConstraintError{row: false, idx: idx, min: min, max: max}
}

func checkPreferenceConstraints(hPref, wPref preferenceGroup) error {
	for row, p := range hPref {
		if p.min > p.max {
			return makeRowViolation(row, p.min, p.max)
		}
	}

	for col, p := range wPref {
		if p.min > p.max {
			return makeColViolation(col, p.min, p.max)
		}
	}

	return nil
}

// expandSpans takes a layout and splits all spans into individual cells. This is a simplification, because
// the span could possibly respect other row/column preferences, but we're discarding the relationship once the
// span has been split to simplify the code.
//
// Here is an example, the parens denote (spanx, spany) overrides:
// ---------------------------------
// | 1        | 2 (2, 2) |    3    |
// ---------------------------------
// | 4 (1, 2) |    5     |
// -----------------------
// |    6     | 7 (1, 2) |
// -----------------------
//
// Turn into this:
// ---------------------------------
// |   1   |       -       |   3   |
// --------- -  -  2  -  - |--------
// |   -   |       -       |   5   |
// | - 4 - -------------------------
// |   -   |   6   |       7       |
// ---------------------------------
//
// In the above example, the 2x2 cell is split into 4 cells, and the 1x2 cells are split into 2 cells.
func expandSpans(layouts Grid) Grid {
	ret := make(Grid, len(layouts))
	for i := 0; i < len(layouts); i++ {
		ret[i] = make([]layout, 0, len(layouts[i]))
		ret[i] = append(ret[i], layouts[i]...)
	}

	longestCol := 0
	for _, row := range ret {
		longestCol = max(longestCol, len(row))
	}

	// spanx and create rows
	for colIdx := 0; colIdx < longestCol; colIdx++ {
		for rowIdx := 0; rowIdx < len(ret); rowIdx++ {
			// skip if this was an empty cell
			if len(ret[rowIdx]) < (colIdx + 1) {
				continue
			}

			// vertical span duplicate handling.
			spanHeight := ret[rowIdx][colIdx].SpanHeight
			if ret[rowIdx][colIdx].SpanHeight > 1 && !ret[rowIdx][colIdx].hDuplicate {
				// TODO: fix rounding errors?
				l := ret[rowIdx][colIdx]
				l.MinHeight /= l.SpanHeight
				l.MaxHeight /= l.SpanHeight
				l.PreferredHeight /= l.SpanHeight

				if ret[rowIdx][colIdx].SpanWidth > 1 && ret[rowIdx][colIdx].wDuplicate {
					// already handled by the horizontal span duplicate handling.
				} else {
					ret[rowIdx][colIdx] = l
					l.hDuplicate = true
					for i := 1; i < spanHeight; i++ {
						// create next row if needed
						if len(ret) <= (rowIdx + i) {
							// pad next row so that we can insert the new cell.
							ret = append(ret, make([]layout, colIdx))
						}
						if len(ret[rowIdx+i]) == 0 {
							// special case for first empty row.
							ret[rowIdx+i] = append(ret[rowIdx+i], l)
						} else {
							ret[rowIdx+i] = append(ret[rowIdx+i][:colIdx+i], ret[rowIdx+i][colIdx+i-1:]...)
							ret[rowIdx+i][colIdx] = l
						}
					}
				}
			}

			// horizontal span duplicate handling.
			spanWidth := ret[rowIdx][colIdx].SpanWidth
			if spanWidth > 1 && !ret[rowIdx][colIdx].wDuplicate {
				l := ret[rowIdx][colIdx]
				// TODO: fix rounding errors?
				l.MinWidth /= l.SpanWidth
				l.MaxWidth /= l.SpanWidth
				l.PreferredWidth /= l.SpanWidth
				ret[rowIdx][colIdx] = l
				l.wDuplicate = true

				for i := 1; i < spanWidth; i++ {
					// make room for a new element
					ret[rowIdx] = append(ret[rowIdx][:colIdx+i], ret[rowIdx][colIdx+i-1:]...)
					ret[rowIdx][colIdx+i] = l
				}
				// grow the longest column if necessary.
				longestCol = max(len(ret[rowIdx]), longestCol)
			}
		}
	}

	return ret
}

// mergeDocks takes a layout and merges the docked layouts. Returns the new layout and width/height deltas.
// This function is called after expandSpans, so it must expand the spans as part of adding the dock.
func mergeDocks(g Grid, docks []layout) Grid {
	if len(g) == 0 {
		return nil
	}

	// Make a copy
	ret := make(Grid, len(g))
	for i := 0; i < len(g); i++ {
		ret[i] = make([]layout, 0, len(g[i]))
		ret[i] = append(ret[i], g[i]...)
	}

	gridHeight := len(g)
	gridWidth := len(g[0])

	// merge docked layouts into the resize cache.
	for _, d := range docks {
		switch d.Cardinal {
		case NORTH:
			// Add it to the first row, spanning the entire width.
			north := layout{
				id: d.id,
				Cell: Cell{
					SpanWidth:       gridWidth,
					MinHeight:       d.Min,
					PreferredHeight: d.Preferred,
					MaxHeight:       d.Max,
				},
			}
			northRow := make([]layout, 0, gridWidth)
			for i := 0; i < gridWidth; i++ {
				northRow = append(northRow, north)
				north.wDuplicate = true // the second and on are duplicate.
			}
			ret = append([][]layout{northRow}, ret...)
			gridHeight++
		case SOUTH:
			// Add it to the final row, spanning the entire width.
			south := layout{
				id: d.id,
				Cell: Cell{
					SpanWidth:       gridWidth,
					MinHeight:       d.Min,
					PreferredHeight: d.Preferred,
					MaxHeight:       d.Max,
				},
			}
			southRow := make([]layout, 0, gridWidth)
			for i := 0; i < gridWidth; i++ {
				southRow = append(southRow, south)
				south.wDuplicate = true // the second and on are duplicate.
			}
			ret = append(ret, southRow)
			gridHeight++
		case EAST:
			// Add it to the end of each row to span the entire height.
			east := layout{
				id: d.id,
				Cell: Cell{
					SpanHeight:     gridHeight,
					MinWidth:       d.Min,
					PreferredWidth: d.Preferred,
					MaxWidth:       d.Max,
				},
			}
			for i := 0; i < gridHeight; i++ {
				ret[i] = append(ret[i], east)
				east.hDuplicate = true
			}
			gridWidth++
		case WEST:
			// Add it to the front of each row to span the entire height.
			west := layout{
				id: d.id,
				Cell: Cell{
					SpanHeight:     gridHeight,
					MinWidth:       d.Min,
					PreferredWidth: d.Preferred,
					MaxWidth:       d.Max,
				},
			}
			for i := 0; i < gridHeight; i++ {
				ret[i] = append([]layout{west}, ret[i]...)
				west.hDuplicate = true
			}
			gridWidth++
		default:
			panic(fmt.Errorf("invalid cardinal"))
		}
	}
	return ret
}

// distillPreferences attempts to normalize the different preferences for cells
// across each row and column. This is done to simplify the resizing algorithm.
//
// For example, the minimum width for a column would be the minimum across all
// layouts in the first column.
func distillPreferences(g Grid) (hPref, wPref preferenceGroup) {
	if len(g) == 0 {
		return
	}

	gridHeight := len(g)
	gridWidth := len(g[0])
	hPref = make(preferenceGroup, gridHeight)
	wPref = make(preferenceGroup, gridWidth)
	for rowIdx := 0; rowIdx < gridHeight; rowIdx++ {
		for colIdx := 0; colIdx < gridWidth; colIdx++ {
			l := g[rowIdx][colIdx]

			// collect height preferences
			if l.MinHeight != 0 {
				hPref[rowIdx].min = max(l.MinHeight, hPref[rowIdx].min)
			}
			if l.MaxHeight != 0 {
				if hPref[rowIdx].max == 0 {
					hPref[rowIdx].max = l.MaxHeight
				} else {
					hPref[rowIdx].max = min(l.MaxHeight, hPref[rowIdx].max)
				}
			}
			if l.PreferredHeight != 0 {
				hPref[rowIdx].preferred = append(hPref[rowIdx].preferred, l.PreferredHeight)
			}
			hPref[rowIdx].grow = hPref[rowIdx].grow || l.GrowHeight

			// collect width preferences
			if l.MinWidth != 0 {
				wPref[colIdx].min = max(l.MinWidth, wPref[colIdx].min)
			}
			if l.MaxWidth != 0 {
				if wPref[colIdx].max == 0 {
					wPref[colIdx].max = l.MaxWidth
				} else {
					wPref[colIdx].max = min(l.MaxWidth, wPref[colIdx].max)
				}
			}
			if l.PreferredWidth != 0 {
				wPref[colIdx].preferred = append(wPref[colIdx].preferred, l.PreferredWidth)
			}
			wPref[colIdx].grow = wPref[colIdx].grow || l.GrowWidth
		}
	}
	return
}

func (bl *bubbleLayout) Validate() error {
	if len(bl.resizeCache) == 0 {
		bl.resizeCache = expandSpans(bl.layouts)
		bl.resizeCache = mergeDocks(bl.resizeCache, bl.docks)

		bl.hPref, bl.wPref = distillPreferences(bl.resizeCache)
		return checkPreferenceConstraints(bl.hPref, bl.wPref)
	}
	return nil
}

// Resize recalculates the layout based on the current terminal size.
func (bl *bubbleLayout) Resize(width, height int) BubbleLayoutMsg {
	if err := bl.Validate(); err != nil {
		panic(err)
	}

	hDims := bl.hPref.computeDims(height)
	wDims := bl.wPref.computeDims(width)

	return bl.resizeCache.makeMessage(wDims, hDims)
}

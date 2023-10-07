package bubblelayout

import (
	"fmt"
)

type ID uint64
type Cardinal string

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
	NORTH Cardinal = "north"
	SOUTH Cardinal = "south"
	EAST  Cardinal = "east"
	WEST  Cardinal = "west"
)

type PreferenceGroup []BoundSize

// computeDims takes a list of BoundSizes and an allocated size and returns the actual size that should be allocated to each component.
// TODO: detect if the minimums add up to more than the allocated size and generate a constraint violation error.
//
//	pass 1: allocate minimums. Priority is given left to right.
//	        -> if minimums fill up the allocated size, everything else remains 0.
//	pass 2: evenly expand up to min(preferred, max), or an even split.
//	        TODO: This is done in a loop. Can it be done in a single iteration?
//	pass 3a: If "grow" is used, allocate remaining space to growers.
//	pass 3b: Otherwise, allocate remaining space to "max" or cells with no max.
//
//	TODO: Grow priorities.
//	TODO: What does it mean to have Grow and Max? Can it go over the Max?
func (pg PreferenceGroup) computeDims(allocated int) []int {
	if len(pg) == 0 {
		return nil
	}

	// cache some useful things
	hasMin := make(map[int]struct{})
	totalMin := 0
	hasMax := make(map[int]struct{})
	hasPref := make(map[int]struct{})
	// totalToPref is the addition needed to get from Min to Preferred.
	totalToPref := 0
	// totalToMax is the addition needed to get from Min/Preferred to Max
	totalToMax := 0
	hasGrow := make(map[int]struct{})
	noGrowNoMax := make(map[int]struct{})

	for idx, p := range pg {
		if p.Min != 0 {
			totalMin += p.Min
			hasMin[idx] = struct{}{}
		}
		if p.Max != 0 {
			hasMax[idx] = struct{}{}
			if p.Preferred != 0 {
				// if there is no preference, the max is already included in 'totalToPref'
				totalToMax += p.Max - p.Preferred
			}
		}
		if p.Preferred != 0 || (p.Preferred == 0 && p.Max != 0) {
			hasPref[idx] = struct{}{}
			totalToPref += max(p.Preferred, p.Max) - p.Min
		}
		if p.Grow {
			hasGrow[idx] = struct{}{}
		}
		if p.Max == 0 && !p.Grow {
			noGrowNoMax[idx] = struct{}{}
		}
	}

	dims := make([]int, len(pg))
	remainder := allocated
	// number of spots left to compute
	numToCompute := len(pg)

	// Pass 1: allocate minimums, exit early if not enough space.
	// don't range over min to avoid nondeterminism
	if len(hasMin) > 0 {
		for idx, p := range pg {
			if p.Min != 0 {
				sz := min(p.Min, remainder)
				dims[idx] = sz
				remainder -= sz
				if remainder == 0 {
					return dims
				}
			}
		}
	}

	// pass 2: even split, stopping at preferred (or max).
	growToPrefferred := func() {
		evenSplit := remainder / numToCompute
		// reallocate a larger even split if totalToPref is reached.
		if evenSplit*len(hasPref) >= totalToPref {
			if numToCompute == len(hasPref) {
				// everyone gets their preference. Set a large evenSplit to avoid a divide by zero.
				evenSplit = 100000
			} else {
				evenSplit = (remainder - totalToPref) / (numToCompute - len(hasPref))
			}
		}
		for idx, p := range pg {
			var sz int
			if _, ok := hasPref[idx]; ok {
				if p.Preferred != 0 {
					// stop at preferred if needed
					sz = min(p.Preferred-dims[idx], evenSplit)
					if sz+dims[idx] >= p.Preferred {
						delete(hasPref, idx)
					}
				} else if p.Max != 0 {
					// stop at max if there is no preference
					sz = min(p.Max-dims[idx], evenSplit)
					if sz+dims[idx] >= p.Max {
						// If max is reached, that item is done. Remove it from numToCompute.
						numToCompute--
						// it is done, remove from lists for later
						delete(hasGrow, idx)
						delete(hasMax, idx)
						delete(hasPref, idx)
					}
				}
				dims[idx] += sz
				remainder -= sz
			}
		}
	}
	for len(hasPref) > 0 && remainder > 0 {
		growToPrefferred()
	}

	// Check if we're done. Either there is no more space, or nothing left to grow.
	if remainder == 0 || numToCompute == 0 {
		return dims
	}

	// pass 3: even split amongst growers OR non-growers with no max.
	growToMax := func() []int {
		var set map[int]struct{}
		var evenSplit int
		if len(hasGrow) > 0 {
			evenSplit = remainder / len(hasGrow)
			set = hasGrow
		}
		// reallocate a larger even split if totalToMax is reached.
		if len(hasGrow) == 0 && evenSplit*len(hasMax) >= totalToMax {
			if len(noGrowNoMax) == 0 {
				evenSplit = 100000
			} else {
				evenSplit = (remainder - totalToMax) / len(noGrowNoMax)
				set = noGrowNoMax
			}
		}
		remainderList := make([]int, 0, len(set))
		for idx := range set {
			var sz int
			if pg[idx].Max != 0 {
				sz = min(pg[idx].Max-dims[idx], evenSplit)
				if sz+dims[idx] >= pg[idx].Max {
					// it is done
					delete(hasGrow, idx)
					delete(hasMax, idx)
				}
			} else {
				sz = evenSplit
			}
			dims[idx] += sz
			remainder -= sz
			if pg[idx].Max != dims[idx] {
				remainderList = append(remainderList, idx)
			}
		}
		return remainderList
	}

	remainderList := growToMax()
	for len(remainderList) > 0 && len(remainderList) < remainder {
		remainderList = growToMax()
	}

	// allocate the remainder if any
	if remainder > 0 {
		for idx := range remainderList {
			dims[idx] += 1
			remainder -= 1
			if remainder == 0 {
				break
			}
		}
	}

	return dims
}

// BoundSize is a size that optionally has a lower and/or upper bound and consists of one to three Unit Values.
// Practically it is a minimum/preferred/maximum size combination but none of the sizes are actually mandatory.
// If a size is missing (e.g. the preferred) it is null and will be replaced by the most appropriate value.
type BoundSize struct {
	Min       int
	Preferred int
	Max       int
	Grow      bool
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
	// PreferredWidth overrides the Preferred width that should be allocated for the view.
	PreferredWidth int
	// MaxWidth overrides the maximum width that should be allocated for the view.
	MaxWidth int

	// MinHeight overrides the minimum height that should be allocated for the view.
	MinHeight int
	// PreferredHeight overrides the Preferred height that should be allocated for the view.
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

	// Preferred overrides the Preferred width or height that should be allocated for the view.
	Preferred int

	// Max overrides the maximum width or height that should be allocated for the view.
	Max int
}

type BubbleLayout interface {
	MaybeAdd(string) (ID, error)
	Add(string) ID
	Cell(Cell) ID
	Dock(Dock) ID
	Wrap()
	Resize(width, height int) BubbleLayoutMsg
	Validate() error
}

// NewWithConstraints creates a new BubbleLayout with the given size constraints.
func NewWithConstraints(width, height PreferenceGroup) BubbleLayout {
	// TODO: Verify these constraints.
	return &bubbleLayout{
		layouts: [][]layout{{}},
		wPref:   width,
		hPref:   height,
	}
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
	hPref       PreferenceGroup
	wPref       PreferenceGroup
}

// MaybeAdd is like Add but returns an error if the string cannot be parsed into a valid Cell or Dock.
func (bl *bubbleLayout) MaybeAdd(str string) (ID, error) {
	l, err := convertToLayout(str)
	if err != nil {
		return 0, err
	}
	if l.Dock == (Dock{}) {
		return bl.add(l), nil
	} else {
		return bl.Dock(l.Dock), nil
	}
}

// Add uses the string notation to define the layout. This is often shorter and easier to read than using the Layout struct.
// If there is an error Add will panic. This is done for code readability, if you want to handle errors use MaybeAdd.
func (bl *bubbleLayout) Add(str string) ID {
	id, err := bl.MaybeAdd(str)
	if err != nil {
		panic(err)
	}
	return id
}

func (bl *bubbleLayout) add(l layout) ID {
	bl.idCounter++
	l.id = bl.idCounter
	idx := len(bl.layouts) - 1
	bl.layouts[idx] = append(bl.layouts[idx], l)

	if l.wrap {
		bl.layouts = append(bl.layouts, []layout{})
	}

	// TODO: Debug mode which panics here as soon as a constraint violation is detected.
	return bl.idCounter
}

// Cell adds a Cell to the Grid. By default, it is placed in the next available cell going left to right top to bottom.
func (bl *bubbleLayout) Cell(c Cell) ID {
	return bl.add(layout{Cell: c})
}

// Wrap inserts a new row into the layout, subsequent calls to Add will place models in the new row.
func (bl *bubbleLayout) Wrap() {
	bl.layouts = append(bl.layouts, []layout{})
}

// Dock places a model on the edge of the layout, spanning the entire width or height.
// For NORTH and SOUTH components, the width is fixed and the height is defined by Min, Preferred and Max.
// For EAST and WEST components, the height is fixed and the width is defined by Min, Preferred and Max.
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
	return fmt.Sprintf("constraint violation: %s %d: Min %s (%d), Max %s (%d)", dir, p.idx, dim, p.min, dim, p.max)
}

func makeRowViolation(idx, min, max int) error {
	return preferenceConstraintError{row: true, idx: idx, min: min, max: max}
}

func makeColViolation(idx, min, max int) error {
	return preferenceConstraintError{row: false, idx: idx, min: min, max: max}
}

func checkPreferenceConstraints(hPref, wPref PreferenceGroup) error {
	for row, p := range hPref {
		if p.Min > p.Max {
			return makeRowViolation(row, p.Min, p.Max)
		}
	}

	for col, p := range wPref {
		if p.Min > p.Max {
			return makeColViolation(col, p.Min, p.Max)
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
		curRow := 0
		for _, l := range row {
			if l.SpanWidth > 0 {
				curRow += l.SpanWidth
			} else {
				curRow++
			}
		}
	}

	// spanx and create rows
	for colIdx := 0; colIdx < longestCol; colIdx++ {
		var rowIdx int
		for rowIdx = 0; rowIdx < len(ret); rowIdx++ {
			// pad empty cells
			if len(ret[rowIdx]) < (colIdx + 1) {
				ret[rowIdx] = append(ret[rowIdx], layout{})
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
							// pad next row to colIdx so that we can append the new cell.
							ret = append(ret, make([]layout, colIdx))
						}
						if len(ret[rowIdx+i]) == colIdx {
							// special case for new rows
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
				// Grow the longest column if necessary.
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
			// Cell it to the first row, spanning the entire width.
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
			// Cell it to the final row, spanning the entire width.
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
			// Cell it to the end of each row to span the entire height.
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
			// Cell it to the front of each row to span the entire height.
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
// across each row and column.
//
// For example, the minimum width for a column would be the largest minimum
// preference across all cells in the first column.
//
// This function is only used if row and column constraints are not defined.
func distillPreferences(g Grid) (hPref, wPref PreferenceGroup) {
	if len(g) == 0 {
		return
	}

	gridHeight := len(g)
	gridWidth := len(g[0])
	hPref = make(PreferenceGroup, gridHeight)
	wPref = make(PreferenceGroup, gridWidth)
	for rowIdx := 0; rowIdx < gridHeight; rowIdx++ {
		for colIdx := 0; colIdx < gridWidth; colIdx++ {
			l := g[rowIdx][colIdx]

			// collect height preferences
			if l.MinHeight != 0 {
				hPref[rowIdx].Min = max(l.MinHeight, hPref[rowIdx].Min)
			}
			if l.MaxHeight != 0 {
				if hPref[rowIdx].Max == 0 {
					hPref[rowIdx].Max = l.MaxHeight
				} else {
					hPref[rowIdx].Max = min(l.MaxHeight, hPref[rowIdx].Max)
				}
			}
			if l.PreferredHeight != 0 {
				hPref[rowIdx].Preferred = max(hPref[rowIdx].Preferred, l.PreferredHeight)
			}
			hPref[rowIdx].Grow = hPref[rowIdx].Grow || l.GrowHeight

			// collect width preferences
			if l.MinWidth != 0 {
				wPref[colIdx].Min = max(l.MinWidth, wPref[colIdx].Min)
			}
			if l.MaxWidth != 0 {
				if wPref[colIdx].Max == 0 {
					wPref[colIdx].Max = l.MaxWidth
				} else {
					wPref[colIdx].Max = min(l.MaxWidth, wPref[colIdx].Max)
				}
			}
			if l.PreferredWidth != 0 {
				wPref[colIdx].Preferred = max(wPref[colIdx].Preferred, l.PreferredWidth)
			}
			wPref[colIdx].Grow = wPref[colIdx].Grow || l.GrowWidth
		}
	}
	return
}

func (bl *bubbleLayout) Validate() error {
	if len(bl.resizeCache) == 0 {
		bl.resizeCache = expandSpans(bl.layouts)
		bl.resizeCache = mergeDocks(bl.resizeCache, bl.docks)

		hPref, wPref := distillPreferences(bl.resizeCache)

		// If the user provided constraints are shorter than the auto generated ones, append the distilled ones.
		// TODO: in the future, cell width/height make this more complicated.
		// 		 These distilled preferences would need to be merged with the user provided ones.
		// TODO: this flexibility may not be needed. Should constraints be more strict?
		appendPref := func(user, distilled PreferenceGroup) PreferenceGroup {
			if len(user) < len(distilled) {
				return append(user, distilled[len(user):]...)
			}
			return user
		}
		bl.hPref = appendPref(bl.hPref, hPref)
		bl.wPref = appendPref(bl.wPref, wPref)

		if len(bl.hPref) != len(bl.resizeCache) {
			return fmt.Errorf("height preferences do not match the cell height")
		}

		if len(bl.resizeCache) > 0 && len(bl.wPref) != len(bl.resizeCache[0]) {
			return fmt.Errorf("width preferences do not match the cell height")
		}

		return checkPreferenceConstraints(bl.hPref, bl.wPref)
	}
	return nil
}

// Resize recalculates the layout based on the current terminal size.
// This function will panic if there is a validation error. If you would like to
// handle errors, use Validate() before calling Resize().
func (bl *bubbleLayout) Resize(width, height int) BubbleLayoutMsg {
	if err := bl.Validate(); err != nil {
		panic(err)
	}

	hDims := bl.hPref.computeDims(height)
	wDims := bl.wPref.computeDims(width)

	return bl.resizeCache.makeMessage(wDims, hDims)
}

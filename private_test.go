package bubblelayout

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPreferenceConstraints(t *testing.T) {
	testcases := []struct {
		name string
		row  []BoundSize
		col  []BoundSize
		err  error
	}{
		{
			name: "no violations 1",
			row:  nil,
			col:  nil,
			err:  nil,
		}, {
			name: "no violations 2",
			row: []BoundSize{
				{Min: 1, Preferred: 1, Max: 10},
				{Min: 5, Preferred: 1, Max: 10},
			},
			col: []BoundSize{
				{Min: 1, Preferred: 1, Max: 10},
				{Min: 5, Preferred: 1, Max: 10},
			},
			err: nil,
		}, {
			name: "ignore zero value max",
			row: []BoundSize{
				{Min: 1, Preferred: 1, Max: 0},
			},
			col: []BoundSize{
				{Min: 1, Preferred: 1, Max: 0},
			},
			err: nil,
		}, {
			name: "row violation",
			row: []BoundSize{
				{Min: 10, Preferred: 0, Max: 1},
			},
			col: []BoundSize{},
			err: makeRowViolation(0, 10, 0, 1),
		}, {
			name: "col violation",
			row:  []BoundSize{},
			col: []BoundSize{
				{Min: 10, Preferred: 0, Max: 1},
			},
			err: makeColViolation(0, 10, 0, 1),
		}, {
			name: "violation index",
			row:  []BoundSize{},
			col: []BoundSize{
				{}, {}, {}, {}, {Min: 10, Preferred: 0, Max: 1},
			},
			err: makeColViolation(4, 10, 0, 1),
		}, {
			name: "preferred row violation 1",
			row: []BoundSize{
				{Min: 10, Preferred: 1, Max: 0},
			},
			col: []BoundSize{},
			err: makeRowViolation(0, 10, 1, 0),
		}, {
			name: "preferred row violation 2",
			row: []BoundSize{
				{Min: 1, Preferred: 3, Max: 2},
			},
			col: []BoundSize{},
			err: makeRowViolation(0, 1, 3, 2),
		}, {
			name: "preferred col violation 1",
			row:  []BoundSize{},
			col: []BoundSize{
				{Min: 10, Preferred: 1, Max: 0},
			},
			err: makeColViolation(0, 10, 1, 0),
		}, {
			name: "preferred col violation 2",
			row:  []BoundSize{},
			col: []BoundSize{
				{Min: 1, Preferred: 3, Max: 2},
			},
			err: makeColViolation(0, 1, 3, 2),
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := checkPreferenceConstraints(tc.row, tc.col)
			assert.Equal(t, tc.err, err)
		})
	}

}

func TestPreferenceGroup(t *testing.T) {
	testcases := []struct {
		name      string
		pg        PreferenceGroup
		allocated int
		expected  []int
	}{
		{
			name:      "empty",
			pg:        PreferenceGroup{},
			allocated: 80,
			expected:  nil,
		}, {
			name: "one",
			pg: PreferenceGroup{
				{Min: 1, Preferred: 5, Max: 10},
			},
			allocated: 80,
			expected:  []int{5},
		}, {
			name: "one Grow to Max",
			pg: PreferenceGroup{
				{Min: 1, Preferred: 5, Max: 10, Grow: true},
			},
			allocated: 80,
			expected:  []int{10},
		}, {
			name: "one Grow to allocated size",
			pg: PreferenceGroup{
				{Min: 1, Preferred: 1, Max: 1000, Grow: true},
			},
			allocated: 80,
			expected:  []int{80},
		}, {
			name: "two even empty Grow",
			pg: PreferenceGroup{
				{},
				{},
			},
			allocated: 80,
			expected:  []int{40, 40},
		}, {
			name: "two uneven to Max",
			pg: PreferenceGroup{
				{Max: 10},
				{Max: 70},
			},
			allocated: 80,
			expected:  []int{10, 70},
		}, {
			name: "two uneven Grow and Max",
			pg: PreferenceGroup{
				{Max: 10},
				{},
			},
			allocated: 80,
			expected:  []int{10, 70},
		}, {
			name: "two growers",
			pg: PreferenceGroup{
				{Preferred: 10, Max: 40, Grow: true},
				{Grow: true},
			},
			allocated: 80,
			expected:  []int{40, 40},
		}, {
			name: "three uneven to Max",
			pg: PreferenceGroup{
				{Max: 10},
				{Max: 20},
				{Max: 30},
			},
			allocated: 80,
			expected:  []int{10, 20, 30},
		}, {
			name: "three uneven preference and Max",
			pg: PreferenceGroup{
				{Max: 10},
				{Preferred: 15, Max: 20},
				{Max: 30},
			},
			allocated: 80,
			expected:  []int{10, 15, 30},
		}, {
			name: "three uneven over allocated",
			pg: PreferenceGroup{
				{Max: 10},
				{Preferred: 95},
				{Max: 30},
			},
			allocated: 80,
			expected:  []int{10, 40, 30},
		}, {
			name: "two growers and uneven max",
			pg: PreferenceGroup{
				{Preferred: 10, Max: 20, Grow: true},
				{Preferred: 10, Max: 80, Grow: true},
			},
			allocated: 60,
			expected:  []int{20, 40},
		}, {
			name: "two growers and uneven max thats reached",
			pg: PreferenceGroup{
				{Preferred: 10, Max: 15, Grow: true},
				{Preferred: 10, Max: 15, Grow: true},
			},
			allocated: 60,
			expected:  []int{15, 15},
		}, {
			name: "just enough for Min",
			pg: PreferenceGroup{
				{Min: 20, Preferred: 30, Max: 40},
				{Min: 20, Preferred: 30, Max: 40},
				{Min: 20, Preferred: 30, Max: 40},
				{Min: 20, Preferred: 30, Max: 40},
			},
			allocated: 80,
			expected:  []int{20, 20, 20, 20},
		}, {
			name: "just enough for Min 2",
			pg: PreferenceGroup{
				{Min: 10, Preferred: 30, Max: 40},
				{Min: 20, Preferred: 30, Max: 40},
				{Min: 30, Preferred: 60, Max: 100},
				{Min: 40, Preferred: 80, Max: 100},
			},
			allocated: 100,
			expected:  []int{10, 20, 30, 40},
		}, {
			name: "go below Min when over allocated",
			pg: PreferenceGroup{
				{Min: 25, Preferred: 30, Max: 40},
				{Min: 25, Preferred: 30, Max: 40},
				{Min: 25, Preferred: 30, Max: 40},
				{Min: 25, Preferred: 30, Max: 40},
			},
			allocated: 80,
			expected:  []int{25, 25, 25, 5},
		}, {
			name: "do not go over max",
			pg: PreferenceGroup{
				{Max: 30, Grow: true},
				{Max: 30, Grow: true},
			},
			allocated: 61,
			expected:  []int{30, 30},
		}, {
			name: "remainder",
			pg: PreferenceGroup{
				{},
				{},
			},
			allocated: 81,
			expected:  []int{41, 40},
		}, {
			name: "even split above preferred",
			pg: PreferenceGroup{
				{Preferred: 10, Max: 100, Grow: true},
				{Preferred: 50, Max: 70, Grow: true},
			},
			allocated: 80,
			expected:  []int{20, 60},
		}, {
			name: "No max no grow",
			pg: PreferenceGroup{
				{Preferred: 10, Max: 15},
				{Preferred: 10},
				{},
			},
			allocated: 60,
			expected:  []int{10, 10, 40},
		}, {
			name: "growToMax",
			pg: PreferenceGroup{
				{Preferred: 10, Max: 100, Grow: true},
				{Preferred: 50, Max: 55, Grow: true},
			},
			allocated: 80,
			expected:  []int{25, 55},
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := tc.pg.computeDims(tc.allocated)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestExpandSpans(t *testing.T) {
	testcases := []struct {
		name        string
		inputLayout Grid
		inputString func(*bubbleLayout) Grid
		expected    Grid
	}{
		{
			name: "no expand",
			inputLayout: [][]layout{
				{{id: 1}, {id: 2, wrap: true}},
				{{id: 3}, {id: 4}},
			},
			inputString: func(bl *bubbleLayout) Grid {
				bl.Add("")
				bl.Add("wrap")
				bl.Add("")
				bl.Add("")
				return bl.layouts
			},
			expected: [][]layout{
				{{id: 1}, {id: 2, wrap: true}},
				{{id: 3}, {id: 4}},
			},
		}, {
			name: "horizontal span",
			inputLayout: [][]layout{
				{{id: 1, Cell: Cell{SpanWidth: 2}, wrap: true}},
				{{id: 2}, {id: 3}},
			},
			inputString: func(bl *bubbleLayout) Grid {
				bl.Add("spanw 2, wrap")
				bl.Add("")
				bl.Add("")
				return bl.layouts
			},
			expected: [][]layout{
				{{id: 1, Cell: Cell{SpanWidth: 2}, wrap: true}, {id: 1, Cell: Cell{SpanWidth: 2, wDuplicate: true}, wrap: true}},
				{{id: 2}, {id: 3}},
			},
		}, {
			name: "horizontal span and expand, with padding",
			inputLayout: [][]layout{
				{{id: 1}, {id: 2, Cell: Cell{SpanWidth: 3}, wrap: true}},
				{{id: 3}, {id: 4}},
			},
			inputString: func(bl *bubbleLayout) Grid {
				bl.Add("")
				bl.Add("spanx 3, wrap")
				bl.Add("")
				bl.Add("")
				return bl.layouts
			},
			expected: [][]layout{
				{{id: 1}, {id: 2, Cell: Cell{SpanWidth: 3}, wrap: true}, {id: 2, Cell: Cell{SpanWidth: 3, wDuplicate: true}, wrap: true}, {id: 2, Cell: Cell{SpanWidth: 3, wDuplicate: true}, wrap: true}},
				{{id: 3}, {id: 4}, {id: 0}, {id: 0}},
			},
		}, {
			name: "vertical span",
			inputLayout: [][]layout{
				{{id: 1, Cell: Cell{SpanHeight: 2}}, {id: 2, wrap: true}},
				{{id: 3}},
			},
			inputString: func(bl *bubbleLayout) Grid {
				bl.Add("spanh 2")
				bl.Add("wrap")
				bl.Add("")
				return bl.layouts
			},
			expected: [][]layout{
				{{id: 1, Cell: Cell{SpanHeight: 2}}, {id: 2, wrap: true}},
				{{id: 1, Cell: Cell{SpanHeight: 2, hDuplicate: true}}, {id: 3}},
			},
		}, {
			name: "vertical span and expand",
			inputLayout: [][]layout{
				{{id: 1, Cell: Cell{SpanHeight: 2}}},
			},
			inputString: func(bl *bubbleLayout) Grid {
				bl.Add("spanh 2")
				return bl.layouts
			},
			expected: [][]layout{
				{{id: 1, Cell: Cell{SpanHeight: 2}}},
				{{id: 1, Cell: Cell{SpanHeight: 2, hDuplicate: true}}},
			},
		}, {
			name: "test empty space",
			inputLayout: [][]layout{
				{{id: 1, wrap: true}},
				{{id: 2}, {id: 3}},
			},
			inputString: func(bl *bubbleLayout) Grid {
				bl.Add("wrap")
				bl.Add("")
				bl.Add("")
				return bl.layouts
			},
			expected: [][]layout{
				{{id: 1, wrap: true}, {id: 0}},
				{{id: 2}, {id: 3}},
			},
		}, {
			name: "empty space 2d expand at horizontal end",
			inputLayout: [][]layout{
				{{id: 1}, {id: 2, Cell: Cell{SpanHeight: 2, SpanWidth: 2}}},
			},
			inputString: func(bl *bubbleLayout) Grid {
				bl.Add("")
				bl.Add("span 2 2")
				return bl.layouts
			},
			expected: [][]layout{
				{{id: 1}, {id: 2, Cell: Cell{SpanHeight: 2, SpanWidth: 2}}, {id: 2, Cell: Cell{SpanHeight: 2, SpanWidth: 2, wDuplicate: true}}},
				{{id: 0}, {id: 2, Cell: Cell{SpanHeight: 2, SpanWidth: 2, hDuplicate: true}}, {id: 2, Cell: Cell{SpanHeight: 2, SpanWidth: 2, hDuplicate: true, wDuplicate: true}}},
			},
		}, {
			name: "empty space 2d expand at vertical end",
			inputLayout: [][]layout{
				{{id: 1, wrap: true}},
				{{id: 2, Cell: Cell{SpanHeight: 2, SpanWidth: 2}}},
			},
			inputString: func(bl *bubbleLayout) Grid {
				bl.Add("wrap")
				bl.Add("span 2 2")
				return bl.layouts
			},
			expected: [][]layout{
				{{id: 1, wrap: true}, {id: 0}},
				{{id: 2, Cell: Cell{SpanHeight: 2, SpanWidth: 2}}, {id: 2, Cell: Cell{SpanHeight: 2, SpanWidth: 2, wDuplicate: true}}},
				{{id: 2, Cell: Cell{SpanHeight: 2, SpanWidth: 2, hDuplicate: true}}, {id: 2, Cell: Cell{SpanHeight: 2, SpanWidth: 2, hDuplicate: true, wDuplicate: true}}},
			},
		}, {
			name: "vertical and horizontal span and overflow",
			inputLayout: [][]layout{
				{{id: 1, Cell: Cell{SpanWidth: 2, SpanHeight: 2}}},
			},
			inputString: func(bl *bubbleLayout) Grid {
				bl.Add("span 2 2")
				return bl.layouts
			},
			expected: [][]layout{
				{{id: 1, Cell: Cell{SpanWidth: 2, SpanHeight: 2}}, {id: 1, Cell: Cell{SpanWidth: 2, SpanHeight: 2, wDuplicate: true}}},
				{{id: 1, Cell: Cell{SpanWidth: 2, SpanHeight: 2, hDuplicate: true}}, {id: 1, Cell: Cell{SpanWidth: 2, SpanHeight: 2, hDuplicate: true, wDuplicate: true}}},
			},
		}, {
			// inputLayout:
			// ---------------------------------
			// | 1        | 2 (2, 2) |    3    |
			// ---------------------------------
			// | 4 (1, 2) |    5     |
			// -----------------------
			// |    6     | 7 (1, 2) |
			// -----------------------
			// expected
			// ---------------------------------
			// |   1   |       -       |   3   |
			// --------- -  -  2  -  - |--------
			// |   -   |       -       |   5   |
			// | - 4 - -------------------------
			// |   -   |   6   |       7       |
			// ---------------------------------
			name: "complex",
			inputLayout: [][]layout{
				{{id: 1}, {id: 2, Cell: Cell{SpanWidth: 2, SpanHeight: 2}}, {id: 3, wrap: true}},
				{{id: 4, Cell: Cell{SpanHeight: 2}}, {id: 5, wrap: true}},
				{{id: 6}, {id: 7, Cell: Cell{SpanWidth: 2}}},
			},
			inputString: func(bl *bubbleLayout) Grid {
				bl.Add("")
				bl.Add("span 2 2")
				bl.Add("wrap")
				bl.Add("spanh 2")
				bl.Add("wrap")
				bl.Add("")
				bl.Add("spanw 2")
				return bl.layouts
			},
			expected: [][]layout{
				{{id: 1}, {id: 2, Cell: Cell{SpanWidth: 2, SpanHeight: 2}}, {id: 2, Cell: Cell{SpanWidth: 2, SpanHeight: 2, wDuplicate: true}}, {id: 3, wrap: true}},
				{{id: 4, Cell: Cell{SpanHeight: 2}}, {id: 2, Cell: Cell{SpanWidth: 2, SpanHeight: 2, hDuplicate: true}}, {id: 2, Cell: Cell{SpanWidth: 2, SpanHeight: 2, hDuplicate: true, wDuplicate: true}}, {id: 5, wrap: true}},
				{{id: 4, Cell: Cell{SpanHeight: 2, hDuplicate: true}}, {id: 6}, {id: 7, Cell: Cell{SpanWidth: 2}}, {id: 7, Cell: Cell{SpanWidth: 2, wDuplicate: true}}},
			},
		}, {
			name: "prefs even split",
			inputLayout: [][]layout{
				{{id: 1, Cell: Cell{SpanWidth: 2, SpanHeight: 2, MinHeight: 10, MaxHeight: 100, PreferredHeight: 50, MinWidth: 10, MaxWidth: 100, PreferredWidth: 50}}},
			},
			inputString: func(bl *bubbleLayout) Grid {
				bl.Add("span 2 2, width 10:50:100, height 10:50:100")
				return bl.layouts
			},
			expected: [][]layout{
				{{id: 1, Cell: Cell{SpanWidth: 2, SpanHeight: 2, MinHeight: 5, MaxHeight: 50, PreferredHeight: 25, MinWidth: 5, MaxWidth: 50, PreferredWidth: 25}},
					{id: 1, Cell: Cell{SpanWidth: 2, SpanHeight: 2, MinHeight: 5, MaxHeight: 50, PreferredHeight: 25, MinWidth: 5, MaxWidth: 50, PreferredWidth: 25, wDuplicate: true}}},
				{{id: 1, Cell: Cell{SpanWidth: 2, SpanHeight: 2, MinHeight: 5, MaxHeight: 50, PreferredHeight: 25, MinWidth: 5, MaxWidth: 50, PreferredWidth: 25, hDuplicate: true}},
					{id: 1, Cell: Cell{SpanWidth: 2, SpanHeight: 2, MinHeight: 5, MaxHeight: 50, PreferredHeight: 25, MinWidth: 5, MaxWidth: 50, PreferredWidth: 25, hDuplicate: true, wDuplicate: true}},
				},
			},
		}, {
			// these are all 1 larger than "perfs even split", the odd number is discarded.
			name: "prefs odd split rounding error",
			inputLayout: [][]layout{
				{{id: 1, Cell: Cell{SpanWidth: 2, SpanHeight: 2, MinHeight: 11, MaxHeight: 101, PreferredHeight: 51, MinWidth: 11, MaxWidth: 101, PreferredWidth: 51}}},
			},
			inputString: func(bl *bubbleLayout) Grid {
				bl.Add("span 2 2, width 11:51:101, height 11:51:101")
				return bl.layouts
			},
			expected: [][]layout{
				{{id: 1, Cell: Cell{SpanWidth: 2, SpanHeight: 2, MinHeight: 5, MaxHeight: 50, PreferredHeight: 25, MinWidth: 5, MaxWidth: 50, PreferredWidth: 25}},
					{id: 1, Cell: Cell{SpanWidth: 2, SpanHeight: 2, MinHeight: 5, MaxHeight: 50, PreferredHeight: 25, MinWidth: 5, MaxWidth: 50, PreferredWidth: 25, wDuplicate: true}}},
				{{id: 1, Cell: Cell{SpanWidth: 2, SpanHeight: 2, MinHeight: 5, MaxHeight: 50, PreferredHeight: 25, MinWidth: 5, MaxWidth: 50, PreferredWidth: 25, hDuplicate: true}},
					{id: 1, Cell: Cell{SpanWidth: 2, SpanHeight: 2, MinHeight: 5, MaxHeight: 50, PreferredHeight: 25, MinWidth: 5, MaxWidth: 50, PreferredWidth: 25, hDuplicate: true, wDuplicate: true}},
				},
			},
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result1 := expandSpans(tc.inputLayout)
			assert.Len(t, result1, len(tc.expected), "number of rows mismatch")
			for i, row := range result1 {
				assert.Len(t, row, len(tc.expected[i]), "row %d: number of columns mismatch", i)
			}
			assert.Equal(t, tc.expected, result1)

			if tc.inputString != nil {
				bl := New().(*bubbleLayout)

				grid := tc.inputString(bl)
				assert.Equal(t, tc.inputLayout, grid)
			}
		})
	}
}

func TestMergeDocks_Empty(t *testing.T) {
	assert.Nil(t, mergeDocks(nil, nil))
}

func TestMergeDocks_InvalidCardinal(t *testing.T) {
	panicFunc := func() {
		mergeDocks([][]layout{{{id: 0}}}, []layout{{id: 1, Dock: Dock{Cardinal: "north-north-west"}}})
	}
	assert.PanicsWithError(t, "invalid cardinal", panicFunc)
}

func TestMergeDocks(t *testing.T) {
	testcases := []struct {
		name     string
		start    Grid
		docks    []layout
		expected Grid
	}{
		{
			name: "simple north",
			start: [][]layout{
				{{id: 0}},
			},
			docks: []layout{
				{id: 1, Dock: Dock{Cardinal: NORTH, Min: 10, Preferred: 10, Max: 10}},
			},
			expected: [][]layout{
				{{id: 1, Cell: Cell{SpanWidth: 1, MinHeight: 10, PreferredHeight: 10, MaxHeight: 10}}},
				{{id: 0}},
			},
		}, {
			name: "simple south",
			start: [][]layout{
				{{id: 0}},
			},
			docks: []layout{
				{id: 1, Dock: Dock{Cardinal: SOUTH, Min: 10, Preferred: 10, Max: 10}},
			},
			expected: [][]layout{
				{{id: 0}},
				{{id: 1, Cell: Cell{SpanWidth: 1, MinHeight: 10, PreferredHeight: 10, MaxHeight: 10}}},
			},
		}, {
			name: "simple west",
			start: [][]layout{
				{{id: 0}},
			},
			docks: []layout{
				{id: 1, Dock: Dock{Cardinal: WEST, Min: 10, Preferred: 10, Max: 10}},
			},
			expected: [][]layout{
				{{id: 1, Cell: Cell{SpanHeight: 1, MinWidth: 10, PreferredWidth: 10, MaxWidth: 10}}, {id: 0}},
			},
		}, {
			name: "simple east",
			start: [][]layout{
				{{id: 0}},
			},
			docks: []layout{
				{id: 1, Dock: Dock{Cardinal: EAST, Min: 10, Preferred: 10, Max: 10}},
			},
			expected: [][]layout{
				{{id: 0}, {id: 1, Cell: Cell{SpanHeight: 1, MinWidth: 10, PreferredWidth: 10, MaxWidth: 10}}},
			},
		}, {
			name: "double north",
			start: [][]layout{
				{{id: 0}, {id: 1}},
			},
			docks: []layout{
				{id: 2, Dock: Dock{Cardinal: NORTH, Min: 10, Preferred: 10, Max: 10}},
			},
			expected: [][]layout{
				{{id: 2, Cell: Cell{SpanWidth: 2, MinHeight: 10, PreferredHeight: 10, MaxHeight: 10}}, {id: 2, Cell: Cell{SpanWidth: 2, MinHeight: 10, PreferredHeight: 10, MaxHeight: 10, wDuplicate: true}}},
				{{id: 0}, {id: 1}},
			},
		}, {
			name: "double south",
			start: [][]layout{
				{{id: 0}, {id: 1}},
			},
			docks: []layout{
				{id: 2, Dock: Dock{Cardinal: SOUTH, Min: 10, Preferred: 10, Max: 10}},
			},
			expected: [][]layout{
				{{id: 0}, {id: 1}},
				{{id: 2, Cell: Cell{SpanWidth: 2, MinHeight: 10, PreferredHeight: 10, MaxHeight: 10}}, {id: 2, Cell: Cell{SpanWidth: 2, MinHeight: 10, PreferredHeight: 10, MaxHeight: 10, wDuplicate: true}}},
			},
		}, {
			name: "double west",
			start: [][]layout{
				{{id: 0}},
				{{id: 1}},
			},
			docks: []layout{
				{id: 2, Dock: Dock{Cardinal: WEST, Min: 10, Preferred: 10, Max: 10}},
			},
			expected: [][]layout{
				{{id: 2, Cell: Cell{SpanHeight: 2, MinWidth: 10, PreferredWidth: 10, MaxWidth: 10}}, {id: 0}},
				{{id: 2, Cell: Cell{SpanHeight: 2, MinWidth: 10, PreferredWidth: 10, MaxWidth: 10, hDuplicate: true}}, {id: 1}},
			},
		}, {
			name: "double east",
			start: [][]layout{
				{{id: 0}},
				{{id: 1}},
			},
			docks: []layout{
				{id: 2, Dock: Dock{Cardinal: EAST, Min: 10, Preferred: 10, Max: 10}},
			},
			expected: [][]layout{
				{{id: 0}, {id: 2, Cell: Cell{SpanHeight: 2, MinWidth: 10, PreferredWidth: 10, MaxWidth: 10}}},
				{{id: 1}, {id: 2, Cell: Cell{SpanHeight: 2, MinWidth: 10, PreferredWidth: 10, MaxWidth: 10, hDuplicate: true}}},
			},
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := mergeDocks(tc.start, tc.docks)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestDistillPrefs(t *testing.T) {
	testcases := []struct {
		name      string
		input     Grid
		hExpected PreferenceGroup
		wExpected PreferenceGroup
	}{
		{
			name: "1x1",
			input: [][]layout{
				{{id: 1, Cell: Cell{MinHeight: 10, MaxHeight: 100, PreferredHeight: 50, MinWidth: 10, MaxWidth: 100, PreferredWidth: 50}}},
			},
			hExpected: []BoundSize{{Min: 10, Preferred: 50, Max: 100}},
			wExpected: []BoundSize{{Min: 10, Preferred: 50, Max: 100}},
		}, {
			name: "2x2 one defaults",
			input: [][]layout{
				{{id: 1, Cell: Cell{MinHeight: 10, MaxHeight: 100, PreferredHeight: 50, MinWidth: 10, MaxWidth: 100, PreferredWidth: 50}}, {id: 2}},
				{{id: 3}, {id: 4, Cell: Cell{MinHeight: 10, MaxHeight: 100, PreferredHeight: 50, MinWidth: 10, MaxWidth: 100, PreferredWidth: 50}}},
			},
			hExpected: []BoundSize{{Min: 10, Preferred: 50, Max: 100}, {Min: 10, Preferred: 50, Max: 100}},
			wExpected: []BoundSize{{Min: 10, Preferred: 50, Max: 100}, {Min: 10, Preferred: 50, Max: 100}},
		}, {
			name: "1x3 Max(Min)) and Min(Max)",
			input: [][]layout{
				{
					{id: 1, Cell: Cell{MinHeight: 5, MaxHeight: 50, PreferredHeight: 25}},
					{id: 2, Cell: Cell{MinHeight: 10, MaxHeight: 100, PreferredHeight: 50}},
					{id: 3},
				},
			},
			hExpected: []BoundSize{{Min: 10, Preferred: 50, Max: 50}},
			wExpected: []BoundSize{{}, {}, {}},
		}, {
			name: "3x1 Max(Min)) and Min(Max)",
			input: [][]layout{
				{{id: 1, Cell: Cell{MinWidth: 5, MaxWidth: 50, PreferredWidth: 25}}},
				{{id: 2, Cell: Cell{MinWidth: 10, MaxWidth: 100, PreferredWidth: 25}}},
				{{id: 3}},
			},
			hExpected: []BoundSize{{}, {}, {}},
			wExpected: []BoundSize{{Min: 10, Preferred: 25, Max: 50}},
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			hPref, wPref := distillPreferences(tc.input)
			assert.Equal(t, tc.hExpected, hPref)
			assert.Equal(t, tc.wExpected, wPref)
		})
	}
}

func TestDistillPrefs_Empty(t *testing.T) {
	hPref, wPref := distillPreferences(nil)
	assert.Nil(t, hPref)
	assert.Nil(t, wPref)
}

func TestDock(t *testing.T) {
	l := New()
	l.Dock(Dock{NORTH, 1, 2, 3})
	bl := l.(*bubbleLayout)
	require.Equal(t, []layout{{id: 1, Dock: Dock{Cardinal: NORTH, Min: 1, Preferred: 2, Max: 3}}}, bl.docks)
}

func TestValidate_FailureHeight(t *testing.T) {
	l := New()
	l.Cell(Cell{MinHeight: 100})
	l.Cell(Cell{MaxHeight: 10})
	require.ErrorContains(t, l.Validate(), makeRowViolation(0, 100, 0, 10).Error())
}

func TestValidate_FailureWidth(t *testing.T) {
	l := New()
	l.Cell(Cell{MinWidth: 100})
	l.Wrap()
	l.Cell(Cell{MaxWidth: 10})
	require.ErrorContains(t, l.Validate(), makeColViolation(0, 100, 0, 10).Error())
}

func TestTooManayConstraints(t *testing.T) {
	// 0x1 layout with 1x1 constraint
	h := PreferenceGroup{{}}
	w := PreferenceGroup{{}}
	l := NewWithConstraints(w, h)
	require.ErrorContains(t, l.Validate(), "width preferences do not match the cell height")

	// 0x1 layout with 1x2 constraint
	h = append(h, BoundSize{})
	l = NewWithConstraints(w, h)
	require.ErrorContains(t, l.Validate(), "height preferences do not match the cell height")
}

func TestProvideConstraints(t *testing.T) {
	col := PreferenceGroup{{Min: 1, Preferred: 2, Max: 3}}
	row := PreferenceGroup{{Min: 4, Preferred: 5, Max: 6}}
	l := NewWithConstraints(row, col)
	l.Add("")
	require.NoError(t, l.Validate())
	bl := l.(*bubbleLayout)
	require.Equal(t, col, bl.hPref)
	require.Equal(t, row, bl.wPref)
}

func TestConstraintExtension(t *testing.T) {
	col := PreferenceGroup{{Min: 1, Preferred: 2, Max: 3}}
	row := PreferenceGroup{{Min: 4, Preferred: 5, Max: 6}}
	l := NewWithConstraints(row, col)
	l.Add("")
	l.Add("grow, wrap")
	l.Add("")
	l.Add("grow")
	require.NoError(t, l.Validate())
	bl := l.(*bubbleLayout)

	// A "grow" bound from the distilled constraints should be added to each.
	addedBound := BoundSize{Grow: true}
	require.Equal(t, append(col, addedBound), bl.hPref)
	require.Equal(t, append(row, addedBound), bl.wPref)
}

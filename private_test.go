package layout

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPreferenceConstraints(t *testing.T) {
	testcases := []struct {
		name string
		row  []sizePreference
		col  []sizePreference
		err  error
	}{
		{
			name: "no violations 1",
			row:  nil,
			col:  nil,
			err:  nil,
		}, {
			name: "no violations 2",
			row: []sizePreference{
				{min: 1, preferred: []int{1}, max: 10},
				{min: 5, preferred: []int{1}, max: 10},
			},
			col: []sizePreference{
				{min: 1, preferred: []int{1}, max: 10},
				{min: 5, preferred: []int{1}, max: 10},
			},
			err: nil,
		}, {
			name: "row violation",
			row: []sizePreference{
				{min: 10, preferred: []int{1}, max: 1},
			},
			col: []sizePreference{},
			err: makeRowViolation(0, 10, 1),
		}, {
			name: "col violation",
			row:  []sizePreference{},
			col: []sizePreference{
				{min: 10, preferred: []int{1}, max: 1},
			},
			err: makeColViolation(0, 10, 1),
		}, {
			name: "violation index",
			row:  []sizePreference{},
			col: []sizePreference{
				{}, {}, {}, {}, {min: 10, preferred: []int{1}, max: 1},
			},
			err: makeColViolation(4, 10, 1),
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
		pg        preferenceGroup
		allocated int
		expected  []int
	}{
		{
			name:      "empty",
			pg:        preferenceGroup{},
			allocated: 80,
			expected:  nil,
		}, {
			name: "one",
			pg: preferenceGroup{
				{min: 1, preferred: []int{5}, max: 10},
			},
			allocated: 80,
			expected:  []int{5},
		}, {
			name: "one grow to max",
			pg: preferenceGroup{
				{max: 10, grow: false},
			},
			allocated: 80,
			expected:  []int{10},
		}, {
			name: "one grow to allocated size",
			pg: preferenceGroup{
				{min: 1, preferred: []int{1}, max: 1000, grow: true},
			},
			allocated: 80,
			expected:  []int{80},
		}, {
			name: "two even empty grow",
			pg: preferenceGroup{
				{},
				{},
			},
			allocated: 80,
			expected:  []int{40, 40},
		}, {
			name: "two uneven to max",
			pg: preferenceGroup{
				{max: 10},
				{max: 70},
			},
			allocated: 80,
			expected:  []int{10, 70},
		}, {
			name: "two uneven grow and max",
			pg: preferenceGroup{
				{max: 10},
				{},
			},
			allocated: 80,
			expected:  []int{10, 70},
		}, {
			name: "three uneven to max",
			pg: preferenceGroup{
				{max: 10},
				{max: 20},
				{max: 30},
			},
			allocated: 80,
			expected:  []int{10, 20, 30},
		}, {
			name: "three uneven preference and max",
			pg: preferenceGroup{
				{max: 10},
				{preferred: []int{5, 15}, max: 20},
				{max: 30},
			},
			allocated: 80,
			expected:  []int{10, 15, 30},
		}, {
			name: "three uneven over allocated",
			pg: preferenceGroup{
				{max: 10},
				{preferred: []int{95, 15}},
				{max: 30},
			},
			allocated: 80,
			expected:  []int{10, 40, 30},
		}, {
			name: "just enough for min",
			pg: preferenceGroup{
				{min: 20, preferred: []int{30}, max: 40},
				{min: 20, preferred: []int{30}, max: 40},
				{min: 20, preferred: []int{30}, max: 40},
				{min: 20, preferred: []int{30}, max: 40},
			},
			allocated: 80,
			expected:  []int{20, 20, 20, 20},
		}, {
			name: "go below min when over allocated",
			pg: preferenceGroup{
				{min: 25, preferred: []int{30}, max: 40},
				{min: 25, preferred: []int{30}, max: 40},
				{min: 25, preferred: []int{30}, max: 40},
				{min: 25, preferred: []int{30}, max: 40},
			},
			allocated: 80,
			expected:  []int{20, 20, 20, 20},
		}, {
			name: "remainder",
			pg: preferenceGroup{
				{max: 30, grow: true},
				{max: 30, grow: true},
			},
			allocated: 61,
			expected:  []int{30, 31},
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
		name     string
		input    Grid
		expected Grid
	}{
		{
			name: "no expand",
			input: [][]Layout{
				{{id: 1}, {id: 2}},
				{{id: 3}, {id: 4}},
			},
			expected: [][]Layout{
				{{id: 1}, {id: 2}},
				{{id: 3}, {id: 4}},
			},
		}, {
			name: "horizontal span",
			input: [][]Layout{
				{{id: 1, SpanWidth: 2}},
				{{id: 2}, {id: 3}},
			},
			expected: [][]Layout{
				{{id: 1, SpanWidth: 2}, {id: 1, SpanWidth: 2, wDuplicate: true}},
				{{id: 2}, {id: 3}},
			},
		}, {
			name: "horizontal span and expand",
			input: [][]Layout{
				{{id: 1}, {id: 2, SpanWidth: 3}},
				{{id: 3}, {id: 4}},
			},
			expected: [][]Layout{
				{{id: 1}, {id: 2, SpanWidth: 3}, {id: 2, SpanWidth: 3, wDuplicate: true}, {id: 2, SpanWidth: 3, wDuplicate: true}},
				{{id: 3}, {id: 4}},
			},
		}, {
			name: "vertical span",
			input: [][]Layout{
				{{id: 1, SpanHeight: 2}, {id: 2}},
				{{id: 3}},
			},
			expected: [][]Layout{
				{{id: 1, SpanHeight: 2}, {id: 2}},
				{{id: 1, SpanHeight: 2, hDuplicate: true}, {id: 3}},
			},
		}, {
			name: "vertical span and expand",
			input: [][]Layout{
				{{id: 1, SpanHeight: 2}},
			},
			expected: [][]Layout{
				{{id: 1, SpanHeight: 2}},
				{{id: 1, SpanHeight: 2, hDuplicate: true}},
			},
		}, {
			name: "vertical and horizontal span and overflow",
			input: [][]Layout{
				{{id: 1, SpanWidth: 2, SpanHeight: 2}},
			},
			expected: [][]Layout{
				{{id: 1, SpanWidth: 2, SpanHeight: 2}, {id: 1, SpanWidth: 2, SpanHeight: 2, wDuplicate: true}},
				{{id: 1, SpanWidth: 2, SpanHeight: 2, hDuplicate: true}, {id: 1, SpanWidth: 2, SpanHeight: 2, hDuplicate: true, wDuplicate: true}},
			},
		}, {
			// input:
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
			input: [][]Layout{
				{{id: 1}, {id: 2, SpanWidth: 2, SpanHeight: 2}, {id: 3}},
				{{id: 4, SpanHeight: 2}, {id: 5}},
				{{id: 6}, {id: 7, SpanWidth: 2}},
			},
			expected: [][]Layout{
				{{id: 1}, {id: 2, SpanWidth: 2, SpanHeight: 2}, {id: 2, SpanWidth: 2, SpanHeight: 2, wDuplicate: true}, {id: 3}},
				{{id: 4, SpanHeight: 2}, {id: 2, SpanWidth: 2, SpanHeight: 2, hDuplicate: true}, {id: 2, SpanWidth: 2, SpanHeight: 2, hDuplicate: true, wDuplicate: true}, {id: 5}},
				{{id: 4, SpanHeight: 2, hDuplicate: true}, {id: 6}, {id: 7, SpanWidth: 2}, {id: 7, SpanWidth: 2, wDuplicate: true}},
			},
		}, {
			name: "prefs even split",
			input: [][]Layout{
				{{id: 1, SpanWidth: 2, SpanHeight: 2, MinHeight: 10, MaxHeight: 100, PreferredHeight: 50, MinWidth: 10, MaxWidth: 100, PreferredWidth: 50}},
			},
			expected: [][]Layout{
				{{id: 1, SpanWidth: 2, SpanHeight: 2, MinHeight: 5, MaxHeight: 50, PreferredHeight: 25, MinWidth: 5, MaxWidth: 50, PreferredWidth: 25},
					{id: 1, SpanWidth: 2, SpanHeight: 2, MinHeight: 5, MaxHeight: 50, PreferredHeight: 25, MinWidth: 5, MaxWidth: 50, PreferredWidth: 25, wDuplicate: true}},
				{{id: 1, SpanWidth: 2, SpanHeight: 2, MinHeight: 5, MaxHeight: 50, PreferredHeight: 25, MinWidth: 5, MaxWidth: 50, PreferredWidth: 25, hDuplicate: true},
					{id: 1, SpanWidth: 2, SpanHeight: 2, MinHeight: 5, MaxHeight: 50, PreferredHeight: 25, MinWidth: 5, MaxWidth: 50, PreferredWidth: 25, hDuplicate: true, wDuplicate: true},
				},
			},
		}, {
			// these are all 1 larger than "perfs even split", the odd number is discarded.
			name: "prefs odd split rounding error",
			input: [][]Layout{
				{{id: 1, SpanWidth: 2, SpanHeight: 2, MinHeight: 11, MaxHeight: 101, PreferredHeight: 51, MinWidth: 11, MaxWidth: 101, PreferredWidth: 51}},
			},
			expected: [][]Layout{
				{{id: 1, SpanWidth: 2, SpanHeight: 2, MinHeight: 5, MaxHeight: 50, PreferredHeight: 25, MinWidth: 5, MaxWidth: 50, PreferredWidth: 25},
					{id: 1, SpanWidth: 2, SpanHeight: 2, MinHeight: 5, MaxHeight: 50, PreferredHeight: 25, MinWidth: 5, MaxWidth: 50, PreferredWidth: 25, wDuplicate: true}},
				{{id: 1, SpanWidth: 2, SpanHeight: 2, MinHeight: 5, MaxHeight: 50, PreferredHeight: 25, MinWidth: 5, MaxWidth: 50, PreferredWidth: 25, hDuplicate: true},
					{id: 1, SpanWidth: 2, SpanHeight: 2, MinHeight: 5, MaxHeight: 50, PreferredHeight: 25, MinWidth: 5, MaxWidth: 50, PreferredWidth: 25, hDuplicate: true, wDuplicate: true},
				},
			},
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := expandSpans(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestMergeDocks_Empty(t *testing.T) {
	assert.Nil(t, mergeDocks(nil, nil))
}

func TestMergeDocks_InvalidCardinal(t *testing.T) {
	panicFunc := func() {
		mergeDocks([][]Layout{{{id: 0}}}, []dock{{id: 1, cardinal: 99}})
	}
	assert.PanicsWithError(t, "invalid cardinal", panicFunc)
}

func TestMergeDocks(t *testing.T) {
	testcases := []struct {
		name     string
		start    Grid
		docks    []dock
		expected Grid
	}{
		{
			name: "simple north",
			start: [][]Layout{
				{{id: 0}},
			},
			docks: []dock{
				{id: 1, cardinal: NORTH, min: 10, preferred: 10, max: 10},
			},
			expected: [][]Layout{
				{{id: 1, SpanWidth: 1, MinHeight: 10, PreferredHeight: 10, MaxHeight: 10}},
				{{id: 0}},
			},
		}, {
			name: "simple south",
			start: [][]Layout{
				{{id: 0}},
			},
			docks: []dock{
				{id: 1, cardinal: SOUTH, min: 10, preferred: 10, max: 10},
			},
			expected: [][]Layout{
				{{id: 0}},
				{{id: 1, SpanWidth: 1, MinHeight: 10, PreferredHeight: 10, MaxHeight: 10}},
			},
		}, {
			name: "simple west",
			start: [][]Layout{
				{{id: 0}},
			},
			docks: []dock{
				{id: 1, cardinal: WEST, min: 10, preferred: 10, max: 10},
			},
			expected: [][]Layout{
				{{id: 1, SpanHeight: 1, MinWidth: 10, PreferredWidth: 10, MaxWidth: 10}, {id: 0}},
			},
		}, {
			name: "simple east",
			start: [][]Layout{
				{{id: 0}},
			},
			docks: []dock{
				{id: 1, cardinal: EAST, min: 10, preferred: 10, max: 10},
			},
			expected: [][]Layout{
				{{id: 0}, {id: 1, SpanHeight: 1, MinWidth: 10, PreferredWidth: 10, MaxWidth: 10}},
			},
		}, {
			name: "double north",
			start: [][]Layout{
				{{id: 0}, {id: 1}},
			},
			docks: []dock{
				{id: 2, cardinal: NORTH, min: 10, preferred: 10, max: 10},
			},
			expected: [][]Layout{
				{{id: 2, SpanWidth: 2, MinHeight: 5, PreferredHeight: 5, MaxHeight: 5}, {id: 2, SpanWidth: 2, MinHeight: 5, PreferredHeight: 5, MaxHeight: 5, wDuplicate: true}},
				{{id: 0}, {id: 1}},
			},
		}, {
			name: "double south",
			start: [][]Layout{
				{{id: 0}, {id: 1}},
			},
			docks: []dock{
				{id: 2, cardinal: SOUTH, min: 10, preferred: 10, max: 10},
			},
			expected: [][]Layout{
				{{id: 0}, {id: 1}},
				{{id: 2, SpanWidth: 2, MinHeight: 5, PreferredHeight: 5, MaxHeight: 5}, {id: 2, SpanWidth: 2, MinHeight: 5, PreferredHeight: 5, MaxHeight: 5, wDuplicate: true}},
			},
		}, {
			name: "double west",
			start: [][]Layout{
				{{id: 0}},
				{{id: 1}},
			},
			docks: []dock{
				{id: 2, cardinal: WEST, min: 10, preferred: 10, max: 10},
			},
			expected: [][]Layout{
				{{id: 2, SpanHeight: 2, MinWidth: 5, PreferredWidth: 5, MaxWidth: 5}, {id: 0}},
				{{id: 2, SpanHeight: 2, MinWidth: 5, PreferredWidth: 5, MaxWidth: 5, hDuplicate: true}, {id: 1}},
			},
		}, {
			name: "double east",
			start: [][]Layout{
				{{id: 0}},
				{{id: 1}},
			},
			docks: []dock{
				{id: 2, cardinal: EAST, min: 10, preferred: 10, max: 10},
			},
			expected: [][]Layout{
				{{id: 0}, {id: 2, SpanHeight: 2, MinWidth: 5, PreferredWidth: 5, MaxWidth: 5}},
				{{id: 1}, {id: 2, SpanHeight: 2, MinWidth: 5, PreferredWidth: 5, MaxWidth: 5, hDuplicate: true}},
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
		hExpected preferenceGroup
		wExpected preferenceGroup
	}{
		{
			name: "1x1",
			input: [][]Layout{
				{{id: 1, MinHeight: 10, MaxHeight: 100, PreferredHeight: 50, MinWidth: 10, MaxWidth: 100, PreferredWidth: 50}},
			},
			hExpected: []sizePreference{{min: 10, preferred: []int{50}, max: 100}},
			wExpected: []sizePreference{{min: 10, preferred: []int{50}, max: 100}},
		}, {
			name: "2x2 one defaults",
			input: [][]Layout{
				{{id: 1, MinHeight: 10, MaxHeight: 100, PreferredHeight: 50, MinWidth: 10, MaxWidth: 100, PreferredWidth: 50}, {id: 2}},
				{{id: 3}, {id: 4, MinHeight: 10, MaxHeight: 100, PreferredHeight: 50, MinWidth: 10, MaxWidth: 100, PreferredWidth: 50}},
			},
			hExpected: []sizePreference{{min: 10, preferred: []int{50}, max: 100}, {min: 10, preferred: []int{50}, max: 100}},
			wExpected: []sizePreference{{min: 10, preferred: []int{50}, max: 100}, {min: 10, preferred: []int{50}, max: 100}},
		}, {
			name: "1x3 max(min)) and min(max)",
			input: [][]Layout{
				{
					{id: 1, MinHeight: 5, MaxHeight: 50, PreferredHeight: 25},
					{id: 2, MinHeight: 10, MaxHeight: 100, PreferredHeight: 50},
					{id: 3},
				},
			},
			hExpected: []sizePreference{{min: 10, preferred: []int{25, 50}, max: 50}},
			wExpected: []sizePreference{{}, {}, {}},
		}, {
			name: "3x1 max(min)) and min(max)",
			input: [][]Layout{
				{{id: 1, MinWidth: 5, MaxWidth: 50, PreferredWidth: 25}},
				{{id: 2, MinWidth: 10, MaxWidth: 100, PreferredWidth: 50}},
				{{id: 3}},
			},
			hExpected: []sizePreference{{}, {}, {}},
			wExpected: []sizePreference{{min: 10, preferred: []int{25, 50}, max: 50}},
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
	l.Dock(NORTH, 1, 2, 3)
	bl := l.(*bubbleLayout)
	require.Equal(t, []dock{{cardinal: NORTH, id: 1, min: 1, preferred: 2, max: 3}}, bl.docks)
}

func TestValidate_FailureHeight(t *testing.T) {
	l := New()
	l.Add(Layout{MinHeight: 100})
	l.Add(Layout{MaxHeight: 10})
	require.ErrorContains(t, l.Validate(), makeRowViolation(0, 100, 10).Error())
}

func TestValidate_FailureWidth(t *testing.T) {
	l := New()
	l.Add(Layout{MinWidth: 100})
	l.Wrap()
	l.Add(Layout{MaxWidth: 10})
	require.ErrorContains(t, l.Validate(), makeColViolation(0, 100, 10).Error())
}
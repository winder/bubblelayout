package bubblelayout

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrorUnwrap(t *testing.T) {
	inner := fmt.Errorf("inner")
	var err error = makeErrStringLayout("one", "two", inner)
	assert.ErrorIs(t, err, inner)
}

func TestIsCardinal(t *testing.T) {
	for _, c := range []Cardinal{NORTH, SOUTH, EAST, WEST} {
		assert.True(t, isCardinal(string(c)))
	}
	assert.False(t, isCardinal("not a cardinal"))
}

func TestBorderSizeRegexp(t *testing.T) {
	testcases := []struct {
		in  string
		out []string
	}{
		{
			in:  "",
			out: nil,
		}, {
			in:  "!",
			out: nil,
		}, {
			in:  "1",
			out: []string{"1", "1", "", "", ""},
		}, {
			in:  "1!",
			out: []string{"1!", "1", "", "", "!"},
		}, {
			in:  "1:2",
			out: []string{"1:2", "1", "2", "", ""},
		}, {
			in:  "1:2!",
			out: []string{"1:2!", "1", "2", "", "!"},
		}, {
			in:  "1:2:3",
			out: []string{"1:2:3", "1", "2", "3", ""},
		}, {
			in:  "1:2:3!",
			out: []string{"1:2:3!", "1", "2", "3", "!"},
		}, {
			in:  "1:2:3:4",
			out: nil,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.in, func(t *testing.T) {
			t.Parallel()
			out := borderSizePattern.FindStringSubmatch(tc.in)
			assert.Equal(t, tc.out, out)
		})
	}
}

func TestPartSize(t *testing.T) {
	testcases := []struct {
		in  string
		out BoundSize
		err string
	}{
		{
			in:  "",
			err: "invalid bound size '': did not match pattern",
		}, {
			in:  "1",
			out: BoundSize{Preferred: 1},
		}, {
			in:  "1!",
			out: BoundSize{Min: 1, Preferred: 1, Max: 1},
		}, {
			in:  "1:2",
			out: BoundSize{Min: 1, Preferred: 2},
		}, {
			in:  "1:2!",
			err: "invalid bound size '1:2!': use '!' with only one number",
		}, {
			in:  "1:2:3",
			out: BoundSize{Min: 1, Preferred: 2, Max: 3},
		}, {
			in:  "1:2:3!",
			err: "invalid bound size '1:2:3!': use '!' with only one number",
		}, {
			in:  "n:10:n",
			out: BoundSize{Preferred: 10},
		}, {
			in:  "null:10:null",
			out: BoundSize{Preferred: 10},
		}, {
			in:  "nil:10:nil",
			out: BoundSize{Preferred: 10},
		}, {
			in:  "n:10:null",
			out: BoundSize{Preferred: 10},
		}, {
			in:  "n:10:nil",
			out: BoundSize{Preferred: 10},
		}, {
			in:  "10::",
			out: BoundSize{Min: 10},
		}, {
			in:  "::10",
			out: BoundSize{Max: 10},
		}, {
			in:  ":10:",
			out: BoundSize{Preferred: 10},
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.in, func(t *testing.T) {
			t.Parallel()
			out, err := parseSize(tc.in)

			if tc.err != "" {
				assert.ErrorContains(t, err, tc.err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.out, out)
			}
		})
	}
}

func TestConvertString(t *testing.T) {
	testcases := []struct {
		name  string
		in    string
		inArr []string
		out   layout
		err   string
	}{
		{
			name: "wrap",
			in:   "wrap",
			out:  layout{wrap: true},
		}, {
			name: "empty",
			in:   "",
			out:  layout{},
		}, {
			name: "grow",
			in:   "grow",
			out:  layout{Cell: Cell{GrowWidth: true, GrowHeight: true}},
		}, {
			name:  "growx,groww",
			inArr: []string{"growx", "groww"},
			out:   layout{Cell: Cell{GrowWidth: true}},
		}, {
			name:  "growy,growh",
			inArr: []string{"growy", "growh"},
			out:   layout{Cell: Cell{GrowHeight: true}},
		}, {
			name:  "invalid dock-wrong direction",
			inArr: []string{"dock left", "dock 1"},
			err:   "invalid cardinal direction",
		}, {
			name: "invalid dock-no direction",
			in:   "dock",
			err:  "string api conversion error for inputLayout 'dock': dock direction is missing",
		}, {
			name: "invalid span-no numbers",
			in:   "span",
			err:  "string api conversion error for inputLayout 'span': wrong number of inputs to span, expected 1 or 2 received '[]'",
		}, {
			name:  "spanw",
			inArr: []string{"span 13", "spanw 13", "spanx 13", "sx 13"},
			out:   layout{Cell: Cell{SpanWidth: 13}},
		}, {
			name:  "spanh",
			inArr: []string{"spanh 13", "spany 13", "sy 13"},
			out:   layout{Cell: Cell{SpanHeight: 13}},
		}, {
			name:  "invalid spanw and spanh",
			inArr: []string{"spanw", "sx", "spanw notanumber", "spanh", "sy", "spanh notanumber"},
			err:   "wrong number of inputs, expected 1 received '[]'",
		}, {
			name:  "span 2d",
			inArr: []string{"span 13 14", "spanw 13, spanh 14", "spanh 14, spanw 13"},
			out:   layout{Cell: Cell{SpanWidth: 13, SpanHeight: 14}},
		}, {
			name:  "dock north",
			inArr: []string{"dock north", "north"},
			out:   layout{Dock: Dock{Cardinal: NORTH}},
		}, {
			name:  "dock south",
			inArr: []string{"dock south", "south"},
			out:   layout{Dock: Dock{Cardinal: SOUTH}},
		}, {
			name:  "dock east",
			inArr: []string{"dock east", "east"},
			out:   layout{Dock: Dock{Cardinal: EAST}},
		}, {
			name:  "dock west",
			inArr: []string{"dock west", "west"},
			out:   layout{Dock: Dock{Cardinal: WEST}},
		}, {
			name: "dock west 1:2:3",
			in:   "dock west 1:2:3",
			out:  layout{Dock: Dock{Cardinal: WEST, Min: 1, Preferred: 2, Max: 3}},
		}, {
			name:  "width and height error-missing",
			inArr: []string{"width", "w", "height", "h"},
			err:   "bound size is missing",
		}, {
			name:  "width and height error-invalid",
			inArr: []string{"width not-a-number", "w 1:2!", "height 1:2:3!", "h !"},
			err:   "invalid bound size",
		}, {
			name:  "width",
			inArr: []string{"width 2:2:2", "width 2!", "w 2!"},
			out:   layout{Cell: Cell{MinWidth: 2, PreferredWidth: 2, MaxWidth: 2}},
		}, {
			name:  "height",
			inArr: []string{"height 2:2:2", "height 2!", "h 2!"},
			out:   layout{Cell: Cell{MinHeight: 2, PreferredHeight: 2, MaxHeight: 2}},
		}, {
			name: "command after multi-token command",
			in:   "width 1:2:3, grow, span 1 2, wrap",
			out:  layout{wrap: true, Cell: Cell{SpanWidth: 1, SpanHeight: 2, GrowWidth: true, GrowHeight: true, MinWidth: 1, PreferredWidth: 2, MaxWidth: 3}},
		}, {
			name:  "unknown constraint",
			inArr: []string{"unknown constraint", "100"},
			err:   "unknown constraint",
		},
	}

	for _, tc := range testcases {
		tc := tc
		if tc.in != "" {
			tc.inArr = append(tc.inArr, tc.in)
		}
		for _, in := range tc.inArr {
			in := in
			t.Run(fmt.Sprintf("%s:%s", tc.name, in), func(t *testing.T) {
				t.Parallel()

				out, err := convertToLayout(in)
				if tc.err != "" {
					require.ErrorContains(t, err, tc.err)
				} else {
					require.NoError(t, err)
					assert.Equal(t, tc.out, out, "inputLayout: %s", in)
				}
			})
		}
	}
}

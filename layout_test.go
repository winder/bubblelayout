package bubblelayout_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	bl "github.com/winder/bubblelayout"
)

func TestUnknownComponent(t *testing.T) {
	l := bl.New()
	msg := l.Resize(10, 10)
	_, err := msg.Size(1)
	require.Error(t, err)
	require.Equal(t, "view not registered", err.Error())
}

func TestOneComponent(t *testing.T) {
	l := bl.New()
	id1 := l.Cell(bl.Cell{})
	msg := l.Resize(10, 10)

	size, err := msg.Size(id1)
	require.NoError(t, err)
	require.Equal(t, bl.Size{Width: 10, Height: 10}, size)
}

func TestHorizontalComponents(t *testing.T) {
	l := bl.New()
	id1 := l.Cell(bl.Cell{})
	id2 := l.Cell(bl.Cell{})
	msg := l.Resize(10, 10)

	{
		size, err := msg.Size(id1)
		require.NoError(t, err)
		require.Equal(t, bl.Size{Width: 5, Height: 10}, size)
	}

	{
		size, err := msg.Size(id2)
		require.NoError(t, err)
		require.Equal(t, bl.Size{Width: 5, Height: 10}, size)
	}
}

func TestVerticalComponents(t *testing.T) {
	l := bl.New()
	id1 := l.Cell(bl.Cell{})
	l.Wrap()
	id2 := l.Cell(bl.Cell{})
	msg := l.Resize(10, 10)

	{
		size, err := msg.Size(id1)
		require.NoError(t, err)
		require.Equal(t, bl.Size{Width: 10, Height: 5}, size)
	}

	{
		size, err := msg.Size(id2)
		require.NoError(t, err)
		require.Equal(t, bl.Size{Width: 10, Height: 5}, size)
	}
}

func TestComplex(t *testing.T) {
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
	l := bl.New()

	var ids []bl.ID
	ids = append(ids, l.Cell(bl.Cell{}))
	ids = append(ids, l.Cell(bl.Cell{SpanWidth: 2, SpanHeight: 2}))
	ids = append(ids, l.Cell(bl.Cell{}))
	l.Wrap()
	ids = append(ids, l.Cell(bl.Cell{SpanHeight: 2}))
	ids = append(ids, l.Cell(bl.Cell{}))
	l.Wrap()
	ids = append(ids, l.Cell(bl.Cell{}))
	ids = append(ids, l.Cell(bl.Cell{SpanWidth: 2}))

	msg := l.Resize(100, 75)

	{
		for _, id := range []bl.ID{ids[0], ids[2], ids[4], ids[5]} {
			size, err := msg.Size(id)
			require.NoError(t, err)
			assert.Equal(t, bl.Size{Width: 25, Height: 25}, size, "1x1 cell should be 25x25")
		}
	}
	{
		size, err := msg.Size(ids[1])
		require.NoError(t, err)
		assert.Equal(t, bl.Size{Width: 50, Height: 50}, size, "2x2 cell should be 50x50")
	}
	{
		size, err := msg.Size(ids[3])
		require.NoError(t, err)
		assert.Equal(t, bl.Size{Width: 25, Height: 50}, size, "1x2 cell should be 25x50")
	}
	{
		size, err := msg.Size(ids[6])
		require.NoError(t, err)
		assert.Equal(t, bl.Size{Width: 50, Height: 25}, size, "2x1 cell should be 50x25")
	}
}

func TestValidateCache(t *testing.T) {
	l := bl.New()
	l.Cell(bl.Cell{})
	require.NoError(t, l.Validate())
	require.NoError(t, l.Validate())
}

func TestResize_Panic(t *testing.T) {
	l := bl.New()
	l.Cell(bl.Cell{MinHeight: 100})
	l.Cell(bl.Cell{MaxHeight: 10})
	panicFunc := func() {
		l.Resize(100, 100)
	}
	require.Panicsf(t, panicFunc, "constraint violation")
}

func TestAdd_Panic(t *testing.T) {
	l := bl.New()
	panicFunc := func() { l.Add("invalid constraint options") }
	require.Panics(t, panicFunc)
}

func TestMaybeAdd(t *testing.T) {
	l := bl.New()
	_, err := l.MaybeAdd("invalid constraint options")
	require.ErrorContains(t, err, "invalid constraint")

}

func TestDockAPI(t *testing.T) {
	l := bl.New()
	id1 := l.Cell(bl.Cell{})
	l.Wrap()
	id2 := l.Cell(bl.Cell{SpanHeight: 2})
	id3 := l.Dock(bl.Dock{Cardinal: bl.EAST})

	msg := l.Resize(100, 75)
	{
		size, err := msg.Size(id1)
		require.NoError(t, err)
		assert.Equal(t, bl.Size{Width: 50, Height: 25}, size)
	}
	{
		size, err := msg.Size(id2)
		require.NoError(t, err)
		assert.Equal(t, bl.Size{Width: 50, Height: 50}, size)
	}
	{
		size, err := msg.Size(id3)
		require.NoError(t, err)
		assert.Equal(t, bl.Size{Width: 50, Height: 75}, size)
	}
}

func TestDockAdd(t *testing.T) {
	l := bl.New()
	id1 := l.Add("wrap")
	id2 := l.Add("spanh 2")
	id3 := l.Add("dock east")

	msg := l.Resize(100, 75)
	{
		size, err := msg.Size(id1)
		require.NoError(t, err)
		assert.Equal(t, bl.Size{Width: 50, Height: 25}, size)
	}
	{
		size, err := msg.Size(id2)
		require.NoError(t, err)
		assert.Equal(t, bl.Size{Width: 50, Height: 50}, size)
	}
	{
		size, err := msg.Size(id3)
		require.NoError(t, err)
		assert.Equal(t, bl.Size{Width: 50, Height: 75}, size)
	}
}

func TestResize(t *testing.T) {
	const width = 80
	const height = 40
	testcases := []struct {
		name      string
		wOverride int
		hOverride int
		in        func() bl.BubbleLayout
		out       map[bl.ID]bl.Size
	}{
		{
			name: "simple",
			in: func() bl.BubbleLayout {
				l := bl.New()
				l.Add("width 10")
				l.Add("grow")
				return l
			},
			out: map[bl.ID]bl.Size{
				1: {Width: 10, Height: height},
				2: {Width: width - 10, Height: height},
			},
		}, {
			name: "dock",
			in: func() bl.BubbleLayout {
				l := bl.New()
				l.Add("")
				l.Add("span 2")
				l.Add("north 5!")
				return l
			},
			out: map[bl.ID]bl.Size{
				1: {Width: 27, Height: 35},
				2: {Width: 53, Height: 35},
				3: {Width: 80, Height: 5},
			},
		}, {
			name:      "empty remainder",
			hOverride: 9,
			in: func() bl.BubbleLayout {
				l := bl.New()
				l.Add("width 9!")
				l.Add("width 9!")
				l.Add("")
				return l
			},
			out: map[bl.ID]bl.Size{
				1: {Width: 9, Height: 9},
				2: {Width: 9, Height: 9},
				3: {Width: 62, Height: 9},
			},
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			l := tc.in()
			w := width
			h := height
			if tc.wOverride != 0 {
				w = tc.wOverride
			}
			if tc.hOverride != 0 {
				h = tc.hOverride
			}
			msg := l.Resize(w, h)
			for id, size := range tc.out {
				actual, err := msg.Size(id)
				require.NoError(t, err)
				assert.Equal(t, size, actual)
			}
		})
	}
}

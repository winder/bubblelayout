package layout_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/winder/layout"
)

func TestUnknownComponent(t *testing.T) {
	l := layout.New()
	msg := l.Resize(10, 10)
	_, err := msg.Size(1)
	require.Error(t, err)
	require.Equal(t, "view not registered", err.Error())
}

func TestOneComponent(t *testing.T) {
	l := layout.New()
	id1 := l.Add(layout.Layout{})
	msg := l.Resize(10, 10)

	size, err := msg.Size(id1)
	require.NoError(t, err)
	require.Equal(t, layout.Size{Width: 10, Height: 10}, size)
}

func TestHorizontalComponents(t *testing.T) {
	l := layout.New()
	id1 := l.Add(layout.Layout{})
	id2 := l.Add(layout.Layout{})
	msg := l.Resize(10, 10)

	{
		size, err := msg.Size(id1)
		require.NoError(t, err)
		require.Equal(t, layout.Size{Width: 5, Height: 10}, size)
	}

	{
		size, err := msg.Size(id2)
		require.NoError(t, err)
		require.Equal(t, layout.Size{Width: 5, Height: 10}, size)
	}
}

func TestVerticalComponents(t *testing.T) {
	l := layout.New()
	id1 := l.Add(layout.Layout{})
	l.Wrap()
	id2 := l.Add(layout.Layout{})
	msg := l.Resize(10, 10)

	{
		size, err := msg.Size(id1)
		require.NoError(t, err)
		require.Equal(t, layout.Size{Width: 10, Height: 5}, size)
	}

	{
		size, err := msg.Size(id2)
		require.NoError(t, err)
		require.Equal(t, layout.Size{Width: 10, Height: 5}, size)
	}
}

func TestComplex(t *testing.T) {
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
	l := layout.New()

	// layout
	// ---------------------------------
	// |   1   |       -       |   3   |
	// --------- -  -  2  -  - |--------
	// |   -   |       -       |   5   |
	// | - 4 - -------------------------
	// |   -   |   6   |       7       |
	// ---------------------------------
	var ids []layout.ID
	ids = append(ids, l.Add(layout.Layout{}))
	ids = append(ids, l.Add(layout.Layout{SpanWidth: 2, SpanHeight: 2}))
	ids = append(ids, l.Add(layout.Layout{}))
	l.Wrap()
	ids = append(ids, l.Add(layout.Layout{SpanHeight: 2}))
	ids = append(ids, l.Add(layout.Layout{}))
	l.Wrap()
	ids = append(ids, l.Add(layout.Layout{}))
	ids = append(ids, l.Add(layout.Layout{SpanWidth: 2}))

	msg := l.Resize(100, 75)

	{
		for _, id := range []layout.ID{ids[0], ids[2], ids[4], ids[5]} {
			size, err := msg.Size(id)
			require.NoError(t, err)
			assert.Equal(t, layout.Size{Width: 25, Height: 25}, size, "1x1 cell should be 25x25")
		}
	}
	{
		size, err := msg.Size(ids[1])
		require.NoError(t, err)
		assert.Equal(t, layout.Size{Width: 50, Height: 50}, size, "2x2 cell should be 50x50")
	}
	{
		size, err := msg.Size(ids[3])
		require.NoError(t, err)
		assert.Equal(t, layout.Size{Width: 25, Height: 50}, size, "1x2 cell should be 25x50")
	}
	{
		size, err := msg.Size(ids[6])
		require.NoError(t, err)
		assert.Equal(t, layout.Size{Width: 50, Height: 25}, size, "2x1 cell should be 50x25")
	}
}

func TestValidateCache(t *testing.T) {
	l := layout.New()
	l.Add(layout.Layout{})
	require.NoError(t, l.Validate())
	require.NoError(t, l.Validate())
}

func TestResize_Panic(t *testing.T) {
	l := layout.New()
	l.Add(layout.Layout{MinHeight: 100})
	l.Add(layout.Layout{MaxHeight: 10})
	panicFunc := func() {
		l.Resize(100, 100)
	}
	require.Panicsf(t, panicFunc, "constraint violation")
}

func TestAddStr(t *testing.T) {
	l := layout.New()
	panicFunc := func() { l.AddStr("hello") }
	require.Panics(t, panicFunc)
}

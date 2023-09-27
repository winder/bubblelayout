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
	id1 := l.Add(bl.Cell{})
	msg := l.Resize(10, 10)

	size, err := msg.Size(id1)
	require.NoError(t, err)
	require.Equal(t, bl.Size{Width: 10, Height: 10}, size)
}

func TestHorizontalComponents(t *testing.T) {
	l := bl.New()
	id1 := l.Add(bl.Cell{})
	id2 := l.Add(bl.Cell{})
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
	id1 := l.Add(bl.Cell{})
	l.Wrap()
	id2 := l.Add(bl.Cell{})
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
	l := bl.New()

	var ids []bl.ID
	ids = append(ids, l.Add(bl.Cell{}))
	ids = append(ids, l.Add(bl.Cell{SpanWidth: 2, SpanHeight: 2}))
	ids = append(ids, l.Add(bl.Cell{}))
	l.Wrap()
	ids = append(ids, l.Add(bl.Cell{SpanHeight: 2}))
	ids = append(ids, l.Add(bl.Cell{}))
	l.Wrap()
	ids = append(ids, l.Add(bl.Cell{}))
	ids = append(ids, l.Add(bl.Cell{SpanWidth: 2}))

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
	l.Add(bl.Cell{})
	require.NoError(t, l.Validate())
	require.NoError(t, l.Validate())
}

func TestResize_Panic(t *testing.T) {
	l := bl.New()
	l.Add(bl.Cell{MinHeight: 100})
	l.Add(bl.Cell{MaxHeight: 10})
	panicFunc := func() {
		l.Resize(100, 100)
	}
	require.Panicsf(t, panicFunc, "constraint violation")
}

func TestAddStr(t *testing.T) {
	l := bl.New()
	panicFunc := func() { l.AddStr("hello") }
	require.Panics(t, panicFunc)
}

Bubble Layout
=============
<p align="center">
  <a href="https://github.com/winder/bubblelayout/actions?query=workflow%3Atests%20and%20coverage+event%3Apush">
    <img title="GitHub Workflow Status (test @ main)" src="https://img.shields.io/github/actions/workflow/status/winder/bubblelayout/test.yml?branch=main&label=test&style=flat-square">
  </a>
  <a href="https://codecov.io/gh/winder/bubblelayout">
    <img title="Code Coverage" src="https://img.shields.io/codecov/c/github/winder/bubblelayout/main?style=flat-square">
  </a>
  <a href="https://pkg.go.dev/github.com/winder/bubblelayout">
    <img title="Go Documentation" src="https://pkg.go.dev/badge/github.com/winder/bubblelayout?style=flat-square">
  </a>
  <a href="https://goreportcard.com/report/github.com/winder/bubblelayout">
    <img title="Go Report Card" src="https://goreportcard.com/badge/github.com/winder/bubblelayout?style=flat-square">
  </a>
</p>

Declarative layout manager for [BubbleTea](https://github.com/charmbracelet/bubbletea/).

BubbleLayout provides a powerful API without sacrificing readability. Inspired by [MiG Layout](http://miglayout.com/).

## Usage
```
go get -u github.com/winder/bubblelayout@latest
```

BubbleLayout uses a declared layout to create `bl.BubbleLayoutMsg` events which are used to provide exact model dimensions. These are created with BubbleLayout's `Resize` function, which translates a `tea.WindowSizeMsg` nto a `bl.BubbleLayoutMsg`.

The dependency should be imported, by convention it is renamed to `bl`:
```
import (
    bl "github.com/winder/bubblelayout"
)
```

The conversion should be done once by calling the Resize function. If it is a top level model, the converted message can be dispatched to child models. Alternatively it can be fed back into the event loop as can be seen below:

```go
func (m SomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  switch msg := msg.(type) {
  case tea.WindowSizeMsg:
    // Convert WindowSizeMsg to BubbleLayoutMsg.
    return m, func() tea.Msg {
      return m.layout.Resize(msg.Width, msg.Height)
    }
  }
  return m, nil
}
```

Window size handling would now be a matter of processing `bl.BubbleMayoutMsg` updates:
```go
func (m aModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  switch msg := msg.(type) {
  case bl.BubbleLayoutMsg:
    sz, _ := msg.Size(m.id)
    m.width = sz.Width
    m.height = sz.Height
  }
}
```

### Layout Declaration

The layout is typically defined during root component initialization. It defines all constrains for sizing the different components using the `Add` function and a StringAPI. For details about how layout works, see the [MiG Layout Quick Start Quide (pdf)](http://www.miglayout.com/mavensite/docs/QuickStart.pdf). Not all options are supported, but most of the basics are.

An alternative to the StringAPI is available by adding raw layout objects directly. This is probably more idiomatic for go APIs, but is significantly more verbose. For more on this refer to the `Cell` and `Dock` methods.


#### **Add** components to the grid
Components are added to the layout with the `Add` function. In the following example two components are added. The first has a preferred width of 10, the second is instructed to grow to fill whatever space remains. In this example, the grow constraint is optional because any component without a size preference attempts to fill all available space.

[Simple example code](./examples/simple/main.go)

```go
layoutModel := layoutModel{layout: bl.New()}
layoutModel.leftID = layoutModel.layout.Add("width 10")
layoutModel.rightID = layoutModel.layout.Add("grow")
```

![Simple example image](./examples/simple/simple.png)

#### **Span** components across multiple cells
In many cases you may not want all cells to be a uniform grid. When this happens you can make use of the `span` constraints. They are used to define components made up of multiple cells. Spans can be made horizontally or vertically.

[Spans example code](./examples/spans/main.go)

```go
layout := bl.New()
layout.Add("")
layout.Add("span 2 2")
layout.Add("wrap")

layout.Add("spanh 2")
layout.Add("wrap")

layout.Add("")
layout.Add("spanw 2")
```

![Spans example image](./examples/spans/spans.png)

#### **Dock** components for common overrides

It is often useful to define certain components by their absolute location. With dock's you can specify things like a header that should always be placed at the top of the UI or a status bar which is always at the bottom. Note that if you have multiple overlapping docs, the order that they are defined determines which one is drawn over the corner.

[Docking example code](./examples/docking/main.go)

```go
layout := bl.New()
layout.Add("")
layout.Add("wrap")
layout.Add("span 2 2")

layout.Add("dock north 1!")
layout.Add("dock south 1!")
layout.Add("dock east 1:10")
layout.Add("dock west 1:10")
```

![Docking example image](./examples/docking/docking.png)

## Comments About Cell Sizes

When defining a layout, width and height `BoundSize` preferences may be provided for each cell. The preferences can be set globally by using `bl.NewWithConstraints(width, height PreferenceGroup)` or on each cell by using **BoundSize** notation. The string definition is compatible with MiGLayout:

> A **bound size** is a size that optionally has a lower and/or upper bound and consists of one to three Unit Values. Practically it is a minimum/preferred/maximum size combination but none of the sizes are actually mandatory. If a size is missing (e.g. the preferred) it is null and will be replaced by the most appropriate value.
>
> The format is **"min:preferred:max"**, however there are shorter versions since for instance it is seldom needed to specify the maximum size.
>
> * A single value (E.g. **"10"**) sets only the preferred size and is exactly the same as "null:10:null" and **":10:"** and **"n:10:n"**.
> * Two values (E.g. **"10:20"**) means minimum and preferred size and is exactly the same as **"10:20:null"** and **"10:20:"** and **"10:20:n"**
> * The use a of an exclamation mark (E.g. **"20!"**) means that the value should be used for all size types and no colon may then be used in the string. It is the same as **"20:20:20"**.

All of this to say: yes, I have brought **null** to go. I've taken the liberty of supporting **nil** as well.

## Future Development

MiGLayout defines many features beyond what is currently supported by bubble layout. What follows is an incomplete list of features which may be added in the future:
* "pad" and "margin" to manage spacing.
* "split" cells to allow cells that do not align with the overall grid.
* "hidden" / "visible" and a way to toggle visibility and whether they still take up space.
* "flow" order to allow defining layouts vertically or from right to left.
* "shrink" to indicate how readily cells should be reduced from their preferred size.
* "priority" for shrink/grow to add finer control over how space is allocated when there is too much or not enough.
* [so many more.](http://www.miglayout.com/whitepaper.html)

Other cool features:
* BubbleTea utilities - currently omitted to avoid a BubbleTea dependency:
  * `ResizeCmd`: helper so that you don't have to wrap `layout.Resize` in an anonymous function.
  * `LayoutModel`: the auto `tea.WindowSizeMsg` translator model used in examples.
* BubbleTea auto renderer: use something like `lipgloss.Place` to render views in place.
* Constraint events: `bl.SpaceOverallocated`, `bl.UnallocatedSpace`, ...
* Borders: automatically fill in border characters when margins and padding is defined.
* What else would you like to see?

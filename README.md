Bubble Layout
=============

Declarative layout manager for [BubbleTea](https://github.com/charmbracelet/bubbletea/).

BubbleLayout provides a powerful API without sacrificing readability. Inspired by [MiG Layout](http://miglayout.com/).

## Usage

BubbleLayout uses a declared layout to translate `tea.WindowSizeMsg` events into `bl.BubbleLayoutMsg` events. Individual parts of the UI would handle the `bl.BubbleLayoutMsg` to retrieve their absolute dimensions using a unique ID.

The dependency should be imported, by convention it is renamed to `bl`:
```
import (
    bl "github.com/winder/bubblelayout"
)
```

### Layout Declaration

The layout is typically defined during root component initialization. It defines all of the constrains which should be used when sizing different components by adding `Cell`s and `Dock`s. You can use `Cell` and `Dock` to add the raw objects, or define the layout using the string API.

For more details about how layout works, see the [MiG Layout Quick Start Quide (pdf)](http://www.miglayout.com/mavensite/docs/QuickStart.pdf).

#### **Add** components to the grid
Here is a simple example which with two side-by-side sections. Notice that the second component doesn't define `MaxWidth` and neither component defines `MaxHeight`, BubbleLayout takes this to mean that the singular row should fill the available height and the second component should fill the available width. `MinWidth` and `PreferredWidth` are also available, these are all considered by the layout engine when calculating view dimensions.

[Simple example code](./examples/simple/main.go)

```go
layout: bl.NewWithConstraints(bl.PreferenceGroup{{Max: 10}, {Grow: true}}, nil),
layoutModel.leftID = layoutModel.layout.Add("")
layoutModel.rightID = layoutModel.layout.Add("")
```

![Simple example image](./examples/simple/simple.png)

#### **Span** components across multiple cells
Here is a more complex example. It defines a layout utilizing horizontal and vertical spans, these allow you to define your grid with components that take up multiple cells.

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

Here is an example that has fixed size components at the top and bottom of the layout. Note that if you have multiple overlapping docs, the order that they are defined determines which one is drawn over the corner.

[Docking example code](./examples/docking/main.go)

```go
layout := bl.New()
layout.Add("")
layout.Add("wrap")
layout.Add("span 2 2")

layout.Add("dock north 1:1:1")
layout.Add("dock south 1:1:1")
layout.Add("dock east 1:10:10")
layout.Add("dock west 1:10:10")
```

![Docking example image](./examples/docking/docking.png)

### Layout model

Somewhere in your program, you'll need to capture the `tea.WindowSizeMsg` and feed it back into BubbleTea as a `BubbleLayoutMsg`.

```go
func (m layoutModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  switch msg := msg.(type) {
  case tea.WindowSizeMsg:
    // Convert WindowSizeMsg to BubbleLayoutMsg.
    return m, func() tea.Msg {
      return m.layout.Resize(msg.Width, msg.Height)
    }
  return m, nil
}
```

### Resizing the view

Each of your views should be initialized with the layout ID emitted from the layout definition. From that point forward, simply listen for the `BubbleLayoutMsg`.

```go
func (m myModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  switch msg := msg.(type) {
  case bl.BubbleLayoutMsg:
    sz, _ := msg.Size(m.id)
    m.width = sz.Width
    m.height = sz.Height
  }

  return m, nil
}
```

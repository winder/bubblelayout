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

## Comments About Cell Sizes

When defining a layout, width and height `BoundSize` preferences may be provided for each cell. The preferences can be set globally by using `bl.NewWithConstraints(width, height PreferenceGroup)` or on each cell by using **BoundSize** notation. The string definition is compatible with MigLayout:

> A **bound size** is a size that optionally has a lower and/or upper bound and consists of one to three Unit Values. Practically it is a minimum/preferred/maximum size combination but none of the sizes are actually mandatory. If a size is missing (e.g. the preferred) it is null and will be replaced by the most appropriate value.
>
> The format is **"min:preferred:max"**, however there are shorter versions since for instance it is seldom needed to specify the maximum size.
>
> * A single value (E.g. **"10"**) sets only the preferred size and is exactly the same as "null:10:null" and **":10:"** and **"n:10:n"**.
> * Two values (E.g. **"10:20"**) means minimum and preferred size and is exactly the same as **"10:20:null"** and **"10:20:"** and **"10:20:n"**
> * The use a of an exclamation mark (E.g. **"20!"**) means that the value should be used for all size types and no colon may then be used in the string. It is the same as **"20:20:20"**.

All of this to say: yes, I have brought **null** to go. I've taken the liberty of supporting **nil** as well.

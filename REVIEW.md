# BubbleLayout Code Review & Refactoring Recommendations

This document summarizes findings and recommendations for `bubblelayout` based on idiomatic Go standards, function purity, performance, and code quality.

---

### 1. Function Purity & State Isolation

#### A. Closure Mutations in `computeDims` (`layout.go`)
- **Current Behavior**: `computeDims` uses nested closures (`growToPrefferred`, `growToMax`) that mutate state in outer function scope (`remainder`, `numToCompute`, `hasPref`, `hasGrow`, `noGrowNoPref`).
- **Problem**: Impure control flow makes tracking dimension allocation logic difficult and error-prone.
- **Recommendation**: Extract allocation passes into pure helper functions that accept explicit inputs (e.g. remaining size, current dimensions, item bounds) and return updated state explicitly.

#### B. Side-Effect Mutating Validation (`layout.go`)
- **Current Behavior**: `bl.Validate()` checks layout validity but also mutates struct fields (`bl.resizeCache`, `bl.hPref`, `bl.wPref`) as a side effect.
- **Problem**: `Validate()` is not a pure inspection method. If validation fails halfway, `bubbleLayout` is left in a partially mutated state.
- **Recommendation**: Decouple validation and compilation. Create a pure transformation function:
  ```go
  func compileGrid(g Grid, docks []layout) (Grid, PreferenceGroup, PreferenceGroup, error)
  ```
  `Validate()` should execute `compileGrid` without mutating state unless validation succeeds.

#### C. In-Place Slice Manipulation in `expandSpans` (`layout.go`)
- **Current Behavior**: `expandSpans` performs in-place slice slicing and insertions on row slices (`append(ret[rowIdx][:colIdx+i], ...)`).
- **Problem**: Slicing in-place during row iteration is state-heavy and susceptible to off-by-one errors.
- **Recommendation**: Construct target grid coordinate map explicitly rather than mutating slices mid-traversal.

---

### 2. Performance & Memory Allocations

#### A. Map Allocations on Layout Calculation (`layout.go`)
- **Current Behavior**: Every invocation of `computeDims` allocates 5 maps on the heap:
  ```go
  hasMin := make(map[int]struct{})
  hasMax := make(map[int]struct{})
  hasPref := make(map[int]struct{})
  hasGrow := make(map[int]struct{})
  noGrowNoPref := make(map[int]struct{})
  ```
- **Problem**: Layout calculation (triggered on every terminal `Resize`) creates heavy GC pressure.
- **Recommendation**: Replace maps with direct evaluation on `pg[idx]` or single boolean slices (`[]bool`). Eliminates heap allocations during size calculation.

#### B. Replace Custom `min`/`max` with Go 1.21+ Standard Library (`util.go`)
- **Current Behavior**: `util.go` implements custom variadic `min` and `max` functions.
- **Problem**: Redundant code that shadows built-ins available in modern Go.
- **Recommendation**: Update `go.mod` to Go 1.21+ (`go 1.21`). Remove `util.go` and use standard library built-in `min()` and `max()`.

---

### 3. Idiomatic Go Standards

#### A. Constant Naming (`layout.go`)
- **Current Behavior**: `NORTH`, `SOUTH`, `EAST`, `WEST` use C/Java `ALL_CAPS`.
- **Problem**: Non-idiomatic in Go.
- **Recommendation**: Use PascalCase: `North`, `South`, `East`, `West`.

#### B. Error Handling & Sentinel Errors (`layout.go`, `stringapi.go`)
- **Exported Sentinel Error**: Define `ErrViewNotRegistered = errors.New("view not registered")` for `BubbleLayoutMsg.Size()`.
- **Invalid `Sprintf` usage**: Replace `fmt.Sprintf("unknown constraint")` with `errors.New("unknown constraint")`.
- **Formatting**: Ensure error messages start with lowercase letters and omit trailing punctuation per Go style rules.

#### C. Copy-Paste Error String Bug (`layout.go`)
- **Location**: `layout.go:538`
- **Current Message**: `"width preferences do not match the cell height"`
- **Fix**: Correct message to `"width preferences do not match the cell width"`.

#### D. Typo Cleanup (`layout.go`)
- **Location**: `layout.go:88`
- **Fix**: Rename `growToPrefferred` to `growToPreferred`.

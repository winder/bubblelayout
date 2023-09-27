package bubblelayout

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var borderSizePattern = regexp.MustCompile(`^([\d]+):?([\d]+)?:?([\d]+)?(!)?$`)

// getNumbers returns all numbers from the slice until a non-numeric string is reached.
func getNumbers(str []string) []int {
	var result []int
	for _, str := range str {
		if num, err := strconv.Atoi(str); err == nil {
			result = append(result, num)
		} else {
			return result
		}
	}
	return result
}

func isCardinal(str string) bool {
	switch Cardinal(str) {
	case NORTH, SOUTH, EAST, WEST:
		return true
	default:
		return false
	}
}

// parseSize parses the BoundSize string.
// The format is "min:preferred:max", however there are shorter versions since for instance it is seldom needed to specify the maximum size.
//
// A single value (E.g. "10") sets only the preferred size and is exactly the same as "null:10:null" and ":10:" and "n:10:n".
// Two values (E.g. "10:20") means minimum and preferred size and is exactly the same as "10:20:null" and "10:20:" and "10:20:n"
// The use a of an exclamation mark (E.g. "20!") means that the value should be used for all size types and no colon may then be used in the string. It is the same as "20:20:20".
func parseSize(sz string) (BoundSize, error) {
	// normalize the inputLayout for some of the weirder options
	sz = strings.ReplaceAll(sz, "null", "0")
	sz = strings.ReplaceAll(sz, "n", "0")
	sz = strings.ReplaceAll(sz, "::", ":0:")
	if strings.HasPrefix(sz, ":") {
		sz = "0" + sz
	}
	if strings.HasSuffix(sz, ":") {
		sz = sz + "0"
	}

	// Parse out the parts. Results vary based on what matches there are.
	parts := borderSizePattern.FindStringSubmatch(sz)
	if parts == nil || parts[0] == "!" {
		return BoundSize{}, fmt.Errorf("invalid bound size '%s': did not match pattern", sz)
	}

	nums := getNumbers(parts[1:])
	exp := parts[4] == "!"

	if exp && len(nums) != 1 {
		return BoundSize{}, fmt.Errorf("invalid bound size '%s': use '!' with only one number", sz)
	}

	if exp {
		return BoundSize{Min: nums[0], Preferred: nums[0], Max: nums[0]}, nil
	}

	if len(nums) == 1 {
		return BoundSize{Preferred: nums[0]}, nil
	}

	if len(nums) == 2 {
		return BoundSize{Min: nums[0], Preferred: nums[1]}, nil
	}

	// The regex doesn't allow more than 3 numbers.
	return BoundSize{Min: nums[0], Preferred: nums[1], Max: nums[2]}, nil
}

type ErrStringLayout struct {
	msg   string
	input string
	err   error
}

func (e ErrStringLayout) Error() string {
	var suffix string
	if e.err != nil {
		suffix = fmt.Sprintf(": %s", e.err)
	}
	return fmt.Sprintf("string api conversion error for inputLayout '%s': %s%s", e.input, e.msg, suffix)
}

func (e ErrStringLayout) Unwrap() error {
	return e.err
}

func makeErrStringLayout(input, msg string, err error) ErrStringLayout {
	return ErrStringLayout{msg: msg, input: input, err: err}
}

func convertToLayout(input string) (layout, error) {
	if input == "" {
		return layout{}, nil
	}

	var result layout
	declarations := strings.Split(input, ",")
	for _, declaration := range declarations {
		parts := strings.Fields(declaration)
		last := len(parts) == 1
		part := parts[0]
		switch part {
		case "wrap":
			result.wrap = true
		case "span":
			nums := getNumbers(parts[1:])
			if len(nums) > 2 || len(nums) == 0 {
				return layout{}, makeErrStringLayout(input, fmt.Sprintf("wrong number of inputs to span, expected 1 or 2 received '%v'", nums), nil)
			}
			if len(nums) > 0 {
				result.SpanWidth = nums[0]
			}
			if len(nums) > 1 {
				result.SpanHeight = nums[1]
			}
		case "spanw", "spanx", "sx":
			nums := getNumbers(parts[1:])
			if len(nums) != 1 {
				return layout{}, makeErrStringLayout(input, fmt.Sprintf("wrong number of inputs, expected 1 received '%v'", nums), nil)
			}
			result.SpanWidth = nums[0]
		case "spanh", "spany", "sy":
			nums := getNumbers(parts[1:])
			if len(nums) != 1 {
				return layout{}, makeErrStringLayout(input, fmt.Sprintf("wrong number of inputs, expected 1 received '%v'", nums), nil)
			}
			result.SpanHeight = nums[0]
		case "grow":
			result.GrowWidth = true
			result.GrowHeight = true
		case "groww", "growx":
			result.GrowWidth = true
		case "growh", "growy":
			result.GrowHeight = true
		case "dock", string(NORTH), string(SOUTH), string(EAST), string(WEST):
			offset := 0
			// dock is optional
			if part == "dock" {
				if last {
					return layout{}, makeErrStringLayout(input, "dock direction is missing", nil)
				}
				if !isCardinal(parts[1]) {
					return layout{}, makeErrStringLayout(input, "invalid cardinal direction", nil)
				}
				offset++
			}
			result.Cardinal = Cardinal(parts[offset])
			offset++
			// size is optional
			if offset < len(parts) {
				bound, err := parseSize(parts[offset])
				if err == nil {
					result.Min = bound.Min
					result.Preferred = bound.Preferred
					result.Max = bound.Max
				}
			}
		case "width", "w":
			if last {
				return layout{}, makeErrStringLayout(input, "width bound size is missing", nil)
			}
			bound, err := parseSize(parts[1])
			if err != nil {
				return layout{}, makeErrStringLayout(input, "unable to parse bound", err)
			}
			result.MinWidth = bound.Min
			result.PreferredWidth = bound.Preferred
			result.MaxWidth = bound.Max
		case "height", "h":
			if last {
				return layout{}, makeErrStringLayout(input, "height bound size is missing", nil)
			}
			bound, err := parseSize(parts[1])
			if err != nil {
				return layout{}, makeErrStringLayout(input, "unable to parse bound", err)
			}
			result.MinHeight = bound.Min
			result.PreferredHeight = bound.Preferred
			result.MaxHeight = bound.Max
		default:
			return layout{}, makeErrStringLayout(input, fmt.Sprintf("unknown constraint"), nil)
		}
	}

	// TODO: is it an error to have a Cell and a Dock?
	return result, nil
}

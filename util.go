package bubblelayout

func min(left int, args ...int) int {
	ret := left
	for _, v := range args {
		if v < ret {
			ret = v
		}
	}
	return ret
}

func max(left int, args ...int) int {
	ret := left
	for _, v := range args {
		if v > ret {
			ret = v
		}
	}
	return ret
}

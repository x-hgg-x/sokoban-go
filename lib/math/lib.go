package math

// Min returns the minimum between 2 integers
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Max returns the maximum between 2 integers
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

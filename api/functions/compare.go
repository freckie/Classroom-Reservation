package functions

// InRange returns if the target is in the `range (start, end)
func InRange(target, start, end int) bool {
	return (start <= target) && (target <= end)
}

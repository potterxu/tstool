package algo

// Find largest index that data[index] <= target
func BinarySearchNoLargerThan(data []int64, target int64) int {
	start := 0
	end := len(data) - 1
	for start <= end {
		mid := (end-start)/2 + start
		if data[mid] > target {
			end = mid - 1
		} else {
			start = mid + 1
		}
	}
	return end
}

package algo

// Find index with largest data[index] with data[index] <= target
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

// Find index with smallest data[index] with data[index] >= target
func BinarySearchNoSmallerThan(data []int64, target int64) int {
	start := 0
	end := len(data) - 1
	for start <= end {
		mid := (end-start)/2 + start
		if data[mid] < target {
			start = mid + 1
		} else {
			end = mid - 1
		}
	}
	return start
}

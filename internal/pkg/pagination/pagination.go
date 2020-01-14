package pagination

import "math"

func Pages(count int64, size int64) int {
	pages := int(math.Ceil(float64(count) / float64(size)))
	if count == 0 {
		return 1
	}
	return pages
}

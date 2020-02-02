package mock

import "time"

// TestTime is used for testing time fields
func TestTime(year int) time.Time {
	return time.Date(year, time.May, 19, 1, 2, 3, 4, time.UTC)
}

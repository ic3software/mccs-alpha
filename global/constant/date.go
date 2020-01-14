package constant

import "time"

var Date = struct {
	DefaultFrom time.Time
	DefaultTo   time.Time
}{
	DefaultFrom: time.Date(2000, 1, 1, 00, 00, 0, 0, time.UTC),
	DefaultTo:   time.Date(2100, 1, 1, 00, 00, 0, 0, time.UTC),
}

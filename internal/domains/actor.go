package domains

import "time"

type Actor struct {
	ID       uint32
	FullName string
	Gender   string
	Birthday time.Time
}

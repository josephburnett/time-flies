package types

import (
	"time"
)

type Log []*Week

type Week struct {
	Date time.Time
	Body []*Entry
}

type Entry struct {
	Line   string
	Labels map[string]string
}

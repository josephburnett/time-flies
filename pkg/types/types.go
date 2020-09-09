package types

import (
	"time"
)

type Log []*Week

type Week struct {
	Date   time.Time
	Header map[string][]string
	Done   []*Entry
	Todo   []*Entry
}

type Entry struct {
	Line   string
	Labels map[string]string
}

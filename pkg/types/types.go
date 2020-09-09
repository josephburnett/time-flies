package types

import (
	"time"
)

type Log []*Week

type Week struct {
	Date   time.Time
	Header map[string][]string
	Body   []*Entry
}

type Entry struct {
	Line   string
	Labels map[string]string
}

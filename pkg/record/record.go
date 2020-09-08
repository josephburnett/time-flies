package record

import (
	"fmt"
	"net/mail"
	"strings"
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

func ParseLog(s string) (Log, error) {
	log := []*Week{}
	for _, stanza := range strings.Split(s, "%%\n") {
		week, err := ParseWeek(stanza)
		if err != nil {
			return nil, err
		}
		log = append(log, week)
	}
	return log, nil
}

func parseDate(s string) (time.Time, error) {
	layout := "Jan 2, 2006"
	return time.Parse(layout, s)
}

func ParseWeek(s string) (*Week, error) {
	message, err := mail.ReadMessage(s.NewReader(stanza))
	if err != nil {
		return err
	}
	date, ok := message.Header["Date"]
	if err != nil {
		return fmt.Errorf("missing required 'Date' header:\n%v", stanza)
	}
	t, err := parseDate(date)
	if err != nil {
		return err
	}
	week := &Week{
		Date: t,
		Body: []*Entry{},
	}
	for _, line := range strings.Split(s, "\n") {
		entry, err := ParseEntry(line)
		if err != nil {
			return nil, err
		}
		week.Body = append(week.Body, entry)
	}
	return week, nil
}

func ParseEntry(s string) (*Entry, error) {
	// Collapse whitespace
	entry := &Entry{
		Labels: map[string]string{},
	}
	components := strings.Split(s, "##")
	if len(components) == 1 {
		entry.Labels = s
		return entry, nil
	}
	// Slice out labels and parse
}

package parse

import (
	"fmt"
	"io/ioutil"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"github.com/josephburnett/time-flies/pkg/types"
)

func ParseLog(recordJar string) (types.Log, error) {
	log := []*types.Week{}
	for _, record := range strings.Split(recordJar, "%%\n") {
		week, err := ParseWeek(record)
		if err != nil {
			return nil, err
		}
		log = append(log, week)
	}
	return log, nil
}

func ParseWeek(record string) (*types.Week, error) {
	message, err := mail.ReadMessage(strings.NewReader(record))
	if err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(message.Body)
	date, ok := message.Header["Date"]
	if !ok {
		return nil, fmt.Errorf("missing required 'Date' header:\n%v", record)
	}
	if len(date) > 1 {
		return nil, fmt.Errorf("duplicate 'Date' header:\n%v", record)
	}
	header := message.Header
	delete(header, "Date")
	t, err := time.Parse("January 2, 2006", date[0])
	if err != nil {
		return nil, fmt.Errorf("invalid date 'January 2, 2006' date format: %v", date)
	}
	week := &types.Week{
		Date:   t,
		Header: message.Header,
		Done:   []*types.Entry{},
		Todo:   []*types.Entry{},
	}
	for _, line := range strings.Split(string(body), "\n") {
		line = dewhite(line)
		if line == "" {
			continue
		}
		entry, done, err := ParseEntry(line)
		if err != nil {
			return nil, err
		}
		if done {
			week.Done = append(week.Done, entry)
		} else {
			week.Todo = append(week.Todo, entry)
		}
	}
	return week, nil
}

func ParseEntry(line string) (*types.Entry, bool, error) {
	done := true
	line = dewhite(line)
	if len(line) > 1 && line[0] == '#' {
		done = false
		line = strings.TrimSpace(line)
	}
	components := strings.Split(line, "##")
	if len(components) == 1 {
		return &types.Entry{
			Line:   line,
			Labels: map[string]string{},
		}, done, nil
	}
	last := components[len(components)-1]
	cut := len(line) - len(last)
	labels, err := parseLabels(line[cut:])
	if err != nil {
		return nil, false, err
	}
	line = line[:cut-2]
	line = dewhite(line)
	return &types.Entry{
		Line:   line,
		Labels: labels,
	}, done, nil
}

func parseLabels(line string) (map[string]string, error) {
	line = dewhite(line)
	labels := map[string]string{}
	for _, pair := range strings.Split(line, " ") {
		if pair == "" {
			continue
		}
		kv := strings.Split(pair, "=")
		if len(kv) != 2 {
			return nil, fmt.Errorf("malformed 'k=v' labels: %q", pair)
		}
		labels[kv[0]] = kv[1]
	}
	return labels, nil
}

func dewhite(s string) string {
	s = regexp.MustCompile(`\s+`).ReplaceAllString(s, " ")
	return strings.TrimSpace(s)
}

package file

import (
	"fmt"
	"io/ioutil"
	"net/mail"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/josephburnett/time-flies/pkg/types"
)

const (
	defaultLogFile = ".tf/log"
)

type FileConfig struct {
	LogFile  *string
	OrgFiles []string
}

func (c *FileConfig) GetLogFile() string {
	if c == nil || c.LogFile == nil {
		home := os.Getenv("HOME")
		return fmt.Sprintf("%v/%v", home, defaultLogFile)
	}
	return *c.LogFile
}

func (c *FileConfig) Read() (types.Log, error) {
	l, err := c.ReadLog()
	if err != nil {
		return nil, err
	}
	o, err := c.ReadOrg()
	if err != nil {
		return nil, err
	}
	return mergeLogs(l, o), nil
}

func (c *FileConfig) ReadLog() (types.Log, error) {
	bs, err := ioutil.ReadFile(c.GetLogFile())
	if err != nil {
		return nil, err
	}
	return c.ParseLog(string(bs))
}

func (c *FileConfig) ParseLog(recordJar string) (types.Log, error) {
	log := []*types.Week{}
	for _, record := range strings.Split(recordJar, "%%\n") {
		week, err := c.ParseWeek(record)
		if err != nil {
			return nil, err
		}
		log = append(log, week)
	}
	return log, nil
}

func (c *FileConfig) ParseWeek(record string) (*types.Week, error) {
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
	t, err := c.parseDate(date[0])
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
		line = c.dewhite(line)
		if line == "" {
			continue
		}
		entry, done, err := c.ParseEntry(line)
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

func (c *FileConfig) parseDate(s string) (time.Time, error) {
	t, err := time.Parse("January 2, 2006", s)
	if err == nil {
		return t, nil
	}
	t, err = time.Parse("Jan 2, 2006", s)
	if err == nil {
		return t, nil
	}
	t, err = time.Parse("Jan 2 2006", s)
	if err == nil {
		return t, nil
	}
	return time.Time{}, fmt.Errorf("could not parse date: %v", s)
}

func (c *FileConfig) isDone(line string) (todo bool, rest string) {
	if strings.HasPrefix(line, "[ ]") {
		return false, strings.TrimSpace(line[3:])
	}
	if strings.HasPrefix(line, "[x]") {
		return true, strings.TrimSpace(line[3:])
	}
	if strings.HasPrefix(line, "#") {
		return false, strings.TrimSpace(line[1:])
	}
	return true, line
}

func (c *FileConfig) ParseEntry(line string) (*types.Entry, bool, error) {
	done, line := c.isDone(c.dewhite(line))
	components := strings.Split(line, "##")
	if len(components) == 1 {
		return &types.Entry{
			Line:   line,
			Labels: map[string]string{},
		}, done, nil
	}
	last := components[len(components)-1]
	cut := len(line) - len(last)
	labels, err := c.parseLabels(line[cut:])
	if err != nil {
		return nil, false, err
	}
	line = line[:cut-2]
	line = c.dewhite(line)
	return &types.Entry{
		Line:   line,
		Labels: labels,
	}, done, nil
}

func (c *FileConfig) parseLabels(line string) (map[string]string, error) {
	line = c.dewhite(line)
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

func (c *FileConfig) dewhite(s string) string {
	s = regexp.MustCompile(`\s+`).ReplaceAllString(s, " ")
	return strings.TrimSpace(s)
}

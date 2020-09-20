package tidy

import (
	"fmt"
	"sort"
	"strings"

	"github.com/josephburnett/time-flies/pkg/types"
)

type TidyConfig struct{}

func (c *TidyConfig) SprintLog(log types.Log) (string, error) {
	sort.Slice(log, func(i, j int) bool { return log[i].Date.After(log[j].Date) })
	out := ""
	for i, week := range log {
		s, err := c.printWeek(week)
		if err != nil {
			return "", err
		}
		out += s
		if i < len(log)-1 {
			out += "%%\n"
		}
	}
	return out, nil
}

func (c *TidyConfig) printWeek(week *types.Week) (string, error) {
	out := fmt.Sprintf("Date: %v\n", week.Date.Format("Jan 02 2006"))
	for k, vs := range week.Header {
		for _, v := range vs {
			out += fmt.Sprintf("%v: %v\n", k, v)
		}
	}
	out += "\n"
	maxWidth := 0
	for _, entry := range week.Done {
		if len(entry.Line) > maxWidth {
			maxWidth = len(entry.Line)
		}
	}
	for _, entry := range week.Todo {
		if len(entry.Line)+2 > maxWidth {
			maxWidth = len(entry.Line) + 2
		}
	}
	for _, entry := range week.Done {
		out += entry.Line
		out += strings.Repeat(" ", maxWidth-len(entry.Line))
		s, err := c.printLabels(entry.Labels)
		if err != nil {
			return "", err
		}
		out += fmt.Sprintf("  ##%v\n", s)
	}
	for _, entry := range week.Todo {
		out += fmt.Sprintf("# %v", entry.Line)
		out += strings.Repeat(" ", maxWidth-len(entry.Line)-2)
		s, err := c.printLabels(entry.Labels)
		if err != nil {
			return "", err
		}
		out += fmt.Sprintf("  ##%v\n", s)
	}
	out += "\n"
	return out, nil
}

func (c *TidyConfig) printLabels(labels map[string]string) (string, error) {
	out := ""
	for k, v := range labels {
		out += fmt.Sprintf(" %v=%v", k, v)
	}
	return out, nil
}

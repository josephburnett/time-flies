package view

import (
	"fmt"
	"sort"
	"strings"

	"github.com/josephburnett/time-flies/pkg/budget"
	"github.com/josephburnett/time-flies/pkg/types"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
	colorGrey   = "\033[90m"
)

var colorIndex = []string{
	colorRed,
	colorGreen,
	colorYellow,
	colorBlue,
	colorPurple,
	colorCyan,
}

const (
	defaultScreenWidth = 100
	minScreenWidth     = 20
)

type ViewConfig struct {
	ScreenWidth *int
	FocusGroup  *string
}

func (c *ViewConfig) screenWidth() int {
	if c == nil || c.ScreenWidth == nil {
		return defaultScreenWidth
	}
	w := *c.ScreenWidth
	if w < minScreenWidth {
		return minScreenWidth
	}
	return w
}

func (c *ViewConfig) focusGroup() string {
	if c == nil || c.FocusGroup == nil {
		return ""
	}
	return *c.FocusGroup
}

func (c *ViewConfig) SprintTodo(log types.Log) (string, error) {
	out := ""
	for _, week := range log {
		if len(week.Todo) > 0 {
			out += fmt.Sprintf("%v\n", week.Date.Format("Jan 02 2006"))
			for _, entry := range week.Todo {
				out += fmt.Sprintf("[ ] %v\n", entry.Line)
			}
		}
	}
	return out, nil
}

func (c *ViewConfig) SprintTotals(totals budget.Totals) (string, error) {
	topLevelTotals := totals
	if c.focusGroup() != "" {
		focusedTotals, err := totals.Focus(c.focusGroup())
		if err != nil {
			return "", err
		}
		topLevelTotals = totals
		totals = focusedTotals
	}
	uniqueValues := map[string]bool{}
	for _, t := range totals {
		for _, s := range t.SubTotals {
			uniqueValues[s.Value] = true
		}
	}
	sortedValues := []string{}
	for v := range uniqueValues {
		sortedValues = append(sortedValues, v)
	}
	sort.Strings(sortedValues)
	out := ""
	for i, total := range totals {
		topTotal := topLevelTotals[i]
		line, err := c.sprintTotal(total, topTotal, sortedValues)
		if err != nil {
			return "", nil
		}
		out += line + "\n"
	}
	return out, nil
}

func (c *ViewConfig) sprintTotal(total, topTotal *budget.Total, values []string) (string, error) {
	screenWidth := float64(c.screenWidth())
	if c.focusGroup() != "" {
		screenWidth = screenWidth / 2
	}
	widthByValue := map[string]float64{}
	for _, sub := range total.SubTotals {
		widthByValue[sub.Value] = sub.Relative * screenWidth
	}
	out := fmt.Sprintf(" %v   |", total.Date.Format("Jan 02 2006"))
	var cursor float64
	i := 0
	for _, value := range values {
		width := widthByValue[value]
		var color string
		if value == "" {
			color = colorGrey
			value = "?"
		} else {
			color = colorIndex[i%len(colorIndex)]
			i++
		}
		chars := int(cursor+width) - int(cursor)
		cursor += width
		if len(value) > chars {
			value = value[:chars]
		}
		pad := chars - len(value)
		out += color
		out += strings.Repeat("-", pad/2)
		out += value
		out += strings.Repeat("-", pad/2)
		if pad%2 == 1 {
			out += "-"
		}
		out += colorReset
		out += "|"
	}
	if len(out)%2 == 1 {
		// TODO: fix the floating point error that makes this necessary
		out += " "
	}
	out += fmt.Sprintf("  (%.1fd) ", total.Absolute.Hours()/8)
	if total != topTotal {
		out += " |"
		var topTotalWidth int
		for _, s := range topTotal.SubTotals {
			if s.Value == c.focusGroup() {
				topTotalWidth = int(s.Relative * float64(screenWidth))
				out += strings.Repeat("-", topTotalWidth)
			}
		}
		out += "|"
		out += strings.Repeat(" ", int(screenWidth)-topTotalWidth)
		out += "|"
	}
	return out, nil
}

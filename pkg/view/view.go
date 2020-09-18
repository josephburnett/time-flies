package view

import (
	"fmt"
	"sort"
	"strings"

	"github.com/josephburnett/time-flies/pkg/budget"
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
)

var colorIndex = []string{
	colorRed,
	colorGreen,
	colorYellow,
	colorBlue,
	colorPurple,
	colorCyan,
	colorWhite,
}

type Config struct {
	ScreenWidth int
}

var defaultConfig = &Config{
	ScreenWidth: 80,
}

func PrintTotals(totals budget.Totals) (string, error) {
	return defaultConfig.PrintTotals(totals)
}

func (c *Config) PrintTotals(totals budget.Totals) (string, error) {
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
	for _, total := range totals {
		line, err := c.printTotal(total, sortedValues)
		if err != nil {
			return "", nil
		}
		out += line + "\n"
	}
	return out, nil
}

func (c *Config) printTotal(total *budget.Total, values []string) (string, error) {
	screenWidth := float64(c.ScreenWidth)
	widthByValue := map[string]float64{}
	for _, sub := range total.SubTotals {
		widthByValue[sub.Value] = sub.Relative * screenWidth
	}
	out := fmt.Sprintf(" %v |", total.Date.Format("Jan 02 2006"))
	var cursor float64
	for i, value := range values {
		width := widthByValue[value]
		chars := int(cursor+width) - int(cursor)
		cursor += width
		if len(value) > chars {
			value = value[:chars]
		}
		pad := chars - len(value)
		color := colorIndex[i%len(colorIndex)]
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
	return out, nil
}

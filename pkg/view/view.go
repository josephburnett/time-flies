package view

import (
	"sort"
	"strings"

	"github.com/josephburnett/time-flies/pkg/budget"
)

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
	out := ""
	for _, total := range totals {
		line, err := c.printTotal(total)
		if err != nil {
			return "", nil
		}
		out += line + "\n"
	}
	return out, nil
}

func (c *Config) printTotal(total *budget.Total) (string, error) {
	screenWidth := float64(c.ScreenWidth)
	widthByValue := map[string]int{}
	sortedValues := []string{}
	for _, sub := range total.SubTotals {
		value := sub.Value
		sortedValues = append(sortedValues, value)
		width := int(sub.Relative * screenWidth)
		widthByValue[value] = width
	}
	sort.Strings(sortedValues)
	out := "|"
	for _, value := range sortedValues {
		width := widthByValue[value]
		if len(value) > width {
			value = value[:width]
		}
		pad := width - len(value)
		out += strings.Repeat("-", pad/2)
		out += value
		out += strings.Repeat("-", pad/2)
		out += "|"
	}
	return out, nil
}

package view

import (
	"fmt"
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
	uniqueValues := map[string]bool{}
	for _, sub := range total.SubTotals {
		value := sub.Value
		uniqueValues[value] = true
		width := int(sub.Relative * screenWidth)
		widthByValue[value] = width
	}
	sortedValues := []string{}
	for value := range uniqueValues {
		sortedValues = append(sortedValues, value)
	}
	sort.Strings(sortedValues)
	out := fmt.Sprintf(" %v |", total.Date.Format("Jan 02 2006"))
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

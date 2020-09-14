package budget

import (
	"time"

	"github.com/josephburnett/time-flies/pkg/types"
)

type Totals []*Total

type Total struct {
	Date      time.Time
	Period    Period
	Absolute  time.Duration
	SubTotals []*SubTotal
}

type SubTotal struct {
	Label     string
	Value     string
	Relative  float64
	Absolute  time.Duration
	SubTotals []*SubTotal
}

type Config struct {
	Period      Period
	DaysPerWeek int
	HoursPerDay int
	Grouping    []string
}

type Period string

const (
	Weekly    Period = "Weekly"
	Monthly          = "Monthly"
	Quarterly        = "Quarterly"
)

var defaultConfig = &Config{
	Period:      Weekly,
	DaysPerWeek: 5,
	HoursPerDay: 8,
	Grouping: []string{
		"cat",
		"sub",
	},
}

func GetTotals(log types.Log) (Totals, error) {
	return defaultConfig.GetTotals(log)
}

func (c *Config) GetTotals(log types.Log) (Totals, error) {
	totals := make(Totals, 0)
	for _, week := range log {
		total, err := c.getTotal(week)
		if err != nil {
			return nil, err
		}
		totals = append(totals, total)
	}
	return totals, nil
}

func (c *Config) getTotal(week *types.Week) (*Total, error) {
	total := &Total{
		Date: week.Date,
		// assuming weekly period
		Period:   Weekly,
		Absolute: time.Duration(c.DaysPerWeek) * time.Duration(c.HoursPerDay) * time.Hour,
	}
	subTotals, err := c.getSubTotals(1, week)
	if err != nil {
		return nil, err
	}
	total.SubTotals = subTotals
	return total, nil
}

func (c *Config) getSubTotals(groupingLevel int, week *types.Week) ([]*SubTotal, error) {
	if groupingLevel > len(c.Grouping) {
		return []*SubTotal{}, nil
	}
	key := c.Grouping[groupingLevel-1]
	// TODO: size by grouping above (for sub-subtotals)
	var relativePerEntry float64
	var absolutePerEntry time.Duration
	if len(week.Done) > 0 {
		relativePerEntry = 1.0 / float64(len(week.Done))
		absolutePerEntry = time.Duration(c.DaysPerWeek) * time.Duration(c.HoursPerDay) * time.Hour / time.Duration(len(week.Done))
	}
	subTotalsByValue := map[string]*SubTotal{}
	for _, entry := range week.Done {
		value, ok := entry.Labels[key]
		if !ok {
			value = ""
		}
		s, ok := subTotalsByValue[value]
		if !ok {
			s = &SubTotal{
				Label:     key,
				Value:     value,
				SubTotals: []*SubTotal{},
			}
			subTotalsByValue[value] = s
		}
		s.Relative += relativePerEntry
		s.Absolute += absolutePerEntry
	}
	// TODO: calculate sub-subtotals
	subTotals := []*SubTotal{}
	for _, s := range subTotalsByValue {
		subTotals = append(subTotals, s)
	}
	return subTotals, nil
}

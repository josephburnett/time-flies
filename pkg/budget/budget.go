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
	Period   Period
	Grouping []string
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

func GetTotals(log Log) (Totals, error) {
	return defaultConfig.GetTotals(log)
}

func (c *Config) GetTotals(log Log) (Totals, error) {
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
	subTotals, err := c.getSubTotals(0, week*types.Week)
	if err != nil {
		return nil, err
	}
	total.SubTotals = subTotals
	return total, nil
}

func (c *Config) getSubTotals(groupingLevel int, week *types.Week) ([]*SubTotal, error) {
	key := c.Grouping[groupingLevel]
	subTotalsByValue := map[string]*SubTotal{}
	for _, entry := range week.Done {
		value, ok := entry.Labels[key]
		if !ok {
			value = ""
		}
		s, ok := subTotalsByValue[value]
		if !ok {
			s = &SubTotal{
				Label: key,
				Value: value,
			}
			subTotalsByValue[value] = s
		}
		s.Relative += relativePerEntry // TODO: calculate this
		s.Absolute += relativePerEntry //
	}
	subTotals := []*SubTotal{}
	for _, s := range subTotalsByValue {
		subTotals = append(subTotals, s)
	}
	return subTotals, nil
}

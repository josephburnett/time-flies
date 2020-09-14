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
	Count     int
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
	subTotals, err := c.getSubTotals(1, 1.0, total.Absolute, week.Done)
	if err != nil {
		return nil, err
	}
	total.SubTotals = subTotals
	return total, nil
}

func (c *Config) getSubTotals(groupingLevel int, relative float64, absolute time.Duration, done []*types.Entry) ([]*SubTotal, error) {
	if groupingLevel > len(c.Grouping) {
		return []*SubTotal{}, nil
	}
	key := c.Grouping[groupingLevel-1]
	var relativePerEntry float64
	var absolutePerEntry time.Duration
	if len(done) > 0 {
		relativePerEntry = relative / float64(len(done))
		absolutePerEntry = absolute / time.Duration(len(done))
	}
	subTotalsByValue := map[string]*SubTotal{}
	doneByValue := map[string][]*types.Entry{}
	for _, entry := range done {
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
		s.Count += 1
		doneByValue[value] = append(doneByValue[value], entry)
	}
	subTotals := []*SubTotal{}
	for _, s := range subTotalsByValue {
		ss, err := c.getSubTotals(groupingLevel+1, s.Relative, s.Absolute, doneByValue[s.Value])
		if err != nil {
			return nil, err
		}
		s.SubTotals = ss
		subTotals = append(subTotals, s)
	}
	return subTotals, nil
}

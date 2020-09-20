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

type Period string

const (
	Weekly    Period = "Weekly"
	Monthly          = "Monthly"
	Quarterly        = "Quarterly"

	defaultEntryPeriod = Weekly
	defaultDaysPerWeek = 5
	defaultHoursPerDay = 8
)

var (
	defaultLabelGrouping = []string{
		"cat",
		"sub",
	}
)

type BudgetConfig struct {
	EntryPeriod   *Period
	DaysPerWeek   *int
	HoursPerDay   *int
	LabelGrouping []string
}

func (c *BudgetConfig) entryPeriod() Period {
	if c == nil || c.EntryPeriod == nil {
		return defaultEntryPeriod
	}
	return *c.EntryPeriod
}

func (c *BudgetConfig) daysPerWeek() int {
	if c == nil || c.DaysPerWeek == nil {
		return defaultDaysPerWeek
	}
	return *c.DaysPerWeek
}

func (c *BudgetConfig) hoursPerDay() int {
	if c == nil || c.HoursPerDay == nil {
		return defaultHoursPerDay
	}
	return *c.HoursPerDay
}

func (c *BudgetConfig) labelGrouping() []string {
	if c == nil || len(c.LabelGrouping) == 0 {
		return defaultLabelGrouping
	}
	return c.LabelGrouping
}

func (c *BudgetConfig) GetTotals(log types.Log) (Totals, error) {
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

func (c *BudgetConfig) getTotal(week *types.Week) (*Total, error) {
	total := &Total{
		Date: week.Date,
		// assuming weekly period
		Period:   Weekly,
		Absolute: time.Duration(c.daysPerWeek()) * time.Duration(c.hoursPerDay()) * time.Hour,
	}
	subTotals, err := c.getSubTotals(1, 1.0, total.Absolute, week.Done)
	if err != nil {
		return nil, err
	}
	total.SubTotals = subTotals
	return total, nil
}

func (c *BudgetConfig) getSubTotals(groupingLevel int, relative float64, absolute time.Duration, done []*types.Entry) ([]*SubTotal, error) {
	if groupingLevel > len(c.labelGrouping()) {
		return []*SubTotal{}, nil
	}
	key := c.labelGrouping()[groupingLevel-1]
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

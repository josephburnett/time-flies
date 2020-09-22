package budget

import (
	"fmt"
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

	defaultAggregationPeriod = Weekly
	defaultDaysPerWeek       = 5
	defaultHoursPerDay       = 8
)

var (
	defaultLabelGrouping = []string{
		"cat",
		"sub",
	}
)

type BudgetConfig struct {
	AggregationPeriod *Period
	DaysPerWeek       *int
	HoursPerDay       *int
	LabelGrouping     []string
}

func (c *BudgetConfig) aggregationPeriod() Period {
	if c == nil || c.AggregationPeriod == nil {
		return defaultAggregationPeriod
	}
	return *c.AggregationPeriod
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
	groups, err := c.groupTotals(totals)
	if err != nil {
		return nil, err
	}
	totals := make(Totals, 0)
	for _, ts := range groups {
		t, err := ts.Merge()
		if err != nil {
			return nil, err
		}
		totals = append(totals, t)
	}
	return totals, nil
}

func (c *BudgetConfig) getTotal(week *types.Week) (*Total, error) {
	total := &Total{
		Date:     week.Date,
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

func (c *BudgetConfig) groupTotals(totals Totals) ([]Totals, error) {
	totalsByTime := map[time.Time]Totals{}
	for _, total := range totals {
		t := c.roundToPeriod(total.Date)
		totalsByTime[t] = append(totalsByTime[t], total)
	}
	groups := []Totals{}
	for _, ts := range totalsByTime {
		groups = append(groups, ts)
	}
	return groups, nil
}

func (c *BudgetConfig) roundToPeriod(t time.Time) time.Time {
	var d time.Duration
	switch c.aggregationPeriod() {
	case Weekly:
		return t.Truncate(7 * 24 * time.Hour)
	case Monthly:
		return t.Truncate(30 * 24 * time.Hour)
	case Quarterly:
		return t.Truncate(90 * 24 * time.Hour)
	default:
		return t
	}
}

func (ts Totals) merge() (*Total, error) {
	if len(ts) == 0 {
		return nil, fmt.Errorf("Cannot merge empty Totals.")
	}
	label, value := totals[0].Label, totals[0].Value
	var (
		absoluteTotal time.Duration
		countTotal    int
	)
	for _, st := range subtotals {
		if st.Label != label {
			return nil, fmt.Errorf("Cannot merge SubTotals with different labels: %v and %v.", label, st.Label)
		}
		if st.Value != value {
			return nil, fmt.Errorf("Cannot merge SubTotals with different label values: %v and %v for %v.", value, st.Value, label)
		}
		absoluteTotal = absoluteTotal.Add(st.Absolute)
		countTotal += st.Count
	}
	relativePerSubtotal := relativeTotal / len(subs)
	s := &SubTotal{
		Relative: 1.0,
		Absolute: absoluteTotal,
		Count:    countTotal,
		// Punting on Sub-SubTotals.
	}
	return s, nil
}

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

type SubTotals []*SubTotal

type SubTotal struct {
	Label     string
	Value     string
	Relative  float64
	Absolute  time.Duration
	Count     int
	SubTotals SubTotals
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
	if c.aggregationPeriod() != Weekly {
		groups, err := c.groupTotals(totals)
		if err != nil {
			return nil, err
		}
		totals = make(Totals, 0)
		for date, ts := range groups {
			t, err := ts.mergeOn(date, c.aggregationPeriod())
			if err != nil {
				return nil, err
			}
			totals = append(totals, t)
		}
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

func (c *BudgetConfig) groupTotals(totals Totals) (map[time.Time]Totals, error) {
	totalsByTime := map[time.Time]Totals{}
	for _, total := range totals {
		t := c.roundToPeriod(total.Date)
		totalsByTime[t] = append(totalsByTime[t], total)
	}
	return totalsByTime, nil
}

func (c *BudgetConfig) roundToPeriod(t time.Time) time.Time {
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

func (ts Totals) mergeOn(date time.Time, period Period) (*Total, error) {
	if len(ts) == 0 {
		return nil, fmt.Errorf("Cannot merge empty Totals.")
	}
	total := &Total{
		Date:   date,
		Period: period,
	}
	ss := []SubTotals{}
	for _, t := range ts {
		total.Absolute = total.Absolute + t.Absolute
		ss = append(ss, t.SubTotals)
	}
	s, err := mergeByValue(ss, len(ts))
	if err != nil {
		return nil, err
	}
	total.SubTotals = s
	return total, nil
}

func mergeByValue(sss []SubTotals, lenTotals int) (SubTotals, error) {
	subTotals := make(SubTotals, 0)
	var label string
	subTotalsByValue := map[string]SubTotals{}
	for _, ss := range sss {
		for _, s := range ss {
			if label == "" {
				label = s.Label
			}
			if label != s.Label {
				return nil, fmt.Errorf("Cannot merge SubTotals with different labels: %v and %v", label, s.Label)
			}
			value := s.Value
			subTotalsByValue[value] = append(subTotalsByValue[value], s)
		}
	}
	for _, ss := range subTotalsByValue {
		s, err := ss.merge(lenTotals)
		if err != nil {
			return nil, err
		}
		subTotals = append(subTotals, s)
	}
	return subTotals, nil
}

func (ss SubTotals) merge(lenTotals int) (*SubTotal, error) {
	if len(ss) == 0 {
		return nil, fmt.Errorf("Cannot merge empty SubTotals.")
	}
	label, value := ss[0].Label, ss[0].Value
	subTotal := &SubTotal{
		Label: label,
		Value: value,
	}
	subSubTotals := make([]SubTotals, 0)
	for _, s := range ss {
		if s.Label != label {
			return nil, fmt.Errorf("Cannot merge SubTotals with different labels: %v and %v", s.Label, label)
		}
		if s.Value != value {
			return nil, fmt.Errorf("Cannot merge SubTotals with different values: %v and %v", s.Value, value)
		}
		subTotal.Absolute += s.Absolute
		subTotal.Count += s.Count
		subTotal.Relative += s.Relative
		subSubTotals = append(subSubTotals, s.SubTotals)
	}
	subTotals, err := mergeByValue(subSubTotals, lenTotals)
	if err != nil {
		return nil, err
	}
	subTotal.Relative = subTotal.Relative / float64(lenTotals)
	subTotal.SubTotals = subTotals
	return subTotal, nil
}

func (ts Totals) Focus(value string) (Totals, error) {
	focusedTotals := make(Totals, 0)
	for _, t := range ts {
		focusedTotal := &Total{
			Date:   t.Date,
			Period: t.Period,
		}
		focusedSubTotals := make(SubTotals, 0)
		for _, s := range t.SubTotals {
			if s.Value == value {
				for _, ss := range s.SubTotals {
					fs := &SubTotal{
						Label:    ss.Label,
						Value:    ss.Value,
						Relative: ss.Relative / s.Relative,
						Absolute: ss.Absolute,
						Count:    ss.Count,
					}
					focusedTotal.Absolute += ss.Absolute
					focusedSubTotals = append(focusedSubTotals, fs)
				}
			}
		}
		focusedTotal.SubTotals = focusedSubTotals
		focusedTotals = append(focusedTotals, focusedTotal)
	}
	return focusedTotals, nil
}

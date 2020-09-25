package budget

import (
	"testing"
	"time"

	"github.com/josephburnett/time-flies/pkg/types"
)

func TestGetTotal(t *testing.T) {
	cases := []struct {
		name string
		week *types.Week
		want *Total
	}{{
		name: "merge a week",
		week: &types.Week{
			Date:   time.Unix(0, 0),
			Header: map[string][]string{},
			Done: []*types.Entry{{
				Line: "thing one",
				Labels: map[string]string{
					"cat": "a",
					"sub": "1",
				},
			}, {
				Line: "thing two",
				Labels: map[string]string{
					"cat": "a",
					"sub": "1",
				},
			}, {
				Line: "thing three",
				Labels: map[string]string{
					"cat": "a",
					"sub": "2",
				},
			}, {
				Line: "thing four",
				Labels: map[string]string{
					"cat": "b",
					"sub": "2",
				},
			}},
			Todo: []*types.Entry{{}},
		},
		want: &Total{
			Date:     time.Unix(0, 0),
			Period:   Weekly,
			Absolute: 40 * time.Hour,
			SubTotals: []*SubTotal{{
				Label:    "cat",
				Value:    "a",
				Relative: 0.75,
				Absolute: 30 * time.Hour,
				Count:    3,
				SubTotals: []*SubTotal{{
					Label:    "sub",
					Value:    "1",
					Relative: 0.5,
					Absolute: 20 * time.Hour,
					Count:    2,
				}, {
					Label:    "sub",
					Value:    "2",
					Relative: 0.25,
					Absolute: 10 * time.Hour,
					Count:    1,
				}},
			}, {
				Label:    "cat",
				Value:    "b",
				Relative: 0.25,
				Absolute: 10 * time.Hour,
				Count:    1,
				SubTotals: []*SubTotal{{
					Label:    "sub",
					Value:    "2",
					Relative: 0.25,
					Absolute: 10 * time.Hour,
					Count:    1,
				}},
			}},
		},
	}}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var bc *BudgetConfig = nil
			got, err := bc.getTotal(c.week)
			if err != nil {
				t.Errorf("Wanted no error. Got %q", err)
			}
			if !got.equal(c.want) {
				t.Errorf("Wanted %+v. Got %+v", c.want, got)
			}
		})
	}
}

func TestTotalsMergeOn(t *testing.T) {
	cases := []struct {
		name      string
		totals    []*Total
		date      time.Time
		period    Period
		wantTotal *Total
		wantError bool
	}{{
		name:      "merge nil",
		totals:    nil,
		period:    Weekly,
		wantError: true,
	}, {
		name:      "merge empty",
		totals:    []*Total{},
		period:    Weekly,
		wantError: true,
	}, {
		name: "merge different labels",
		totals: []*Total{{
			Date:     time.Unix(0, 1),
			Period:   Weekly,
			Absolute: 2 * time.Hour,
			SubTotals: []*SubTotal{{
				Label:    "a",
				Value:    "1",
				Relative: 1.0,
				Absolute: time.Hour,
				Count:    1,
			}, {
				Label:    "b",
				Value:    "1",
				Relative: 1.0,
				Absolute: time.Hour,
				Count:    1,
			}},
		}},
		wantError: true,
	}, {
		name: "merge two totals with two subtotals each",
		totals: []*Total{{
			Date:     time.Unix(0, 1),
			Period:   Weekly,
			Absolute: 5 * 8 * time.Hour,
			SubTotals: []*SubTotal{{
				Label:    "a",
				Value:    "1",
				Relative: 0.5,
				Absolute: 20 * time.Hour,
				Count:    1,
			}, {
				Label:    "a",
				Value:    "2",
				Relative: 0.5,
				Absolute: 20 * time.Hour,
				Count:    1,
			}},
		}, {
			Date:     time.Unix(0, 2),
			Period:   Weekly,
			Absolute: 5 * 8 * time.Hour,
			SubTotals: []*SubTotal{{
				Label:    "a",
				Value:    "1",
				Relative: 0.2,
				Absolute: 8 * time.Hour,
				Count:    1,
			}, {
				Label:    "a",
				Value:    "2",
				Relative: 0.8,
				Absolute: 32 * time.Hour,
				Count:    1,
			}},
		}},
		date:   time.Unix(0, 3),
		period: Monthly,
		wantTotal: &Total{
			Date:     time.Unix(0, 3),
			Period:   Monthly,
			Absolute: 10 * 8 * time.Hour,
			SubTotals: []*SubTotal{{
				Label:    "a",
				Value:    "1",
				Relative: 0.35,
				Absolute: 28 * time.Hour,
				Count:    2,
			}, {
				Label:    "a",
				Value:    "2",
				Relative: 0.65,
				Absolute: 52 * time.Hour,
				Count:    2,
			}},
		},
	}}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			totals := Totals(c.totals)
			got, err := totals.mergeOn(c.date, c.period)
			if c.wantError && err == nil {
				t.Errorf("Wanted error. Got none.")
			}
			if !c.wantError && err != nil {
				t.Errorf("Wanted no error. Got %q", err)
			}
			if !got.equal(c.wantTotal) {
				t.Errorf("Wanted Total %+v. Got %+v", c.wantTotal, got)
			}
		})
	}
}

func TestSubTotalsMerge(t *testing.T) {
	cases := []struct {
		name         string
		subTotals    []*SubTotal
		wantSubTotal *SubTotal
		wantError    bool
	}{{
		name:      "merge nil",
		subTotals: nil,
		wantError: true,
	}, {
		name:      "merge empty",
		subTotals: []*SubTotal{},
		wantError: true,
	}, {
		name: "merge different labels",
		subTotals: []*SubTotal{
			{Label: "a"},
			{Label: "b"},
		},
		wantError: true,
	}, {
		name: "merge different values",
		subTotals: []*SubTotal{
			{Label: "a", Value: "1"},
			{Label: "a", Value: "2"},
		},
		wantError: true,
	}, {
		name: "merge compatible sub totals",
		subTotals: []*SubTotal{{
			Label:    "a",
			Value:    "1",
			Relative: 0.4,
			Absolute: time.Hour,
			Count:    1,
		}, {
			Label:    "a",
			Value:    "1",
			Relative: 0.6,
			Absolute: time.Hour,
			Count:    1,
		}},
		wantSubTotal: &SubTotal{
			Label:    "a",
			Value:    "1",
			Relative: 0.5,
			Absolute: 2 * time.Hour,
			Count:    2,
		},
	}}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			subTotals := SubTotals(c.subTotals)
			got, err := subTotals.merge()
			if c.wantError && err == nil {
				t.Errorf("Wanted error. Got none.")
			}
			if !c.wantError && err != nil {
				t.Errorf("Wanted no error. Got %q", err)
			}
			if !got.equal(c.wantSubTotal) {
				t.Errorf("Wanted SubTotal %+v. Got %+v", c.wantSubTotal, got)
			}
		})
	}
}

func (t1 *Total) equal(t2 *Total) bool {
	switch {
	case t1 == nil && t2 == nil:
		return true
	case t1 == nil && t2 != nil:
		return false
	case t1 != nil && t2 == nil:
		return false
	case !t1.Date.Equal(t2.Date):
		return false
	case t1.Period != t2.Period:
		return false
	case t1.Absolute != t2.Absolute:
		return false
	case len(t1.SubTotals) != len(t2.SubTotals):
		return false
	default:
		for i, ts1 := range t1.SubTotals {
			ts2 := t2.SubTotals[i]
			if !ts1.equal(ts2) {
				return false
			}
		}
		return true
	}
}

func (s1 *SubTotal) equal(s2 *SubTotal) bool {
	switch {
	case s1 == nil && s2 == nil:
		return true
	case s1 == nil && s2 != nil:
		return false
	case s2 != nil && s2 == nil:
		return false
	case s1.Label != s2.Label:
		return false
	case s1.Value != s2.Value:
		return false
	case s1.Relative != s2.Relative:
		return false
	case s1.Absolute != s2.Absolute:
		return false
	case s1.Count != s2.Count:
		return false
	case len(s1.SubTotals) != len(s2.SubTotals):
		return false
	default:
		for i, ss1 := range s1.SubTotals {
			ss2 := s2.SubTotals[i]
			if !ss1.equal(ss2) {
				return false
			}
		}
		return true
	}
}

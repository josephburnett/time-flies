package budget

import (
	"encoding/json"
	"sort"
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
				wantString := "nil"
				if c.want != nil {
					b, _ := json.Marshal(c.want)
					wantString = string(b)
				}
				gotString := "nil"
				if got != nil {
					b, _ := json.Marshal(got)
					gotString = string(b)
				}
				t.Errorf("Wanted %+v. Got %+v", wantString, gotString)
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
				wantString := "nil"
				if c.wantTotal != nil {
					b, _ := json.Marshal(c.wantTotal)
					wantString = string(b)
				}
				gotString := "nil"
				if got != nil {
					b, _ := json.Marshal(got)
					gotString = string(b)
				}
				t.Errorf("Wanted Total %+v. Got %+v", wantString, gotString)
			}
		})
	}
}

func TestSubTotalsMerge(t *testing.T) {
	cases := []struct {
		name         string
		lenTotals    int
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
			Relative: 1.0,
			Absolute: 2 * time.Hour,
			Count:    2,
		},
	}}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			subTotals := SubTotals(c.subTotals)
			lenTotals := c.lenTotals
			if c.lenTotals == 0 {
				lenTotals = 1
			}
			got, err := subTotals.merge(lenTotals)
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

func TestEntryTimes(t *testing.T) {
	cases := []struct {
		name     string
		relative float64
		absolute time.Duration
		done     []*types.Entry
		want     []*entryTime
		wantErr  bool
	}{{
		name:     "fuzzy expands",
		relative: 0.5,
		absolute: 8 * time.Hour,
		done: []*types.Entry{
			entry("", "1h"),
			entry("", "1h"),
		},
		want: []*entryTime{
			entryT(0.25, 0, 4*time.Hour),
			entryT(0.25, 0, 4*time.Hour),
		},
	}, {
		name:     "fuzzy contracts",
		relative: 0.5,
		absolute: 8 * time.Hour,
		done: []*types.Entry{
			entry("", "6h"),
			entry("", "6h"),
		},
		want: []*entryTime{
			entryT(0.25, 0, 4*time.Hour),
			entryT(0.25, 0, 4*time.Hour),
		},
	}, {
		name:     "strict expands",
		relative: 0.5,
		absolute: 8 * time.Hour,
		done: []*types.Entry{
			entry("1h", ""),
			entry("1h", ""),
		},
		want: []*entryTime{
			entryT(0.25, 4*time.Hour, 0),
			entryT(0.25, 4*time.Hour, 0),
		},
	}, {
		name:     "strict contracts",
		relative: 0.5,
		absolute: 8 * time.Hour,
		done: []*types.Entry{
			entry("6h", ""),
			entry("6h", ""),
		},
		want: []*entryTime{
			entryT(0.25, 4*time.Hour, 0),
			entryT(0.25, 4*time.Hour, 0),
		},
	}, {
		name:     "expands only fuzzy",
		relative: 1.0,
		absolute: 10 * time.Hour,
		done: []*types.Entry{
			entry("1h", ""),
			entry("", "1h"),
		},
		want: []*entryTime{
			entryT(0.1, time.Hour, 0),
			entryT(0.9, 0, 9*time.Hour),
		},
	}, {
		name:     "contracts only fuzzy",
		relative: 1.0,
		absolute: 10 * time.Hour,
		done: []*types.Entry{
			entry("9h", ""),
			entry("", "9h"),
		},
		want: []*entryTime{
			entryT(0.9, 9*time.Hour, 0),
			entryT(0.1, 0, time.Hour),
		},
	}}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			entryTimes, err := (*BudgetConfig)(nil).entryTimes(c.relative, c.absolute, c.done)
			if err != nil && !c.wantErr {
				t.Errorf("wanted no error. got %v", err)
			}
			if err == nil && c.wantErr {
				t.Errorf("wanted error. got nil")
			}
			if len(c.want) != len(entryTimes) {
				t.Errorf("wanted %v entries. got %v", len(c.want), len(entryTimes))
				return
			}
			for i, e := range entryTimes {
				if want := c.want[i].relative; want != e.relative {
					t.Errorf("[%v] wanted relative %v. got %v", i, want, e.relative)
				}
				if want := c.want[i].strict; want != e.strict {
					t.Errorf("[%v] wanted strict %v. got %v", i, want, e.strict)
				}
				if want := c.want[i].fuzzy; want != e.fuzzy {
					t.Errorf("[%v] wanted fuzzy %v. got %v", i, want, e.fuzzy)
				}
			}
		})
	}
}

func entry(strict, fuzzy string) *types.Entry {
	labels := map[string]string{}
	if strict != "" {
		labels["t"] = strict
	}
	if fuzzy != "" {
		labels["f"] = fuzzy
	}
	return &types.Entry{
		Labels: labels,
	}
}

func entryT(relative float64, strict, fuzzy time.Duration) *entryTime {
	return &entryTime{
		relative: relative,
		strict:   strict,
		fuzzy:    fuzzy,
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
		sort.Slice(t1.SubTotals, func(i, j int) bool { return t1.SubTotals[i].Label < t1.SubTotals[j].Label })
		sort.Slice(t2.SubTotals, func(i, j int) bool { return t2.SubTotals[i].Label < t2.SubTotals[j].Label })
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

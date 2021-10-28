package file

import (
	"time"

	"github.com/josephburnett/time-flies/pkg/types"
)

func mergeLogs(a, b types.Log) types.Log {
	weeks := map[time.Time]*types.Week{}
	for _, w1 := range a {
		key := w1.Date.Round(0)
		w2 := weeks[key]
		weeks[key] = mergeWeeks(w1, w2)
	}
	for _, w1 := range b {
		key := w1.Date.Round(0)
		w2 := weeks[key]
		weeks[key] = mergeWeeks(w1, w2)
	}
	var log types.Log
	for _, w := range weeks {
		log = append(log, w)
	}
	return log
}

func mergeWeeks(a, b *types.Week) *types.Week {
	ab := &types.Week{}
	if a != nil {
		ab.Header = mergeHeaders(ab.Header, a.Header)
		ab.Date = a.Date
		for _, d := range a.Done {
			ab.Done = append(ab.Done, d)
		}
		for _, t := range a.Todo {
			ab.Todo = append(ab.Todo, t)
		}
	}
	if b != nil {
		ab.Header = mergeHeaders(ab.Header, b.Header)
		ab.Date = b.Date
		for _, d := range b.Done {
			ab.Done = append(ab.Done, d)
		}
		for _, t := range b.Todo {
			ab.Todo = append(ab.Todo, t)
		}
	}
	return ab
}

func mergeHeaders(a, b map[string][]string) map[string][]string {
	ab := map[string][]string{}
	for ka, vsa := range a {
		var vsab []string
		copy(vsab, vsa)
		vsb := b[ka]
		for _, vb := range vsb {
			var have bool
			for _, vab := range vsab {
				if vab == vb {
					have = true
				}
			}
			if !have {
				vsab = append(vsab, vb)
			}
		}
		ab[ka] = vsab
	}
	return ab
}

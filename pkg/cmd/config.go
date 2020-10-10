package cmd

import (
	"github.com/josephburnett/time-flies/pkg/budget"
	"github.com/josephburnett/time-flies/pkg/file"
	"github.com/josephburnett/time-flies/pkg/tidy"
	"github.com/josephburnett/time-flies/pkg/view"
	flag "github.com/spf13/pflag"
)

type Config struct {
	budget.BudgetConfig
	file.FileConfig
	tidy.TidyConfig
	view.ViewConfig
}

var (
	focus  = flag.StringP("focus", "f", "", "Focus on a particular label group.")
	group  = flag.StringSliceP("group", "g", []string{}, "Group entries by labels.")
	log    = flag.StringP("log", "l", "", "Log file.")
	period = flag.StringP("period", "p", "", "Aggregation period.")
)

func getConfig() *Config {
	cfg := &Config{}
	if *focus != "" {
		cfg.ViewConfig.FocusGroup = focus
	}
	if len(*group) > 0 {
		cfg.BudgetConfig.LabelGrouping = *group
	}
	if *log != "" {
		cfg.FileConfig.LogFile = log
	}
	if *period != "" {
		budgetPeriod := budget.Period(*period)
		cfg.BudgetConfig.AggregationPeriod = &budgetPeriod
	}
	return cfg
}

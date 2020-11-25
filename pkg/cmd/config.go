package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

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
	config = flag.StringP("config", "c", "", "Config file. JSON serialization of pkg/cmd/Config.")
	focus  = flag.StringP("focus", "f", "", "Focus on a particular label group.")
	group  = flag.StringSliceP("group", "g", []string{}, "Group entries by labels.")
	log    = flag.StringP("log", "l", "", "Log file.")
	period = flag.StringP("period", "p", "", "Aggregation period.")
)

const (
	defaultConfigFile = ".tf/config"
)

func getConfig() (*Config, error) {
	home := os.Getenv("HOME")
	configFile := fmt.Sprintf("%v/%v", home, defaultConfigFile)
	mustExist := false
	if *config != "" {
		configFile = *config
		mustExist = true
	}
	cfg := &Config{}
	b, err := ioutil.ReadFile(configFile)
	if mustExist && os.IsNotExist(err) {
		// We don't have a file and we need one
		return nil, err
	}
	if err == nil {
		// We have a file
		err = json.Unmarshal(b, cfg)
		if err != nil {
			return nil, fmt.Errorf("unable to read config file %q: %v", configFile, err)
		}
	}

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
	return cfg, nil
}

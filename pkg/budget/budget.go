package budget

import "time"

type Totals []*Total

type Total struct {
	Label     string
	Date      time.Time
	Relative  float64
	Absolute  time.Duration
	SubTotals map[string]*Total
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
	Period: Weekly,
	Grouping: []string{
		"cat",
		"sub",
	},
}

func GetTotals(log Log) (Totals, error) {
	return defaultConfig.GetTotals(log)
}

func (c *Config) GetTotals(log Log) (Totals, error) {

}

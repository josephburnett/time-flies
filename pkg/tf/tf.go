package tf

import (
	"fmt"

	"github.com/josephburnett/time-flies/pkg/budget"
	"github.com/josephburnett/time-flies/pkg/file"
	"github.com/josephburnett/time-flies/pkg/tidy"
	"github.com/josephburnett/time-flies/pkg/view"
	"github.com/spf13/cobra"
)

type Config struct {
	budget.BudgetConfig
	file.FileConfig
	tidy.TidyConfig
	view.ViewConfig
}

var cfg = &Config{}

var CmdTidy = &cobra.Command{
	Use:   "tidy [log file]",
	Short: "Reformats log to spark joy.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg.LogFile = &args[0]
		log, err := cfg.FileConfig.Read()
		if err != nil {
			return err
		}
		s, err := cfg.TidyConfig.SprintLog(log)
		if err != nil {
			return err
		}
		fmt.Print(s)
		return nil
	},
}

var CmdTotals = &cobra.Command{
	Use:   "tots [log file]",
	Short: "Output weekly focus totals.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg.LogFile = &args[0]
		log, err := cfg.FileConfig.Read()
		if err != nil {
			return err
		}
		tots, err := cfg.BudgetConfig.GetTotals(log)
		if err != nil {
			return err
		}
		s, err := cfg.ViewConfig.SprintTotals(tots)
		if err != nil {
			return err
		}
		fmt.Printf("%v\n", s)
		return nil
	},
}

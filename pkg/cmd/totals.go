package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
)

var CmdTotals = &cobra.Command{
	Use:   "tots",
	Short: "Output weekly focus totals.",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := getConfig()
		log, err := cfg.FileConfig.Read()
		if err != nil {
			return err
		}
		tots, err := cfg.BudgetConfig.GetTotals(log)
		if err != nil {
			return err
		}
		sort.Slice(tots, func(i, j int) bool { return tots[i].Date.Before(tots[j].Date) })
		s, err := cfg.ViewConfig.SprintTotals(tots)
		if err != nil {
			return err
		}
		fmt.Printf("%v\n", s)
		return nil
	},
}

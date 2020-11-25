package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
)

var CmdTodo = &cobra.Command{
	Use:   "todo",
	Short: "List TODO entries.",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := getConfig()
		if err != nil {
			return err
		}
		log, err := cfg.FileConfig.Read()
		if err != nil {
			return err
		}
		sort.Slice(log, func(i, j int) bool { return log[i].Date.Before(log[j].Date) })
		s, err := cfg.ViewConfig.SprintTodo(log)
		if err != nil {
			return err
		}
		fmt.Printf("%v", s)
		return nil
	},
}

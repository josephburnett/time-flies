package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var CmdTidy = &cobra.Command{
	Use:   "tidy",
	Short: "Reformats log to spark joy.",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := getConfig()
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

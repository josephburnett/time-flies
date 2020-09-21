package tf

import (
	"fmt"
	"os"
	"os/exec"
	"sort"

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

var CmdTidy = &cobra.Command{
	Use:   "tidy",
	Short: "Reformats log to spark joy.",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := &Config{}
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
	Use:   "tots",
	Short: "Output weekly focus totals.",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := &Config{}
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

var CmdEdit = &cobra.Command{
	Use:   "edit",
	Short: "Edit the log file.",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		editor := os.Getenv("EDITOR")
		if editor == "" {
			return fmt.Errorf("no EDITOR set")
		}
		cfg := &Config{}
		filename := cfg.FileConfig.GetLogFile()
		execCmd := exec.Command(editor, filename)
		execCmd.Stdin = os.Stdin
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr
		return execCmd.Run()
	},
}

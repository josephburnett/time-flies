package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/josephburnett/time-flies/pkg/budget"
	"github.com/josephburnett/time-flies/pkg/parse"
	"github.com/josephburnett/time-flies/pkg/view"

	"github.com/josephburnett/time-flies/pkg/tidy"
	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:   "tf",
		Short: "Time Flies (tf) is a tool for budgeting focus time.",
	}
	root.AddCommand(cmdTidy)
	root.AddCommand(cmdTotals)
	root.AddCommand(cmdTodo)
	root.AddCommand(cmdJson)
	root.Execute()
}

var cmdTidy = &cobra.Command{
	Use:   "tidy [log file]",
	Short: "Reformats log to spark joy.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bytes, err := ioutil.ReadFile(args[0])
		if err != nil {
			return err
		}
		log, err := parse.ParseLog(string(bytes))
		if err != nil {
			return err
		}
		s, err := tidy.PrintLog(log)
		if err != nil {
			return err
		}
		fmt.Print(s)
		return nil
	},
}

var cmdTotals = &cobra.Command{
	Use:   "tots [log file]",
	Short: "Output weekly focus totals.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bytes, err := ioutil.ReadFile(args[0])
		if err != nil {
			return err
		}
		log, err := parse.ParseLog(string(bytes))
		if err != nil {
			return err
		}
		tots, err := budget.GetTotals(log)
		if err != nil {
			return err
		}
		out, err := view.PrintTotals(tots)
		if err != nil {
			return err
		}
		fmt.Printf("%v\n", string(out))
		return nil
	},
}

var cmdTodo = &cobra.Command{
	Use:   "todo [log file]",
	Short: "Output the current week's TODO list.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bytes, err := ioutil.ReadFile(args[0])
		if err != nil {
			return err
		}
		log, err := parse.ParseLog(string(bytes))
		if err != nil {
			return err
		}
		if len(log) == 0 {
			return nil
		}
		last := log[len(log)-1]
		for _, entry := range last.Todo {
			fmt.Println(entry.Line)
		}
		return nil
	},
}

var cmdJson = &cobra.Command{
	Use:   "json [log file]",
	Short: "Output log file in JSON format.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		bytes, err := ioutil.ReadFile(args[0])
		if err != nil {
			return err
		}
		log, err := parse.ParseLog(string(bytes))
		if err != nil {
			return err
		}
		s, err := json.MarshalIndent(log, "", "  ")
		if err != nil {
			return err
		}
		fmt.Printf("%v\n", string(s))
		return nil
	},
}

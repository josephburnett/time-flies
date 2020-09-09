package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/josephburnett/time-flies/pkg/parse"
	"github.com/spf13/cobra"
)

func main() {
	root.AddCommand(cmdJson)
	root.Execute()
}

var root = &cobra.Command{
	Use:   "tf",
	Short: "Time Flies (tf) is a tool for budgeting focus time.",
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

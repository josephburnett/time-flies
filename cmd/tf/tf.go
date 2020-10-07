package main

import (
	"github.com/spf13/cobra"

	"github.com/josephburnett/time-flies/pkg/cmd"
)

func main() {
	root := &cobra.Command{
		Use:   "tf",
		Short: "Time Flies (tf) is a tool for budgeting focus time.",
	}
	root.AddCommand(cmd.CmdTidy)
	root.AddCommand(cmd.CmdTotals)
	root.AddCommand(cmd.CmdEdit)
	root.AddCommand(cmd.CmdTodo)
	root.Execute()
}

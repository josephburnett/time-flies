package main

import (
	"github.com/spf13/cobra"

	"github.com/josephburnett/time-flies/pkg/tf"
)

func main() {
	root := &cobra.Command{
		Use:   "tf",
		Short: "Time Flies (tf) is a tool for budgeting focus time.",
	}
	root.AddCommand(tf.CmdTidy)
	root.AddCommand(tf.CmdTotals)
	root.Execute()
}

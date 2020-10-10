package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var CmdEdit = &cobra.Command{
	Use:   "edit",
	Short: "Edit the log file.",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		editor := os.Getenv("EDITOR")
		if editor == "" {
			return fmt.Errorf("no EDITOR set")
		}
		cfg := getConfig()
		filename := cfg.FileConfig.GetLogFile()
		execCmd := exec.Command(editor, filename)
		execCmd.Stdin = os.Stdin
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr
		return execCmd.Run()
	},
}

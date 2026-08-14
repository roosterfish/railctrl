package cmd

import (
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the interlocking daemon",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	interlockingCmd.AddCommand(runCmd)
}

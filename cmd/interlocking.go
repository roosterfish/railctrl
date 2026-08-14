package cmd

import (
	"github.com/spf13/cobra"
)

var interlockingCmd = &cobra.Command{
	Use:   "interlocking",
	Short: "Run the interlocking components for track control",
}

func init() {
	rootCmd.AddCommand(interlockingCmd)

	interlockingCmd.PersistentFlags().String("config", "", "Layout configuration file")
}

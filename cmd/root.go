package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "railctrl",
	Short:        "Plan, schedule, control and drive model railway infrastructure",
	SilenceUsage: true,
}

// Execute the root command.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.railctrl.yaml)")
}

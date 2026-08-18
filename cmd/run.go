package cmd

import (
	"github.com/roosterfish/railctrl/internal/interlocking"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the interlocking daemon",
	RunE: func(cmd *cobra.Command, args []string) error {
		address, err := cmd.Flags().GetString("address")
		if err != nil {
			return err
		}

		configPath, err := cmd.Flags().GetString("config")
		if err != nil {
			return err
		}

		d, err := interlocking.NewDaemon(address, configPath)
		if err != nil {
			return err
		}

		return d.Start()
	},
}

func init() {
	interlockingCmd.AddCommand(runCmd)
}

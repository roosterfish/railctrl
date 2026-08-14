package cmd

import (
	"fmt"

	"githu.com/roosterfish/railctrl/internal/config"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate the interlocking layout",
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath, err := cmd.Flags().GetString("config")
		if err != nil {
			return err
		}

		layout, err := config.DecodeLayout(configPath)
		if err != nil {
			return err
		}

		printConfig, err := cmd.Flags().GetBool("print")
		if err != nil {
			return err
		}

		if printConfig {
			encodedLayout, err := config.EncodeLayout(layout)
			if err != nil {
				return err
			}

			fmt.Print(string(encodedLayout))
		}

		return nil
	},
}

func init() {
	interlockingCmd.AddCommand(validateCmd)

	validateCmd.Flags().Bool("print", false, "Print layout configuration file")
}

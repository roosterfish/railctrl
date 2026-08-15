package cmd

import (
	"fmt"

	"github.com/roosterfish/railctrl/internal/config"
	"github.com/spf13/cobra"
	"go.yaml.in/yaml/v4"
)

var routeCmd = &cobra.Command{
	Use:   "route <start_track> <end_track>",
	Args:  cobra.ExactArgs(2),
	Short: "Calculate a route from start to end track",
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath, err := cmd.Flags().GetString("config")
		if err != nil {
			return err
		}

		layout, err := config.DecodeLayout(configPath)
		if err != nil {
			return err
		}

		route, cost, err := layout.Dijkstra(args[0], args[1])
		if err != nil {
			return err
		}

		output := map[string]any{
			"cost":  cost,
			"route": route,
		}

		bytes, err := yaml.Marshal(output)
		if err != nil {
			return fmt.Errorf("failed marshalling route: %w", err)
		}

		fmt.Print(string(bytes))
		return nil
	},
}

func init() {
	interlockingCmd.AddCommand(routeCmd)
}

package cmd

import (
	"context"
	"fmt"

	"github.com/roosterfish/railctrl/api"
	"github.com/spf13/cobra"
	"go.yaml.in/yaml/v4"
)

var routeCmd = &cobra.Command{
	Use:   "route <start_track> <end_track>",
	Args:  cobra.ExactArgs(2),
	Short: "Calculate a route from start to end track",
	RunE: func(cmd *cobra.Command, args []string) error {
		address, err := cmd.Flags().GetString("address")
		if err != nil {
			return err
		}

		client := api.NewClient(address)
		route, err := client.GetRouteFind(context.Background(), args[0], args[1])
		if err != nil {
			return fmt.Errorf("failed finding route: %w", err)
		}

		bytes, err := yaml.Marshal(route)
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

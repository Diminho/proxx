package cmd

import (
	"github.com/spf13/cobra"
)

// Execute executes proxx
func Execute() error {
	command := &cobra.Command{
		Use:   "proxx",
		Short: "Proxx game",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}

	command.AddCommand(start())

	return command.Execute()
}

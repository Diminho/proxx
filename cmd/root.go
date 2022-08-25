package cmd

import (
	"github.com/spf13/cobra"
)

// Execute executes proxx
func Execute() error {
	command := &cobra.Command{
		Use:          "proxx",
		Short:        "Proxx game",
		SilenceUsage: true,
	}

	command.AddCommand(start())

	return command.Execute()
}

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	cobra.OnInitialize(initEntries)
}

// Execute starts Extensibility Vendor Service
func Execute() error {
	command := &cobra.Command{
		Use:   "proxx",
		Short: "proxx",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("ARGS: ", args)
			return cmd.Usage()
		},
	}

	command.AddCommand(start())

	return command.Execute()
}

func initEntries() {

	fmt.Println("STARTED!")

}

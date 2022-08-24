package cmd

import (
	"fmt"
	"os"

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

	//command.PersistentFlags().StringVarP(&configPath, "config", "c", "", "config file path")
	command.AddCommand()

	return command.Execute()
}

func initEntries() {

	fmt.Println("STARTED!")

}

func webServer() *cobra.Command {
	command := &cobra.Command{
		Use:   "web-server",
		Short: "Manage Web Server",
		Run: func(cmd *cobra.Command, _ []string) {
			var cellNumber int
			var blackholes int
			fmt.Print("Enter board size\n")

			fmt.Print("Enter N size: ")
			_, err := fmt.Scanln(&cellNumber)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(0)
			}

			fmt.Print("Enter number of black holes: ")
			_, err = fmt.Scanln(&blackholes)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(0)
			}

			for {

			}
		},
	}

	return command
}

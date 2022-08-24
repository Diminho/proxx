package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

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

func start() *cobra.Command {
	command := &cobra.Command{
		Use:   "start",
		Short: "Start the game",
		Run: func(cmd *cobra.Command, _ []string) {
			var cellNumber int
			var blackHoles int
			fmt.Print("Enter board size\n")

			fmt.Print("Enter N size: ")
			_, err := fmt.Scanln(&cellNumber)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(0)
			}

			fmt.Print("Enter number of black holes: ")
			_, err = fmt.Scanln(&blackHoles)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(0)
			}

			b := newBoard(cellNumber, blackHoles)
			b.printBoard()
			fmt.Println()
			b.printBoardStateLess()
			return
			for {
				var coordinatesInput string
				fmt.Fprintln(os.Stdout, "Enter board coordinates - row and column (two digits with space):")
				_, err = fmt.Scanln(&coordinatesInput)
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(0)
				}

				clickInput := strings.Split(coordinatesInput, " ")
				// check if user types more that two digits
				if len(clickInput) > 2 {
					fmt.Fprintln(os.Stdout, "Enter board coordinates - row and column (two digits with space):")
					continue
				}

				click := make([]int, 0, len(clickInput))
				for _, c := range clickInput {
					clickInt, err := strconv.Atoi(c)
					if err != nil {
						return
					}
					click = append(click, clickInt)
				}

				err := b.click(click)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					return
				}

				b.printBoard()
				fmt.Println()
				b.printBoardStateLess()
				if b.IsFinished() {
					fmt.Fprintf(os.Stdout, "You %s", b.getBoardState())
					return
				}
			}
		},
	}

	return command
}

package cmd

import (
	"fmt"
	"github.com/proxx/game"
	"github.com/spf13/cobra"
	"os"
)

func start() *cobra.Command {
	command := &cobra.Command{
		Use:   "start",
		Short: "Start the game",
		Run: func(cmd *cobra.Command, _ []string) {
			var cellNumber int
			var blackHoles int

			fmt.Fprintln(os.Stdout, "Enter board N size:")
			_, err := fmt.Scanln(&cellNumber)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(0)
			}

			fmt.Fprintln(os.Stdout, "Enter number of black holes:")
			_, err = fmt.Scanln(&blackHoles)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(0)
			}

			b := game.NewBoard(cellNumber, blackHoles)
			gameInstance := game.NewGame(b)
			b.Print()
			fmt.Println()
			b.PrintStateless()
			for {
				var coordinateX, coordinateY int
				fmt.Fprintln(os.Stdout, "Enter board coordinates - row and column (two digits with space):")
				_, err = fmt.Scan(&coordinateX, &coordinateY)
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(0)
				}

				// subtracting one since user types from 1 to 5 and to align with 0-indexed slices subtracting is done
				err := b.Click([]int{coordinateX - 1, coordinateY - 1})
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					return
				}

				b.Print()
				fmt.Println()
				b.PrintStateless()
				if gameInstance.IsFinished() {
					fmt.Fprintf(os.Stdout, "You %v", gameInstance.GetState())
					return
				}
			}
		},
	}

	return command
}

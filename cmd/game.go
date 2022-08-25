package cmd

import (
	"fmt"
	"os"

	"github.com/proxx/game"
	"github.com/spf13/cobra"
)

func start() *cobra.Command {
	command := &cobra.Command{
		Use:   "start",
		Short: "Start the game",
		RunE: func(cmd *cobra.Command, _ []string) error {
			var (
				cellNumber, blackHoles int
			)

			fmt.Println("Enter board N size:")
			_, err := fmt.Scanln(&cellNumber)
			if err != nil {
				return err
			}

			fmt.Println("Enter number of black holes:")
			_, err = fmt.Scanln(&blackHoles)
			if err != nil {
				return err
			}

			b, err := game.NewBoard(cellNumber, blackHoles)
			if err != nil {
				return err
			}
			gameInstance := game.NewGame(b)

			return gameInstance.Start(os.Stdin)
		},
	}

	return command
}

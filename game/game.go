package game

import (
	"fmt"
	"os"
)

//go:generate mockgen -destination=./mocks/playground.go -package=mocks github.com/proxx/game Playground

// State represent game state
type State string

// list of game states
const (
	win  State = "win"
	lose State = "lose"
	// this state is not used in this game but for clarity I added inProgress
	// since state is implicitly inProgress exists when player not lost or won yet
	inProgress State = "inProgress"
)

// Playground interface that represents methods of playground
type Playground interface {
	Click(click []int) error
	Print()
	WinState() bool
	LoseState() bool
	SetOnStateChangeHook(func())
}

// Game represents game data
type Game struct {
	playground Playground
	state      State
}

// setState sets game state
func (g *Game) setState(state State) {
	g.state = state
}

// GetState get game state
func (g *Game) GetState() State {
	return g.state
}

// IsFinished checks whether game finished
func (g *Game) IsFinished() bool {
	return g.state == win || g.state == lose
}

// NewGame inits new game.
// accepts playground interface
func NewGame(playground Playground) *Game {
	g := &Game{
		playground: playground,
		state:      inProgress,
	}

	g.playground.SetOnStateChangeHook(g.gameStateChangeHook)
	return g
}

// Start starts the game
func (g *Game) Start(in *os.File) error {
	// initial playground print
	g.playground.Print()
	for {
		var coordinateX, coordinateY int
		fmt.Print("Enter board coordinates - row and column (two digits with space):")
		_, err := fmt.Fscan(in, &coordinateX, &coordinateY)
		if err != nil {
			return err
		}

		// subtracting one since user types from 1 to n and to align with 0-indexed slices subtracting is done
		err = g.playground.Click([]int{coordinateX - 1, coordinateY - 1})
		if err != nil {
			fmt.Printf("Notice: %v. Repeat please.", err)
			continue
		}

		g.playground.Print()
		if g.IsFinished() {
			fmt.Printf("You %v \n", g.GetState())
			return nil
		}
	}
}

// gameStateChangeHook hook that observes playground change when game is over
func (g *Game) gameStateChangeHook() {
	switch {
	case g.playground.LoseState():
		g.setState(lose)
	case g.playground.WinState():
		g.setState(win)
	default:
	}
}

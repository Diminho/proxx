package game

// State represent game state
type State string

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
	PrintStateless()
	WinCondition() bool
	LoseCondition() bool
	SetOnStateChangeHook(func())
}

type Game struct {
	playground Playground
	state      State
}

func (g *Game) setState(state State) {
	g.state = state
}

func (g *Game) GetState() State {
	return g.state
}

func (g *Game) IsFinished() bool {
	return g.state == win || g.state == lose
}

func NewGame(playground Playground) *Game {
	g := &Game{
		playground: playground,
		state:      inProgress,
	}

	g.playground.SetOnStateChangeHook(g.gameStateChangeHook)
	return g
}

func (g *Game) Start() {

}

func (g *Game) gameStateChangeHook() {
	switch {
	case g.playground.LoseCondition():
		g.setState(lose)
	case g.playground.WinCondition():
		g.setState(win)
	default:
	}
}

package game

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

const (
	keyCoordinatesFmt   = "%d_%d"
	clickOutOfBoundsFmt = "click coordinate [%d %d] is out of board bounds %d x %d"

	//padding when printing
	paddingLen = "2"
)

var (
	errCellOpened = errors.New("cell already opened")
)

// icon is cell (vertex) view for board printing. E.g. if cell is closed then "c" will be displayed when printed
func stateToIconMapping() map[cellState]string {
	return map[cellState]string{
		openedState:     "o",
		closedState:     "c",
		blackHoledState: "H",
	}
}

// defining directions (neighbors) of given node on the board.
func directions() [][]int {
	return [][]int{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}
}

type boardState string

const (
	blackHoled boardState = "blackHoled"
	cleared    boardState = "cleared"
)

// Board represents board playground
type Board struct {
	boardState boardState
	// represents relations between vertexes (cells)
	adjacencyList map[string][]*cell
	// represents board as two-dimensional slice. this is for printing the board
	board           [][]*cell
	cellList        map[string]*cell
	sideCellsNumber int
	// number of cells to be revealed in order to win
	toBeRevealed     int
	stateChangeHooks []func()
	rows, cols       int
}

// NewBoard init new board as playground
func NewBoard(sideCellsNumber, blackHolesNumber int) (*Board, error) {
	totalCellNumber := sideCellsNumber * sideCellsNumber
	b := &Board{
		adjacencyList:   make(map[string][]*cell),
		sideCellsNumber: sideCellsNumber,
		cellList:        make(map[string]*cell, totalCellNumber),
		toBeRevealed:    totalCellNumber - blackHolesNumber,
		rows:            sideCellsNumber,
		cols:            sideCellsNumber,
	}

	if totalCellNumber < blackHolesNumber {
		return nil, fmt.Errorf(
			"number of blackholes [%d] is bigger than max amount of board cells [%d]. Quiting game",
			blackHolesNumber,
			totalCellNumber)
	}

	blackHolesLocations := distributeBlackHoles(sideCellsNumber, blackHolesNumber)
	b.board = b.generateBoard(blackHolesLocations)

	b.buildGraph(b.board)
	return b, nil
}

func (b *Board) setBoardState(boardState boardState) {
	b.boardState = boardState
	b.execStateChangeHooks()
}

// WinState represents winning state
func (b *Board) WinState() bool {
	return b.boardState == cleared
}

// LoseState represents lose state
func (b *Board) LoseState() bool {
	return b.boardState == blackHoled
}

// decrementToBeRevealed decrements field toBeRevealed to track cells that is yet to be revealed to identify user win
func (b *Board) decrementToBeRevealed() {
	b.toBeRevealed--
	if b.toBeRevealed == 0 {
		b.setBoardState(cleared)
	}
}

func (b *Board) execStateChangeHooks() {
	for _, exec := range b.stateChangeHooks {
		exec()
	}
}

// SetOnStateChangeHook set hook on every board state change
func (b *Board) SetOnStateChangeHook(hookFn func()) {
	b.stateChangeHooks = append(b.stateChangeHooks, hookFn)
}

func isClickValid(click []int, sideCellsNumber int) bool {
	return click[0] > sideCellsNumber || click[1] > sideCellsNumber ||
		click[0] < 0 || click[1] < 0
}

// Click executes click on the given cell. click parameter is x,y coordinates ([]int{x,y})
func (b *Board) Click(click []int) error {
	if isClickValid(click, b.sideCellsNumber) {
		return fmt.Errorf(clickOutOfBoundsFmt, click[0], click[1], b.sideCellsNumber, b.sideCellsNumber)
	}

	currentCell := b.cellList[cellIdentificationKey(click[0], click[1])]
	if currentCell.state.isOpened() {
		return errCellOpened
	}
	if currentCell.value.isBlackHole() {
		b.setBoardState(blackHoled)
		b.revealEntireBoard()
		return nil
	}

	return b.revealCells(cellIdentificationKey(click[0], click[1]))
}

// revealCells uses breadth-first-search to get connected cells with void value.
// BFS is used since it better suits for finding the closest connections (siblings/neighbors)
// and during revealing connected neighbors this is exactly what we need
func (b *Board) revealCells(cellID string) error {
	currentCell := b.cellList[cellID]

	currentCell.state.setToOpened()
	// if cell touches black hole - exit immediately and open just this cell
	if currentCell.value.isTouchingBlackHoles() {
		b.decrementToBeRevealed()
		return nil
	}
	visited := make(map[string]struct{})
	queue := make([]*cell, 0)
	queue = append(queue, currentCell)

	for len(queue) > 0 {
		currentNode := queue[0]
		currentNodeID := cellIdentificationKey(currentNode.x, currentNode.y)
		queue = queue[1:]
		_, ok := visited[currentNodeID]
		if ok {
			continue
		}
		// visit cell
		visited[currentNodeID] = struct{}{}
		currentNode.state.setToOpened()
		b.decrementToBeRevealed()
		// skip revealing neighbors since current cell is touching to the black hole
		if !currentNode.value.isVoid() {
			continue
		}
		neighbors, ok := b.adjacencyList[currentNodeID]
		if !ok {
			continue
		}

		for _, neighbor := range neighbors {
			_, ok := visited[cellIdentificationKey(neighbor.x, neighbor.y)]
			if !ok {
				queue = append(queue, neighbor)
			}
		}
	}

	return nil
}

func (b *Board) revealEntireBoard() {
	for _, c := range b.cellList {
		if c.value == blackHole {
			c.state.setToBlackHoled()
			continue
		}
		c.state.setToOpened()
	}
}

func (b *Board) buildGraph(board [][]*cell) {
	b.addVertexes(board)
	b.addEdges(board)
}

func (b *Board) addVertexes(board [][]*cell) {
	for i := 0; i < b.rows; i++ {
		for j := 0; j < b.cols; j++ {
			b.addVertex(board[i][j])
		}
	}
}

func (b *Board) addEdges(board [][]*cell) {
	// adding edges that connect vertices based of neighbor placement
	for i := 0; i < b.rows; i++ {
		for j := 0; j < b.cols; j++ {
			for _, direction := range directions() {
				dirI := i + direction[0]
				dirJ := j + direction[1]
				if (0 <= dirI && dirI < b.rows) &&
					(0 <= dirJ && dirJ < b.cols) {
					node1 := board[i][j]
					node2 := board[dirI][dirJ]
					b.addEdge(node1, node2)
				}
			}
		}
	}
}

func (b *Board) addEdge(node1, node2 *cell) {
	if node1 == nil || node2 == nil {
		return
	}
	node1Key := cellIdentificationKey(node1.x, node1.y)
	list1 := b.adjacencyList[node1Key]
	list1 = append(list1, node2)
	b.adjacencyList[node1Key] = list1
}

func (b *Board) addVertex(node *cell) {
	if node == nil {
		return
	}
	key := cellIdentificationKey(node.x, node.y)
	_, ok := b.adjacencyList[key]
	if !ok {
		b.adjacencyList[key] = []*cell{}
	}
}

func distributeBlackHoles(sideCount, blackHolesTargetNumber int) [][]int {
	//bh - black hole.
	bhLocations := make([][]int, 0, blackHolesTargetNumber)

	occupiedPositions := make(map[string]struct{}, blackHolesTargetNumber)

	var blackHolesPlaced int
	for blackHolesPlaced < blackHolesTargetNumber {
		// inits rand providing source to create uniformly-distributed number
		r := rand.New(rand.NewSource(time.Now().UnixNano()))

		// excluding this from linter check since it for game purposes it is acceptable to use it
		//nolint: gosec
		x := r.Intn(sideCount)
		//nolint: gosec
		y := r.Intn(sideCount)

		position := cellIdentificationKey(x, y)
		_, ok := occupiedPositions[position]
		if ok {
			continue
		}
		bhLocations = append(bhLocations, []int{x, y})
		occupiedPositions[position] = struct{}{}
		blackHolesPlaced++
	}

	return bhLocations
}

// cellIdentificationKey builds key identify cells in the board. id key is basically x and y coordinates.
func cellIdentificationKey(x, y int) string {
	return fmt.Sprintf(keyCoordinatesFmt, x, y)
}

// cell represents cell data
type cell struct {
	state cellState
	value cellValue
	x, y  int
}

// cellState represents state of cell
type cellState int

const (
	openedState     cellState = 1
	closedState     cellState = 0
	blackHoledState cellState = -1
)

func (cs *cellState) setToOpened() {
	*cs = openedState
}

func (cs *cellState) setToBlackHoled() {
	*cs = blackHoledState
}

func (cs cellState) isOpened() bool {
	return cs == openedState
}

func (cs cellState) isBlackHoled() bool {
	return cs == blackHoledState
}

func (cs cellState) isClosed() bool {
	return cs == closedState
}

type cellValue int

const (
	void      cellValue = 0
	blackHole cellValue = -1
)

func (c cellValue) isBlackHole() bool {
	return c == blackHole
}

func (c cellValue) isVoid() bool {
	return c == void
}

func (c cellValue) isTouchingBlackHoles() bool {
	return !c.isVoid() && !c.isBlackHole()
}

func (b *Board) generateBoard(blackHoles [][]int) [][]*cell {
	emptyBoard := b.initBoard()

	return b.setItems(blackHoles, emptyBoard)
}

// initiates board will default cell values
func (b *Board) initBoard() [][]*cell {
	board := make([][]*cell, b.rows)
	for i := 0; i < b.rows; i++ {
		board[i] = make([]*cell, b.cols)
		for j := 0; j < b.cols; j++ {
			c := &cell{
				value: void,
				x:     i,
				y:     j,
			}
			board[i][j] = c
			b.cellList[cellIdentificationKey(i, j)] = c
		}
	}

	return board
}

// setItems sets black holes and cell counters that touch cells with black holes
func (b *Board) setItems(blackHoles [][]int, board [][]*cell) [][]*cell {
	for _, r := range blackHoles {
		rowI, colI := r[0], r[1]
		board[rowI][colI].value = blackHole
		for i := rowI - 1; i <= rowI+1; i++ {
			for j := colI - 1; j <= colI+1; j++ {
				if (0 <= i && i < b.rows) && (0 <= j && j < b.cols) && board[i][j].value != blackHole {
					board[i][j].value++
				}
			}
		}
	}

	return board
}

// Print prints current state of board
func (b *Board) Print() {
	for i, row := range b.board {
		for col := range row {
			var (
				cellView string
			)

			switch {
			case b.board[i][col].state.isClosed(), b.board[i][col].state.isBlackHoled():
				cellView = stateToIconMapping()[b.board[i][col].state]
			default:
				cellView = fmt.Sprintf("%d", b.board[i][col].value)
			}

			fmt.Printf("%v %"+paddingLen+"s", cellView, "")
			fmt.Print(" ")
		}
		fmt.Println()
	}
}

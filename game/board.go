package game

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

const (
	keyCoordinatesFmt   = "%d_%d"
	clickOutOfBoundsFmt = "click coordinate [%d %d] is out of board bounds %d x %d"
)

// icon is cell (vertex) view for board printing. E.g. if cell is closed then "c" will be displayed when printed
func stateToIconMapping() map[cellState]string {
	return map[cellState]string{
		openedState:     "o",
		closedState:     "c",
		blackHoledState: "H",
	}
}

type boardState string

const (
	blackHoled boardState = "blackHoled"
	cleared    boardState = "cleared"

	//padding when printing
	paddingLen = "2"
)

// Board represents board playground
type Board struct {
	boardState boardState
	numNodes   []*cell
	// represents relations between vertexes (cells)
	adjacencyList map[string][]*cell
	// represents board as two-dimensional slice. this is for printing the board
	board           [][]*cell
	cellList        map[string]*cell
	sidecellsNumber int
	// number of cells to be revealed in order to win
	toBeRevealed     int
	stateChangeHooks []func()
}

// NewBoard init new board as playground
func NewBoard(sideCellsNumber, blackHolesNumber int) *Board {
	totalCellNumber := sideCellsNumber * sideCellsNumber
	b := &Board{
		numNodes:        make([]*cell, 0, totalCellNumber),
		adjacencyList:   make(map[string][]*cell),
		sidecellsNumber: sideCellsNumber,
		cellList:        make(map[string]*cell, totalCellNumber),
		toBeRevealed:    totalCellNumber - blackHolesNumber,
	}

	blackHolesLocations := distributeBlackHoles(sideCellsNumber, blackHolesNumber)
	b.board = b.generateBoard(blackHolesLocations, sideCellsNumber, sideCellsNumber)

	b.buildGraph(b.board, sideCellsNumber, sideCellsNumber)
	return b
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

// Click executes click on the given cell. click parameter is x,y coordinates ([]int{x,y})
func (b *Board) Click(click []int) error {
	if click[0] > b.sidecellsNumber || click[1] > b.sidecellsNumber {
		return fmt.Errorf(clickOutOfBoundsFmt, click[0], click[1], b.sidecellsNumber, b.sidecellsNumber)
	}

	currentCell := b.cellList[cellIdentificationKey(click[0], click[1])]
	if currentCell.state.isOpened() {
		_, err := fmt.Fprintln(os.Stdout, "cell already opened:")
		return err
	}
	if currentCell.value.isBlackHole() {
		b.setBoardState(blackHoled)
		b.revealEntireBoard()
		return nil
	}

	return b.revealCells(cellIdentificationKey(click[0], click[1]))
}

// revealCells uses breadth first search to get connected cells with void value
func (b *Board) revealCells(cellID string) error {
	currentCell := b.cellList[cellID]

	currentCell.state.setToOpened()
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
		// visit
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

func (b *Board) buildGraph(artifacts [][]*cell, rows, cols int) {
	b.addVertexes(artifacts, rows, cols)
	b.addEdges(artifacts, rows, cols)
}

func (b *Board) addVertexes(artifacts [][]*cell, rows, cols int) {
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			b.addVertex(artifacts[i][j])
		}
	}
}

func (b *Board) addEdges(artifacts [][]*cell, rows, cols int) {
	// defining directions (neighbors) places related to given node
	directions := [][]int{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}
	// adding edges that connect vertices based of neighbor placement
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			for _, direction := range directions {
				dirI := i + direction[0]
				dirJ := j + direction[1]
				if (0 <= dirI && dirI < rows) &&
					(0 <= dirJ && dirJ < cols) {
					node1 := artifacts[i][j]
					node2 := artifacts[dirI][dirJ]
					b.addEdge(node1, node2)
				}
			}
		}
	}
}

func (b *Board) addEdge(node1, node2 *cell) {
	node1Key := cellIdentificationKey(node1.x, node1.y)
	list1, ok1 := b.adjacencyList[node1Key]
	if !ok1 {
		return
	}
	list1 = append(list1, node2)
	b.adjacencyList[node1Key] = list1
}

func (b *Board) addVertex(node *cell) {
	key := cellIdentificationKey(node.x, node.y)
	b.numNodes = append(b.numNodes, node)
	_, ok := b.adjacencyList[key]
	if !ok {
		b.adjacencyList[key] = []*cell{}
	}
}

// @TODO make distribution even
func distributeBlackHoles(sideCount, blackHolesTargetNumber int) [][]int {
	//bh - black hole.
	bhLocations := make([][]int, 0, blackHolesTargetNumber)

	occupiedPositions := make(map[string]struct{}, blackHolesTargetNumber)

	var blackHolesPlaced int
	for blackHolesPlaced < blackHolesTargetNumber {
		rand.Seed(time.Now().UnixNano())

		// excluding this since it for game purposes it is acceptable to use it
		//nolint: gosec
		x := rand.Intn(sideCount)
		//nolint: gosec
		y := rand.Intn(sideCount)

		position := fmt.Sprintf("%d_%d", x, y)
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

func (b *Board) generateBoard(blackHoles [][]int, rows, cols int) [][]*cell {
	artifacts := make([][]*cell, rows)
	for i := 0; i < rows; i++ {
		artifacts[i] = make([]*cell, cols)
		for j := 0; j < cols; j++ {
			c := &cell{
				value: void,
				x:     i,
				y:     j,
			}
			artifacts[i][j] = c
			b.cellList[cellIdentificationKey(i, j)] = c
		}
	}

	return setArtifacts(blackHoles, artifacts, rows, cols)
}

// setArtifacts sets black holes and cell counter that touch black holes
func setArtifacts(blackHoles [][]int, artifacts [][]*cell, rows, cols int) [][]*cell {
	for _, r := range blackHoles {
		rowI, colI := r[0], r[1]
		artifacts[rowI][colI].value = blackHole
		for i := rowI - 1; i <= rowI+1; i++ {
			for j := colI - 1; j <= colI+1; j++ {
				if (0 <= i && i < rows) && (0 <= j && j < cols) && artifacts[i][j].value != blackHole {
					artifacts[i][j].value++
				}
			}
		}
	}

	return artifacts
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

// PrintStateless is for debugging (or verifying that game works correctly) purposes.
// shows all values on the board. Blackhole will be as -1
func (b *Board) PrintStateless() {
	for i, row := range b.board {
		for col := range row {
			fmt.Printf("%v %"+paddingLen+"s", b.board[i][col].value, "")
			fmt.Print(" ")
		}
		fmt.Println()
	}
}

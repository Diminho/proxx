package game

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	keyCoordinatesFmt   = "%d_%d"
	clickOutOfBoundsFmt = "click coordinate [%d %d] is out of board bounds %d x %d"
)

type boardState string

const (
	blackHoled boardState = "blackHoled"
	cleared    boardState = "cleared"
)

type board struct {
	boardState boardState
	numNodes   []*Cell
	// represents relations between vertexes (cells)
	adjacencyList map[string][]*Cell
	// represents board as two-dimensional slice. this is for printing the board
	board           [][]*Cell
	cellList        map[string]*Cell
	sideCellsNumber int
	// number of cells to be revealed in order to win
	toBeRevealed     int
	stateChangeHooks []func()
}

func NewBoard(sideCellsNumber, blackHolesNumber int) *board {
	b := &board{
		numNodes:        make([]*Cell, 0, sideCellsNumber*sideCellsNumber),
		adjacencyList:   make(map[string][]*Cell),
		sideCellsNumber: sideCellsNumber,
		cellList:        make(map[string]*Cell, sideCellsNumber*sideCellsNumber),
		toBeRevealed:    sideCellsNumber*sideCellsNumber - blackHolesNumber,
	}

	blackHolesLocations := distributeBlackHoles(sideCellsNumber, blackHolesNumber)
	b.board = b.generateArtifacts(blackHolesLocations, sideCellsNumber, sideCellsNumber)

	b.buildGraph(b.board, sideCellsNumber, sideCellsNumber)
	return b
}

func (b *board) setBoardState(boardState boardState) {
	b.boardState = boardState
	b.execStateChangeHooks()
}

func (b *board) WinCondition() bool {
	return b.boardState == cleared
}

func (b *board) LoseCondition() bool {
	return b.boardState == blackHoled
}

func (b *board) decrementToBeRevealed() {
	fmt.Println("BEFORE ", b.toBeRevealed)

	b.toBeRevealed--
	fmt.Println("AFTER ", b.toBeRevealed)

	if b.toBeRevealed == 0 {
		b.setBoardState(cleared)
	}
}

func (b *board) execStateChangeHooks() {
	for _, exec := range b.stateChangeHooks {
		exec()
	}
}

func (b *board) SetOnStateChangeHook(hookFn func()) {
	b.stateChangeHooks = append(b.stateChangeHooks, hookFn)
}

// Click executes click on the given cell. click parameter is x,y coordinates ([]int{x,y})
func (b *board) Click(click []int) error {
	if click[0] > b.sideCellsNumber || click[1] > b.sideCellsNumber {
		return fmt.Errorf(clickOutOfBoundsFmt, click[0], click[1], b.sideCellsNumber, b.sideCellsNumber)
	}

	cell := b.cellList[cellIdentificationKey(click[0], click[1])]
	if cell.state.isOpened() {
		fmt.Println("already opened cell ", click[0], click[1])
		return nil
	}
	if cell.value.isBlackHole() {
		b.setBoardState(blackHoled)
		b.revealEntireBoard()
		return nil
	}

	return b.revealCells(cellIdentificationKey(click[0], click[1]))
}

// revealCells uses breadth first search to get connected cells with void value
func (b *board) revealCells(cellID string) error {
	cell := b.cellList[cellID]

	cell.state.setToOpened()
	if cell.value.isTouchingBlackHoles() {
		b.decrementToBeRevealed()
		return nil
	}
	visited := make(map[string]struct{})
	queue := make([]*Cell, 0)

	queue = append(queue, cell)

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
		fmt.Println("OPENING ", currentNodeID)
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

func (b *board) revealEntireBoard() {
	for _, cell := range b.cellList {
		cell.state.setToOpened()
	}
}

func (b *board) buildGraph(artifacts [][]*Cell, rows, cols int) {
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			b.addVertex(artifacts[i][j])
		}
	}

	// defining directions (neighbours) places related to given node
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

// @TODO make distribution even
func distributeBlackHoles(sideCount, blackHolesTargetNumber int) [][]int {
	//bh - black hole.
	bhLocations := make([][]int, 0, blackHolesTargetNumber)

	occupiedPositions := make(map[string]struct{}, blackHolesTargetNumber)

	var blackHolesPlaced int
	for blackHolesPlaced < blackHolesTargetNumber {
		rand.Seed(time.Now().UnixNano())

		x := rand.Intn(sideCount)
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

func (b *board) addEdge(node1, node2 *Cell) {
	node1Key := cellIdentificationKey(node1.x, node1.y)
	list1, ok1 := b.adjacencyList[node1Key]
	if !ok1 {
		return
	}
	list1 = append(list1, node2)
	b.adjacencyList[node1Key] = list1
}

func (b *board) addVertex(node *Cell) {
	key := cellIdentificationKey(node.x, node.y)
	b.numNodes = append(b.numNodes, node)
	_, ok := b.adjacencyList[key]
	if !ok {
		b.adjacencyList[key] = []*Cell{}
	}
}

type Cell struct {
	state CellState
	value CellValue
	x, y  int
}

type CellState int

const (
	opened CellState = 1
	closed CellState = 0
)

func (cs *CellState) setToOpened() {
	*cs = opened
}

func (cs CellState) isOpened() bool {
	return cs == opened
}

func (cs CellState) isClosed() bool {
	return cs == closed
}

type CellValue int

const (
	void      CellValue = 0
	blackHole CellValue = -1
)

func (c CellValue) isBlackHole() bool {
	return c == blackHole
}

func (c CellValue) isVoid() bool {
	return c == void
}

func (c CellValue) isTouchingBlackHoles() bool {
	return !c.isVoid() && !c.isBlackHole()
}

// icon is cell (vertex) view for board printing. E.g. if cell is closed then "c" will be displayed when printed
var stateToIconMapping = map[CellState]string{
	opened: "o",
	closed: "c",
}

func (b *board) generateArtifacts(blackHoles [][]int, rows, cols int) [][]*Cell {
	artifacts := make([][]*Cell, rows)
	for i := 0; i < rows; i++ {
		artifacts[i] = make([]*Cell, cols)
		for j := 0; j < cols; j++ {
			cell := &Cell{
				value: void,
				x:     i,
				y:     j,
			}
			artifacts[i][j] = cell
			b.cellList[cellIdentificationKey(i, j)] = cell
		}
	}

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
func (b *board) Print() {
	for i, row := range b.board {
		for col := range row {
			var (
				cellView, paddingLen string
			)
			if b.LoseCondition() && b.board[i][col].state.isOpened() {
				paddingLen = fmt.Sprintf("%d", 2)
			} else {
				paddingLen = fmt.Sprintf("%d", 3)
			}

			if b.board[i][col].state.isClosed() {
				cellView = stateToIconMapping[closed]
			} else {
				cellView = fmt.Sprintf("%d", b.board[i][col].value)
			}
			fmt.Printf("%v %"+paddingLen+"s", cellView, "")
			fmt.Print(" ")
		}
		fmt.Println()
	}
}

// PrintStateless is for debugging (or verifying that game works correctly) purposes
func (b *board) PrintStateless() {
	for i, row := range b.board {
		for col := range row {
			var (
				paddingLen string
			)
			if b.LoseCondition() {
				paddingLen = fmt.Sprintf("%d", 2)
			} else {
				paddingLen = fmt.Sprintf("%d", 3)
			}

			fmt.Printf("%v %"+paddingLen+"s", b.board[i][col].value, "")
			fmt.Print(" ")
		}
		fmt.Println()
	}
}

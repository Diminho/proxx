package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

const (
	keyCoordinatesFmt   = "%d_%d"
	clickOutOfBoundsFmt = "click coordinate [%d %d] is out of board bounds %d x %d"
)

var (
	GameOver = errors.New("you have been sucked by mighty black-black hole")
)

func main() {
	//err := cmd.Execute()
	//if err != nil {
	//	log.Fatalf("%+v\n", err)
	//}

	b := newBoard(5, 3)
	b.printBoard()
	fmt.Println()
	b.printBoardStateLess()

	clicks := [][]int{
		{1, 0},
		{4, 3},
		{2, 2},
		{0, 2},
		{1, 4},
	}
	//fmt.Println(b)
	for _, click := range clicks {
		err := b.click(click)
		if err != nil {
			fmt.Println("ERROR: ", err)
			return
		}

		b.printBoard()
		fmt.Println()
		b.printBoardStateLess()
	}

}

type board struct {
	numNodes []*Cell
	// represents relations between vertexes (cells)
	adjacencyList map[string][]*Cell
	// represents board as two-dimensional slice
	board           [][]*Cell
	cellList        map[string]*Cell
	sideCellsNumber int
}

func newBoard(sideCellsNumber, blackHolesNumber int) *board {
	b := &board{
		numNodes:        make([]*Cell, 0, sideCellsNumber*sideCellsNumber),
		adjacencyList:   make(map[string][]*Cell),
		sideCellsNumber: sideCellsNumber,
		cellList:        make(map[string]*Cell, sideCellsNumber*sideCellsNumber),
	}

	blackHolesLocations := distributeBlackHoles(sideCellsNumber, blackHolesNumber)
	b.board = b.generateArtifacts(blackHolesLocations, sideCellsNumber, sideCellsNumber)

	b.buildGraph(b.board, sideCellsNumber, sideCellsNumber)
	return b
}

// click executes click on the given cell. click parameter is x,y coordinates ([]int{x,y})
func (b *board) click(click []int) error {
	if click[0] > b.sideCellsNumber || click[1] > b.sideCellsNumber {
		return fmt.Errorf(clickOutOfBoundsFmt, click[0], click[1], b.sideCellsNumber, b.sideCellsNumber)
	}

	cell := b.cellList[cellIdentificationKey(click[0], click[1])]
	if cell.state.isOpened(){
		fmt.Println("already opened cell ", click[0], click[1])
		return nil
	}
	if cell.value.isBlackHole() {
		return GameOver
	}
	cell.state.setToOpened()
	if cell.value.isTouchingBlackHoles(){
		fmt.Println("touching cell ", cell.value)
		return nil
	}

	return b.bfs(cellIdentificationKey(click[0], click[1]))
}

func (b *board) bfs(cellID string) error {
	cell := b.cellList[cellID]
	fmt.Println("cell value: ", cell.value)

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
		if !currentNode.value.isVoid() {
			continue
		}
		neighbors, ok := b.adjacencyList[currentNodeID]
		if !ok {
			continue
		}

		for _, neighbor := range neighbors {
			neighborCellID := cellIdentificationKey(neighbor.x, neighbor.y)
			_, ok := visited[neighborCellID]
			if !ok {
				queue = append(queue, neighbor)
			}
		}
	}

	return nil
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

func (cs *CellState) setToOpened() {
	*cs = opened
}

const (
	opened CellState = 1
	closed CellState = 0
)

func (cs CellState) isOpened() bool {
	return cs == opened
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

func (b *board) printBoard() {
	for i, row := range b.board {
		for col := range row {
			var (
				cellView, paddingLen string
			)
			if b.board[i][col].value == blackHole && b.board[i][col].state == opened {
				paddingLen = fmt.Sprintf("%d", 2)
			} else {
				paddingLen = fmt.Sprintf("%d", 3)
			}

			if b.board[i][col].state == closed {
				cellView = stateToIconMapping[closed]
			} else {
				cellView = fmt.Sprintf("%d", b.board[i][col].value)
			}
			fmt.Print(fmt.Sprintf("%v %"+paddingLen+"s", cellView, ""))
			fmt.Print(" ")
		}
		fmt.Println()
	}
}

//
func (b *board) printBoardStateLess() {
	for i, row := range b.board {
		for col := range row {
			var (
				paddingLen string
			)
			if b.board[i][col].value == blackHole {
				paddingLen = fmt.Sprintf("%d", 2)
			} else {
				paddingLen = fmt.Sprintf("%d", 3)
			}

			fmt.Print(fmt.Sprintf("%v %"+paddingLen+"s", b.board[i][col].value, ""))
			fmt.Print(" ")
		}
		fmt.Println()
	}
}

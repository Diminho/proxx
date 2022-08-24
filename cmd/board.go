package cmd

import (
	"fmt"
	"math/rand"
	"time"
)


func newGraph(sideCellsCount, blackHolesNumber int) *board {
	b := &board{
		numNodes:      make([]*Cell, sideCellsCount*sideCellsCount),
		adjacencyList: make(map[int][]*Cell),
	}

	blackHolesLocations := distributeBlackHoles(sideCellsCount, blackHolesNumber)
	generateArtifacts(blackHolesLocations, sideCellsCount, sideCellsCount)
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

func (g *board) addEdge(node1, node2 int) {
	list1, ok1 := g.adjacencyList[node1]
	list2, ok2 := g.adjacencyList[node2]
	if !ok1 || !ok2 {
		return
	}
	list1 = append(list1, node2)
	list2 = append(list2, node1)

	g.adjacencyList[node2] = list2
	g.adjacencyList[node1] = list1
}

func (g *board) addVertex(node int) {
	g.numNodes = append(g.numNodes, node)
	_, ok := g.adjacencyList[node]
	if !ok {
		g.adjacencyList[node] = []int{}
	}
}

type board struct {
	numNodes      []*Cell
	adjacencyList map[int][]*Cell
}

type Cell struct {
	state CellState
	value CellValue
}

type CellState string

const (
	opened CellState = "opened"
	closed CellState = "closed"
)

type CellValue int

const (
	void      CellValue = 0
	blackHole CellValue = -1
)

func generateArtifacts(bombs [][]int, rows, cols int) {
	board := make([][]CellValue, rows)
	for i := range board {
		board[i] = make([]CellValue, cols)
	}

	for _, r := range bombs {
		row_i, col_i := r[0], r[1]
		board[row_i][col_i] = blackHole
		for i := row_i - 1; i <= row_i+1; i++ {
			for j := col_i - 1; j <= col_i+1; j++ {
				if (0 <= i && i < rows) && (0 <= j && j < cols) && board[i][j] != -1 {
					board[i][j]++
				}
			}
		}
	}

	fmt.Println(board)
}

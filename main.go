package main

import (
	"fmt"
	"math/rand"
	"time"
)

const keyCoordinatesFmt = "%d_%d"

func main() {
	//err := cmd.Execute()
	//if err != nil {
	//	log.Fatalf("%+v\n", err)
	//}

	b := newBoard(3, 3)
	fmt.Println(b)

}

func newBoard(sideCellsCount, blackHolesNumber int) *board {
	b := &board{
		numNodes:      make([]Cell, 0, sideCellsCount*sideCellsCount),
		adjacencyList: make(map[string][]Cell),
	}

	blackHolesLocations := distributeBlackHoles(sideCellsCount, blackHolesNumber)
	artifacts := generateArtifacts(blackHolesLocations, sideCellsCount, sideCellsCount)
	b.buildGraph(artifacts, sideCellsCount, sideCellsCount)
	return b
}

func (b *board) buildGraph(artifacts [][]Cell, rows, cols int) {
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			b.addVertex(artifacts[i][j])
		}
	}

	directions := [][]int{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}
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

func (b *board) addEdge(node1, node2 Cell) {
	node1Key := fmt.Sprintf(keyCoordinatesFmt, node1.x, node1.y)
	//node2Key := fmt.Sprintf(keyCoordinatesFmt, node2.x, node2.y)
	list1, ok1 := b.adjacencyList[node1Key]
	//list2, ok2 := b.adjacencyList[node2Key]
	if !ok1 {
		return
	}
	list1 = append(list1, node2)
	//list2 = append(list2, node1)

	//b.adjacencyList[node2Key] = list2
	b.adjacencyList[node1Key] = list1
}

func (b *board) addVertex(node Cell) {
	key := fmt.Sprintf(keyCoordinatesFmt, node.x, node.y)
	b.numNodes = append(b.numNodes, node)
	_, ok := b.adjacencyList[key]
	if !ok {
		b.adjacencyList[key] = []Cell{}
	}
}

type board struct {
	numNodes      []Cell
	adjacencyList map[string][]Cell
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

type CellValue int

const (
	void      CellValue = 0
	blackHole CellValue = -1
)

func generateArtifacts(bombs [][]int, rows, cols int) [][]Cell {
	artifacts := make([][]Cell, rows)
	for i := 0; i < rows; i++ {
		artifacts[i] = make([]Cell, cols)
		for j := 0; j < cols; j++ {
			value := Cell{
				value: void,
				x:     i,
				y:     j,
			}
			artifacts[i][j] = value
		}
	}

	for _, r := range bombs {
		rowI, colI := r[0], r[1]
		artifacts[rowI][colI].value = blackHole
		for i := rowI - 1; i <= rowI+1; i++ {
			for j := colI - 1; j <= colI+1; j++ {
				if (0 <= i && i < rows) && (0 <= j && j < cols) && artifacts[rowI][colI].value != blackHole {
					artifacts[rowI][colI].value++
				}
			}
		}
	}

	return artifacts
}

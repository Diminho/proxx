package main

import (
	"github.com/proxx/cmd"
	"log"
)

const (
	keyCoordinatesFmt   = "%d_%d"
	clickOutOfBoundsFmt = "click coordinate [%d %d] is out of board bounds %d x %d"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatalf("%+v\n", err)
	}

	//b := newBoard(5, 3)
	//b.printBoard()
	//fmt.Println()
	//b.printBoardStateLess()
	//
	//clicks := [][]int{
	//	{1, 0},
	//	{4, 3},
	//	{2, 2},
	//	{0, 2},
	//	{1, 4},
	//}
	//for _, click := range clicks {
	//	err := b.click(click)
	//	if err != nil {
	//		fmt.Println("ERROR: ", err)
	//		return
	//	}
	//
	//	b.printBoard()
	//	fmt.Println()
	//	b.printBoardStateLess()
	//	if b.IsFinished() {
	//		fmt.Println("you ", b.getBoardState())
	//		return
	//	}
	//	fmt.Println("===============")
	//}
	//fmt.Println("stop")
}

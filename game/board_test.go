package game

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBoard_addEdge(t *testing.T) {
	type fields struct {
		rows, cols, blackHolesNumber int
		blackHoleLocations           [][]int
	}
	type args struct {
		node1 *cell
		node2 *cell
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		verifyFn func(b *Board, node *cell) bool
	}{
		{
			name: "success",
			verifyFn: func(b *Board, node *cell) bool {
				// verify whether node is present
				id := cellIdentificationKey(node.x, node.y)
				_, ok := b.adjacencyList[id]
				return ok
			},
			fields: fields{
				rows: 3,
				cols: 3,
				blackHoleLocations: [][]int{
					{0, 1},
				},
			},
			args: args{
				node1: &cell{
					x: 1,
					y: 1,
				},
				node2: &cell{
					x: 1,
					y: 0,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			totalCellNumber := tt.fields.rows * tt.fields.cols
			b := &Board{
				adjacencyList:   make(map[string][]*cell),
				sideCellsNumber: tt.fields.rows,
				cellList:        make(map[string]*cell, totalCellNumber),
				toBeRevealed:    totalCellNumber - tt.fields.blackHolesNumber,
				rows:            tt.fields.rows,
				cols:            tt.fields.cols,
			}
			b.board = b.generateBoard(tt.fields.blackHoleLocations)

			b.addEdge(tt.args.node1, tt.args.node2)

			assert.True(t, tt.verifyFn(b, tt.args.node1))
		})
	}
}

func TestBoard_addEdges(t *testing.T) {
	type fields struct {
		boardState       boardState
		adjacencyList    map[string][]*cell
		board            [][]*cell
		cellList         map[string]*cell
		sideCellsNumber  int
		toBeRevealed     int
		stateChangeHooks []func()
	}
	type args struct {
		rows               int
		cols               int
		blackHolesNumber   int
		blackHoleLocations [][]int
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		verifyFn func(b *Board) bool
	}{
		{
			name: "success",
			verifyFn: func(b *Board) bool {
				// 9 since rows and cols 3x3
				return len(b.adjacencyList) == 9
			},
			args: args{
				blackHolesNumber: 2,
				rows:             3,
				cols:             3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			totalCellNumber := tt.args.rows * tt.args.cols
			b := &Board{
				adjacencyList:   make(map[string][]*cell),
				sideCellsNumber: tt.args.rows,
				cellList:        make(map[string]*cell, totalCellNumber),
				toBeRevealed:    totalCellNumber - tt.args.blackHolesNumber,
				rows:            tt.args.rows,
				cols:            tt.args.cols,
			}
			b.board = b.generateBoard(tt.args.blackHoleLocations)

			b.addEdges(b.board)
			assert.True(t, tt.verifyFn(b))
		})
	}
}

func TestBoard_addVertexes(t *testing.T) {
	type args struct {
		rows               int
		cols               int
		blackHolesNumber   int
		blackHoleLocations [][]int
	}
	tests := []struct {
		name     string
		args     args
		verifyFn func(b *Board) bool
	}{
		{
			name: "success",
			verifyFn: func(b *Board) bool {
				for _, row := range b.board {
					for _, c := range row {
						key := cellIdentificationKey(c.x, c.y)
						_, ok := b.adjacencyList[key]
						// return false if key is not present
						if !ok {
							return false
						}
					}
				}
				// 9 since rows and cols 3x3
				return len(b.adjacencyList) == 9
			},
			args: args{
				blackHolesNumber: 2,
				rows:             3,
				cols:             3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			totalCellNumber := tt.args.rows * tt.args.cols
			b := &Board{
				adjacencyList:   make(map[string][]*cell),
				sideCellsNumber: tt.args.rows,
				cellList:        make(map[string]*cell, totalCellNumber),
				toBeRevealed:    totalCellNumber - tt.args.blackHolesNumber,
				rows:            tt.args.rows,
				cols:            tt.args.cols,
			}
			b.board = b.generateBoard(tt.args.blackHoleLocations)

			b.addVertexes(b.board)
			assert.True(t, tt.verifyFn(b))
		})
	}
}

func Test_distributeBlackHoles(t *testing.T) {
	type args struct {
		sideCount              int
		blackHolesTargetNumber int
	}
	tests := []struct {
		name        string
		args        args
		expectedLen int
		want        [][]int
	}{
		{
			name: "1_black_hole",
			args: args{
				sideCount:              3,
				blackHolesTargetNumber: 1,
			},
			expectedLen: 1,
		},
		{
			name: "2_black_holes",
			args: args{
				sideCount:              3,
				blackHolesTargetNumber: 2,
			},
			expectedLen: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := distributeBlackHoles(tt.args.sideCount, tt.args.blackHolesTargetNumber)
			assert.Equal(t, tt.expectedLen, len(actual))
		})
	}
}

func Test_setArtifacts(t *testing.T) {
	type args struct {
		blackHoleLocations [][]int
		rows               int
		cols               int
		blackHolesNumber   int
	}
	tests := []struct {
		name string
		args args
		want [][]*cell
	}{
		{
			name: "one_black_hole",
			args: args{
				blackHoleLocations: [][]int{
					{0, 0},
				},
				rows: 3,
				cols: 3,
			},
			want: [][]*cell{
				{
					{
						state: closedState,
						value: blackHole,
						x:     0,
						y:     0,
					},
					{
						state: closedState,
						value: 1,
						x:     0,
						y:     1,
					},
					{
						state: closedState,
						value: void,
						x:     0,
						y:     2,
					},
				},
				{
					{
						state: closedState,
						value: 1,
						x:     1,
						y:     0,
					},
					{
						state: closedState,
						value: 1,
						x:     1,
						y:     1,
					},
					{
						state: closedState,
						value: void,
						x:     1,
						y:     2,
					},
				},
				{
					{
						state: closedState,
						value: void,
						x:     2,
						y:     0,
					},
					{
						state: closedState,
						value: void,
						x:     2,
						y:     1,
					},
					{
						state: closedState,
						value: void,
						x:     2,
						y:     2,
					},
				},
			},
		},
		{
			name: "two_black_holes",
			args: args{
				blackHoleLocations: [][]int{
					{0, 0},
					{2, 2},
				},
				rows: 3,
				cols: 3,
			},
			want: [][]*cell{
				{
					{
						state: closedState,
						value: blackHole,
						x:     0,
						y:     0,
					},
					{
						state: closedState,
						value: 1,
						x:     0,
						y:     1,
					},
					{
						state: closedState,
						value: void,
						x:     0,
						y:     2,
					},
				},
				{
					{
						state: closedState,
						value: 1,
						x:     1,
						y:     0,
					},
					{
						state: closedState,
						value: 2,
						x:     1,
						y:     1,
					},
					{
						state: closedState,
						value: 1,
						x:     1,
						y:     2,
					},
				},
				{
					{
						state: closedState,
						value: void,
						x:     2,
						y:     0,
					},
					{
						state: closedState,
						value: 1,
						x:     2,
						y:     1,
					},
					{
						state: closedState,
						value: blackHole,
						x:     2,
						y:     2,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			totalCellNumber := tt.args.rows * tt.args.cols
			b := &Board{
				adjacencyList:   make(map[string][]*cell),
				sideCellsNumber: tt.args.rows,
				cellList:        make(map[string]*cell, totalCellNumber),
				toBeRevealed:    totalCellNumber - tt.args.blackHolesNumber,
				rows:            tt.args.rows,
				cols:            tt.args.cols,
			}
			emptyBoard := b.initBoard()

			actual := b.setItems(tt.args.blackHoleLocations, emptyBoard)
			assert.Equal(t, tt.want, actual)
		})
	}
}

// since this func is just initialization so only to verify if cells are not nil
func Test_initBoard(t *testing.T) {
	type args struct {
		rows             int
		cols             int
		blackHolesNumber int
	}
	tests := []struct {
		name     string
		args     args
		verifyFn func([][]*cell) bool
		want     [][]*cell
	}{
		{
			name: "success",
			verifyFn: func(b [][]*cell) bool {
				for _, row := range b {
					for _, col := range row {
						if col == nil {
							return false
						}
					}
				}

				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			totalCellNumber := tt.args.rows * tt.args.cols
			b := &Board{
				adjacencyList:   make(map[string][]*cell),
				sideCellsNumber: tt.args.rows,
				cellList:        make(map[string]*cell, totalCellNumber),
				toBeRevealed:    totalCellNumber - tt.args.blackHolesNumber,
				rows:            tt.args.rows,
				cols:            tt.args.cols,
			}
			actual := b.initBoard()

			assert.True(t, tt.verifyFn(actual))
		})
	}
}

func TestBoard_revealCells(t *testing.T) {
	type fields struct {
		rows               int
		cols               int
		blackHolesNumber   int
		blackHoleLocations [][]int
	}
	type args struct {
		cellID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		want    [][]*cell
	}{
		{
			name: "success_one_opened_cell_that_touches_black_hole",
			args: args{
				cellID: cellIdentificationKey(1, 2),
			},
			fields: fields{
				cols: 3,
				rows: 3,
				blackHoleLocations: [][]int{
					{0, 1},
					{2, 2},
				},
			},
			want: [][]*cell{
				{
					{
						state: closedState,
						value: 1,
						x:     0,
						y:     0,
					},
					{
						state: closedState,
						value: blackHole,
						x:     0,
						y:     1,
					},
					{
						state: closedState,
						value: 1,
						x:     0,
						y:     2,
					},
				},
				{
					{
						state: closedState,
						value: 1,
						x:     1,
						y:     0,
					},
					{
						state: closedState,
						value: 2,
						x:     1,
						y:     1,
					},
					{
						state: openedState,
						value: 2,
						x:     1,
						y:     2,
					},
				},
				{
					{
						state: closedState,
						value: void,
						x:     2,
						y:     0,
					},
					{
						state: closedState,
						value: 1,
						x:     2,
						y:     1,
					},
					{
						state: closedState,
						value: blackHole,
						x:     2,
						y:     2,
					},
				},
			},
		},
		{
			name: "success_one_opened_void_cell",
			args: args{
				cellID: cellIdentificationKey(2, 2),
			},
			fields: fields{
				cols: 3,
				rows: 3,
				blackHoleLocations: [][]int{
					{0, 1},
				},
			},
			want: [][]*cell{
				{
					{
						state: closedState,
						value: 1,
						x:     0,
						y:     0,
					},
					{
						state: closedState,
						value: blackHole,
						x:     0,
						y:     1,
					},
					{
						state: closedState,
						value: 1,
						x:     0,
						y:     2,
					},
				},
				{
					{
						state: openedState,
						value: 1,
						x:     1,
						y:     0,
					},
					{
						state: openedState,
						value: 1,
						x:     1,
						y:     1,
					},
					{
						state: openedState,
						value: 1,
						x:     1,
						y:     2,
					},
				},
				{
					{
						state: openedState,
						value: void,
						x:     2,
						y:     0,
					},
					{
						state: openedState,
						value: void,
						x:     2,
						y:     1,
					},
					{
						state: openedState,
						value: void,
						x:     2,
						y:     2,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			totalCellNumber := tt.fields.rows * tt.fields.rows
			b := &Board{
				adjacencyList:   make(map[string][]*cell),
				sideCellsNumber: tt.fields.rows,
				cellList:        make(map[string]*cell, totalCellNumber),
				toBeRevealed:    totalCellNumber - tt.fields.blackHolesNumber,
				rows:            tt.fields.rows,
				cols:            tt.fields.cols,
			}
			b.board = b.generateBoard(tt.fields.blackHoleLocations)
			b.buildGraph(b.board)

			err := b.revealCells(tt.args.cellID)
			if tt.wantErr {
				return
			}

			assert.Equal(t, b.board, tt.want)
			require.NoError(t, err)
		})
	}
}

func TestBoard_Click(t *testing.T) {
	type fields struct {
		rows               int
		cols               int
		blackHolesNumber   int
		blackHoleLocations [][]int
	}
	type args struct {
		click []int
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantErr     bool
		expectedErr error
		want        [][]*cell
		setupFn     func(b *Board)
	}{
		{
			name: "success",
			fields: fields{
				cols: 3,
				rows: 3,
				blackHoleLocations: [][]int{
					{0, 1},
				},
			},
			args: args{
				click: []int{2, 2},
			},
			want: [][]*cell{
				{
					{
						state: closedState,
						value: 1,
						x:     0,
						y:     0,
					},
					{
						state: closedState,
						value: blackHole,
						x:     0,
						y:     1,
					},
					{
						state: closedState,
						value: 1,
						x:     0,
						y:     2,
					},
				},
				{
					{
						state: openedState,
						value: 1,
						x:     1,
						y:     0,
					},
					{
						state: openedState,
						value: 1,
						x:     1,
						y:     1,
					},
					{
						state: openedState,
						value: 1,
						x:     1,
						y:     2,
					},
				},
				{
					{
						state: openedState,
						value: void,
						x:     2,
						y:     0,
					},
					{
						state: openedState,
						value: void,
						x:     2,
						y:     1,
					},
					{
						state: openedState,
						value: void,
						x:     2,
						y:     2,
					},
				},
			},
			setupFn: func(b *Board) {
				return
			},
		},
		{
			name: "error_click_out_of_bounds",
			fields: fields{
				cols: 3,
				rows: 3,
				blackHoleLocations: [][]int{
					{0, 1},
				},
			},
			args: args{
				click: []int{5, 2},
			},
			wantErr:     true,
			expectedErr: errors.New("click coordinate [5 2] is out of board bounds 3 x 3"),
			setupFn: func(b *Board) {
				return
			},
		},
		{
			name: "error_already_opened",
			fields: fields{
				cols: 3,
				rows: 3,
				blackHoleLocations: [][]int{
					{0, 1},
				},
			},
			args: args{
				click: []int{2, 2},
			},
			wantErr:     true,
			expectedErr: errCellOpened,
			setupFn: func(b *Board) {
				// simulation of click that already took place (in some previous clicks) in order to get error of already clicked cell
				b.board[2][2].state.setToOpened()
				return
			},
		},
		{
			name: "click_on_black_hole",
			fields: fields{
				cols: 3,
				rows: 3,
				blackHoleLocations: [][]int{
					{0, 1},
				},
			},
			args: args{
				click: []int{0, 1},
			},
			setupFn: func(b *Board) {
				return
			},
			want: [][]*cell{
				{
					{
						state: openedState,
						value: 1,
						x:     0,
						y:     0,
					},
					{
						state: blackHoledState,
						value: blackHole,
						x:     0,
						y:     1,
					},
					{
						state: openedState,
						value: 1,
						x:     0,
						y:     2,
					},
				},
				{
					{
						state: openedState,
						value: 1,
						x:     1,
						y:     0,
					},
					{
						state: openedState,
						value: 1,
						x:     1,
						y:     1,
					},
					{
						state: openedState,
						value: 1,
						x:     1,
						y:     2,
					},
				},
				{
					{
						state: openedState,
						value: void,
						x:     2,
						y:     0,
					},
					{
						state: openedState,
						value: void,
						x:     2,
						y:     1,
					},
					{
						state: openedState,
						value: void,
						x:     2,
						y:     2,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			totalCellNumber := tt.fields.rows * tt.fields.rows
			b := &Board{
				adjacencyList:   make(map[string][]*cell),
				sideCellsNumber: tt.fields.rows,
				cellList:        make(map[string]*cell, totalCellNumber),
				toBeRevealed:    totalCellNumber - tt.fields.blackHolesNumber,
				rows:            tt.fields.rows,
				cols:            tt.fields.cols,
			}
			b.board = b.generateBoard(tt.fields.blackHoleLocations)
			b.buildGraph(b.board)
			tt.setupFn(b)

			err := b.Click(tt.args.click)
			if tt.wantErr {
				assert.Equal(t, tt.expectedErr, err)
				return
			}

			assert.Equal(t, b.board, tt.want)
			require.NoError(t, err)
		})
	}
}

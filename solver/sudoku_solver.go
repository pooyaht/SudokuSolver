package sudokusolver

import (
	"math"
)

type SudokuSolver struct {
	board    SudokuBoard
	frontier []SudokuBoard
	explored []SudokuBoard
}

type SudokuBoard [][]Cell

type Cell struct {
	Value          int
	PossibleValues []int
	Row            int
	Column         int
}

func (c *Cell) removeValue(value int) {
	for i, v := range c.PossibleValues {
		if v == value {
			c.PossibleValues = append(c.PossibleValues[:i], c.PossibleValues[i+1:]...)
			break
		}
	}
}

func (sb SudokuBoard) getSubBoardBoundingBox(row int, column int) (int, int, int, int) {
	n := len(sb)
	subBoardSize := int(math.Sqrt(float64(n)))

	subRowStart := (row / subBoardSize) * subBoardSize
	subColStart := (column / subBoardSize) * subBoardSize
	subRowEnd := subRowStart + subBoardSize - 1
	subColEnd := subColStart + subBoardSize - 1

	return subRowStart, subColStart, subRowEnd, subColEnd
}

func NewSudokuSolver(board [][]int) SudokuSolver {
	sb := convertBoardToCells(board)
	return SudokuSolver{
		board: sb,
		frontier: []SudokuBoard{
			sb,
		},
		explored: []SudokuBoard{},
	}
}

func convertBoardToCells(board [][]int) SudokuBoard {
	sudokuBoard := make(SudokuBoard, len(board))
	for rowIndex, row := range board {
		sudokuBoard[rowIndex] = make([]Cell, len(row))
		for colIndex, value := range row {
			var possibleValues []int
			if value == 0 {
				possibleValues = make([]int, len(board))
				for i := 1; i <= len(board); i++ {
					possibleValues[i-1] = i
				}
			}
			sudokuBoard[rowIndex][colIndex] = Cell{
				Value:          value,
				PossibleValues: possibleValues,
				Row:            rowIndex,
				Column:         colIndex,
			}
		}
	}
	return sudokuBoard
}

func (ss SudokuSolver) GetBoard() SudokuBoard {
	return ss.board
}

func (ss *SudokuSolver) ReducePossibleValues(cell Cell, value int) {
	// remove from same column
	for _, row := range ss.board {
		if row[cell.Column].Row != cell.Row {
			row[cell.Column].removeValue(value)
		}
	}
	// remove from same row
	for i, col := range ss.board[cell.Row] {
		if i != cell.Column {
			col.removeValue(value)
		}
	}
	// remove from same subgrid
	subRowStart, subColStart, subRowEnd, subColEnd := ss.board.getSubBoardBoundingBox(cell.Row, cell.Column)
	for r := subRowStart; r <= subRowEnd; r++ {
		for c := subColStart; c <= subColEnd; c++ {
			if ss.board[r][c].Row != cell.Row || ss.board[r][c].Column != cell.Column {
				ss.board[r][c].removeValue(value)
			}
		}
	}
}

package solver

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
	value          int
	possibleValues []int
	row            int
	column         int
}

func (c *Cell) removeValue(value int) {
	for i, v := range c.possibleValues {
		if v == value {
			c.possibleValues = append(c.possibleValues[:i], c.possibleValues[i+1:]...)
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
				value:          value,
				possibleValues: possibleValues,
				row:            rowIndex,
				column:         colIndex,
			}
		}
	}
	return sudokuBoard
}

func (ss SudokuSolver) GetBoard() SudokuBoard {
	return ss.board
}

func (ss *SudokuSolver) reducePossibleValues(cell Cell, value int) {
	// remove from same column
	for _, row := range ss.board {
		if row[cell.column].row != cell.row {
			row[cell.column].removeValue(value)
		}
	}
	// remove from same row
	for i, col := range ss.board[cell.row] {
		if i != cell.column {
			col.removeValue(value)
		}
	}
	// remove from same subgrid
	subRowStart, subColStart, subRowEnd, subColEnd := ss.board.getSubBoardBoundingBox(cell.row, cell.column)
	for r := subRowStart; r <= subRowEnd; r++ {
		for c := subColStart; c <= subColEnd; c++ {
			if ss.board[r][c].row != cell.row || ss.board[r][c].column != cell.column {
				ss.board[r][c].removeValue(value)
			}
		}
	}
}

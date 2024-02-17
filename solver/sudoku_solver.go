package sudokusolver

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

package solver

type SequentialSudokuSolver struct {
	initialBoard SudokuBoard
}

func NewSequentialSudokuSolver(board [][]int) *SequentialSudokuSolver {
	cells := ConvertBoardToCells(board)
	sb := NewSudokuBoard(cells)
	return &SequentialSudokuSolver{
		initialBoard: sb,
	}
}

func (ss *SequentialSudokuSolver) Solve() [][]int {
	frontier := make([]SudokuBoard, 0)
	explored := NewStandardMap[string, bool]()
	frontier = append(frontier, ss.initialBoard)

	result := solveUtil(&frontier, &explored)
	if result == nil {
		return nil
	}
	return result
}

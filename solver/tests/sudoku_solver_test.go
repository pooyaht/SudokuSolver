package solver_test

import (
	_ "embed"
	"testing"

	solver "github.com/pooyaht/SudokuSolver/solver"
)

//go:embed sudoku_test_input.txt
var testInput string

func TestSudokuSolver(t *testing.T) {
	inputTestGrids, _ := solver.ParseSudoku(testInput, 9)
	for i := range inputTestGrids {
		ss := solver.NewSudokuSolver(inputTestGrids[i])
		output := ss.Solve()
		sb := solver.NewSudokuBoard(solver.ConvertBoardToCells(output))
		if !sb.IsSolved() {
			t.Errorf("%v is not a valid answer", output)
		}
	}

}

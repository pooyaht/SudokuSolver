package solver_test

import (
	_ "embed"
	"testing"

	solver "github.com/pooyaht/SudokuSolver/solver"
)

//go:embed sudoku_test_input_passes.txt
var testInputPasses string

//go:embed sudoku_test_input_fails.txt
var testInputFails string

func TestSequentialSudokuSolver(t *testing.T) {
	solvableTestGrids, _ := solver.ParseSudoku(testInputPasses)
	unsolvableTestGrids, _ := solver.ParseSudoku(testInputFails)
	for i := range solvableTestGrids {
		ss := solver.NewSequentialSudokuSolver(solvableTestGrids[i])
		output := ss.Solve()
		sb := solver.NewSudokuBoard(solver.ConvertBoardToCells(output))
		if !sb.IsSolved() {
			t.Errorf("%v is not a valid answer for %v", output, solvableTestGrids[i])
		}
	}

	for i := range unsolvableTestGrids {
		ss := solver.NewSequentialSudokuSolver(unsolvableTestGrids[i])
		output := ss.Solve()
		if output != nil {
			t.Errorf("%v is solved but it is unsolvable", unsolvableTestGrids[i])
		}
	}
}

func TestParallelSudokuSolver(t *testing.T) {
	solvableTestGrids, _ := solver.ParseSudoku(testInputPasses)
	unsolvableTestGrids, _ := solver.ParseSudoku(testInputFails)
	for i := range solvableTestGrids {
		ps := solver.NewParallelSudokuSolver(solvableTestGrids[i], 4)
		output := ps.Solve()
		sb := solver.NewSudokuBoard(solver.ConvertBoardToCells(output))
		if !sb.IsSolved() {
			t.Errorf("%v is not a valid answer for %v", output, solvableTestGrids[i])
		}
	}

	for i := range unsolvableTestGrids {
		ps := solver.NewParallelSudokuSolver(unsolvableTestGrids[i], 4)
		output := ps.Solve()
		if output != nil {
			t.Errorf("%v is solved but it is unsolvable", unsolvableTestGrids[i])
		}
	}
}

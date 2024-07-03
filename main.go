package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/pooyaht/SudokuSolver/solver"
)

func main() {
	solverType := flag.String("solver", "sequential", "type of sudoku solver [sequential, parallelbfs, paralleldfs]")
	numWorkers := flag.Int("num_workers", 1, "num of concurrent workers")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: command -solver [sequential, parallelbfs, paralleldfs] -num_workers path_to_sudoku_file\n")
		flag.CommandLine.PrintDefaults()
	}
	flag.Parse()
	if len(flag.Args()) != 1 {
		return
	}

	sudokuFileContent, err := os.ReadFile(flag.Args()[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Sudoku's filepath is invalid")
		return
	}
	sudokuFileContentString := string(sudokuFileContent)
	sudokus, err := solver.ParseSudoku(sudokuFileContentString)
	if err != nil {
		return
	}

	sudokuSolver := solverFactory(solverType, numWorkers)
	for i, sudoku := range sudokus {
		solvedBoard := sudokuSolver(sudoku).Solve()
		if solvedBoard == nil {
			fmt.Printf("Sudoku %d is unsolvable\n", i)
		} else {
			fmt.Printf("Sudoku %d is solved\n", i)
			printBoard(solvedBoard)
		}
	}
}

type sodokuSolver interface {
	Solve() [][]int
}

func solverFactory(solverType *string, numWorkers *int) func([][]int) sodokuSolver {
	switch *solverType {
	case "sequential":
		return func(board [][]int) sodokuSolver { return solver.NewSequentialSudokuSolver(board) }
	case "parallelbfs":
		return func(board [][]int) sodokuSolver { return solver.NewParallelBFSSudokuSolver(board, *numWorkers) }
	case "paralleldfs":
		return func(board [][]int) sodokuSolver { return solver.NewParallelDFSSudokuSolver(board, *numWorkers) }
	}

	return func(board [][]int) sodokuSolver { return solver.NewSequentialSudokuSolver(board) }
}
func printBoard(board [][]int) {
	for _, row := range board {
		for _, cell := range row {
			fmt.Printf("%d ", cell)
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")
}

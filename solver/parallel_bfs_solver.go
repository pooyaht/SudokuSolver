package solver

import "fmt"

type ParallelBFSSudokuSolver struct {
	initialBoard SudokuBoard
	explored     ConcurrentMap[string, bool]
	numWorkers   int
}

func NewParallelBFSSudokuSolver(board [][]int, numWorkers int) ParallelBFSSudokuSolver {
	cells := ConvertBoardToCells(board)
	sb := NewSudokuBoard(cells)
	return ParallelBFSSudokuSolver{
		initialBoard: sb,
		explored:     NewConcurrentMap[string, bool](),
		numWorkers:   numWorkers,
	}
}

func (ps *ParallelBFSSudokuSolver) Solve() [][]int {
	board := ps.initialBoard
	result := make(chan [][]int, ps.numWorkers)

	for i := 0; i < ps.numWorkers; i++ {
		frontier := make([]SudokuBoard, 0)
		index := board.selectCellWithMaximumPossibleValues()
		updatedFrontier := extendFrontier(&frontier, board, index)
		if len(*updatedFrontier) == 0 {
			result <- nil
			continue
		}
		board = (*updatedFrontier)[0]
		go func() {
			res := solveUtil(updatedFrontier, &ps.explored)
			result <- res
		}()
	}

	for i := 0; i < ps.numWorkers; i++ {
		solvedBoard := <-result
		if solvedBoard != nil {
			fmt.Println(solvedBoard)
			return solvedBoard
		}
	}
	fmt.Println("No solution found!")
	return nil
}

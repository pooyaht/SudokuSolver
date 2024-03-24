package solver

import (
	"fmt"
)

type SequentialSudokuSolver struct {
	initialBoard SudokuBoard
}

func NewSequentialSudokuSolver(board [][]int) SequentialSudokuSolver {
	cells := ConvertBoardToCells(board)
	sb := NewSudokuBoard(cells)
	return SequentialSudokuSolver{
		initialBoard: sb,
	}
}

func (ss *SequentialSudokuSolver) Solve() [][]int {
	frontier := make([]SudokuBoard, 0)
	explored := NewStandardMap[string, bool]()
	frontier = append(frontier, ss.initialBoard)

	result := solve(&frontier, &explored)
	if result == nil {
		fmt.Println("No solution found!")
		return nil
	}
	fmt.Println(result)
	return result
}

type ParallelSudokuSolver struct {
	initialBoard SudokuBoard
	explored     ConcurrentMap[string, bool]
	numWorkers   int
}

func NewParallelSudokuSolver(board [][]int, numWorkers int) ParallelSudokuSolver {
	cells := ConvertBoardToCells(board)
	sb := NewSudokuBoard(cells)
	return ParallelSudokuSolver{
		initialBoard: sb,
		explored:     NewConcurrentMap[string, bool](),
		numWorkers:   numWorkers,
	}
}

func (ps *ParallelSudokuSolver) Solve() [][]int {
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
			res := solve(updatedFrontier, &ps.explored)
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

func solve(frontier *[]SudokuBoard, explored Map[string, bool]) [][]int {
	for len(*frontier) > 0 {
		currentBoard := (*frontier)[len(*frontier)-1]
		*frontier = (*frontier)[:len(*frontier)-1]

		currentBoardStr := currentBoard.string()
		if _, ok := explored.Load(currentBoardStr); ok {
			continue
		}
		explored.Store(currentBoardStr, true)

		if len(currentBoard.unsolved) == 0 {
			return currentBoard.toArray()
		}

		index := currentBoard.selectCellWithFewestPossibleValues()
		frontier = extendFrontier(frontier, currentBoard, index)
	}
	return nil
}

func extendFrontier(frontier *[]SudokuBoard, board SudokuBoard, index Index) *[]SudokuBoard {
	cell := board.grid[index.row][index.column]
	for _, value := range cell.possibleValues {
		newBoard := DeepCopy(board)
		newBoard.grid[index.row][index.column].value = value
		delete(newBoard.unsolved, index)
		newBoard.reducePossibleValues(index, value)
		if newBoard.isValid() {
			*frontier = append(*frontier, newBoard)
		}
	}
	return frontier
}

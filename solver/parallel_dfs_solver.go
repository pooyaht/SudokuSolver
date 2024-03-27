package solver

import (
	"fmt"
	"sync"
)

type ParallelDFSsSudokuSolver struct {
	initialBoard SudokuBoard
	explored     ConcurrentMap[string, bool]
	frontier     *[]SudokuBoard
	numWorkers   int
}

func NewParallelDFSSudokuSolver(board [][]int, numWorkers int) ParallelDFSsSudokuSolver {
	cells := ConvertBoardToCells(board)
	sb := NewSudokuBoard(cells)
	return ParallelDFSsSudokuSolver{
		initialBoard: sb,
		explored:     NewConcurrentMap[string, bool](),
		frontier:     &[]SudokuBoard{sb},
		numWorkers:   numWorkers,
	}
}

func (ps *ParallelDFSsSudokuSolver) Solve() [][]int {
	result := make(chan [][]int, ps.numWorkers)
	wg := sync.WaitGroup{}
	cv := sync.NewCond(&sync.Mutex{})
	wg.Add(ps.numWorkers)
	for i := 0; i < ps.numWorkers; i++ {
		go ps.solveUtil(cv, &wg, result)
	}

	wg.Wait()
	cv.Broadcast()
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

func (ps *ParallelDFSsSudokuSolver) solveUtil(cv *sync.Cond, wg *sync.WaitGroup, res chan<- [][]int) {
	for {
		cv.L.Lock()
		if len(*ps.frontier) == 0 {
			wg.Done()
			cv.Wait()
			wg.Add(1)
		}
		// Hack to findout whether all goroutines were blocked on the condvar
		// TODO find a better solution
		if len(*ps.frontier) == 0 {
			wg.Done()
			cv.L.Unlock()
			res <- nil
			return
		}
		currentBoard := (*ps.frontier)[len(*ps.frontier)-1]
		(*ps.frontier) = (*ps.frontier)[:len(*ps.frontier)-1]
		cv.L.Unlock()

		currentBoardStr := currentBoard.string()
		if _, ok := ps.explored.Load(currentBoardStr); ok {
			continue
		}
		ps.explored.Store(currentBoardStr, true)

		if len(currentBoard.unsolved) == 0 {
			res <- currentBoard.toArray()
			wg.Done()
			return
		}

		index := currentBoard.selectCellWithFewestPossibleValues()

		cv.L.Lock()
		ps.frontier = extendFrontier(ps.frontier, currentBoard, index)
		if len(*ps.frontier) > 0 {
			cv.Signal()
		}
		cv.L.Unlock()
	}
}

package solver

import "fmt"

type SudokuSolver struct {
	frontier []SudokuBoard
	explored map[string]bool
}

func NewSudokuSolver(board [][]int) SudokuSolver {
	cells := ConvertBoardToCells(board)
	sb := NewSudokuBoard(cells)
	return SudokuSolver{
		frontier: []SudokuBoard{
			sb,
		},
		explored: make(map[string]bool),
	}
}

func (ss *SudokuSolver) Solve() [][]int {
	for len(ss.frontier) > 0 {
		currentBoard := ss.frontier[len(ss.frontier)-1]
		ss.frontier = ss.frontier[:len(ss.frontier)-1]

		currentBoardStr := currentBoard.string()
		if ss.explored[currentBoardStr] {
			continue
		}
		ss.explored[currentBoardStr] = true

		if len(currentBoard.unsolved) == 0 {
			fmt.Println(currentBoard.toArray())
			return currentBoard.toArray()
		}

		index := currentBoard.selectCellWithFewestPossibleValues()
		cell := currentBoard.grid[index.row][index.column]
		for _, value := range cell.possibleValues {
			newBoard := DeepCopy(currentBoard)
			newBoard.grid[index.row][index.column].value = value
			delete(newBoard.unsolved, index)
			newBoard.reducePossibleValues(index, value)
			if newBoard.isValid() {
				ss.frontier = append(ss.frontier, newBoard)
			}
		}
	}
	return nil
}

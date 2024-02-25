package solver

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func DeepCopy(sb SudokuBoard) SudokuBoard {
	newGrid := make([][]Cell, len(sb.grid))
	for rowIndex, row := range sb.grid {
		newGrid[rowIndex] = make([]Cell, len(row))
		for colIndex, cell := range row {
			newGrid[rowIndex][colIndex] = Cell{
				value:          cell.value,
				possibleValues: append([]int(nil), cell.possibleValues...),
			}
		}
	}

	newUnsolved := make(map[Index]Cell, len(sb.unsolved))
	for key, value := range sb.unsolved {
		newUnsolved[key] = value
	}

	return SudokuBoard{
		grid:     newGrid,
		unsolved: newUnsolved,
	}
}

func ConvertBoardToCells(board [][]int) [][]Cell {
	if board == nil {
		return nil
	}

	sudokuBoard := make([][]Cell, len(board))
	for rowIndex, row := range board {
		sudokuBoard[rowIndex] = make([]Cell, len(row))
		for colIndex, value := range row {
			sudokuBoard[rowIndex][colIndex] = Cell{
				value:          value,
				possibleValues: CalcPossibleValues(board, Index{rowIndex, colIndex}),
			}
		}
	}
	return sudokuBoard
}

func CalcPossibleValues(grid [][]int, index Index) []int {
	if grid[index.row][index.column] != 0 {
		return nil
	}

	gridSize := len(grid)
	possibleValues := make([]int, gridSize+1)
	for i := 1; i <= gridSize; i++ {
		possibleValues[i] = i
	}

	for i := 0; i < gridSize; i++ {
		if grid[index.row][i] != 0 {
			possibleValues[grid[index.row][i]] = 0
		}
		if grid[i][index.column] != 0 {
			possibleValues[grid[i][index.column]] = 0
		}
	}

	subBoardSize := int(math.Sqrt(float64(gridSize)))
	subRowStart, subColStart := (index.row/subBoardSize)*subBoardSize, (index.column/subBoardSize)*subBoardSize
	for r := subRowStart; r < subRowStart+subBoardSize; r++ {
		for c := subColStart; c < subColStart+subBoardSize; c++ {
			if grid[r][c] != 0 {
				possibleValues[grid[r][c]] = 0
			}
		}
	}

	var result []int
	for i := 1; i <= gridSize; i++ {
		if possibleValues[i] != 0 {
			result = append(result, i)
		}
	}

	return result
}

func ParseSudoku(sudokuString string) ([][][]int, error) {
	var grids [][][]int
	var currentGrid [][]int

	lines := strings.Split(sudokuString, "\n")

	for _, line := range lines {
		if line == "" {
			if len(currentGrid) > 0 {
				grids = append(grids, currentGrid)
				currentGrid = nil
			}
			continue
		}

		cells := strings.Fields(line)
		row := make([]int, len(cells))
		for i, cell := range cells {
			if cell == "." || cell == "0" {
				row[i] = 0
			} else {
				value, err := strconv.Atoi(cell)
				if err != nil || value < 1 || value > len(cells) {
					return nil, fmt.Errorf("failed to parse cell value: %v", err)
				}
				row[i] = value
			}
		}
		currentGrid = append(currentGrid, row)
	}

	if len(currentGrid) > 0 {
		grids = append(grids, currentGrid)
	}

	return grids, nil
}

package solver

import (
	"math"
	"strconv"
	"strings"
)

type SudokuBoard struct {
	grid     [][]Cell
	unsolved map[Index]Cell
}

type Cell struct {
	value          int
	possibleValues []int
}

type Index struct {
	row    int
	column int
}

func (c *Cell) removeValue(value int) {
	for i := len(c.possibleValues) - 1; i >= 0; i-- {
		if c.possibleValues[i] == value {
			c.possibleValues = append(c.possibleValues[:i], c.possibleValues[i+1:]...)
		}
	}
}

func NewSudokuBoard(cells [][]Cell) SudokuBoard {
	unsolved := make(map[Index]Cell)
	for rowIndex, row := range cells {
		for colIndex, cell := range row {
			if cell.value == 0 {
				unsolved[Index{rowIndex, colIndex}] = cell
			}
		}
	}
	return SudokuBoard{
		grid:     cells,
		unsolved: unsolved,
	}
}

func (sb SudokuBoard) getSubBoardBoundingBox(row int, column int) (int, int, int, int) {
	n := len(sb.grid)
	subBoardSize := int(math.Sqrt(float64(n)))

	subRowStart := (row / subBoardSize) * subBoardSize
	subColStart := (column / subBoardSize) * subBoardSize
	subRowEnd := subRowStart + subBoardSize - 1
	subColEnd := subColStart + subBoardSize - 1

	return subRowStart, subColStart, subRowEnd, subColEnd
}

func (sb *SudokuBoard) reducePossibleValues(index Index, value int) {
	// remove from same column
	for i := range sb.grid {
		if i != index.row {
			sb.grid[i][index.column].removeValue(value)
		}
	}
	// remove from same row
	for i := range sb.grid[index.row] {
		if i != index.column {
			sb.grid[index.row][i].removeValue(value)
		}
	}
	// remove from same subgrid
	subRowStart, subColStart, subRowEnd, subColEnd := sb.getSubBoardBoundingBox(index.row, index.column)
	for r := subRowStart; r <= subRowEnd; r++ {
		for c := subColStart; c <= subColEnd; c++ {
			if r != index.row || c != index.column {
				sb.grid[r][c].removeValue(value)
			}
		}
	}
}

func (sb SudokuBoard) toArray() [][]int {
	var sbArray [][]int
	for _, row := range sb.grid {
		var sbRow []int
		for _, cell := range row {
			sbRow = append(sbRow, cell.value)
		}
		sbArray = append(sbArray, sbRow)
	}
	return sbArray
}

func (sb SudokuBoard) string() string {
	var sbStr strings.Builder
	for _, row := range sb.grid {
		for _, cell := range row {
			sbStr.WriteString(strconv.Itoa(cell.value))
		}
	}
	return sbStr.String()
}

func (sb SudokuBoard) IsSolved() bool {
	return len(sb.unsolved) == 0 && sb.isValid()
}

func (sb SudokuBoard) isValid() bool {
	size := len(sb.grid)
	subGridSize := int(math.Sqrt(float64(size)))

	// Check uniqueness of values in each row
	for _, row := range sb.grid {
		rowSet := make(map[int]bool)
		for _, cell := range row {
			if cell.value != 0 {
				if rowSet[cell.value] {
					return false
				}
				rowSet[cell.value] = true
			}
		}
	}

	// Check uniqueness of values in each column
	for col := 0; col < size; col++ {
		colSet := make(map[int]bool)
		for row := 0; row < size; row++ {
			cell := sb.grid[row][col]
			if cell.value != 0 {
				if colSet[cell.value] {
					return false
				}
				colSet[cell.value] = true
			}
		}
	}

	// Check uniqueness of values in each sub-grid
	for boxStartRow := 0; boxStartRow < size; boxStartRow += subGridSize {
		for boxStartCol := 0; boxStartCol < size; boxStartCol += subGridSize {
			boxSet := make(map[int]bool)
			for row := boxStartRow; row < boxStartRow+subGridSize; row++ {
				for col := boxStartCol; col < boxStartCol+subGridSize; col++ {
					cell := sb.grid[row][col]
					if cell.value != 0 {
						if boxSet[cell.value] {
							return false
						}
						boxSet[cell.value] = true
					}
				}
			}
		}
	}

	return true
}

func (sb SudokuBoard) selectCellWithFewestPossibleValues() Index {
	minPossibleValues := len(sb.grid) + 1
	var minIndex Index
	for index, cell := range sb.unsolved {
		if len(cell.possibleValues) < minPossibleValues {
			minPossibleValues = len(cell.possibleValues)
			minIndex = index
		}
	}
	return minIndex
}

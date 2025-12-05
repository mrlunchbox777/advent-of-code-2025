package main

import (
	"bufio"
	"os"
)

// Position represents a coordinate in the grid with 1-based indexing (bottom-left is [1,1])
type Position struct {
	X int // column (1-based, left to right)
	Y int // row (1-based, bottom to top)
}

// Grid represents a 2D grid of symbols
type Grid struct {
	Cells  [][]rune // cells[row][col] where row 0 is the top
	Width  int
	Height int
}

// NewGridFromFile reads a file and creates a Grid
func NewGridFromFile(filepath string) (*Grid, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cells [][]rune
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		row := []rune(line)
		cells = append(cells, row)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	height := len(cells)
	width := 0
	if height > 0 {
		width = len(cells[0])
	}

	return &Grid{
		Cells:  cells,
		Width:  width,
		Height: height,
	}, nil
}

// GetCell returns the rune at the given position (1-based coordinates)
func (g *Grid) GetCell(pos Position) rune {
	// Convert 1-based position to 0-based array indices
	// Y=1 is bottom (last row), Y=Height is top (first row)
	row := g.Height - pos.Y
	col := pos.X - 1

	if row < 0 || row >= g.Height || col < 0 || col >= g.Width {
		return '.' // out of bounds is treated as empty
	}

	return g.Cells[row][col]
}

// CountAdjacentAt counts how many '@' symbols are adjacent to the given position
func (g *Grid) CountAdjacentAt(pos Position) int {
	// Define the 8 adjacent directions
	directions := []struct{ dx, dy int }{
		{-1, 0},  // left
		{-1, 1},  // left-up
		{0, 1},   // up
		{1, 1},   // right-up
		{1, 0},   // right
		{1, -1},  // right-down
		{0, -1},  // down
		{-1, -1}, // left-down
	}

	count := 0
	for _, dir := range directions {
		neighbor := Position{X: pos.X + dir.dx, Y: pos.Y + dir.dy}
		if g.GetCell(neighbor) == '@' {
			count++
		}
	}

	return count
}

// FindSelectedPositions finds all '@' positions with fewer than 4 adjacent '@' symbols
func (g *Grid) FindSelectedPositions() []Position {
	var selected []Position

	// Iterate through all positions in the grid (1-based coordinates)
	for y := 1; y <= g.Height; y++ {
		for x := 1; x <= g.Width; x++ {
			pos := Position{X: x, Y: y}
			if g.GetCell(pos) == '@' {
				adjacentCount := g.CountAdjacentAt(pos)
				if adjacentCount < 4 {
					selected = append(selected, pos)
				}
			}
		}
	}

	return selected
}

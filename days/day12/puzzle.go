package main

import (
	"strings"
)

// PieceSpec specifies which piece and how many to use
type PieceSpec struct {
	PieceID int
	Count   int
}

// Puzzle represents a puzzle to solve
type Puzzle struct {
	Width      int
	Height     int
	PieceSpecs []PieceSpec
}

// Solution represents a solved puzzle
type Solution struct {
	Width  int
	Height int
	Grid   [][]byte
}

// String converts the solution to a printable string
func (s *Solution) String() string {
	var sb strings.Builder
	for i := 0; i < s.Height; i++ {
		for j := 0; j < s.Width; j++ {
			sb.WriteByte(s.Grid[i][j])
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

// PlacedPiece represents a piece placed at a specific position and orientation
type PlacedPiece struct {
	Piece       *Piece
	X           int
	Y           int
	DisplayChar byte
}

// PieceToPlace represents a piece ready to be placed
type PieceToPlace struct {
	piece        *Piece
	displayChar  byte
	orientations []*Piece
	filledCount  int
}

// Solve attempts to solve the puzzle
func (p *Puzzle) Solve(pieces map[int]*Piece) *Solution {
	// Create empty grid
	grid := make([][]byte, p.Height)
	for i := 0; i < p.Height; i++ {
		grid[i] = make([]byte, p.Width)
		for j := 0; j < p.Width; j++ {
			grid[i][j] = '.'
		}
	}

	// Build list of pieces to place with their display characters
	var piecesToPlace []PieceToPlace
	displayChar := byte('A')
	
	for _, spec := range p.PieceSpecs {
		piece, exists := pieces[spec.PieceID]
		if !exists {
			return nil
		}
		
		// Pre-compute orientations once per unique piece
		orientations := piece.AllOrientations()
		filledCount := countFilledCells(piece)
		
		for i := 0; i < spec.Count; i++ {
			piecesToPlace = append(piecesToPlace, PieceToPlace{
				piece:        piece,
				displayChar:  displayChar,
				orientations: orientations,
				filledCount:  filledCount,
			})
			displayChar++
		}
	}

	// Try to solve using backtracking
	if p.backtrackOptimized(grid, piecesToPlace, 0) {
		return &Solution{
			Width:  p.Width,
			Height: p.Height,
			Grid:   grid,
		}
	}

	return nil
}

// countFilledCells counts the number of filled cells in a piece
func countFilledCells(piece *Piece) int {
	count := 0
	for i := 0; i < piece.Height; i++ {
		for j := 0; j < piece.Width; j++ {
			if piece.Grid[i][j] {
				count++
			}
		}
	}
	return count
}

// backtrackOptimized recursively tries to place pieces with optimizations
func (p *Puzzle) backtrackOptimized(grid [][]byte, piecesToPlace []PieceToPlace, index int) bool {
	// Base case: all pieces placed
	if index >= len(piecesToPlace) {
		return true
	}

	displayChar := piecesToPlace[index].displayChar

	// Pruning: check if remaining pieces could possibly fit in remaining space
	emptyCount := p.countEmptyCells(grid)
	remainingFilledNeeded := 0
	for i := index; i < len(piecesToPlace); i++ {
		remainingFilledNeeded += piecesToPlace[i].filledCount
	}
	
	// If we need more filled cells than we have empty cells, it's impossible
	if remainingFilledNeeded > emptyCount {
		return false
	}

	// Try all orientations of the current piece (pre-computed)
	for _, oriented := range piecesToPlace[index].orientations {
		// Try all valid positions for this orientation
		for y := 0; y <= p.Height-oriented.Height; y++ {
			for x := 0; x <= p.Width-oriented.Width; x++ {
				// Check if piece can be placed at this position
				if p.canPlace(grid, oriented, x, y) {
					// Place the piece
					p.place(grid, oriented, x, y, displayChar)
					
					// Recurse to place next piece
					if p.backtrackOptimized(grid, piecesToPlace, index+1) {
						return true
					}
					
					// Backtrack: remove the piece
					p.remove(grid, oriented, x, y)
				}
			}
		}
	}

	return false
}

// coversCell checks if a piece at position (px, py) would cover cell (cx, cy)
func (p *Puzzle) coversCell(piece *Piece, px, py, cx, cy int) bool {
	// Check if the target cell is within the piece's bounding box
	if cx < px || cx >= px+piece.Width || cy < py || cy >= py+piece.Height {
		return false
	}
	
	// Check if the piece has a filled cell at that position
	return piece.Grid[cy-py][cx-px]
}

// countEmptyCells counts empty cells in the grid
func (p *Puzzle) countEmptyCells(grid [][]byte) int {
	count := 0
	for i := 0; i < p.Height; i++ {
		for j := 0; j < p.Width; j++ {
			if grid[i][j] == '.' {
				count++
			}
		}
	}
	return count
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// findFirstEmpty finds the first empty cell in the grid (reading left-to-right, top-to-bottom)
func (p *Puzzle) findFirstEmpty(grid [][]byte) (int, int) {
	for y := 0; y < p.Height; y++ {
		for x := 0; x < p.Width; x++ {
			if grid[y][x] == '.' {
				return x, y
			}
		}
	}
	return -1, -1
}

// canPlace checks if a piece can be placed at position (x, y)
func (p *Puzzle) canPlace(grid [][]byte, piece *Piece, x, y int) bool {
	for i := 0; i < piece.Height; i++ {
		for j := 0; j < piece.Width; j++ {
			if piece.Grid[i][j] {
				// Check if this filled cell overlaps with existing filled cell
				if grid[y+i][x+j] != '.' {
					return false
				}
			}
		}
	}
	return true
}

// place places a piece on the grid
func (p *Puzzle) place(grid [][]byte, piece *Piece, x, y int, displayChar byte) {
	for i := 0; i < piece.Height; i++ {
		for j := 0; j < piece.Width; j++ {
			if piece.Grid[i][j] {
				grid[y+i][x+j] = displayChar
			}
		}
	}
}

// remove removes a piece from the grid
func (p *Puzzle) remove(grid [][]byte, piece *Piece, x, y int) {
	for i := 0; i < piece.Height; i++ {
		for j := 0; j < piece.Width; j++ {
			if piece.Grid[i][j] {
				grid[y+i][x+j] = '.'
			}
		}
	}
}

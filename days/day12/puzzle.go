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
	piecesToPlace := []struct {
		piece       *Piece
		displayChar byte
	}{}
	
	displayChar := byte('A')
	for _, spec := range p.PieceSpecs {
		piece, exists := pieces[spec.PieceID]
		if !exists {
			return nil
		}
		for i := 0; i < spec.Count; i++ {
			piecesToPlace = append(piecesToPlace, struct {
				piece       *Piece
				displayChar byte
			}{piece, displayChar})
			displayChar++
		}
	}

	// Try to solve using backtracking
	if p.backtrack(grid, piecesToPlace, 0) {
		return &Solution{
			Width:  p.Width,
			Height: p.Height,
			Grid:   grid,
		}
	}

	return nil
}

// backtrack recursively tries to place pieces
func (p *Puzzle) backtrack(grid [][]byte, piecesToPlace []struct {
	piece       *Piece
	displayChar byte
}, index int) bool {
	// Base case: all pieces placed
	if index >= len(piecesToPlace) {
		return true
	}

	piece := piecesToPlace[index].piece
	displayChar := piecesToPlace[index].displayChar

	// Try all orientations of the current piece
	for _, oriented := range piece.AllOrientations() {
		// Try all positions in the grid
		for y := 0; y <= p.Height-oriented.Height; y++ {
			for x := 0; x <= p.Width-oriented.Width; x++ {
				// Check if piece can be placed at this position
				if p.canPlace(grid, oriented, x, y) {
					// Place the piece
					p.place(grid, oriented, x, y, displayChar)
					
					// Recurse to place next piece
					if p.backtrack(grid, piecesToPlace, index+1) {
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

package main

import (
	"fmt"
	"strings"
)

// Piece represents a puzzle piece with its shape
type Piece struct {
	ID     int
	Grid   [][]bool
	Width  int
	Height int
}

// NewPiece creates a new piece from text lines
func NewPiece(id int, lines []string) (*Piece, error) {
	if len(lines) == 0 {
		return nil, fmt.Errorf("piece has no lines")
	}

	height := len(lines)
	width := len(lines[0])
	
	grid := make([][]bool, height)
	for i, line := range lines {
		grid[i] = make([]bool, width)
		for j, ch := range line {
			if j >= width {
				break
			}
			if ch == '#' {
				grid[i][j] = true
			}
		}
	}

	return &Piece{
		ID:     id,
		Grid:   grid,
		Width:  width,
		Height: height,
	}, nil
}

// AllOrientations returns all 8 possible orientations (4 rotations × 2 flips)
func (p *Piece) AllOrientations() []*Piece {
	orientations := []*Piece{}
	
	current := p
	for rotation := 0; rotation < 4; rotation++ {
		orientations = append(orientations, current)
		orientations = append(orientations, current.Flip())
		current = current.Rotate90()
	}
	
	// Deduplicate identical orientations
	unique := []*Piece{}
	seen := make(map[string]bool)
	for _, orient := range orientations {
		key := orient.Key()
		if !seen[key] {
			seen[key] = true
			unique = append(unique, orient)
		}
	}
	
	return unique
}

// Rotate90 rotates the piece 90 degrees clockwise
func (p *Piece) Rotate90() *Piece {
	newHeight := p.Width
	newWidth := p.Height
	newGrid := make([][]bool, newHeight)
	
	for i := 0; i < newHeight; i++ {
		newGrid[i] = make([]bool, newWidth)
		for j := 0; j < newWidth; j++ {
			newGrid[i][j] = p.Grid[newWidth-1-j][i]
		}
	}
	
	return &Piece{
		ID:     p.ID,
		Grid:   newGrid,
		Width:  newWidth,
		Height: newHeight,
	}
}

// Flip flips the piece horizontally
func (p *Piece) Flip() *Piece {
	newGrid := make([][]bool, p.Height)
	for i := 0; i < p.Height; i++ {
		newGrid[i] = make([]bool, p.Width)
		for j := 0; j < p.Width; j++ {
			newGrid[i][j] = p.Grid[i][p.Width-1-j]
		}
	}
	
	return &Piece{
		ID:     p.ID,
		Grid:   newGrid,
		Width:  p.Width,
		Height: p.Height,
	}
}

// Key returns a string key for deduplication
func (p *Piece) Key() string {
	var sb strings.Builder
	for i := 0; i < p.Height; i++ {
		for j := 0; j < p.Width; j++ {
			if p.Grid[i][j] {
				sb.WriteByte('#')
			} else {
				sb.WriteByte('.')
			}
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

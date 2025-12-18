package main

import (
	"fmt"
	"strconv"
	"strings"
)

// PuzzleData contains all pieces and puzzles from the input file
type PuzzleData struct {
	Pieces  map[int]*Piece
	Puzzles []*Puzzle
}

// ParseInput parses the input lines into pieces and puzzles
func ParseInput(lines []string) (*PuzzleData, error) {
	data := &PuzzleData{
		Pieces:  make(map[int]*Piece),
		Puzzles: []*Puzzle{},
	}

	i := 0
	// Parse pieces
	for i < len(lines) {
		line := strings.TrimSpace(lines[i])
		
		// Check if this is a puzzle line (contains 'x')
		if strings.Contains(line, "x") {
			break
		}
		
		// Check if this is a piece definition
		if strings.Contains(line, ":") && !strings.Contains(line, "x") {
			parts := strings.SplitN(line, ":", 2)
			idStr := strings.TrimSpace(parts[0])
			id, err := strconv.Atoi(idStr)
			if err != nil {
				return nil, fmt.Errorf("invalid piece id: %s", idStr)
			}
			
			// Collect piece lines
			i++
			var pieceLines []string
			for i < len(lines) {
				pieceLine := lines[i]
				if strings.TrimSpace(pieceLine) == "" {
					break
				}
				if strings.Contains(pieceLine, ":") {
					break
				}
				pieceLines = append(pieceLines, pieceLine)
				i++
			}
			
			piece, err := NewPiece(id, pieceLines)
			if err != nil {
				return nil, fmt.Errorf("error parsing piece %d: %v", id, err)
			}
			data.Pieces[id] = piece
		} else {
			i++
		}
	}

	// Parse puzzles
	for i < len(lines) {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			i++
			continue
		}
		
		if strings.Contains(line, "x") {
			puzzle, err := ParsePuzzle(line)
			if err != nil {
				return nil, fmt.Errorf("error parsing puzzle: %v", err)
			}
			data.Puzzles = append(data.Puzzles, puzzle)
		}
		i++
	}

	return data, nil
}

// ParsePuzzle parses a puzzle line like "4x4: 0 0 0 0 2 0"
func ParsePuzzle(line string) (*Puzzle, error) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid puzzle format: %s", line)
	}

	// Parse dimensions
	dimParts := strings.Split(strings.TrimSpace(parts[0]), "x")
	if len(dimParts) != 2 {
		return nil, fmt.Errorf("invalid dimensions: %s", parts[0])
	}
	
	width, err := strconv.Atoi(dimParts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid width: %s", dimParts[0])
	}
	
	height, err := strconv.Atoi(dimParts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid height: %s", dimParts[1])
	}

	// Parse piece specifications
	// Each number represents the count for piece 0, 1, 2, etc. in order
	pieceSpecs := []PieceSpec{}
	numStrs := strings.Fields(strings.TrimSpace(parts[1]))
	for i, numStr := range numStrs {
		count, err := strconv.Atoi(numStr)
		if err != nil {
			return nil, fmt.Errorf("invalid count: %s", numStr)
		}
		
		pieceSpecs = append(pieceSpecs, PieceSpec{
			PieceID: i, // Piece ID is the position in the list
			Count:   count,
		})
	}

	return &Puzzle{
		Width:      width,
		Height:     height,
		PieceSpecs: pieceSpecs,
	}, nil
}

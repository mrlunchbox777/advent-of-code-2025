package main

import (
	"strings"
	"testing"
)

func TestParseInput(t *testing.T) {
	input := `0:
###
##.
##.

1:
###
##.
.##

4x4: 0 0 1`

	lines := strings.Split(input, "\n")
	data, err := ParseInput(lines)
	if err != nil {
		t.Fatalf("ParseInput failed: %v", err)
	}

	if len(data.Pieces) != 2 {
		t.Errorf("Expected 2 pieces, got %d", len(data.Pieces))
	}

	if len(data.Puzzles) != 1 {
		t.Errorf("Expected 1 puzzle, got %d", len(data.Puzzles))
	}

	puzzle := data.Puzzles[0]
	if puzzle.Width != 4 || puzzle.Height != 4 {
		t.Errorf("Expected 4x4 puzzle, got %dx%d", puzzle.Width, puzzle.Height)
	}
	
	// Should have 3 specs (piece 0: 0, piece 1: 0, piece 2: 1)
	if len(puzzle.PieceSpecs) != 3 {
		t.Errorf("Expected 3 piece specs, got %d", len(puzzle.PieceSpecs))
	}
}

func TestPieceRotate(t *testing.T) {
	lines := []string{"##.", ".#."}
	piece, err := NewPiece(0, lines)
	if err != nil {
		t.Fatalf("NewPiece failed: %v", err)
	}

	rotated := piece.Rotate90()
	if rotated.Width != 2 || rotated.Height != 3 {
		t.Errorf("Expected 2x3 after rotation, got %dx%d", rotated.Width, rotated.Height)
	}
}

func TestPieceFlip(t *testing.T) {
	lines := []string{"##.", ".#."}
	piece, err := NewPiece(0, lines)
	if err != nil {
		t.Fatalf("NewPiece failed: %v", err)
	}

	flipped := piece.Flip()
	if flipped.Width != 3 || flipped.Height != 2 {
		t.Errorf("Expected 3x2 after flip, got %dx%d", flipped.Width, flipped.Height)
	}

	// Original: "##." becomes ".##" when flipped horizontally
	if flipped.Grid[0][0] || !flipped.Grid[0][1] || !flipped.Grid[0][2] {
		t.Errorf("First row not flipped correctly: got %v %v %v", 
			flipped.Grid[0][0], flipped.Grid[0][1], flipped.Grid[0][2])
	}
	// Original: ".#." becomes ".#." when flipped horizontally (symmetric)
	if flipped.Grid[1][0] || !flipped.Grid[1][1] || flipped.Grid[1][2] {
		t.Errorf("Second row not flipped correctly")
	}
}

func TestPuzzleSolve(t *testing.T) {
	input := `0:
##
##

2x2: 1`

	lines := strings.Split(input, "\n")
	data, err := ParseInput(lines)
	if err != nil {
		t.Fatalf("ParseInput failed: %v", err)
	}

	if len(data.Puzzles) != 1 {
		t.Fatalf("Expected 1 puzzle, got %d", len(data.Puzzles))
	}

	solution := data.Puzzles[0].Solve(data.Pieces)
	if solution == nil {
		t.Error("Expected to find a solution, got nil")
	}
}

func TestExampleData(t *testing.T) {
	input := `0:
###
##.
##.

1:
###
##.
.##

2:
.##
###
##.

3:
##.
###
##.

4:
###
#..
###

5:
###
.#.
###

4x4: 0 0 0 0 2 0
12x5: 1 0 1 0 2 2
12x5: 1 0 1 0 3 2`

	lines := strings.Split(input, "\n")
	data, err := ParseInput(lines)
	if err != nil {
		t.Fatalf("ParseInput failed: %v", err)
	}

	if len(data.Pieces) != 6 {
		t.Errorf("Expected 6 pieces, got %d", len(data.Pieces))
	}

	if len(data.Puzzles) != 3 {
		t.Errorf("Expected 3 puzzles, got %d", len(data.Puzzles))
	}

	solvedCount := 0
	for _, puzzle := range data.Puzzles {
		// Skip puzzles with no pieces to place
		hasPieces := false
		for _, spec := range puzzle.PieceSpecs {
			if spec.Count > 0 {
				hasPieces = true
				break
			}
		}
		
		if !hasPieces {
			continue
		}
		
		solution := puzzle.Solve(data.Pieces)
		if solution != nil {
			solvedCount++
		}
	}

	if solvedCount != 2 {
		t.Errorf("Expected 2 solved puzzles, got %d", solvedCount)
	}
}

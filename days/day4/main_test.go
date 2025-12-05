package main

import (
	"testing"
)

func TestCountAdjacentAt(t *testing.T) {
	grid := &Grid{
		Cells: [][]rune{
			{'@', '@', '@'},
			{'@', '@', '@'},
			{'@', '@', '@'},
		},
		Width:  3,
		Height: 3,
	}

	tests := []struct {
		pos      Position
		expected int
	}{
		{Position{X: 2, Y: 2}, 8}, // center (row 1, col 1 in 0-based) has 8 neighbors
		{Position{X: 1, Y: 3}, 3}, // top-left corner has 3 neighbors
		{Position{X: 3, Y: 3}, 3}, // top-right corner has 3 neighbors
		{Position{X: 1, Y: 1}, 3}, // bottom-left corner has 3 neighbors
		{Position{X: 3, Y: 1}, 3}, // bottom-right corner has 3 neighbors
		{Position{X: 2, Y: 3}, 5}, // top-middle edge has 5 neighbors
		{Position{X: 2, Y: 1}, 5}, // bottom-middle edge has 5 neighbors
	}

	for _, tt := range tests {
		count := grid.CountAdjacentAt(tt.pos)
		if count != tt.expected {
			t.Errorf("CountAdjacentAt(%v) = %d, expected %d", tt.pos, count, tt.expected)
		}
	}
}

func TestFindSelectedPositions(t *testing.T) {
	grid := &Grid{
		Cells: [][]rune{
			{'.', '@', '@'},
			{'@', '@', '.'},
			{'.', '.', '@'},
		},
		Width:  3,
		Height: 3,
	}

	selected := grid.FindSelectedPositions()

	// Position [2,3] (top-middle '@') has 3 adjacent '@' (left, down, down-left)
	// Position [3,3] (top-right '@') has 2 adjacent '@' (left, down-left)
	// Position [1,2] (middle-left '@') has 3 adjacent '@' (right, right-up, up)
	// Position [2,2] (center '@') has 4 adjacent '@' - NOT selected
	// Position [3,1] (bottom-right '@') has 0 adjacent '@'
	// 4 should be selected (< 4)

	if len(selected) != 4 {
		t.Errorf("Expected 4 selected positions, got %d", len(selected))
		for _, pos := range selected {
			t.Logf("  [%d,%d]", pos.X, pos.Y)
		}
	}
}

func TestExampleData(t *testing.T) {
	grid, err := NewGridFromFile("example-data.txt")
	if err != nil {
		t.Fatalf("Failed to load example-data.txt: %v", err)
	}

	selected := grid.FindSelectedPositions()

	if len(selected) != 13 {
		t.Errorf("Expected 13 selected positions, got %d", len(selected))
		t.Logf("Selected positions:")
		for _, pos := range selected {
			t.Logf("  [%d,%d]", pos.X, pos.Y)
		}
	}
}

func TestGetCell(t *testing.T) {
	grid := &Grid{
		Cells: [][]rune{
			{'@', '.', '@'},
			{'.', '@', '.'},
			{'@', '.', '@'},
		},
		Width:  3,
		Height: 3,
	}

	tests := []struct {
		pos      Position
		expected rune
	}{
		{Position{X: 1, Y: 3}, '@'}, // top-left
		{Position{X: 2, Y: 3}, '.'}, // top-middle
		{Position{X: 3, Y: 3}, '@'}, // top-right
		{Position{X: 2, Y: 2}, '@'}, // center
		{Position{X: 1, Y: 1}, '@'}, // bottom-left
		{Position{X: 2, Y: 1}, '.'}, // bottom-middle
		{Position{X: 3, Y: 1}, '@'}, // bottom-right
		{Position{X: 0, Y: 1}, '.'}, // out of bounds (left)
		{Position{X: 4, Y: 1}, '.'}, // out of bounds (right)
		{Position{X: 1, Y: 0}, '.'}, // out of bounds (below)
		{Position{X: 1, Y: 4}, '.'}, // out of bounds (above)
	}

	for _, tt := range tests {
		cell := grid.GetCell(tt.pos)
		if cell != tt.expected {
			t.Errorf("GetCell(%v) = %c, expected %c", tt.pos, cell, tt.expected)
		}
	}
}

package main

import (
	"os"
	"strings"
	"testing"
)

func TestNewGrid(t *testing.T) {
	lines := []string{
		"..S..",
		".....",
		"..^..",
	}
	
	grid := NewGrid(lines)
	
	if grid.Height != 3 {
		t.Errorf("Expected height 3, got %d", grid.Height)
	}
	
	if grid.Width != 5 {
		t.Errorf("Expected width 5, got %d", grid.Width)
	}
	
	if grid.Cells[0][2] != Start {
		t.Errorf("Expected Start at [0][2], got %c", grid.Cells[0][2])
	}
	
	if grid.Cells[2][2] != Split {
		t.Errorf("Expected Split at [2][2], got %c", grid.Cells[2][2])
	}
}

func TestFindStart(t *testing.T) {
	lines := []string{
		".....",
		"..S..",
		".....",
	}
	
	grid := NewGrid(lines)
	start := grid.FindStart()
	
	if start == nil {
		t.Fatal("Expected to find start position")
	}
	
	if start.Row != 1 || start.Col != 2 {
		t.Errorf("Expected start at (1, 2), got (%d, %d)", start.Row, start.Col)
	}
}

func TestSimpleBeamMovement(t *testing.T) {
	lines := []string{
		"..S..",
		".....",
		".....",
	}
	
	grid := NewGrid(lines)
	
	paths := grid.CountPaths()
	if paths != 1 {
		t.Errorf("Expected 1 path, got %d", paths)
	}
	
	grid = NewGrid(lines)
	grid.ProcessBeams()
	
	if grid.Cells[1][2] != Beam {
		t.Errorf("Expected beam at [1][2], got %c", grid.Cells[1][2])
	}
	
	if grid.Cells[2][2] != Beam {
		t.Errorf("Expected beam at [2][2], got %c", grid.Cells[2][2])
	}
}

func TestBeamSplit(t *testing.T) {
	lines := []string{
		"..S..",
		".....",
		"..^..",
		".....",
	}
	
	grid := NewGrid(lines)
	
	paths := grid.CountPaths()
	if paths != 2 {
		t.Errorf("Expected 2 paths, got %d", paths)
	}
	
	grid = NewGrid(lines)
	grid.ProcessBeams()
	
	if grid.Cells[1][2] != Beam {
		t.Errorf("Expected beam at [1][2], got %c", grid.Cells[1][2])
	}
	
	if grid.Cells[2][1] != Beam {
		t.Errorf("Expected beam at [2][1] (left of split), got %c", grid.Cells[2][1])
	}
	
	if grid.Cells[2][3] != Beam {
		t.Errorf("Expected beam at [2][3] (right of split), got %c", grid.Cells[2][3])
	}
	
	if grid.Cells[3][1] != Beam {
		t.Errorf("Expected beam at [3][1], got %c", grid.Cells[3][1])
	}
	
	if grid.Cells[3][3] != Beam {
		t.Errorf("Expected beam at [3][3], got %c", grid.Cells[3][3])
	}
}

func TestExampleData(t *testing.T) {
	if _, err := os.Stat("example-data-1.txt"); os.IsNotExist(err) {
		t.Skip("example-data-1.txt not found")
	}
	
	grid, err := parseFile("example-data-1.txt")
	if err != nil {
		t.Fatalf("Failed to parse example-data-1.txt: %v", err)
	}
	
	paths := grid.CountPaths()
	if paths != 40 {
		t.Errorf("Expected 40 paths, got %d", paths)
	}
	
	grid, err = parseFile("example-data-1.txt")
	if err != nil {
		t.Fatalf("Failed to parse example-data-1.txt: %v", err)
	}
	
	grid.ProcessBeams()
	
	expectedFile, err := os.ReadFile("example-data-2.txt")
	if err != nil {
		t.Fatalf("Failed to read example-data-2.txt: %v", err)
	}
	
	expected := strings.TrimSpace(string(expectedFile))
	expectedLines := strings.Split(expected, "\n")
	
	if len(expectedLines) != grid.Height {
		t.Fatalf("Expected %d rows, got %d", len(expectedLines), grid.Height)
	}
	
	for i := 0; i < grid.Height; i++ {
		var actual strings.Builder
		for j := 0; j < grid.Width; j++ {
			actual.WriteRune(rune(grid.Cells[i][j]))
		}
		
		if actual.String() != expectedLines[i] {
			t.Errorf("Row %d mismatch:\nExpected: %s\nGot:      %s", 
				i, expectedLines[i], actual.String())
		}
	}
}

func TestMultipleSplits(t *testing.T) {
	lines := []string{
		"...S...",
		".......",
		"..^.^..",
		".......",
	}
	
	grid := NewGrid(lines)
	
	paths := grid.CountPaths()
	if paths != 1 {
		t.Errorf("Expected 1 path, got %d", paths)
	}
	
	grid = NewGrid(lines)
	grid.ProcessBeams()
	
	if grid.Cells[1][3] != Beam {
		t.Errorf("Expected beam at [1][3], got %c", grid.Cells[1][3])
	}
	
	if grid.Cells[2][3] != Beam {
		t.Errorf("Expected beam at [2][3] (between splits), got %c", grid.Cells[2][3])
	}
	
	if grid.Cells[3][3] != Beam {
		t.Errorf("Expected beam at [3][3] (continued down), got %c", grid.Cells[3][3])
	}
}

func TestEdgeSplit(t *testing.T) {
	lines := []string{
		"S....",
		".....",
		"^....",
		".....",
	}
	
	grid := NewGrid(lines)
	
	paths := grid.CountPaths()
	if paths != 1 {
		t.Errorf("Expected 1 path (split at edge, only right side), got %d", paths)
	}
	
	grid = NewGrid(lines)
	grid.ProcessBeams()
	
	if grid.Cells[1][0] != Beam {
		t.Errorf("Expected beam at [1][0], got %c", grid.Cells[1][0])
	}
	
	if grid.Cells[2][1] != Beam {
		t.Errorf("Expected beam at [2][1] (right of edge split), got %c", grid.Cells[2][1])
	}
	
	if grid.Cells[3][1] != Beam {
		t.Errorf("Expected beam at [3][1], got %c", grid.Cells[3][1])
	}
}

func TestCountPathsSimple(t *testing.T) {
	tests := []struct {
		name     string
		lines    []string
		expected int
	}{
		{
			name: "straight line",
			lines: []string{
				"S",
				".",
				".",
			},
			expected: 1,
		},
		{
			name: "single split",
			lines: []string{
				"..S..",
				".....",
				"..^..",
				".....",
			},
			expected: 2,
		},
		{
			name: "double split - tree",
			lines: []string{
				"...S...",
				".......",
				"...^...",
				".......",
				"..^.^..",
				".......",
			},
			expected: 4,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			grid := NewGrid(tt.lines)
			paths := grid.CountPaths()
			if paths != tt.expected {
				t.Errorf("Expected %d paths, got %d", tt.expected, paths)
			}
		})
	}
}

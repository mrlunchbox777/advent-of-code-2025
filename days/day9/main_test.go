package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProcessExampleData(t *testing.T) {
	p := filepath.Join(".", "example-data.txt")
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("failed to read example data: %v", err)
	}

	var lines []string
	for _, l := range splitLines(string(b)) {
		lines = append(lines, l)
	}

	maxArea := processCoordinates(lines)

	expected := 50
	if maxArea != expected {
		t.Fatalf("unexpected max area: got %d want %d", maxArea, expected)
	}
}

func TestParsePoint(t *testing.T) {
	tests := []struct {
		input    string
		expected Point
		wantErr  bool
	}{
		{"7,1", Point{7, 1}, false},
		{"11,7", Point{11, 7}, false},
		{"2,3", Point{2, 3}, false},
		{"", Point{}, true},
		{"invalid", Point{}, true},
		{"1,2,3", Point{}, true},
	}

	for _, tt := range tests {
		got, err := parsePoint(tt.input)
		if tt.wantErr {
			if err == nil {
				t.Errorf("parsePoint(%q) expected error, got nil", tt.input)
			}
		} else {
			if err != nil {
				t.Errorf("parsePoint(%q) unexpected error: %v", tt.input, err)
			}
			if got != tt.expected {
				t.Errorf("parsePoint(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		}
	}
}

func TestRectangleArea(t *testing.T) {
	tests := []struct {
		p1       Point
		p2       Point
		expected int
	}{
		{Point{0, 0}, Point{5, 10}, 66},  // Width=6 (0-5 inclusive), Height=11 (0-10 inclusive)
		{Point{2, 3}, Point{7, 3}, 6},    // Same Y, height = 1 (inclusive)
		{Point{7, 1}, Point{7, 10}, 10},  // Same X, width = 1 (inclusive)
		{Point{7, 1}, Point{11, 7}, 35},  // Width=5, Height=7
		{Point{11, 1}, Point{7, 7}, 35},  // Reversed, should still be 35
		{Point{11, 1}, Point{2, 5}, 50},  // Width=10, Height=5
	}

	for _, tt := range tests {
		rect := NewRectangle(tt.p1, tt.p2)
		got := rect.Area()
		if got != tt.expected {
			t.Errorf("Rectangle%v.Area() = %d, want %d", rect, got, tt.expected)
		}
	}
}

func TestProcessCoordinatesSimple(t *testing.T) {
	lines := []string{
		"0,0",
		"5,10",
		"1,1",
	}

	maxArea := processCoordinates(lines)
	// (0,0) to (5,10) = (5-0+1)*(10-0+1) = 6*11 = 66
	expected := 66
	if maxArea != expected {
		t.Fatalf("unexpected max area: got %d want %d", maxArea, expected)
	}
}

func splitLines(s string) []string {
	var out []string
	cur := ""
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			out = append(out, cur)
			cur = ""
			continue
		}
		if s[i] == '\r' {
			continue
		}
		cur += string(s[i])
	}
	if cur != "" {
		out = append(out, cur)
	}
	return out
}

package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestColumnCalculate(t *testing.T) {
	tests := []struct {
		name     string
		column   Column
		expected int
	}{
		{
			name: "multiplication",
			column: Column{
				Numbers:  []int{123, 45, 6},
				Operator: "*",
			},
			expected: 33210,
		},
		{
			name: "addition",
			column: Column{
				Numbers:  []int{328, 64, 98},
				Operator: "+",
			},
			expected: 490,
		},
		{
			name: "empty column",
			column: Column{
				Numbers:  []int{},
				Operator: "+",
			},
			expected: 0,
		},
		{
			name: "single number",
			column: Column{
				Numbers:  []int{42},
				Operator: "*",
			},
			expected: 42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.column.Calculate()
			if result != tt.expected {
				t.Errorf("Calculate() = %d, want %d", result, tt.expected)
			}
		})
	}
}

func TestGridCalculateTotal(t *testing.T) {
	grid := &Grid{
		Columns: []Column{
			{Numbers: []int{123, 45, 6}, Operator: "*"},
			{Numbers: []int{328, 64, 98}, Operator: "+"},
			{Numbers: []int{51, 387, 215}, Operator: "*"},
			{Numbers: []int{64, 23, 314}, Operator: "+"},
		},
	}

	expected := 4277556
	result := grid.CalculateTotal()
	
	if result != expected {
		t.Errorf("CalculateTotal() = %d, want %d", result, expected)
	}
}

func TestParseFileOriginalMode(t *testing.T) {
	content := `123 328  51 64 
 45 64  387 23 
  6 98  215 314
*   +   *   +  
`
	tmpfile, err := os.CreateTemp("", "test-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	grid, err := parseFile(tmpfile.Name(), "original")
	if err != nil {
		t.Fatalf("parseFile() error = %v", err)
	}

	if len(grid.Columns) != 4 {
		t.Errorf("Expected 4 columns, got %d", len(grid.Columns))
	}

	expectedTotal := 4277556
	result := grid.CalculateTotal()
	if result != expectedTotal {
		t.Errorf("Total = %d, want %d", result, expectedTotal)
	}
}

func TestParseFileAlignedMode(t *testing.T) {
	content := `123 328  51 64 
 45 64  387 23 
  6 98  215 314
*   +   *   +  
`
	tmpfile, err := os.CreateTemp("", "test-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	grid, err := parseFile(tmpfile.Name(), "aligned")
	if err != nil {
		t.Fatalf("parseFile() error = %v", err)
	}

	if len(grid.Columns) != 4 {
		t.Errorf("Expected 4 columns, got %d", len(grid.Columns))
	}

	expectedTotal := 3263827
	result := grid.CalculateTotal()
	if result != expectedTotal {
		t.Errorf("Total = %d, want %d", result, expectedTotal)
	}
}

func TestParseFileWithExampleDataOriginal(t *testing.T) {
	examplePath := filepath.Join(".", "example-data.txt")
	
	if _, err := os.Stat(examplePath); os.IsNotExist(err) {
		t.Skip("example-data.txt not found")
	}

	grid, err := parseFile(examplePath, "original")
	if err != nil {
		t.Fatalf("parseFile() error = %v", err)
	}

	expectedTotal := 4277556
	result := grid.CalculateTotal()
	if result != expectedTotal {
		t.Errorf("Total from example-data.txt in original mode = %d, want %d", result, expectedTotal)
	}
}

func TestParseFileWithExampleDataAligned(t *testing.T) {
	examplePath := filepath.Join(".", "example-data.txt")
	
	if _, err := os.Stat(examplePath); os.IsNotExist(err) {
		t.Skip("example-data.txt not found")
	}

	grid, err := parseFile(examplePath, "aligned")
	if err != nil {
		t.Fatalf("parseFile() error = %v", err)
	}

	expectedTotal := 3263827
	result := grid.CalculateTotal()
	if result != expectedTotal {
		t.Errorf("Total from example-data.txt in aligned mode = %d, want %d", result, expectedTotal)
	}
}

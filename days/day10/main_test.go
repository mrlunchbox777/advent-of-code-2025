package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseMachine(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		wantLen  int
		wantOpts int
		wantErr  bool
	}{
		{
			name:     "line1",
			line:     "[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}",
			wantLen:  4,
			wantOpts: 6,
			wantErr:  false,
		},
		{
			name:     "line2",
			line:     "[...#.] (0,2,3,4) (2,3) (0,4) (0,1,2) (1,2,3,4) {7,5,12,7,2}",
			wantLen:  5,
			wantOpts: 5,
			wantErr:  false,
		},
		{
			name:     "line3",
			line:     "[.###.#] (0,1,2,3,4) (0,3,4) (0,1,2,4,5) (1,2) {10,11,11,5,10,5}",
			wantLen:  6,
			wantOpts: 4,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := ParseMachine(tt.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseMachine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(m.DesiredState) != tt.wantLen {
					t.Errorf("DesiredState length = %d, want %d", len(m.DesiredState), tt.wantLen)
				}
				if len(m.Options) != tt.wantOpts {
					t.Errorf("Options count = %d, want %d", len(m.Options), tt.wantOpts)
				}
			}
		})
	}
}

func TestMachineSolve(t *testing.T) {
	tests := []struct {
		name           string
		line           string
		wantSelections int
	}{
		{
			name:           "line1",
			line:           "[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}",
			wantSelections: 2,
		},
		{
			name:           "line2",
			line:           "[...#.] (0,2,3,4) (2,3) (0,4) (0,1,2) (1,2,3,4) {7,5,12,7,2}",
			wantSelections: 3,
		},
		{
			name:           "line3",
			line:           "[.###.#] (0,1,2,3,4) (0,3,4) (0,1,2,4,5) (1,2) {10,11,11,5,10,5}",
			wantSelections: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := ParseMachine(tt.line)
			if err != nil {
				t.Fatalf("ParseMachine() error = %v", err)
			}
			_, selections := m.Solve()
			if selections != tt.wantSelections {
				t.Errorf("Solve() selections = %d, want %d", selections, tt.wantSelections)
			}
		})
	}
}

func TestProcessExampleData(t *testing.T) {
	p := filepath.Join(".", "example-data.txt")
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("failed to read example data: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(b)), "\n")
	results, totalSelections := ProcessLines(lines)

	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}

	expectedTotal := 7
	if totalSelections != expectedTotal {
		t.Errorf("total selections = %d, want %d", totalSelections, expectedTotal)
	}

	// Verify each line has the correct number of selections
	expectedSelections := []int{2, 3, 2}
	for i, expected := range expectedSelections {
		if !strings.Contains(results[i], fmt.Sprintf("%d selections", expected)) {
			t.Errorf("Line %d: expected %d selections, got %s", i+1, expected, results[i])
		}
	}
}

func TestApplyOption(t *testing.T) {
	m := &Machine{
		DesiredState: []bool{false, true, true, false},
		Options:      [][]int{{3}, {1, 3}, {2}},
	}

	state := []bool{false, false, false, false}
	
	// Apply option 0 (toggle position 3)
	state = m.ApplyOption(state, 0)
	expected := []bool{false, false, false, true}
	if !StatesEqual(state, expected) {
		t.Errorf("After option 0: got %v, want %v", state, expected)
	}

	// Apply option 1 (toggle positions 1,3)
	state = m.ApplyOption(state, 1)
	expected = []bool{false, true, false, false}
	if !StatesEqual(state, expected) {
		t.Errorf("After option 1: got %v, want %v", state, expected)
	}

	// Apply option 2 (toggle position 2)
	state = m.ApplyOption(state, 2)
	expected = []bool{false, true, true, false}
	if !StatesEqual(state, expected) {
		t.Errorf("After option 2: got %v, want %v", state, expected)
	}

	// Verify we reached desired state
	if !StatesEqual(state, m.DesiredState) {
		t.Errorf("Final state %v doesn't match desired %v", state, m.DesiredState)
	}
}

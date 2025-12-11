package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProcessExampleData(t *testing.T) {
	// example-data.txt is included in this package directory
	p := filepath.Join(".", "example-data.txt")
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("failed to read example data: %v", err)
	}
	// split lines
	var lines []string
	for _, l := range stringsSplitLines(string(b)) {
		lines = append(lines, l)
	}

	outs, zeroCount := processEntries(lines, "exact")

	expected := []string{
		"L68 50 -> 82",
		"L30 82 -> 52",
		"R48 52 -> 0",
		"L5 0 -> 95",
		"R60 95 -> 55",
		"L55 55 -> 0",
		"L1 0 -> 99",
		"L99 99 -> 0",
		"R14 0 -> 14",
		"L82 14 -> 32",
	}

	if len(outs) != len(expected) {
		t.Fatalf("unexpected number of outputs: got %d want %d", len(outs), len(expected))
	}
	for i := range expected {
		if outs[i] != expected[i] {
			t.Fatalf("mismatch at %d: got %q want %q", i, outs[i], expected[i])
		}
	}
	if zeroCount != 3 {
		t.Fatalf("unexpected zero count: got %d want %d", zeroCount, 3)
	}
}

func TestProcessExampleDataPasses(t *testing.T) {
	// example-data.txt is included in this package directory
	p := filepath.Join(".", "example-data.txt")
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("failed to read example data: %v", err)
	}
	// split lines
	var lines []string
	for _, l := range stringsSplitLines(string(b)) {
		lines = append(lines, l)
	}

	outs, zeroCount := processEntries(lines, "passes")

	// Same output format
	expected := []string{
		"L68 50 -> 82",
		"L30 82 -> 52",
		"R48 52 -> 0",
		"L5 0 -> 95",
		"R60 95 -> 55",
		"L55 55 -> 0",
		"L1 0 -> 99",
		"L99 99 -> 0",
		"R14 0 -> 14",
		"L82 14 -> 32",
	}

	if len(outs) != len(expected) {
		t.Fatalf("unexpected number of outputs: got %d want %d", len(outs), len(expected))
	}
	for i := range expected {
		if outs[i] != expected[i] {
			t.Fatalf("mismatch at %d: got %q want %q", i, outs[i], expected[i])
		}
	}
	// Passes mode counts crossings including landings
	if zeroCount != 6 {
		t.Fatalf("unexpected zero passes count: got %d want %d", zeroCount, 6)
	}
}

func TestMultiWrapCrossings(t *testing.T) {
	lines := []string{"R250", "L260"} // start 50
	outsExact, cExact := processEntries(lines, "exact")
	if len(outsExact) != 2 {
		t.Fatalf("expected 2 outputs")
	}
	// R250: 50-> (50+250)=300 -> 0 ends at 0 counts 1 in exact mode
	// L260: 0 -> after first move dial at 0; then left 260 from 0 ends at (0-260)=-260 mod100=40 not zero
	if cExact != 1 {
		t.Fatalf("exact mode count mismatch: got %d want %d", cExact, 1)
	}
	outsPass, cPass := processEntries(lines, "passes")
	// R250 crossings: (50+250)/100 = 3
	// L260 crossings from start 0: 260/100 = 2
	// total 5
	if cPass != 5 {
		t.Fatalf("passes mode count mismatch: got %d want %d", cPass, 5)
	}
	if len(outsPass) != 2 {
		t.Fatalf("expected 2 outputs in passes mode")
	}
}

func TestLargeRotationCrossings(t *testing.T) {
	// Single large right rotation R1000 from start 50:
	// crossings = (50 + 1000)/100 = 10, ends back at 50 so exact=0
	linesR := []string{"R1000"}
	_, cExactR := processEntries(linesR, "exact")
	if cExactR != 0 {
		t.Fatalf("R1000 exact mode should be 0, got %d", cExactR)
	}
	outsPassR, cPassR := processEntries(linesR, "passes")
	if cPassR != 10 {
		t.Fatalf("R1000 passes count mismatch: got %d want %d", cPassR, 10)
	}
	if len(outsPassR) != 1 || outsPassR[0] != "R1000 50 -> 50" {
		t.Fatalf("unexpected output for R1000: %v", outsPassR)
	}

	// Single large left rotation L1000 from start 50:
	// crossings formula (left, start>0): 1 + (stepsOrig - start)/100 = 1 + (1000-50)/100 = 10
	linesL := []string{"L1000"}
	_, cExactL := processEntries(linesL, "exact")
	if cExactL != 0 {
		t.Fatalf("L1000 exact mode should be 0, got %d", cExactL)
	}
	outsPassL, cPassL := processEntries(linesL, "passes")
	if cPassL != 10 {
		t.Fatalf("L1000 passes count mismatch: got %d want %d", cPassL, 10)
	}
	if len(outsPassL) != 1 || outsPassL[0] != "L1000 50 -> 50" {
		t.Fatalf("unexpected output for L1000: %v", outsPassL)
	}
}

// stringsSplitLines is a tiny helper avoiding extra imports in the test body.
func stringsSplitLines(s string) []string {
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

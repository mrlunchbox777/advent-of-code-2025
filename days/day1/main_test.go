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

    outs, zeroCount := processEntries(lines)

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

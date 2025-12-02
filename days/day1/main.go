package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <path-to-input-file> <mode>\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "  mode: 'exact' (ends at 0) or 'passes' (crosses or ends at 0)\n")
		os.Exit(2)
	}
	path := os.Args[1]
	mode := os.Args[2]
	if mode != "exact" && mode != "passes" {
		fmt.Fprintf(os.Stderr, "Invalid mode %q. Must be 'exact' or 'passes'\n", mode)
		os.Exit(2)
	}
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer f.Close()

	// Read all lines into a slice and process using the testable function.
	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("read error: %v", err)
	}

	outs, zeroCount := processEntries(lines, mode)
	for _, out := range outs {
		fmt.Println(out)
	}
	if mode == "exact" {
		fmt.Printf("Ended at 0 count: %d\n", zeroCount)
	} else {
		fmt.Printf("Passed through 0 count: %d\n", zeroCount)
	}
}

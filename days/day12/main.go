package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <path-to-input-file>\n", filepath.Base(os.Args[0]))
		os.Exit(2)
	}
	path := os.Args[1]

	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("read error: %v", err)
	}

	data, err := ParseInput(lines)
	if err != nil {
		log.Fatalf("parse error: %v", err)
	}

	solvedCount := 0
	for i, puzzle := range data.Puzzles {
		fmt.Fprintf(os.Stderr, "Solving puzzle %d/%d (%dx%d)...\n", i+1, len(data.Puzzles), puzzle.Width, puzzle.Height)
		
		// Skip puzzles with no pieces to place
		haspieces := false
		for _, spec := range puzzle.PieceSpecs {
			if spec.Count > 0 {
				haspieces = true
				break
			}
		}
		
		if !haspieces {
			fmt.Println("No solution found")
			continue
		}
		
		solution := puzzle.Solve(data.Pieces)
		if solution != nil {
			fmt.Println(solution.String())
			solvedCount++
		} else {
			fmt.Println("No solution found")
		}
	}

	fmt.Printf("Puzzles with solutions found: %d\n", solvedCount)
}

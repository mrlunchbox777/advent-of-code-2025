package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: day4 <filepath> <mode>")
		fmt.Println("  mode: 'initial' for single pass, 'completion' for iterative passes")
		os.Exit(1)
	}

	filepath := os.Args[1]
	mode := os.Args[2]

	if mode != "initial" && mode != "completion" {
		fmt.Println("Error: mode must be 'initial' or 'completion'")
		os.Exit(1)
	}

	// Load the grid from file
	grid, err := NewGridFromFile(filepath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	if mode == "initial" {
		runInitialPass(grid)
	} else {
		runCompletionMode(grid)
	}
}

func runInitialPass(grid *Grid) {
	// Find all selected positions
	selected := grid.FindSelectedPositions()

	// Print each selected position
	fmt.Println("Selected positions:")
	for _, pos := range selected {
		fmt.Printf("[%d,%d]\n", pos.X, pos.Y)
	}

	// Print the total count
	fmt.Printf("\nTotal count: %d\n", len(selected))
}

func runCompletionMode(grid *Grid) {
	allPositions := []Position{}
	runningTotal := 0
	round := 1

	for {
		selected := grid.FindSelectedPositions()
		if len(selected) == 0 {
			break
		}

		// Print round information
		fmt.Printf("Round %d:\n", round)
		fmt.Println("Selected positions:")
		for _, pos := range selected {
			fmt.Printf("[%d,%d]\n", pos.X, pos.Y)
			allPositions = append(allPositions, pos)
		}
		runningTotal += len(selected)
		fmt.Printf("Total from round: %d\n", len(selected))
		fmt.Printf("Running total: %d\n\n", runningTotal)

		// Create new grid with selected positions replaced by '.'
		grid = grid.ReplacePositions(selected)
		round++
	}

	// Print final summary
	fmt.Println("=== Final Summary ===")
	fmt.Println("All selected positions:")
	for _, pos := range allPositions {
		fmt.Printf("[%d,%d]\n", pos.X, pos.Y)
	}
	fmt.Printf("\nNumber of rounds: %d\n", round-1)
	fmt.Printf("Final total: %d\n", runningTotal)
}

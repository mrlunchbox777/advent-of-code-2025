package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: day4 <filepath>")
		os.Exit(1)
	}

	filepath := os.Args[1]

	// Load the grid from file
	grid, err := NewGridFromFile(filepath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

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

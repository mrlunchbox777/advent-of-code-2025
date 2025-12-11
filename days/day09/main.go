package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <path-to-input-file> <mode> [output-file]\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "  mode: 'original' (any pair as corners) or 'contained' (rectangle within shape)\n")
		fmt.Fprintf(os.Stderr, "  output-file: optional, if provided draws visualization for 'contained' mode\n")
		os.Exit(2)
	}
	path := os.Args[1]
	mode := os.Args[2]
	var outputFile string
	if len(os.Args) > 3 {
		outputFile = os.Args[3]
	}
	
	if mode != "original" && mode != "contained" {
		fmt.Fprintf(os.Stderr, "Invalid mode %q. Must be 'original' or 'contained'\n", mode)
		os.Exit(2)
	}
	
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

	if outputFile != "" && mode == "contained" {
		maxArea, rect := processCoordinatesWithResult(lines, mode)
		fmt.Printf("Largest rectangle area: %d\n", maxArea)
		
		if err := drawVisualization(lines, rect, outputFile); err != nil {
			log.Fatalf("failed to draw visualization: %v", err)
		}
		fmt.Printf("Visualization saved to: %s\n", outputFile)
	} else {
		maxArea := processCoordinates(lines, mode)
		fmt.Printf("Largest rectangle area: %d\n", maxArea)
	}
}

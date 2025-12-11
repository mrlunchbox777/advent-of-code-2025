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
		fmt.Fprintf(os.Stderr, "  mode: 'toggle' or 'counter'\n")
		os.Exit(2)
	}
	path := os.Args[1]
	mode := os.Args[2]
	if mode != "toggle" && mode != "counter" {
		fmt.Fprintf(os.Stderr, "Invalid mode %q. Must be 'toggle' or 'counter'\n", mode)
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

	results, totalSelections := ProcessLines(lines, mode)
	for _, result := range results {
		fmt.Println(result)
	}
	fmt.Printf("Total selections: %d\n", totalSelections)
}

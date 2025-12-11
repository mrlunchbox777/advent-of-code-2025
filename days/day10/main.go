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

	results, totalSelections := ProcessLines(lines)
	for _, result := range results {
		fmt.Println(result)
	}
	fmt.Printf("Total selections: %d\n", totalSelections)
}

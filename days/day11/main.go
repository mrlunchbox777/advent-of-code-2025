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
		fmt.Fprintf(os.Stderr, "  mode: 'all' or 'must-visit'\n")
		os.Exit(2)
	}
	path := os.Args[1]
	mode := os.Args[2]
	if mode != "all" && mode != "must-visit" {
		fmt.Fprintf(os.Stderr, "Invalid mode %q. Must be 'all' or 'must-visit'\n", mode)
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

	graph, err := ParseGraph(lines)
	if err != nil {
		log.Fatalf("parse error: %v", err)
	}

	var paths []Path
	if mode == "all" {
		paths = graph.FindAllPaths("you", "out")
	} else {
		paths = graph.FindPathsWithRequiredNodes("svr", "out", []string{"dac", "fft"})
	}
	
	for _, path := range paths {
		fmt.Println(path.String())
	}
	
	fmt.Printf("Total unique paths: %d\n", len(paths))
}

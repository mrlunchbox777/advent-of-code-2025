package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

	scanner := bufio.NewScanner(f)
	dial := 50
	zeroCount := 0
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		raw := strings.TrimSpace(scanner.Text())
		if raw == "" {
			continue
		}
		// normalize (remove spaces and make upper-case) so both "L68" and "L 68" work
		entry := strings.ToUpper(strings.ReplaceAll(raw, " ", ""))
		if len(entry) < 2 {
			fmt.Fprintf(os.Stderr, "invalid entry at line %d: %q\n", lineNum, raw)
			continue
		}
		dir := entry[0]
		numStr := entry[1:]
		n, err := strconv.Atoi(numStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid number at line %d: %q\n", lineNum, numStr)
			continue
		}

		start := dial
		shift := n % 100
		if dir == 'L' {
			// subtract with wrap-around 0..99
			dial = ((dial-shift)%100 + 100) % 100
		} else if dir == 'R' {
			dial = (dial + shift) % 100
		} else {
			fmt.Fprintf(os.Stderr, "invalid direction at line %d: %q\n", lineNum, string(dir))
			continue
		}

		fmt.Printf("%s %d -> %d\n", entry, start, dial)
		if dial == 0 {
			zeroCount++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("read error: %v", err)
	}

	fmt.Printf("Ended at 0 count: %d\n", zeroCount)
}

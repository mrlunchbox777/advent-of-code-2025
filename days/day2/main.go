package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . <filepath>")
		os.Exit(1)
	}

	filePath := os.Args[1]
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	line := strings.TrimSpace(string(data))
	entries := strings.Split(line, ",")

	totalSum := 0
	for _, entry := range entries {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}

		r, err := ParseRange(entry)
		if err != nil {
			fmt.Printf("Error parsing range %q: %v\n", entry, err)
			continue
		}

		invalidIDs := r.FindRepeatedSequenceNumbers()

		if len(invalidIDs) > 0 {
			fmt.Printf("%s has %d invalid ID(s): %v\n", entry, len(invalidIDs), invalidIDs)
		} else {
			fmt.Printf("%s contains no invalid IDs.\n", entry)
		}

		for _, id := range invalidIDs {
			totalSum += id
		}
	}

	fmt.Printf("\nTotal sum of invalid IDs: %d\n", totalSum)
}

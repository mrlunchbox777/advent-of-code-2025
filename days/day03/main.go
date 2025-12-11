package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run . <filepath> <digitCount>")
		fmt.Println("  digitCount: number of digits to concatenate (e.g., 2, 12)")
		os.Exit(1)
	}

	filePath := os.Args[1]
	digitCount := 0
	if _, err := fmt.Sscanf(os.Args[2], "%d", &digitCount); err != nil || digitCount < 1 {
		fmt.Printf("Invalid digit count %q. Must be a positive integer.\n", os.Args[2])
		os.Exit(1)
	}
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	totalSum := 0

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		entry := NewEntry(line)
		digits, result := entry.FindLargestNumber(digitCount)

		fmt.Printf("%s -> %v = %d\n", line, string(digits), result)
		totalSum += result
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nTotal sum: %d\n", totalSum)
}

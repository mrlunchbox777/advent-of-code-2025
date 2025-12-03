package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . <filepath>")
		os.Exit(1)
	}

	filePath := os.Args[1]
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
		digit1, digit2, result := entry.FindLargestTwoDigitNumber()

		fmt.Printf("%s -> %c and %c = %d\n", line, digit1, digit2, result)
		totalSum += result
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nTotal sum: %d\n", totalSum)
}

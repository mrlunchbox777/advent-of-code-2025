package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Range struct {
	Start int
	End   int
}

func (r Range) Contains(n int) bool {
	return n >= r.Start && n <= r.End
}

type RangeList struct {
	Ranges []Range
}

func (rl *RangeList) AddRange(r Range) {
	rl.Ranges = append(rl.Ranges, r)
}

func (rl *RangeList) IsValid(n int) bool {
	for _, r := range rl.Ranges {
		if r.Contains(n) {
			return true
		}
	}
	return false
}

func (rl *RangeList) CountTotalValid() int {
	if len(rl.Ranges) == 0 {
		return 0
	}
	
	// Sort ranges by start position
	ranges := make([]Range, len(rl.Ranges))
	copy(ranges, rl.Ranges)
	
	for i := 0; i < len(ranges); i++ {
		for j := i + 1; j < len(ranges); j++ {
			if ranges[j].Start < ranges[i].Start {
				ranges[i], ranges[j] = ranges[j], ranges[i]
			}
		}
	}
	
	// Merge overlapping ranges and count
	total := 0
	currentStart := ranges[0].Start
	currentEnd := ranges[0].End
	
	for i := 1; i < len(ranges); i++ {
		if ranges[i].Start <= currentEnd+1 {
			// Overlapping or adjacent - merge
			if ranges[i].End > currentEnd {
				currentEnd = ranges[i].End
			}
		} else {
			// No overlap - count current range and start new one
			total += currentEnd - currentStart + 1
			currentStart = ranges[i].Start
			currentEnd = ranges[i].End
		}
	}
	
	// Add the last range
	total += currentEnd - currentStart + 1
	
	return total
}

type NumberList struct {
	Numbers []int
}

func (nl *NumberList) AddNumber(n int) {
	nl.Numbers = append(nl.Numbers, n)
}

func (nl *NumberList) ValidateAgainstRanges(rangeList *RangeList) int {
	count := 0
	for _, num := range nl.Numbers {
		valid := rangeList.IsValid(num)
		fmt.Printf("%d: %t\n", num, valid)
		if valid {
			count++
		}
	}
	return count
}

func parseRange(line string) (Range, error) {
	parts := strings.Split(line, "-")
	if len(parts) != 2 {
		return Range{}, fmt.Errorf("invalid range format: %s", line)
	}
	start, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return Range{}, err
	}
	end, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return Range{}, err
	}
	return Range{Start: start, End: end}, nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <path-to-file> <mode>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Modes:\n")
		fmt.Fprintf(os.Stderr, "  validate - Count valid numbers from second list\n")
		fmt.Fprintf(os.Stderr, "  total    - Count total possible valid numbers from ranges\n")
		os.Exit(1)
	}

	filePath := os.Args[1]
	mode := os.Args[2]
	
	if mode != "validate" && mode != "total" {
		fmt.Fprintf(os.Stderr, "Invalid mode: %s\n", mode)
		fmt.Fprintf(os.Stderr, "Valid modes are: validate, total\n")
		os.Exit(1)
	}
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	rangeList := &RangeList{}
	numberList := &NumberList{}
	
	scanner := bufio.NewScanner(file)
	parsingRanges := true

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		if line == "" {
			parsingRanges = false
			continue
		}

		if parsingRanges {
			r, err := parseRange(line)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error parsing range '%s': %v\n", line, err)
				os.Exit(1)
			}
			rangeList.AddRange(r)
		} else {
			num, err := strconv.Atoi(line)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error parsing number '%s': %v\n", line, err)
				os.Exit(1)
			}
			numberList.AddNumber(num)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	if mode == "validate" {
		count := numberList.ValidateAgainstRanges(rangeList)
		fmt.Printf("\nTotal valid numbers: %d\n", count)
	} else {
		count := rangeList.CountTotalValid()
		fmt.Printf("Total possible valid numbers: %d\n", count)
	}
}

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

func (rl *RangeList) BuildValidSet() map[int]bool {
	validSet := make(map[int]bool)
	for _, r := range rl.Ranges {
		for i := r.Start; i <= r.End; i++ {
			validSet[i] = true
		}
	}
	return validSet
}

type NumberList struct {
	Numbers []int
}

func (nl *NumberList) AddNumber(n int) {
	nl.Numbers = append(nl.Numbers, n)
}

func (nl *NumberList) ValidateAgainstSet(validSet map[int]bool) int {
	count := 0
	for _, num := range nl.Numbers {
		valid := validSet[num]
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
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <path-to-file>\n", os.Args[0])
		os.Exit(1)
	}

	filePath := os.Args[1]
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

	validSet := rangeList.BuildValidSet()
	count := numberList.ValidateAgainstSet(validSet)
	
	fmt.Printf("\nTotal valid numbers: %d\n", count)
}

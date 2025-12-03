package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Range represents a range of numbers with lower and upper bounds
type Range struct {
	Lower int
	Upper int
}

// ParseRange parses a string in the format "lower-upper" into a Range
func ParseRange(s string) (Range, error) {
	parts := strings.Split(s, "-")
	if len(parts) != 2 {
		return Range{}, fmt.Errorf("invalid range format: %s", s)
	}

	lower, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return Range{}, fmt.Errorf("invalid lower bound: %v", err)
	}

	upper, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return Range{}, fmt.Errorf("invalid upper bound: %v", err)
	}

	if lower > upper {
		return Range{}, fmt.Errorf("lower bound %d greater than upper bound %d", lower, upper)
	}

	return Range{Lower: lower, Upper: upper}, nil
}

// FindRepeatedSequenceNumbers finds all numbers in the range that are comprised
// entirely of repeated sequences (e.g., 11, 1010, 222222)
func (r Range) FindRepeatedSequenceNumbers() []int {
	var result []int

	for n := r.Lower; n <= r.Upper; n++ {
		if r.isRepeatedSequence(n) {
			result = append(result, n)
		}
	}

	return result
}

// isRepeatedSequence checks if a number is comprised of a repeated pattern
// Examples: 11 (pattern "1" x2), 1010 (pattern "10" x2), 222222 (pattern "2" x6)
// Not: 101 (no repeating pattern), 111 (pattern "1" x3, but odd total length for single char)
// Rule: Must have even total length OR pattern length >= 2
func (r Range) isRepeatedSequence(n int) bool {
	s := strconv.Itoa(n)
	length := len(s)

	// Try all possible pattern lengths from 1 to length/2
	for patternLen := 1; patternLen <= length/2; patternLen++ {
		// Pattern length must divide evenly into total length
		if length%patternLen != 0 {
			continue
		}

		pattern := s[:patternLen]
		repeats := length / patternLen

		// Build the repeated string
		repeated := strings.Repeat(pattern, repeats)

		if repeated == s {
			// For single-digit patterns, total length must be even
			if patternLen == 1 && length%2 != 0 {
				continue
			}
			return true
		}
	}

	return false
}

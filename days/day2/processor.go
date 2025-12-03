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
// Examples: 11 (pattern "1" x2), 1010 (pattern "10" x2), 222222 (pattern "222" x2)
// Not: 101 (no repeating pattern), 111 (pattern "1" x3), 1111 (pattern "1" x4)
// Rule: The number must be a pattern repeated EXACTLY 2 times
func (r Range) isRepeatedSequence(n int) bool {
	s := strconv.Itoa(n)
	length := len(s)

	// Number must have even length to be repeated exactly 2 times
	if length%2 != 0 {
		return false
	}

	// Check if first half equals second half
	half := length / 2
	pattern := s[:half]
	secondHalf := s[half:]

	return pattern == secondHalf
}

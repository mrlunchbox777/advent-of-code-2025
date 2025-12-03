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
// mode: "exact" for pattern repeated exactly 2 times, "any" for pattern repeated 2+ times
func (r Range) FindRepeatedSequenceNumbers(mode string) []int {
	var result []int

	for n := r.Lower; n <= r.Upper; n++ {
		if r.isRepeatedSequence(n, mode) {
			result = append(result, n)
		}
	}

	return result
}

// isRepeatedSequence checks if a number is comprised of a repeated pattern
// Examples (exact): 11 (pattern "1" x2), 1010 (pattern "10" x2), 222222 (pattern "222" x2)
// Examples (any): 111 (pattern "1" x3), 123123123 (pattern "123" x3)
// Not: 101 (no repeating pattern)
// mode: "exact" requires pattern repeated exactly 2 times, "any" requires 2+ times
func (r Range) isRepeatedSequence(n int, mode string) bool {
	s := strconv.Itoa(n)
	length := len(s)

	if mode == "exact" {
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

	// mode == "any": check all possible pattern lengths
	for patternLen := 1; patternLen <= length/2; patternLen++ {
		// Pattern length must divide evenly into total length
		if length%patternLen != 0 {
			continue
		}

		pattern := s[:patternLen]
		repeats := length / patternLen

		// Must repeat at least 2 times
		if repeats < 2 {
			continue
		}

		// Build the repeated string
		repeated := strings.Repeat(pattern, repeats)

		if repeated == s {
			return true
		}
	}

	return false
}

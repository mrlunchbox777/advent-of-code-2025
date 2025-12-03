package main

import (
	"unicode"
)

// Entry represents a line of digits
type Entry struct {
	Raw string
}

// NewEntry creates a new Entry from a string
func NewEntry(raw string) Entry {
	return Entry{Raw: raw}
}

// FindLargestNumber finds the n digits that form the largest n-digit number
// by selecting n digits in order (without reordering) and concatenating them in that order.
// Returns the selected digits as a slice and the resulting number.
func (e Entry) FindLargestNumber(n int) ([]rune, int) {
	// Extract all digits from the raw string
	var digits []rune
	for _, ch := range e.Raw {
		if unicode.IsDigit(ch) {
			digits = append(digits, ch)
		}
	}

	if len(digits) < n {
		return nil, 0
	}

	// Optimized iterative combination generation
	maxValue := 0
	var maxDigits []rune

	// Pre-allocate arrays to avoid repeated allocations
	indices := make([]int, n)
	selected := make([]rune, n)

	// Initialize indices to first n positions
	for i := 0; i < n; i++ {
		indices[i] = i
	}

	for {
		// Calculate value for current combination
		value := 0
		for i := 0; i < n; i++ {
			selected[i] = digits[indices[i]]
			value = value*10 + int(digits[indices[i]]-'0')
		}

		if value > maxValue {
			maxValue = value
			maxDigits = make([]rune, n)
			copy(maxDigits, selected)
		}

		// Generate next combination (lexicographic order)
		// Find rightmost index that can be incremented
		i := n - 1
		for i >= 0 && indices[i] == len(digits)-n+i {
			i--
		}

		// No more combinations
		if i < 0 {
			break
		}

		// Increment this index and reset all following indices
		indices[i]++
		for j := i + 1; j < n; j++ {
			indices[j] = indices[j-1] + 1
		}
	}

	return maxDigits, maxValue
}

// FindLargestTwoDigitNumber is a convenience wrapper for FindLargestNumber(2)
// Kept for backward compatibility with existing tests
func (e Entry) FindLargestTwoDigitNumber() (rune, rune, int) {
	digits, value := e.FindLargestNumber(2)
	if len(digits) < 2 {
		return '0', '0', 0
	}
	return digits[0], digits[1], value
}

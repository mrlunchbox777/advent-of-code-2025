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
// Uses a greedy algorithm: for each position, select the largest digit that still leaves
// enough remaining digits to complete the n-digit number.
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

	// Greedy algorithm: for each position, pick the largest digit
	// that leaves enough digits remaining to complete the selection
	result := make([]rune, 0, n)
	startPos := 0

	for len(result) < n {
		remaining := n - len(result) // how many more digits we need
		maxDigit := '0' - 1          // invalid value to ensure first digit is always picked
		maxPos := -1

		// Search window: we can only pick from positions that leave enough digits after
		// searchEnd is the last position we can pick from and still have enough digits
		searchEnd := len(digits) - remaining + 1

		for i := startPos; i < searchEnd; i++ {
			if digits[i] > maxDigit {
				maxDigit = digits[i]
				maxPos = i
			}
		}

		result = append(result, maxDigit)
		startPos = maxPos + 1 // next search starts after the position we just picked
	}

	// Calculate the numeric value
	value := 0
	for _, d := range result {
		value = value*10 + int(d-'0')
	}

	return result, value
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

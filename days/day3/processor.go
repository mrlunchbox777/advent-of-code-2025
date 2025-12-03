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

	// Try all combinations of n digits maintaining original order
	maxValue := 0
	var maxDigits []rune

	// Generate all combinations of n positions
	var tryCombo func(start, depth int, indices []int)
	tryCombo = func(start, depth int, indices []int) {
		if depth == n {
			// Calculate the value for this combination
			value := 0
			selected := make([]rune, n)
			for i, idx := range indices {
				selected[i] = digits[idx]
				value = value*10 + int(digits[idx]-'0')
			}
			if value > maxValue {
				maxValue = value
				maxDigits = selected
			}
			return
		}
		for i := start; i < len(digits); i++ {
			tryCombo(i+1, depth+1, append(indices, i))
		}
	}

	tryCombo(0, 0, []int{})
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

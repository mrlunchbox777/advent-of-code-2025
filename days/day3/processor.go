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

// FindLargestTwoDigitNumber finds the two digits that form the largest two-digit number
// by selecting two digits in order (without reordering) and concatenating them in that order.
// Returns the two digits and the resulting number.
func (e Entry) FindLargestTwoDigitNumber() (rune, rune, int) {
	// Extract all digits from the raw string
	var digits []rune
	for _, ch := range e.Raw {
		if unicode.IsDigit(ch) {
			digits = append(digits, ch)
		}
	}

	if len(digits) < 2 {
		return '0', '0', 0
	}

	// Try all pairs of digits maintaining original order
	// Concatenate in the order they appear (i before j)
	maxValue := 0
	maxDigit1 := '0'
	maxDigit2 := '0'

	for i := 0; i < len(digits); i++ {
		for j := i + 1; j < len(digits); j++ {
			d1 := int(digits[i] - '0')
			d2 := int(digits[j] - '0')

			// Concatenate in order: first digit then second digit
			value := d1*10 + d2

			if value > maxValue {
				maxValue = value
				maxDigit1 = digits[i]
				maxDigit2 = digits[j]
			}
		}
	}

	return maxDigit1, maxDigit2, maxValue
}

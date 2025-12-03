package main

import (
	"bufio"
	"os"
	"testing"
)

func TestFindLargestTwoDigitNumber(t *testing.T) {
	tests := []struct {
		input          string
		expectedDigit1 rune
		expectedDigit2 rune
		expectedResult int
		description    string
	}{
		{"12345", '4', '5', 45, "should select 4 and 5 for 45"},
		{"54321", '5', '4', 54, "should select 5 and 4 for 54"},
		{"987654321111111", '9', '8', 98, "should select 9 and 8 for 98"},
		{"811111111111119", '8', '9', 89, "should select 8 and 9 for 89"},
		{"234234234234278", '7', '8', 78, "should select 7 and 8 for 78"},
		{"818181911112111", '9', '2', 92, "should select 9 and 2 for 92"},
		{"99", '9', '9', 99, "both 9s should give 99"},
		{"19", '1', '9', 19, "should select 1 and 9 for 19"},
	}

	for _, tt := range tests {
		entry := NewEntry(tt.input)
		d1, d2, result := entry.FindLargestTwoDigitNumber()

		if result != tt.expectedResult {
			t.Errorf("%s: FindLargestTwoDigitNumber(%q) = %d, want %d",
				tt.description, tt.input, result, tt.expectedResult)
		}

		// Check that the digits are correct (order might vary)
		if !((d1 == tt.expectedDigit1 && d2 == tt.expectedDigit2) ||
			(d1 == tt.expectedDigit2 && d2 == tt.expectedDigit1)) {
			t.Errorf("%s: got digits %c and %c, want %c and %c",
				tt.description, d1, d2, tt.expectedDigit1, tt.expectedDigit2)
		}
	}
}

func TestProcessExampleData(t *testing.T) {
	file, err := os.Open("example-data.txt")
	if err != nil {
		t.Fatalf("Failed to open example-data.txt: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	totalSum := 0
	lineCount := 0

	expectedResults := []int{98, 89, 78, 92}

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		entry := NewEntry(line)
		_, _, result := entry.FindLargestTwoDigitNumber()

		if lineCount < len(expectedResults) {
			if result != expectedResults[lineCount] {
				t.Errorf("Line %d (%s): got %d, want %d",
					lineCount+1, line, result, expectedResults[lineCount])
			}
		}

		totalSum += result
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("Error reading file: %v", err)
	}

	expectedSum := 357
	if totalSum != expectedSum {
		t.Errorf("Total sum = %d, want %d", totalSum, expectedSum)
	}

	if lineCount != 4 {
		t.Errorf("Expected 4 lines, got %d", lineCount)
	}
}

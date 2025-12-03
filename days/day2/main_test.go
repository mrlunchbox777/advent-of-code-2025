package main

import (
	"os"
	"strings"
	"testing"
)

func TestParseRange(t *testing.T) {
	tests := []struct {
		input    string
		expected Range
		wantErr  bool
	}{
		{"11-22", Range{Lower: 11, Upper: 22}, false},
		{"95-115", Range{Lower: 95, Upper: 115}, false},
		{"1-9", Range{Lower: 1, Upper: 9}, false},
		{"invalid", Range{}, true},
		{"50-40", Range{}, true},
	}

	for _, tt := range tests {
		r, err := ParseRange(tt.input)
		if tt.wantErr {
			if err == nil {
				t.Errorf("ParseRange(%q) expected error, got nil", tt.input)
			}
			continue
		}
		if err != nil {
			t.Errorf("ParseRange(%q) unexpected error: %v", tt.input, err)
			continue
		}
		if r != tt.expected {
			t.Errorf("ParseRange(%q) = %v, want %v", tt.input, r, tt.expected)
		}
	}
}

func TestIsRepeatedSequence(t *testing.T) {
	r := Range{Lower: 0, Upper: 100000}

	tests := []struct {
		num      int
		expected bool
	}{
		{11, true},         // "1" repeated
		{22, true},         // "2" repeated
		{99, true},         // "9" repeated
		{1010, true},       // "10" repeated
		{222222, true},     // "2" repeated
		{446446, true},     // "446" repeated
		{101, false},       // no repeating pattern
		{123, false},       // no repeating pattern
		{1188511885, true}, // "11885" repeated
		{38593859, true},   // "3859" repeated
	}

	for _, tt := range tests {
		result := r.isRepeatedSequence(tt.num)
		if result != tt.expected {
			t.Errorf("isRepeatedSequence(%d) = %v, want %v", tt.num, result, tt.expected)
		}
	}
}

func TestFindRepeatedSequenceNumbers(t *testing.T) {
	tests := []struct {
		rangeStr string
		expected []int
	}{
		{"11-22", []int{11, 22}},
		{"95-115", []int{99}},
		{"998-1012", []int{1010}},
		{"222220-222224", []int{222222}},
		{"1698522-1698528", []int{}},
		{"446443-446449", []int{446446}},
	}

	for _, tt := range tests {
		r, err := ParseRange(tt.rangeStr)
		if err != nil {
			t.Fatalf("ParseRange(%q) error: %v", tt.rangeStr, err)
		}
		result := r.FindRepeatedSequenceNumbers()
		if len(result) != len(tt.expected) {
			t.Errorf("Range %s: got %d numbers, want %d: %v", tt.rangeStr, len(result), len(tt.expected), result)
			continue
		}
		for i := range result {
			if result[i] != tt.expected[i] {
				t.Errorf("Range %s: got %v, want %v", tt.rangeStr, result, tt.expected)
				break
			}
		}
	}
}

func TestProcessExampleData(t *testing.T) {
	data, err := os.ReadFile("example-data.txt")
	if err != nil {
		t.Fatalf("Failed to read example-data.txt: %v", err)
	}

	line := strings.TrimSpace(string(data))
	entries := strings.Split(line, ",")

	totalSum := 0
	expectedResults := map[string][]int{
		"11-22":                 {11, 22},
		"95-115":                {99},
		"998-1012":              {1010},
		"1188511880-1188511890": {1188511885},
		"222220-222224":         {222222},
		"1698522-1698528":       {},
		"446443-446449":         {446446},
		"38593856-38593862":     {38593859},
	}

	for _, entry := range entries {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}

		r, err := ParseRange(entry)
		if err != nil {
			t.Fatalf("Error parsing range %q: %v", entry, err)
		}

		invalidIDs := r.FindRepeatedSequenceNumbers()

		expected, ok := expectedResults[entry]
		if !ok {
			t.Errorf("Unexpected range in example data: %s", entry)
			continue
		}

		if len(invalidIDs) != len(expected) {
			t.Errorf("Range %s: got %d invalid IDs %v, want %d: %v",
				entry, len(invalidIDs), invalidIDs, len(expected), expected)
			continue
		}

		for i := range invalidIDs {
			if invalidIDs[i] != expected[i] {
				t.Errorf("Range %s: got %v, want %v", entry, invalidIDs, expected)
				break
			}
		}

		for _, id := range invalidIDs {
			totalSum += id
		}
	}

	expectedSum := 1227775554
	if totalSum != expectedSum {
		t.Errorf("Total sum = %d, want %d", totalSum, expectedSum)
	}
}

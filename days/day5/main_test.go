package main

import (
	"testing"
)

func TestRangeContains(t *testing.T) {
	r := Range{Start: 10, End: 20}
	
	tests := []struct {
		num      int
		expected bool
	}{
		{5, false},
		{10, true},
		{15, true},
		{20, true},
		{25, false},
	}
	
	for _, tt := range tests {
		result := r.Contains(tt.num)
		if result != tt.expected {
			t.Errorf("Range{10, 20}.Contains(%d) = %t, want %t", tt.num, result, tt.expected)
		}
	}
}

func TestRangeListIsValid(t *testing.T) {
	rl := &RangeList{}
	rl.AddRange(Range{Start: 3, End: 5})
	rl.AddRange(Range{Start: 10, End: 14})
	rl.AddRange(Range{Start: 16, End: 20})
	rl.AddRange(Range{Start: 12, End: 18})
	
	expectedValid := []int{3, 4, 5, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	expectedInvalid := []int{1, 2, 6, 7, 8, 9, 21, 32}
	
	for _, num := range expectedValid {
		if !rl.IsValid(num) {
			t.Errorf("Expected %d to be valid", num)
		}
	}
	
	for _, num := range expectedInvalid {
		if rl.IsValid(num) {
			t.Errorf("Expected %d to NOT be valid", num)
		}
	}
}

func TestNumberListValidateAgainstRanges(t *testing.T) {
	rl := &RangeList{}
	rl.AddRange(Range{Start: 3, End: 5})
	rl.AddRange(Range{Start: 10, End: 14})
	rl.AddRange(Range{Start: 16, End: 20})
	
	nl := &NumberList{}
	nl.AddNumber(1)
	nl.AddNumber(5)
	nl.AddNumber(8)
	nl.AddNumber(11)
	nl.AddNumber(17)
	nl.AddNumber(32)
	
	count := nl.ValidateAgainstRanges(rl)
	
	if count != 3 {
		t.Errorf("Expected 3 valid numbers, got %d", count)
	}
}

func TestRangeListCountTotalValid(t *testing.T) {
	rl := &RangeList{}
	rl.AddRange(Range{Start: 3, End: 5})    // 3, 4, 5 = 3 numbers
	rl.AddRange(Range{Start: 10, End: 14})  // 10, 11, 12, 13, 14 = 5 numbers
	rl.AddRange(Range{Start: 16, End: 20})  // 16, 17, 18, 19, 20 = 5 numbers
	rl.AddRange(Range{Start: 12, End: 18})  // 12, 13, 14 already counted, adds 15, 16, 17, 18 = 1 new (15)
	
	count := rl.CountTotalValid()
	
	// 3,4,5,10,11,12,13,14,15,16,17,18,19,20 = 14 unique numbers
	if count != 14 {
		t.Errorf("Expected 14 total possible valid numbers, got %d", count)
	}
}

func TestRangeListCountTotalValidLargeNumbers(t *testing.T) {
	rl := &RangeList{}
	// Test with billion-scale ranges
	rl.AddRange(Range{Start: 1000000000, End: 1000000100})  // 101 numbers
	rl.AddRange(Range{Start: 2000000000, End: 2000000200})  // 201 numbers
	rl.AddRange(Range{Start: 1000000050, End: 1000000150})  // Overlaps with first, adds 50
	
	count := rl.CountTotalValid()
	
	// First range: 101, second: 201, third adds 50 (1000000101-1000000150)
	expected := 101 + 201 + 50
	if count != expected {
		t.Errorf("Expected %d total possible valid numbers, got %d", expected, count)
	}
}

func TestRangeListCountTotalValidNoOverlap(t *testing.T) {
	rl := &RangeList{}
	rl.AddRange(Range{Start: 1, End: 10})
	rl.AddRange(Range{Start: 20, End: 30})
	rl.AddRange(Range{Start: 40, End: 50})
	
	count := rl.CountTotalValid()
	
	// 10 + 11 + 11 = 32
	if count != 32 {
		t.Errorf("Expected 32 total possible valid numbers, got %d", count)
	}
}

func TestRangeListCountTotalValidFullOverlap(t *testing.T) {
	rl := &RangeList{}
	rl.AddRange(Range{Start: 1, End: 100})
	rl.AddRange(Range{Start: 10, End: 50})
	rl.AddRange(Range{Start: 25, End: 75})
	
	count := rl.CountTotalValid()
	
	// All merged into 1-100 = 100
	if count != 100 {
		t.Errorf("Expected 100 total possible valid numbers, got %d", count)
	}
}

func TestParseRange(t *testing.T) {
	tests := []struct {
		input       string
		expected    Range
		expectError bool
	}{
		{"3-5", Range{Start: 3, End: 5}, false},
		{"10-14", Range{Start: 10, End: 14}, false},
		{"100-200", Range{Start: 100, End: 200}, false},
		{"invalid", Range{}, true},
		{"1-2-3", Range{}, true},
	}
	
	for _, tt := range tests {
		result, err := parseRange(tt.input)
		if tt.expectError {
			if err == nil {
				t.Errorf("parseRange(%q) expected error, got nil", tt.input)
			}
		} else {
			if err != nil {
				t.Errorf("parseRange(%q) unexpected error: %v", tt.input, err)
			}
			if result != tt.expected {
				t.Errorf("parseRange(%q) = %+v, want %+v", tt.input, result, tt.expected)
			}
		}
	}
}

func TestPerformance500RangesTrillions(t *testing.T) {
rl := &RangeList{}

// Add 500 ranges in the hundred trillions
for i := 0; i < 500; i++ {
start := int(100000000000000 + int64(i)*2000000000000)
end := start + 50000000
rl.AddRange(Range{Start: start, End: end})
}

// This should complete very quickly
count := rl.CountTotalValid()

// Each range has ~50M numbers, all non-overlapping = 500 * 50M = 25 billion
expected := 500 * 50000001
if count != expected {
t.Errorf("Expected %d, got %d", expected, count)
}
}

func BenchmarkCountTotalValid500Ranges(b *testing.B) {
rl := &RangeList{}

// Add 500 ranges in the hundred trillions with some overlaps
for i := 0; i < 500; i++ {
start := int(100000000000000 + int64(i)*1500000000000)
end := start + 2000000000
rl.AddRange(Range{Start: start, End: end})
}

b.ResetTimer()
for i := 0; i < b.N; i++ {
rl.CountTotalValid()
}
}

func BenchmarkValidate500Ranges(b *testing.B) {
rl := &RangeList{}

for i := 0; i < 500; i++ {
start := int(100000000000000 + int64(i)*2000000000000)
end := start + 50000000
rl.AddRange(Range{Start: start, End: end})
}

testNum := 100000000000000 + 25000000

b.ResetTimer()
for i := 0; i < b.N; i++ {
rl.IsValid(testNum)
}
}

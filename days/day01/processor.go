package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Dial represents a dial with values in 0..99
type Dial struct{ v int }

func NewDial() *Dial       { return &Dial{v: 50} }
func (d *Dial) Value() int { return d.v }

// Entry with original steps (no modulo)
type Entry struct {
	Raw       string
	Dir       rune
	StepsOrig int
}

// countZeroCrossings counts how many times we pass through or land on 0 during a rotation
// Does not count if we start at 0
// countZeroCrossings counts crossings (including landing) of 0 using original steps.
// dirSign: +1 right, -1 left.
func countZeroCrossings(start, stepsOrig, dirSign int) int {
	if stepsOrig == 0 {
		return 0
	}
	if dirSign > 0 { // right
		return (start + stepsOrig) / 100
	}
	// left
	if start == 0 {
		return stepsOrig / 100
	}
	if stepsOrig < start {
		return 0
	}
	return 1 + (stepsOrig-start)/100
}

// ParseEntry parses a raw line into Entry. Returns error on invalid input.
func ParseEntry(raw string) (Entry, error) {
	s := strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(raw), " ", ""))
	if len(s) < 2 {
		return Entry{}, fmt.Errorf("invalid entry")
	}
	dir := rune(s[0])
	n, err := strconv.Atoi(s[1:])
	if err != nil {
		return Entry{}, err
	}
	return Entry{Raw: s, Dir: dir, StepsOrig: n}, nil
}

// processEntries takes raw lines (possibly with spaces) and returns formatted
// output lines and the number of times the dial ended at exactly 0 (mode="exact")
// or passed through 0 (mode="passes").
func processEntries(lines []string, mode string) ([]string, int) {
	d := NewDial()
	zeroCount := 0
	var outs []string
	for _, raw := range lines {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}
		e, err := ParseEntry(raw)
		if err != nil {
			continue
		}
		start := d.Value()
		dirSign := 1
		if e.Dir == 'L' {
			dirSign = -1
		} else if e.Dir != 'R' {
			continue
		}
		crossings := countZeroCrossings(start, e.StepsOrig, dirSign)
		stepsMod := e.StepsOrig % 100
		if dirSign < 0 {
			d.v = ((d.v-stepsMod)%100 + 100) % 100
		} else {
			d.v = (d.v + stepsMod) % 100
		}
		outs = append(outs, fmt.Sprintf("%s %d -> %d", e.Raw, start, d.Value()))
		if mode == "passes" {
			zeroCount += crossings
		} else if mode == "exact" && d.Value() == 0 {
			zeroCount++
		}
	}
	return outs, zeroCount
}

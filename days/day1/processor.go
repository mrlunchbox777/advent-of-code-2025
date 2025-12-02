package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Dial represents a dial with values in 0..99
type Dial struct {
	v int
}

func NewDial() *Dial { return &Dial{v: 50} }

func (d *Dial) Value() int { return d.v }

func (d *Dial) RotateLeft(n int) int {
	shift := n % 100
	oldV := d.v
	d.v = ((d.v-shift)%100 + 100) % 100
	return countZeroCrossings(oldV, d.v, -shift)
}

func (d *Dial) RotateRight(n int) int {
	shift := n % 100
	oldV := d.v
	d.v = (d.v + shift) % 100
	return countZeroCrossings(oldV, d.v, shift)
}

// Entry represents a parsed instruction like L68 or R5
type Entry struct {
	Raw   string // normalized raw form like "L68"
	Dir   rune
	Shift int
}

// countZeroCrossings counts how many times we pass through or land on 0 during a rotation
// Does not count if we start at 0
func countZeroCrossings(start, end, shift int) int {
	if shift == 0 || start == 0 {
		return 0
	}

	// Count if we land on 0
	if end == 0 {
		return 1
	}

	// Count if we pass through 0 (but don't land on it)
	absShift := shift
	if absShift < 0 {
		absShift = -absShift
	}

	// How many complete wraps?
	completeWraps := absShift / 100
	if completeWraps > 0 {
		return completeWraps
	}

	// Check if we cross 0 in the partial rotation (without landing on it)
	if shift > 0 {
		// Moving right: cross 0 if start + shift >= 100 and we didn't land on 0
		if start+shift >= 100 {
			return 1
		}
	} else {
		// Moving left: cross 0 if start + shift < 0 and we didn't land on 0
		if start+shift < 0 {
			return 1
		}
	}
	return 0
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
	return Entry{Raw: s, Dir: dir, Shift: n % 100}, nil
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
		var crossings int
		if e.Dir == 'L' {
			crossings = d.RotateLeft(e.Shift)
		} else if e.Dir == 'R' {
			crossings = d.RotateRight(e.Shift)
		} else {
			continue
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

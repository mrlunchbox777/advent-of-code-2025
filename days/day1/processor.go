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

func (d *Dial) RotateLeft(n int) {
    shift := n % 100
    d.v = ((d.v - shift) % 100 + 100) % 100
}

func (d *Dial) RotateRight(n int) {
    d.v = (d.v + (n % 100)) % 100
}

// Entry represents a parsed instruction like L68 or R5
type Entry struct {
    Raw   string // normalized raw form like "L68"
    Dir   rune
    Shift int
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
// output lines and the number of times the dial ended at exactly 0.
func processEntries(lines []string) ([]string, int) {
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
        if e.Dir == 'L' {
            d.RotateLeft(e.Shift)
        } else if e.Dir == 'R' {
            d.RotateRight(e.Shift)
        } else {
            continue
        }
        outs = append(outs, fmt.Sprintf("%s %d -> %d", e.Raw, start, d.Value()))
        if d.Value() == 0 {
            zeroCount++
        }
    }

    return outs, zeroCount
}

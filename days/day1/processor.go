package main

import (
	"fmt"
	"strconv"
	"strings"
)

// processEntries takes raw lines (possibly with spaces) and returns formatted
// output lines and the number of times the dial ended at exactly 0.
func processEntries(lines []string) ([]string, int) {
	dial := 50
	zeroCount := 0
	var outs []string

	for _, raw := range lines {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}
		entry := strings.ToUpper(strings.ReplaceAll(raw, " ", ""))
		if len(entry) < 2 {
			// skip invalid
			continue
		}
		dir := entry[0]
		numStr := entry[1:]
		n, err := strconv.Atoi(numStr)
		if err != nil {
			continue
		}

		start := dial
		shift := n % 100
		if dir == 'L' {
			dial = ((dial-shift)%100 + 100) % 100
		} else if dir == 'R' {
			dial = (dial + shift) % 100
		} else {
			continue
		}

		outs = append(outs, fmt.Sprintf("%s %d -> %d", entry, start, dial))
		if dial == 0 {
			zeroCount++
		}
	}

	return outs, zeroCount
}

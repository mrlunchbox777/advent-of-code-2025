package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Machine represents a toggle machine with desired state and options
type Machine struct {
	DesiredState []bool
	Options      [][]int
}

// ParseMachine parses a line into a Machine
func ParseMachine(line string) (*Machine, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, fmt.Errorf("empty line")
	}

	// Find sections
	desiredStart := strings.Index(line, "[")
	desiredEnd := strings.Index(line, "]")
	if desiredStart == -1 || desiredEnd == -1 {
		return nil, fmt.Errorf("invalid format: missing desired state")
	}

	desiredStr := line[desiredStart+1 : desiredEnd]
	desired := make([]bool, len(desiredStr))
	for i, ch := range desiredStr {
		if ch == '#' {
			desired[i] = true
		} else if ch != '.' {
			return nil, fmt.Errorf("invalid character in desired state: %c", ch)
		}
	}

	// Parse options (everything between ] and {)
	optionsStr := line[desiredEnd+1:]
	curlyBraceIdx := strings.Index(optionsStr, "{")
	if curlyBraceIdx != -1 {
		optionsStr = optionsStr[:curlyBraceIdx]
	}

	var options [][]int
	optionsStr = strings.TrimSpace(optionsStr)
	if optionsStr != "" {
		// Find all parenthesized groups
		for {
			start := strings.Index(optionsStr, "(")
			if start == -1 {
				break
			}
			end := strings.Index(optionsStr[start:], ")")
			if end == -1 {
				break
			}
			end += start

			numStr := optionsStr[start+1 : end]
			var nums []int
			for _, part := range strings.Split(numStr, ",") {
				part = strings.TrimSpace(part)
				if part == "" {
					continue
				}
				n, err := strconv.Atoi(part)
				if err != nil {
					return nil, fmt.Errorf("invalid number in option: %s", part)
				}
				nums = append(nums, n)
			}
			if len(nums) > 0 {
				options = append(options, nums)
			}
			optionsStr = optionsStr[end+1:]
		}
	}

	return &Machine{
		DesiredState: desired,
		Options:      options,
	}, nil
}

// ApplyOption applies an option to the current state
func (m *Machine) ApplyOption(state []bool, optionIdx int) []bool {
	newState := make([]bool, len(state))
	copy(newState, state)
	for _, pos := range m.Options[optionIdx] {
		if pos >= 0 && pos < len(newState) {
			newState[pos] = !newState[pos]
		}
	}
	return newState
}

// StatesEqual checks if two states are equal
func StatesEqual(a, b []bool) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Solve finds the minimum number of option selections to reach desired state using BFS
func (m *Machine) Solve() ([]int, int) {
	initialState := make([]bool, len(m.DesiredState))

	// BFS to find shortest path
	type SearchNode struct {
		state     []bool
		path      []int
		selection int
	}

	queue := []SearchNode{{state: initialState, path: []int{}, selection: 0}}
	visited := make(map[string]bool)
	visited[stateKey(initialState)] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if StatesEqual(current.state, m.DesiredState) {
			return current.path, current.selection
		}

		// Try each option
		for i := range m.Options {
			newState := m.ApplyOption(current.state, i)
			key := stateKey(newState)
			if !visited[key] {
				visited[key] = true
				newPath := make([]int, len(current.path))
				copy(newPath, current.path)
				newPath = append(newPath, i)
				queue = append(queue, SearchNode{
					state:     newState,
					path:      newPath,
					selection: current.selection + 1,
				})
			}
		}
	}

	return nil, -1
}

// stateKey creates a string key from a state for use in visited map
func stateKey(state []bool) string {
	var sb strings.Builder
	for _, b := range state {
		if b {
			sb.WriteByte('1')
		} else {
			sb.WriteByte('0')
		}
	}
	return sb.String()
}

// ProcessLines processes all lines and returns results
func ProcessLines(lines []string) ([]string, int) {
	var results []string
	totalSelections := 0
	lineNum := 1

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		machine, err := ParseMachine(line)
		if err != nil {
			results = append(results, fmt.Sprintf("Line %d: Error parsing - %v", lineNum, err))
			lineNum++
			continue
		}

		path, selections := machine.Solve()
		if selections == -1 {
			results = append(results, fmt.Sprintf("Line %d: No solution found", lineNum))
		} else {
			// Convert 0-indexed options to 1-indexed for display
			displayPath := make([]int, len(path))
			for i, p := range path {
				displayPath[i] = p + 1
			}
			results = append(results, fmt.Sprintf("Line %d: %d selections - options %v", lineNum, selections, displayPath))
			totalSelections += selections
		}
		lineNum++
	}

	return results, totalSelections
}

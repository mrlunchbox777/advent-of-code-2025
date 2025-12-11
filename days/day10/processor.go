package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Machine represents a toggle machine with desired state and options
type Machine struct {
	DesiredState   []bool
	TargetCounts   []int
	Options        [][]int
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

	// Parse target counts (in curly braces)
	var targetCounts []int
	curlyStart := strings.Index(line, "{")
	curlyEnd := strings.Index(line, "}")
	if curlyStart != -1 && curlyEnd != -1 {
		countsStr := line[curlyStart+1 : curlyEnd]
		for _, part := range strings.Split(countsStr, ",") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			n, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("invalid number in target counts: %s", part)
			}
			targetCounts = append(targetCounts, n)
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
		TargetCounts: targetCounts,
		Options:      options,
	}, nil
}

// ApplyOption applies an option to the current state (toggle mode)
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

// ApplyOptionCounter applies an option to counter state
func (m *Machine) ApplyOptionCounter(counts []int, optionIdx int) []int {
	newCounts := make([]int, len(counts))
	copy(newCounts, counts)
	for _, pos := range m.Options[optionIdx] {
		if pos >= 0 && pos < len(newCounts) {
			newCounts[pos]++
		}
	}
	return newCounts
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

// CountsEqual checks if two count arrays are equal
func CountsEqual(a, b []int) bool {
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

// Solve finds the minimum number of option selections to reach desired state using BFS (toggle mode)
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

// State represents a search state for A*
type State struct {
	counts    []int
	parent    *State
	option    int // which option led to this state
	cost      int // g(n): actual cost (path length)
	heuristic int // h(n): estimated cost to goal
}

// SolveCounter uses A* search with Manhattan distance heuristic
func (m *Machine) SolveCounter() ([]int, int) {
	
	initialCounts := make([]int, len(m.TargetCounts))
	initialState := &State{
		counts:    initialCounts,
		parent:    nil,
		option:    -1,
		cost:      0,
		heuristic: m.manhattanDistance(initialCounts),
	}
	
	// Priority queue (min-heap based on cost + heuristic)
	pq := &PriorityQueue{}
	pq.Push(initialState)
	
	// Use counts as int slices directly for faster comparison
	visited := make(map[string]bool)
	visited[countsKey(initialCounts)] = true
	
	statesExplored := 0
	maxStates := 2000000 // Safety limit
	
	for pq.Len() > 0 {
		current := pq.Pop()
		
		// Check if we reached the goal
		if current.heuristic == 0 {
			// Reconstruct path
			path := []int{}
			for node := current; node.parent != nil; node = node.parent {
				path = append([]int{node.option}, path...)
			}
			return path, len(path)
		}
		
		statesExplored++
		if statesExplored > maxStates {
			return nil, -1
		}
		
		// Try each option
		for i := range m.Options {
			newCounts := m.ApplyOptionCounter(current.counts, i)
			
			// Skip if any position exceeds target
			exceeds := false
			for j := range newCounts {
				if newCounts[j] > m.TargetCounts[j] {
					exceeds = true
					break
				}
			}
			if exceeds {
				continue
			}
			
			newKey := countsKey(newCounts)
			if !visited[newKey] {
				visited[newKey] = true
				newState := &State{
					counts:    newCounts,
					parent:    current,
					option:    i,
					cost:      current.cost + 1,
					heuristic: m.manhattanDistance(newCounts),
				}
				pq.Push(newState)
			}
		}
	}
	
	return nil, -1
}

// manhattanDistance calculates how far current counts are from target
func (m *Machine) manhattanDistance(counts []int) int {
	dist := 0
	for i := range counts {
		diff := m.TargetCounts[i] - counts[i]
		if diff > 0 {
			dist += diff
		}
	}
	return dist
}

// PriorityQueue implements a min-heap for A* search
type PriorityQueue struct {
	items []*State
}

func (pq *PriorityQueue) Len() int {
	return len(pq.items)
}

func (pq *PriorityQueue) Push(state *State) {
	pq.items = append(pq.items, state)
	pq.up(len(pq.items) - 1)
}

func (pq *PriorityQueue) Pop() *State {
	n := len(pq.items)
	pq.swap(0, n-1)
	item := pq.items[n-1]
	pq.items = pq.items[:n-1]
	if len(pq.items) > 0 {
		pq.down(0)
	}
	return item
}

func (pq *PriorityQueue) up(i int) {
	for {
		parent := (i - 1) / 2
		if parent == i || pq.less(parent, i) {
			break
		}
		pq.swap(parent, i)
		i = parent
	}
}

func (pq *PriorityQueue) down(i int) {
	for {
		left := 2*i + 1
		if left >= len(pq.items) {
			break
		}
		j := left
		if right := left + 1; right < len(pq.items) && pq.less(right, left) {
			j = right
		}
		if pq.less(i, j) {
			break
		}
		pq.swap(i, j)
		i = j
	}
}

func (pq *PriorityQueue) less(i, j int) bool {
	// Compare by f(n) = g(n) + h(n)
	fi := pq.items[i].cost + pq.items[i].heuristic
	fj := pq.items[j].cost + pq.items[j].heuristic
	return fi < fj
}

func (pq *PriorityQueue) swap(i, j int) {
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
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

// countsKey creates a string key from counts for use in visited map
func countsKey(counts []int) string {
	// Use fmt.Sprintf for faster string creation
	return fmt.Sprint(counts)
}

// ProcessLines processes all lines and returns results
func ProcessLines(lines []string, mode string) ([]string, int) {
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

		var path []int
		var selections int
		if mode == "counter" {
			path, selections = machine.SolveCounter()
		} else {
			path, selections = machine.Solve()
		}

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

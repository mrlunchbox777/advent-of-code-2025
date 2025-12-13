package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
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

// SolveCounter solves using Gaussian elimination with Dijkstra on free variables
func (m *Machine) SolveCounter() ([]int, int) {
	// Try Gaussian elimination - much faster than BFS
	if path, count := m.solveCounterGaussian(); path != nil {
		return path, count
	}
	
	// Fallback to BFS if Gaussian fails (rare)
	return m.solveCounterBFS()
}

// solveCounterGaussian uses Gaussian elimination to solve the linear system
func (m *Machine) solveCounterGaussian() ([]int, int) {
	numPositions := len(m.TargetCounts)
	numOptions := len(m.Options)
	
	// Build the matrix A where A[i][j] = 1 if option j affects position i
	A := make([][]float64, numPositions)
	for i := range A {
		A[i] = make([]float64, numOptions)
	}
	
	for optIdx, opt := range m.Options {
		for _, pos := range opt {
			if pos >= 0 && pos < numPositions {
				A[pos][optIdx] = 1
			}
		}
	}
	
	// Target vector b
	b := make([]float64, numPositions)
	for i := range b {
		b[i] = float64(m.TargetCounts[i])
	}
	
	// Perform Gaussian elimination with partial pivoting
	// Augment matrix [A | b]
	aug := make([][]float64, numPositions)
	for i := range aug {
		aug[i] = make([]float64, numOptions+1)
		copy(aug[i], A[i])
		aug[i][numOptions] = b[i]
	}
	
	// Forward elimination
	pivotCols := []int{}
	row := 0
	for col := 0; col < numOptions && row < numPositions; col++ {
		// Find pivot
		maxRow := row
		maxVal := abs(aug[row][col])
		for i := row + 1; i < numPositions; i++ {
			if abs(aug[i][col]) > maxVal {
				maxVal = abs(aug[i][col])
				maxRow = i
			}
		}
		
		if maxVal < 1e-10 {
			continue // Skip this column
		}
		
		// Swap rows
		aug[row], aug[maxRow] = aug[maxRow], aug[row]
		pivotCols = append(pivotCols, col)
		
		// Eliminate
		pivot := aug[row][col]
		for j := col; j <= numOptions; j++ {
			aug[row][j] /= pivot
		}
		
		for i := 0; i < numPositions; i++ {
			if i != row {
				factor := aug[i][col]
				for j := col; j <= numOptions; j++ {
					aug[i][j] -= factor * aug[row][j]
				}
			}
		}
		row++
	}
	
	// Identify free variables
	freeVars := []int{}
	isPivot := make([]bool, numOptions)
	for _, col := range pivotCols {
		isPivot[col] = true
	}
	for i := 0; i < numOptions; i++ {
		if !isPivot[i] {
			freeVars = append(freeVars, i)
		}
	}
	
	// Try to find integer solution
	if len(freeVars) == 0 {
		// Direct solution - check if all integer
		solution := make([]int, numOptions)
		for rowIdx, col := range pivotCols {
			val := aug[rowIdx][numOptions]
			if abs(val-float64(int(val+0.5))) > 0.01 {
				return nil, -1
			}
			solution[col] = int(val + 0.5)
			if solution[col] < 0 {
				return nil, -1
			}
		}
		return solutionToPath(solution), sumInts(solution)
	} else if len(freeVars) <= 4 {
		// Search over free variables (increased from 2 to 4)
		return m.searchFreeVars(aug, pivotCols, freeVars, numOptions)
	}
	
	// Too many free variables, fallback
	return nil, -1
}

type freeVarCandidate struct {
	freeVals []int
	cost     int // Heuristic cost (sum of free vars)
}

type PQItem struct {
	vals []int
	cost int
}

func (m *Machine) searchFreeVars(aug [][]float64, pivotCols, freeVars []int, numOptions int) ([]int, int) {
	// Exhaustive enumeration of all valid free variable combinations
	// Find the one with minimum cost
	
	maxTarget := 0
	for _, t := range m.TargetCounts {
		if t > maxTarget {
			maxTarget = t
		}
	}
	
	bestSolution := []int(nil)
	bestCost := int(1e9)
	
	// Use recursive enumeration
	var enumerate func(vals []int, index int)
	enumerate = func(vals []int, index int) {
		if index == len(freeVars) {
			// Try this combination
			if solution, cost := m.trySolution(aug, pivotCols, freeVars, vals, numOptions); solution != nil {
				if cost < bestCost {
					bestSolution = solution
					bestCost = cost
				}
			}
			return
		}
		
		// Try all values for this free variable
		for v := 0; v <= maxTarget; v++ {
			vals[index] = v
			enumerate(vals, index+1)
		}
	}
	
	vals := make([]int, len(freeVars))
	enumerate(vals, 0)
	
	if bestSolution != nil {
		return solutionToPath(bestSolution), bestCost
	}
	return nil, -1
}

func insertSorted(pq []freeVarCandidate, item freeVarCandidate) []freeVarCandidate {
	// Binary search insertion
	left, right := 0, len(pq)
	for left < right {
		mid := (left + right) / 2
		if pq[mid].cost < item.cost {
			left = mid + 1
		} else {
			right = mid
		}
	}
	// Insert at position 'left'
	pq = append(pq, freeVarCandidate{})
	copy(pq[left+1:], pq[left:])
	pq[left] = item
	return pq
}

func (m *Machine) trySolution(aug [][]float64, pivotCols, freeVars, freeVals []int, numOptions int) ([]int, int) {
	solution := make([]int, numOptions)
	
	// Set free variables
	for i, fv := range freeVars {
		solution[fv] = freeVals[i]
	}
	
	// Solve for pivot variables
	for rowIdx, col := range pivotCols {
		val := aug[rowIdx][numOptions] // RHS
		for _, fv := range freeVars {
			val -= aug[rowIdx][fv] * float64(solution[fv])
		}
		
		if abs(val-float64(int(val+0.5))) > 0.01 {
			return nil, -1
		}
		intVal := int(val + 0.5)
		if intVal < 0 {
			return nil, -1
		}
		solution[col] = intVal
	}
	
	cost := sumInts(solution)
	return solution, cost
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func sumInts(arr []int) int {
	sum := 0
	for _, v := range arr {
		sum += v
	}
	return sum
}

func solutionToPath(solution []int) []int {
	var path []int
	for optIdx, count := range solution {
		for i := 0; i < count; i++ {
			path = append(path, optIdx)
		}
	}
	return path
}

// solveCounterGreedy uses iterative refinement
func (m *Machine) solveCounterGreedy() ([]int, int) {
	// Try to solve using a mathematical approach
	// Count how many times each option should be used
	optionCounts := make([]int, len(m.Options))
	
	// Start with a greedy estimate
	remaining := make([]int, len(m.TargetCounts))
	copy(remaining, m.TargetCounts)
	
	maxIterations := 500
	for iteration := 0; iteration < maxIterations; iteration++ {
		// Check if solved
		allZero := true
		for _, r := range remaining {
			if r != 0 {
				allZero = false
				break
			}
		}
		if allZero {
			// Convert option counts to path
			var path []int
			for optIdx, count := range optionCounts {
				for i := 0; i < count; i++ {
					path = append(path, optIdx)
				}
			}
			return path, len(path)
		}
		
		// Find the option that best satisfies the most critical remaining position
		bestOption := -1
		bestScore := float64(-1)
		
		for i := range m.Options {
			score := float64(0)
			canUse := true
			usefulCount := 0
			
			for _, pos := range m.Options[i] {
				if pos >= 0 && pos < len(remaining) {
					if remaining[pos] > 0 {
						score += float64(remaining[pos])
						usefulCount++
					} else if remaining[pos] < 0 {
						canUse = false
						break
					}
				}
			}
			
			// Normalize by the number of positions this option affects
			if canUse && usefulCount > 0 {
				score = score / float64(len(m.Options[i]))
				if score > bestScore {
					bestScore = score
					bestOption = i
				}
			}
		}
		
		if bestOption == -1 {
			return nil, -1
		}
		
		// Use this option once
		optionCounts[bestOption]++
		for _, pos := range m.Options[bestOption] {
			if pos >= 0 && pos < len(remaining) {
				remaining[pos]--
			}
		}
	}
	
	return nil, -1
}

// BeamState represents a state in beam search
type BeamState struct {
	counts []int
	path   []int
	cost   int
}

// solveCounterBFS uses pure BFS for guaranteed optimal solution  
func (m *Machine) solveCounterBFS() ([]int, int) {
	type QueueItem struct {
		counts []int
		path   []int
	}
	
	initialCounts := make([]int, len(m.TargetCounts))
	queue := []QueueItem{{counts: initialCounts, path: []int{}}}
	
	// Use map for visited states - track with compact encoding
	visited := make(map[string]struct{})
	visited[encodeCounts(initialCounts)] = struct{}{}
	
	itemsProcessed := 0
	
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		itemsProcessed++
		
		// Check for solution
		if CountsEqual(current.counts, m.TargetCounts) {
			return current.path, len(current.path)
		}
		
		// Memory safety: if we've processed too many items without finding solution, abort
		// This prevents indefinite memory growth
		if itemsProcessed > 50000000 {
			return nil, -1
		}
		
		// Try each option
		for i := range m.Options {
			newCounts := m.ApplyOptionCounter(current.counts, i)
			
			// Skip if any position exceeds target
			valid := true
			for j := range newCounts {
				if newCounts[j] > m.TargetCounts[j] {
					valid = false
					break
				}
			}
			if !valid {
				continue
			}
			
			key := encodeCounts(newCounts)
			if _, seen := visited[key]; !seen {
				visited[key] = struct{}{}
				newPath := make([]int, len(current.path)+1)
				copy(newPath, current.path)
				newPath[len(current.path)] = i
				queue = append(queue, QueueItem{counts: newCounts, path: newPath})
			}
		}
	}
	
	return nil, -1
}

// encodeCounts encodes counts as a compact string
func encodeCounts(counts []int) string {
	return fmt.Sprint(counts)
}

// BeamPriorityQueue is a min-heap for beam states
type BeamPriorityQueue struct {
	items []*BeamState
}

func (pq *BeamPriorityQueue) Len() int { return len(pq.items) }

func (pq *BeamPriorityQueue) Push(state *BeamState) {
	pq.items = append(pq.items, state)
	pq.up(len(pq.items) - 1)
}

func (pq *BeamPriorityQueue) Pop() *BeamState {
	n := len(pq.items)
	pq.swap(0, n-1)
	item := pq.items[n-1]
	pq.items = pq.items[:n-1]
	if len(pq.items) > 0 {
		pq.down(0)
	}
	return item
}

func (pq *BeamPriorityQueue) up(i int) {
	for {
		parent := (i - 1) / 2
		if parent == i || pq.items[parent].cost <= pq.items[i].cost {
			break
		}
		pq.swap(parent, i)
		i = parent
	}
}

func (pq *BeamPriorityQueue) down(i int) {
	for {
		left := 2*i + 1
		if left >= len(pq.items) {
			break
		}
		j := left
		if right := left + 1; right < len(pq.items) && pq.items[right].cost < pq.items[left].cost {
			j = right
		}
		if pq.items[i].cost <= pq.items[j].cost {
			break
		}
		pq.swap(i, j)
		i = j
	}
}

func (pq *BeamPriorityQueue) swap(i, j int) {
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
}

func sortBeamStatePointers(states []*BeamState) {
	n := len(states)
	if n < 2 {
		return
	}
	quicksortBeamStatePointers(states, 0, n-1)
}

func quicksortBeamStatePointers(states []*BeamState, low, high int) {
	if low < high {
		pi := partitionBeamStatePointers(states, low, high)
		quicksortBeamStatePointers(states, low, pi-1)
		quicksortBeamStatePointers(states, pi+1, high)
	}
}

func partitionBeamStatePointers(states []*BeamState, low, high int) int {
	pivot := states[high].cost
	i := low - 1
	for j := low; j < high; j++ {
		if states[j].cost < pivot {
			i++
			states[i], states[j] = states[j], states[i]
		}
	}
	states[i+1], states[high] = states[high], states[i+1]
	return i + 1
}

// sortBeamStates sorts beam states by cost
func sortBeamStates(states []BeamState) {
	// Simple insertion sort for small arrays, quicksort for large
	n := len(states)
	if n < 20 {
		for i := 1; i < n; i++ {
			key := states[i]
			j := i - 1
			for j >= 0 && states[j].cost > key.cost {
				states[j+1] = states[j]
				j--
			}
			states[j+1] = key
		}
	} else {
		quicksortBeamStates(states, 0, n-1)
	}
}

func quicksortBeamStates(states []BeamState, low, high int) {
	if low < high {
		pi := partitionBeamStates(states, low, high)
		quicksortBeamStates(states, low, pi-1)
		quicksortBeamStates(states, pi+1, high)
	}
}

func partitionBeamStates(states []BeamState, low, high int) int {
	pivot := states[high].cost
	i := low - 1
	for j := low; j < high; j++ {
		if states[j].cost < pivot {
			i++
			states[i], states[j] = states[j], states[i]
		}
	}
	states[i+1], states[high] = states[high], states[i+1]
	return i + 1
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

// ProcessLines processes all lines and prints results as it goes
func ProcessLines(lines []string, mode string) int {
	totalSelections := 0
	lineNum := 1

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		machine, err := ParseMachine(line)
		if err != nil {
			fmt.Printf("Line %d: Error parsing - %v\n", lineNum, err)
			lineNum++
			continue
		}

		startTime := time.Now()
		var path []int
		var selections int
		if mode == "counter" {
			path, selections = machine.SolveCounter()
		} else {
			path, selections = machine.Solve()
		}
		elapsed := time.Since(startTime)

		if selections == -1 {
			fmt.Printf("Line %d: No solution found (%.2fs)\n", lineNum, elapsed.Seconds())
		} else {
			// Convert 0-indexed options to 1-indexed for display
			displayPath := make([]int, len(path))
			for i, p := range path {
				displayPath[i] = p + 1
			}
			fmt.Printf("Line %d: %d selections - options %v (%.2fs)\n", lineNum, selections, displayPath, elapsed.Seconds())
			totalSelections += selections
		}
		lineNum++
	}

	return totalSelections
}

package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Coordinate struct {
	X, Y, Z int
	ID      int
}

func (c *Coordinate) Distance(other *Coordinate) float64 {
	dx := float64(c.X - other.X)
	dy := float64(c.Y - other.Y)
	dz := float64(c.Z - other.Z)
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

type UnionFind struct {
	parent []int
	rank   []int
}

func NewUnionFind(size int) *UnionFind {
	parent := make([]int, size)
	rank := make([]int, size)
	for i := range parent {
		parent[i] = i
	}
	return &UnionFind{parent: parent, rank: rank}
}

func (uf *UnionFind) Find(x int) int {
	if uf.parent[x] != x {
		uf.parent[x] = uf.Find(uf.parent[x])
	}
	return uf.parent[x]
}

func (uf *UnionFind) Union(x, y int) bool {
	rootX := uf.Find(x)
	rootY := uf.Find(y)
	
	if rootX == rootY {
		return false
	}
	
	if uf.rank[rootX] < uf.rank[rootY] {
		uf.parent[rootX] = rootY
	} else if uf.rank[rootX] > uf.rank[rootY] {
		uf.parent[rootY] = rootX
	} else {
		uf.parent[rootY] = rootX
		uf.rank[rootX]++
	}
	return true
}

type Edge struct {
	idx1, idx2 int
	distance   float64
	index      int
}

type EdgeHeap []*Edge

func (h EdgeHeap) Len() int           { return len(h) }
func (h EdgeHeap) Less(i, j int) bool { return h[i].distance < h[j].distance }
func (h EdgeHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *EdgeHeap) Push(x interface{}) {
	n := len(*h)
	edge := x.(*Edge)
	edge.index = n
	*h = append(*h, edge)
}

func (h *EdgeHeap) Pop() interface{} {
	old := *h
	n := len(old)
	edge := old[n-1]
	old[n-1] = nil
	edge.index = -1
	*h = old[0 : n-1]
	return edge
}

type CoordinateSet struct {
	coords      []*Coordinate
	uf          *UnionFind
	connections map[string]bool
	edgeHeap    *EdgeHeap
	useHeap     bool
}

func NewCoordinateSet(coords []*Coordinate) *CoordinateSet {
	return &CoordinateSet{
		coords:      coords,
		uf:          NewUnionFind(len(coords)),
		connections: make(map[string]bool),
		useHeap:     false,
	}
}

func NewCoordinateSetWithHeap(coords []*Coordinate) *CoordinateSet {
	cs := &CoordinateSet{
		coords:      coords,
		uf:          NewUnionFind(len(coords)),
		connections: make(map[string]bool),
		useHeap:     true,
	}
	
	h := make(EdgeHeap, 0)
	heap.Init(&h)
	
	for i := 0; i < len(coords); i++ {
		for j := i + 1; j < len(coords); j++ {
			dist := coords[i].Distance(coords[j])
			heap.Push(&h, &Edge{
				idx1:     i,
				idx2:     j,
				distance: dist,
			})
		}
	}
	
	cs.edgeHeap = &h
	return cs
}

func (cs *CoordinateSet) connectionKey(idx1, idx2 int) string {
	if idx1 > idx2 {
		idx1, idx2 = idx2, idx1
	}
	return fmt.Sprintf("%d-%d", idx1, idx2)
}

func (cs *CoordinateSet) hasConnection(idx1, idx2 int) bool {
	return cs.connections[cs.connectionKey(idx1, idx2)]
}

func (cs *CoordinateSet) FindClosestPair() (int, int, float64) {
	if cs.useHeap {
		for cs.edgeHeap.Len() > 0 {
			edge := heap.Pop(cs.edgeHeap).(*Edge)
			if !cs.hasConnection(edge.idx1, edge.idx2) {
				return edge.idx1, edge.idx2, edge.distance
			}
		}
		return -1, -1, 0
	}
	
	minDist := math.MaxFloat64
	idx1, idx2 := -1, -1
	
	for i := 0; i < len(cs.coords); i++ {
		for j := i + 1; j < len(cs.coords); j++ {
			if cs.hasConnection(i, j) {
				continue
			}
			
			dist := cs.coords[i].Distance(cs.coords[j])
			if dist < minDist {
				minDist = dist
				idx1 = i
				idx2 = j
			}
		}
	}
	
	return idx1, idx2, minDist
}

func (cs *CoordinateSet) Connect(idx1, idx2 int) {
	cs.connections[cs.connectionKey(idx1, idx2)] = true
	cs.uf.Union(idx1, idx2)
}

func (cs *CoordinateSet) GetGroups() map[int][]int {
	groups := make(map[int][]int)
	for i := range cs.coords {
		root := cs.uf.Find(i)
		groups[root] = append(groups[root], i)
	}
	return groups
}

func (cs *CoordinateSet) GetTopGroups(n int) [][]int {
	groups := cs.GetGroups()
	
	var groupList [][]int
	for _, members := range groups {
		groupList = append(groupList, members)
	}
	
	sort.Slice(groupList, func(i, j int) bool {
		return len(groupList[i]) > len(groupList[j])
	})
	
	if len(groupList) > n {
		groupList = groupList[:n]
	}
	
	return groupList
}

func parseFile(filepath string) ([]*Coordinate, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	
	var coords []*Coordinate
	scanner := bufio.NewScanner(file)
	id := 0
	
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		
		parts := strings.Split(line, ",")
		if len(parts) < 3 {
			continue
		}
		
		x, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			continue
		}
		
		y, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			continue
		}
		
		z, err := strconv.Atoi(strings.TrimSpace(parts[2]))
		if err != nil {
			continue
		}
		
		coords = append(coords, &Coordinate{X: x, Y: y, Z: z, ID: id})
		id++
	}
	
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	
	return coords, nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: day8 <mode> <filepath> [max_rounds]")
		fmt.Println("  mode: 'grouping' or 'completion'")
		os.Exit(1)
	}
	
	mode := os.Args[1]
	filepath := os.Args[2]
	maxRounds := 1000
	
	if mode != "grouping" && mode != "completion" {
		fmt.Fprintf(os.Stderr, "Invalid mode: %s (must be 'grouping' or 'completion')\n", mode)
		os.Exit(1)
	}
	
	if len(os.Args) >= 4 {
		rounds, err := strconv.Atoi(os.Args[3])
		if err == nil && rounds > 0 {
			maxRounds = rounds
		}
	}
	
	coords, err := parseFile(filepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing file: %v\n", err)
		os.Exit(1)
	}
	
	if len(coords) == 0 {
		fmt.Println("No coordinates found")
		os.Exit(1)
	}
	
	if mode == "grouping" {
		runGroupingMode(coords, maxRounds)
	} else {
		runCompletionMode(coords)
	}
}

func runGroupingMode(coords []*Coordinate, maxRounds int) {
	fmt.Printf("Loaded %d coordinates\n", len(coords))
	fmt.Printf("Running up to %d rounds\n\n", maxRounds)
	
	cs := NewCoordinateSet(coords)
	
	for round := 1; round <= maxRounds; round++ {
		idx1, idx2, dist := cs.FindClosestPair()
		
		if idx1 == -1 || idx2 == -1 {
			fmt.Printf("\nAll coordinates are connected after %d rounds\n", round-1)
			break
		}
		
		cs.Connect(idx1, idx2)
		
		fmt.Printf("Round %d: Connected (%d,%d,%d) and (%d,%d,%d) - Distance: %.2f\n",
			round,
			coords[idx1].X, coords[idx1].Y, coords[idx1].Z,
			coords[idx2].X, coords[idx2].Y, coords[idx2].Z,
			dist)
		
		topGroups := cs.GetTopGroups(5)
		fmt.Printf("  Top 5 groups: ")
		for i, group := range topGroups {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Printf("%d", len(group))
		}
		fmt.Println()
	}
	
	fmt.Println("\n=== Final Results ===")
	top3 := cs.GetTopGroups(3)
	
	fmt.Println("Top 3 largest groups:")
	product := 1
	for i, group := range top3 {
		fmt.Printf("  Group %d: %d members\n", i+1, len(group))
		product *= len(group)
	}
	
	fmt.Printf("\nProduct of top 3 group sizes: %d\n", product)
}

func runCompletionMode(coords []*Coordinate) {
	fmt.Printf("Loaded %d coordinates\n", len(coords))
	fmt.Println("Running until all coordinates are in a single group")
	fmt.Println()
	
	cs := NewCoordinateSetWithHeap(coords)
	round := 0
	var completionIdx1, completionIdx2 int
	var completionRound int
	
	numGroups := len(coords)
	
	for {
		idx1, idx2, dist := cs.FindClosestPair()
		
		if idx1 == -1 || idx2 == -1 {
			fmt.Printf("\nAll possible connections made after %d rounds\n", round)
			break
		}
		
		round++
		
		root1 := cs.uf.Find(idx1)
		root2 := cs.uf.Find(idx2)
		cs.Connect(idx1, idx2)
		
		if root1 != root2 {
			numGroups--
		}
		
		if round <= 100 || round%1000 == 0 {
			fmt.Printf("Round %d: Connected (%d,%d,%d) and (%d,%d,%d) - Distance: %.2f\n",
				round,
				coords[idx1].X, coords[idx1].Y, coords[idx1].Z,
				coords[idx2].X, coords[idx2].Y, coords[idx2].Z,
				dist)
			
			topGroups := cs.GetTopGroups(5)
			fmt.Printf("  Top 5 groups: ")
			for i, group := range topGroups {
				if i > 0 {
					fmt.Print(", ")
				}
				fmt.Printf("%d", len(group))
			}
			fmt.Printf(" (%d total groups)\n", numGroups)
		}
		
		if numGroups == 1 && completionRound == 0 {
			completionIdx1 = idx1
			completionIdx2 = idx2
			completionRound = round
			fmt.Printf("*** All coordinates now in a single group at round %d! ***\n", round)
		}
	}
	
	fmt.Println("\n=== Final Results ===")
	if completionRound > 0 {
		fmt.Printf("Single group achieved at round %d\n", completionRound)
		fmt.Printf("Completion connection: (%d,%d,%d) and (%d,%d,%d)\n",
			coords[completionIdx1].X, coords[completionIdx1].Y, coords[completionIdx1].Z,
			coords[completionIdx2].X, coords[completionIdx2].Y, coords[completionIdx2].Z)
		
		product := coords[completionIdx1].X * coords[completionIdx2].X
		fmt.Printf("Product of X coordinates: %d × %d = %d\n",
			coords[completionIdx1].X, coords[completionIdx2].X, product)
	} else {
		fmt.Println("Did not reach single group")
	}
}

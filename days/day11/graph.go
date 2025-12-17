package main

import (
	"fmt"
	"strings"
)

// Node represents a graph node with its connections
type Node struct {
	Name        string
	Connections []string
}

// Graph represents a directed graph
type Graph struct {
	Nodes map[string]*Node
}

// Path represents a path through the graph
type Path []string

// String formats a path as node names separated by arrows
func (p Path) String() string {
	return strings.Join(p, "->")
}

// ParseGraph parses lines into a Graph structure
func ParseGraph(lines []string) (*Graph, error) {
	g := &Graph{
		Nodes: make(map[string]*Node),
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid line format: %s", line)
		}

		nodeName := strings.TrimSpace(parts[0])
		connectionsStr := strings.TrimSpace(parts[1])

		node := &Node{
			Name:        nodeName,
			Connections: []string{},
		}

		if connectionsStr != "" {
			for _, conn := range strings.Fields(connectionsStr) {
				node.Connections = append(node.Connections, conn)
			}
		}

		g.Nodes[nodeName] = node
	}

	return g, nil
}

// FindAllPaths finds all unique paths from start to end node using DFS
func (g *Graph) FindAllPaths(start, end string) []Path {
	var allPaths []Path
	visited := make(map[string]bool)
	currentPath := []string{}

	g.dfs(start, end, visited, currentPath, &allPaths)

	return allPaths
}

// CountAllPaths counts all unique paths from start to end without storing them
func (g *Graph) CountAllPaths(start, end string) int {
	visited := make(map[string]bool)
	memo := make(map[string]int)
	return g.dfsCountSimpleMemoized(start, end, visited, memo)
}

// dfsCountSimpleMemoized performs DFS with memoization for simple path counting
func (g *Graph) dfsCountSimpleMemoized(current, end string, visited map[string]bool, memo map[string]int) int {
	// For part 1, memo key is just the current node (no required nodes to track)
	if count, found := memo[current]; found {
		return count
	}
	
	visited[current] = true

	var count int
	if current == end {
		count = 1
	} else {
		node, exists := g.Nodes[current]
		if exists {
			for _, neighbor := range node.Connections {
				if !visited[neighbor] {
					count += g.dfsCountSimpleMemoized(neighbor, end, visited, memo)
				}
			}
		}
	}

	visited[current] = false
	
	memo[current] = count
	return count
}

// dfs performs depth-first search to find all paths
func (g *Graph) dfs(current, end string, visited map[string]bool, currentPath []string, allPaths *[]Path) {
	// Add current node to path
	currentPath = append(currentPath, current)
	visited[current] = true

	// If we reached the end, save this path
	if current == end {
		// Make a copy of the path
		pathCopy := make(Path, len(currentPath))
		copy(pathCopy, currentPath)
		*allPaths = append(*allPaths, pathCopy)
	} else {
		// Continue searching from neighbors
		node, exists := g.Nodes[current]
		if exists {
			for _, neighbor := range node.Connections {
				if !visited[neighbor] {
					g.dfs(neighbor, end, visited, currentPath, allPaths)
				}
			}
		}
	}

	// Backtrack
	visited[current] = false
}

// FindPathsWithRequiredNodes finds all paths from start to end that visit all required nodes
func (g *Graph) FindPathsWithRequiredNodes(start, end string, required []string) []Path {
	requiredSet := make(map[string]bool, len(required))
	for _, node := range required {
		requiredSet[node] = true
	}
	
	var allPaths []Path
	visited := make(map[string]bool)
	currentPath := []string{}
	requiredVisited := make(map[string]bool)
	
	g.dfsWithRequired(start, end, requiredSet, required, visited, currentPath, requiredVisited, &allPaths)
	
	return allPaths
}

// dfsWithRequired performs DFS to find paths with required nodes
func (g *Graph) dfsWithRequired(current, end string, requiredSet map[string]bool, required []string,
	visited map[string]bool, currentPath []string, requiredVisited map[string]bool, allPaths *[]Path) {
	
	currentPath = append(currentPath, current)
	visited[current] = true
	
	wasRequired := false
	if requiredSet[current] && !requiredVisited[current] {
		requiredVisited[current] = true
		wasRequired = true
	}
	
	if current == end {
		if len(requiredVisited) == len(required) {
			pathCopy := make(Path, len(currentPath))
			copy(pathCopy, currentPath)
			*allPaths = append(*allPaths, pathCopy)
		}
	} else {
		node, exists := g.Nodes[current]
		if exists {
			for _, neighbor := range node.Connections {
				if !visited[neighbor] {
					g.dfsWithRequired(neighbor, end, requiredSet, required, visited, currentPath, requiredVisited, allPaths)
				}
			}
		}
	}
	
	visited[current] = false
	if wasRequired {
		delete(requiredVisited, current)
	}
}

// CountPathsWithRequiredNodes counts all paths from start to end that visit all required nodes without storing them
func (g *Graph) CountPathsWithRequiredNodes(start, end string, required []string) int {
	requiredSet := make(map[string]bool, len(required))
	for _, node := range required {
		requiredSet[node] = true
	}
	
	// Use memoization: map from (current_node, required_nodes_visited_set) -> count
	memo := make(map[string]int)
	visited := make(map[string]bool)
	requiredVisited := make(map[string]bool)
	
	return g.dfsCountMemoized(start, end, requiredSet, required, visited, requiredVisited, memo)
}

// dfsCountMemoized performs DFS with memoization
func (g *Graph) dfsCountMemoized(current, end string, requiredSet map[string]bool, required []string,
	visited map[string]bool, requiredVisited map[string]bool, memo map[string]int) int {
	
	// Build memo key from current node and required nodes visited
	memoKey := g.buildMemoKey(current, requiredVisited, required)
	
	// Check memo
	if count, found := memo[memoKey]; found {
		return count
	}
	
	visited[current] = true
	
	wasRequired := false
	if requiredSet[current] && !requiredVisited[current] {
		requiredVisited[current] = true
		wasRequired = true
	}
	
	var count int
	if current == end {
		if len(requiredVisited) == len(required) {
			count = 1
		} else {
			count = 0
		}
	} else {
		node, exists := g.Nodes[current]
		if exists {
			for _, neighbor := range node.Connections {
				if !visited[neighbor] {
					count += g.dfsCountMemoized(neighbor, end, requiredSet, required, visited, requiredVisited, memo)
				}
			}
		}
	}
	
	visited[current] = false
	if wasRequired {
		delete(requiredVisited, current)
	}
	
	// Store in memo
	memo[memoKey] = count
	return count
}

// buildMemoKey creates a unique key for memoization based on current node and required nodes visited
func (g *Graph) buildMemoKey(current string, requiredVisited map[string]bool, required []string) string {
	// Create a sorted, consistent representation of which required nodes have been visited
	var visited []string
	for _, req := range required {
		if requiredVisited[req] {
			visited = append(visited, req)
		}
	}
	// Simple encoding: "node|req1,req2,..."
	if len(visited) == 0 {
		return current + "|"
	}
	return current + "|" + strings.Join(visited, ",")
}



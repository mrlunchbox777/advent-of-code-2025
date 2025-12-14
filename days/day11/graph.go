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
	count := 0
	g.dfsCount(start, end, visited, &count)
	return count
}

// dfsCount performs DFS and only counts paths
func (g *Graph) dfsCount(current, end string, visited map[string]bool, count *int) {
	visited[current] = true

	if current == end {
		*count++
	} else {
		node, exists := g.Nodes[current]
		if exists {
			for _, neighbor := range node.Connections {
				if !visited[neighbor] {
					g.dfsCount(neighbor, end, visited, count)
				}
			}
		}
	}

	visited[current] = false
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
	
	// Pre-compute reachability to prune dead branches
	canReachEnd := g.computeReachability(end, true)    // Reverse: who can reach 'end'
	canReachFromStart := g.computeReachability(start, false) // Forward: who can be reached from 'start'
	
	// Check that all required nodes can be reached
	for _, req := range required {
		if !canReachFromStart[req] {
			// Can't reach a required node from start
			return []Path{}
		}
		if !canReachEnd[req] {
			// Required node can't reach end
			return []Path{}
		}
	}
	
	// Pre-compute which nodes can reach each required node  
	canReachFromNode := make(map[string]map[string]bool)
	for nodeName := range g.Nodes {
		canReachFromNode[nodeName] = g.computeReachability(nodeName, false)
	}
	
	var allPaths []Path
	visited := make(map[string]bool)
	currentPath := make([]string, 0, 100) // Pre-allocate larger size
	requiredVisited := make(map[string]bool)
	
	maxDepth := 50
	g.dfsWithPruning(start, end, requiredSet, required, visited, currentPath, requiredVisited, 
		&allPaths, canReachEnd, canReachFromStart, canReachFromNode, 0, maxDepth)
	
	return allPaths
}

// computeReachability finds all nodes that can reach (reverse=true) or be reached from (reverse=false) a target node
func (g *Graph) computeReachability(target string, reverse bool) map[string]bool {
	reachable := make(map[string]bool)
	visited := make(map[string]bool)
	queue := []string{target}
	reachable[target] = true
	visited[target] = true
	
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		
		if reverse {
			// Find nodes that can reach current (look at who points to current)
			for nodeName, node := range g.Nodes {
				if visited[nodeName] {
					continue
				}
				for _, conn := range node.Connections {
					if conn == current {
						reachable[nodeName] = true
						visited[nodeName] = true
						queue = append(queue, nodeName)
						break
					}
				}
			}
		} else {
			// Find nodes that current can reach (forward search)
			node, exists := g.Nodes[current]
			if !exists {
				continue
			}
			for _, neighbor := range node.Connections {
				if !visited[neighbor] {
					reachable[neighbor] = true
					visited[neighbor] = true
					queue = append(queue, neighbor)
				}
			}
		}
	}
	
	return reachable
}

// dfsWithPruning performs DFS with aggressive dead-end pruning
func (g *Graph) dfsWithPruning(current, end string, requiredSet map[string]bool, required []string,
	visited map[string]bool, currentPath []string, requiredVisited map[string]bool,
	allPaths *[]Path, canReachEnd, canReachFromStart map[string]bool, 
	canReachFromNode map[string]map[string]bool, depth, maxDepth int) {
	
	// Depth limit
	if depth > maxDepth {
		return
	}
	
	// Store current path length to avoid allocating new slice
	pathLen := len(currentPath)
	currentPath = append(currentPath, current)
	visited[current] = true
	
	// Track required nodes
	wasRequired := false
	if requiredSet[current] && !requiredVisited[current] {
		requiredVisited[current] = true
		wasRequired = true
	}
	
	// Check if we reached the end
	if current == end {
		if len(requiredVisited) == len(required) {
			pathCopy := make(Path, len(currentPath))
			copy(pathCopy, currentPath)
			*allPaths = append(*allPaths, pathCopy)
		}
	} else {
		// Continue exploring neighbors
		node, exists := g.Nodes[current]
		if exists {
			for _, neighbor := range node.Connections {
				if visited[neighbor] {
					continue
				}
				
				// PRUNING: Check if this neighbor is a dead end
				// 1. Can this neighbor reach the end?
				if !canReachEnd[neighbor] {
					continue
				}
				
				// 2. Can this neighbor be reached from start? (should always be true, but check anyway)
				if !canReachFromStart[neighbor] {
					continue
				}
				
				// 3. Check if we can still visit all remaining required nodes
				canReachAll := true
				for _, req := range required {
					if !requiredVisited[req] && req != neighbor {
						// Can neighbor reach this required node?
						if reachable, exists := canReachFromNode[neighbor]; !exists || !reachable[req] {
							canReachAll = false
							break
						}
					}
				}
				if !canReachAll {
					continue
				}
				
				g.dfsWithPruning(neighbor, end, requiredSet, required, visited, currentPath, 
					requiredVisited, allPaths, canReachEnd, canReachFromStart, canReachFromNode, depth+1, maxDepth)
			}
		}
	}
	
	// Backtrack - reuse the slice by truncating instead of creating new slices
	currentPath = currentPath[:pathLen]
	visited[current] = false
	if wasRequired {
		requiredVisited[current] = false
	}
}

// encodeRequired creates a string key from required nodes visited
func encodeRequired(required map[string]bool) string {
	// Create a sorted list of required nodes for consistent encoding
	nodes := make([]string, 0, len(required))
	for node := range required {
		nodes = append(nodes, node)
	}
	return fmt.Sprint(nodes)
}

// encodeVisited creates a compact string key from visited nodes
func encodeVisited(visited map[string]bool) string {
	return fmt.Sprint(visited)
}

// CountPathsWithRequiredNodes counts all paths from start to end that visit all required nodes without storing them
func (g *Graph) CountPathsWithRequiredNodes(start, end string, required []string) int {
	requiredSet := make(map[string]bool, len(required))
	for _, node := range required {
		requiredSet[node] = true
	}
	
	// Pre-compute reachability to prune dead branches
	canReachEnd := g.computeReachability(end, true)
	canReachFromStart := g.computeReachability(start, false)
	
	// Check that all required nodes can be reached
	for _, req := range required {
		if !canReachFromStart[req] || !canReachEnd[req] {
			return 0
		}
	}
	
	// Pre-compute which nodes can reach each required node  
	canReachFromNode := make(map[string]map[string]bool)
	for nodeName := range g.Nodes {
		canReachFromNode[nodeName] = g.computeReachability(nodeName, false)
	}
	
	visited := make(map[string]bool)
	requiredVisited := make(map[string]bool)
	count := 0
	maxDepth := 50
	
	g.dfsCountWithPruning(start, end, requiredSet, required, visited, requiredVisited, 
		&count, canReachEnd, canReachFromStart, canReachFromNode, 0, maxDepth)
	
	return count
}

// dfsCountWithPruning performs DFS with pruning and only counts paths
func (g *Graph) dfsCountWithPruning(current, end string, requiredSet map[string]bool, required []string,
	visited map[string]bool, requiredVisited map[string]bool,
	count *int, canReachEnd, canReachFromStart map[string]bool, 
	canReachFromNode map[string]map[string]bool, depth, maxDepth int) {
	
	if depth > maxDepth {
		return
	}
	
	visited[current] = true
	
	wasRequired := false
	if requiredSet[current] && !requiredVisited[current] {
		requiredVisited[current] = true
		wasRequired = true
	}
	
	if current == end {
		if len(requiredVisited) == len(required) {
			*count++
		}
	} else {
		node, exists := g.Nodes[current]
		if exists {
			for _, neighbor := range node.Connections {
				if visited[neighbor] {
					continue
				}
				
				if !canReachEnd[neighbor] || !canReachFromStart[neighbor] {
					continue
				}
				
				canReachAll := true
				for _, req := range required {
					if !requiredVisited[req] && req != neighbor {
						if reachable, exists := canReachFromNode[neighbor]; !exists || !reachable[req] {
							canReachAll = false
							break
						}
					}
				}
				if !canReachAll {
					continue
				}
				
				g.dfsCountWithPruning(neighbor, end, requiredSet, required, visited, requiredVisited, 
					count, canReachEnd, canReachFromStart, canReachFromNode, depth+1, maxDepth)
			}
		}
	}
	
	visited[current] = false
	if wasRequired {
		requiredVisited[current] = false
	}
}



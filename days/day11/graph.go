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
	// Use BFS to find all paths systematically
	type SearchState struct {
		node            string
		path            []string
		visited         map[string]bool
		requiredVisited map[string]bool
	}
	
	requiredSet := make(map[string]bool, len(required))
	for _, node := range required {
		requiredSet[node] = true
	}
	
	initialVisited := make(map[string]bool)
	initialVisited[start] = true
	initialRequired := make(map[string]bool)
	if requiredSet[start] {
		initialRequired[start] = true
	}
	
	queue := []SearchState{{
		node:            start,
		path:            []string{start},
		visited:         initialVisited,
		requiredVisited: initialRequired,
	}}
	
	var allPaths []Path
	
	// Track best depth seen for each (node, required_visited) to prune inefficient paths
	// We only care about which required nodes have been visited, not the full path
	type StateKey string
	bestDepth := make(map[StateKey]int)
	
	maxDepth := 40 // Reasonable limit to prevent infinite search  
	maxQueueSize := 100000000 // Limit queue size to prevent memory exhaustion
	
	for len(queue) > 0 {
		// Memory safety
		if len(queue) > maxQueueSize {
			// Queue is too large, abort
			return allPaths
		}
		current := queue[0]
		queue = queue[1:]
		
		currentDepth := len(current.path)
		
		// Depth limit
		if currentDepth > maxDepth {
			continue
		}
		
		// Check if we reached the end
		if current.node == end {
			// Check if all required nodes were visited
			if len(current.requiredVisited) == len(required) {
				pathCopy := make(Path, len(current.path))
				copy(pathCopy, current.path)
				allPaths = append(allPaths, pathCopy)
			}
			continue
		}
		
		// Explore neighbors
		node, exists := g.Nodes[current.node]
		if !exists {
			continue
		}
		
		for _, neighbor := range node.Connections {
			if current.visited[neighbor] {
				continue
			}
			
			// Create new state for this neighbor
			newVisited := make(map[string]bool, len(current.visited)+1)
			for k, v := range current.visited {
				newVisited[k] = v
			}
			newVisited[neighbor] = true
			
			newRequired := make(map[string]bool, len(current.requiredVisited))
			for k, v := range current.requiredVisited {
				newRequired[k] = v
			}
			if requiredSet[neighbor] {
				newRequired[neighbor] = true
			}
			
			// Create state key for pruning - only use node and required nodes visited
			stateKey := StateKey(neighbor + "|" + encodeRequired(newRequired))
			
			// Only continue if this is a better or equal path to this state
			if prevDepth, seen := bestDepth[stateKey]; seen && currentDepth+1 > prevDepth {
				continue
			}
			bestDepth[stateKey] = currentDepth + 1
			
			newPath := make([]string, len(current.path)+1)
			copy(newPath, current.path)
			newPath[len(current.path)] = neighbor
			
			queue = append(queue, SearchState{
				node:            neighbor,
				path:            newPath,
				visited:         newVisited,
				requiredVisited: newRequired,
			})
		}
	}
	
	return allPaths
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



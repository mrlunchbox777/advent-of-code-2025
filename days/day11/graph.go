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
	var allPaths []Path
	visited := make(map[string]bool)
	currentPath := []string{}
	requiredSet := make(map[string]bool)
	for _, node := range required {
		requiredSet[node] = true
	}

	g.dfsWithRequired(start, end, visited, currentPath, requiredSet, &allPaths)

	return allPaths
}

// dfsWithRequired performs DFS to find paths that visit all required nodes
func (g *Graph) dfsWithRequired(current, end string, visited map[string]bool, currentPath []string, required map[string]bool, allPaths *[]Path) {
	// Add current node to path
	currentPath = append(currentPath, current)
	visited[current] = true

	// If we reached the end, check if all required nodes were visited
	if current == end {
		allVisited := true
		for reqNode := range required {
			found := false
			for _, node := range currentPath {
				if node == reqNode {
					found = true
					break
				}
			}
			if !found {
				allVisited = false
				break
			}
		}
		
		if allVisited {
			// Make a copy of the path
			pathCopy := make(Path, len(currentPath))
			copy(pathCopy, currentPath)
			*allPaths = append(*allPaths, pathCopy)
		}
	} else {
		// Continue searching from neighbors
		node, exists := g.Nodes[current]
		if exists {
			for _, neighbor := range node.Connections {
				if !visited[neighbor] {
					g.dfsWithRequired(neighbor, end, visited, currentPath, required, allPaths)
				}
			}
		}
	}

	// Backtrack
	visited[current] = false
}

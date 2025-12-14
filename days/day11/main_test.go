package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseGraph(t *testing.T) {
	lines := []string{
		"you: bbb ccc",
		"bbb: ddd eee",
		"ccc: ddd eee fff",
		"ddd: ggg",
		"eee: out",
	}

	graph, err := ParseGraph(lines)
	if err != nil {
		t.Fatalf("ParseGraph() error = %v", err)
	}

	// Check that nodes were created
	if len(graph.Nodes) != 5 {
		t.Errorf("Expected 5 nodes, got %d", len(graph.Nodes))
	}

	// Check specific node
	you := graph.Nodes["you"]
	if you == nil {
		t.Fatal("Node 'you' not found")
	}
	if len(you.Connections) != 2 {
		t.Errorf("Expected 'you' to have 2 connections, got %d", len(you.Connections))
	}
	if you.Connections[0] != "bbb" || you.Connections[1] != "ccc" {
		t.Errorf("Expected 'you' connections to be [bbb ccc], got %v", you.Connections)
	}
}

func TestFindAllPathsSimple(t *testing.T) {
	lines := []string{
		"you: bbb",
		"bbb: out",
	}

	graph, err := ParseGraph(lines)
	if err != nil {
		t.Fatalf("ParseGraph() error = %v", err)
	}

	paths := graph.FindAllPaths("you", "out")
	
	expectedCount := 1
	if len(paths) != expectedCount {
		t.Errorf("Expected %d paths, got %d", expectedCount, len(paths))
	}

	if len(paths) > 0 {
		expectedPath := "you->bbb->out"
		if paths[0].String() != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, paths[0].String())
		}
	}
}

func TestFindAllPathsMultiple(t *testing.T) {
	lines := []string{
		"you: bbb ccc",
		"bbb: out",
		"ccc: out",
	}

	graph, err := ParseGraph(lines)
	if err != nil {
		t.Fatalf("ParseGraph() error = %v", err)
	}

	paths := graph.FindAllPaths("you", "out")
	
	expectedCount := 2
	if len(paths) != expectedCount {
		t.Errorf("Expected %d paths, got %d", expectedCount, len(paths))
	}
}

func TestProcessExampleData(t *testing.T) {
	p := filepath.Join(".", "example-data.txt")
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("failed to read example data: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(b)), "\n")
	graph, err := ParseGraph(lines)
	if err != nil {
		t.Fatalf("ParseGraph() error = %v", err)
	}

	paths := graph.FindAllPaths("you", "out")

	expectedCount := 5
	if len(paths) != expectedCount {
		t.Errorf("Expected %d paths, got %d", expectedCount, len(paths))
		for i, path := range paths {
			t.Logf("Path %d: %s", i+1, path.String())
		}
	}
}

func TestPathString(t *testing.T) {
	path := Path{"you", "bbb", "out"}
	expected := "you->bbb->out"
	if path.String() != expected {
		t.Errorf("Expected %s, got %s", expected, path.String())
	}
}

func TestEmptyConnections(t *testing.T) {
	lines := []string{
		"you: bbb",
		"bbb:",
	}

	graph, err := ParseGraph(lines)
	if err != nil {
		t.Fatalf("ParseGraph() error = %v", err)
	}

	bbb := graph.Nodes["bbb"]
	if bbb == nil {
		t.Fatal("Node 'bbb' not found")
	}
	if len(bbb.Connections) != 0 {
		t.Errorf("Expected 'bbb' to have 0 connections, got %d", len(bbb.Connections))
	}
}

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

func TestProcessExampleDataAllMode(t *testing.T) {
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

func TestFindPathsWithRequiredNodes(t *testing.T) {
	lines := []string{
		"start: aaa bbb",
		"aaa: req1",
		"bbb: req1",
		"req1: req2",
		"req2: end",
	}

	graph, err := ParseGraph(lines)
	if err != nil {
		t.Fatalf("ParseGraph() error = %v", err)
	}

	paths := graph.FindPathsWithRequiredNodes("start", "end", []string{"req1", "req2"})

	expectedCount := 2
	if len(paths) != expectedCount {
		t.Errorf("Expected %d paths, got %d", expectedCount, len(paths))
		for i, path := range paths {
			t.Logf("Path %d: %s", i+1, path.String())
		}
	}

	// Verify both required nodes are in each path
	for _, path := range paths {
		hasReq1 := false
		hasReq2 := false
		for _, node := range path {
			if node == "req1" {
				hasReq1 = true
			}
			if node == "req2" {
				hasReq2 = true
			}
		}
		if !hasReq1 || !hasReq2 {
			t.Errorf("Path %s missing required nodes", path.String())
		}
	}
}

func TestFindPathsWithRequiredNodesNoMatch(t *testing.T) {
	lines := []string{
		"start: aaa",
		"aaa: end",
	}

	graph, err := ParseGraph(lines)
	if err != nil {
		t.Fatalf("ParseGraph() error = %v", err)
	}

	paths := graph.FindPathsWithRequiredNodes("start", "end", []string{"missing"})

	if len(paths) != 0 {
		t.Errorf("Expected 0 paths, got %d", len(paths))
	}
}

func TestProcessExampleData2MustVisitMode(t *testing.T) {
	p := filepath.Join(".", "example-data-2.txt")
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("failed to read example data: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(b)), "\n")
	graph, err := ParseGraph(lines)
	if err != nil {
		t.Fatalf("ParseGraph() error = %v", err)
	}

	paths := graph.FindPathsWithRequiredNodes("svr", "out", []string{"dac", "fft"})

	expectedCount := 2
	if len(paths) != expectedCount {
		t.Errorf("Expected %d paths, got %d", expectedCount, len(paths))
		for i, path := range paths {
			t.Logf("Path %d: %s", i+1, path.String())
		}
	}

	// Verify all paths contain both required nodes
	for _, path := range paths {
		hasDac := false
		hasFft := false
		for _, node := range path {
			if node == "dac" {
				hasDac = true
			}
			if node == "fft" {
				hasFft = true
			}
		}
		if !hasDac || !hasFft {
			t.Errorf("Path %s missing required nodes (dac and/or fft)", path.String())
		}
	}
}

func BenchmarkFindPathsWithRequiredNodes(b *testing.B) {
	p := filepath.Join(".", "example-data-2.txt")
	data, err := os.ReadFile(p)
	if err != nil {
		b.Fatalf("failed to read example data: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	graph, err := ParseGraph(lines)
	if err != nil {
		b.Fatalf("ParseGraph() error = %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = graph.FindPathsWithRequiredNodes("svr", "out", []string{"dac", "fft"})
	}
}

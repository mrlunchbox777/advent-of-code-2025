package main

import (
	"math"
	"os"
	"testing"
)

func TestCoordinateDistance(t *testing.T) {
	c1 := &Coordinate{X: 0, Y: 0, Z: 0}
	c2 := &Coordinate{X: 3, Y: 4, Z: 0}
	
	dist := c1.Distance(c2)
	expected := 5.0
	
	if math.Abs(dist-expected) > 0.001 {
		t.Errorf("Distance = %.2f, want %.2f", dist, expected)
	}
}

func TestUnionFind(t *testing.T) {
	uf := NewUnionFind(5)
	
	if uf.Find(0) == uf.Find(1) {
		t.Error("Initially separate elements should have different roots")
	}
	
	uf.Union(0, 1)
	if uf.Find(0) != uf.Find(1) {
		t.Error("After union, elements should have same root")
	}
	
	uf.Union(2, 3)
	uf.Union(1, 2)
	
	if uf.Find(0) != uf.Find(3) {
		t.Error("Transitively connected elements should have same root")
	}
	
	if uf.Find(0) == uf.Find(4) {
		t.Error("Unconnected element should have different root")
	}
}

func TestFindClosestPair(t *testing.T) {
	coords := []*Coordinate{
		{X: 0, Y: 0, Z: 0, ID: 0},
		{X: 1, Y: 0, Z: 0, ID: 1},
		{X: 10, Y: 0, Z: 0, ID: 2},
	}
	
	cs := NewCoordinateSet(coords)
	idx1, idx2, dist := cs.FindClosestPair()
	
	if (idx1 != 0 || idx2 != 1) && (idx1 != 1 || idx2 != 0) {
		t.Errorf("Expected closest pair to be 0 and 1, got %d and %d", idx1, idx2)
	}
	
	if math.Abs(dist-1.0) > 0.001 {
		t.Errorf("Distance = %.2f, want 1.0", dist)
	}
}

func TestGetGroups(t *testing.T) {
	coords := []*Coordinate{
		{X: 0, Y: 0, Z: 0, ID: 0},
		{X: 1, Y: 0, Z: 0, ID: 1},
		{X: 2, Y: 0, Z: 0, ID: 2},
		{X: 10, Y: 0, Z: 0, ID: 3},
	}
	
	cs := NewCoordinateSet(coords)
	cs.Connect(0, 1)
	cs.Connect(1, 2)
	
	groups := cs.GetGroups()
	
	if len(groups) != 2 {
		t.Errorf("Expected 2 groups, got %d", len(groups))
	}
	
	largestGroupSize := 0
	for _, members := range groups {
		if len(members) > largestGroupSize {
			largestGroupSize = len(members)
		}
	}
	
	if largestGroupSize != 3 {
		t.Errorf("Expected largest group size 3, got %d", largestGroupSize)
	}
}

func TestGetTopGroups(t *testing.T) {
	coords := []*Coordinate{
		{X: 0, Y: 0, Z: 0, ID: 0},
		{X: 1, Y: 0, Z: 0, ID: 1},
		{X: 2, Y: 0, Z: 0, ID: 2},
		{X: 10, Y: 0, Z: 0, ID: 3},
		{X: 11, Y: 0, Z: 0, ID: 4},
	}
	
	cs := NewCoordinateSet(coords)
	cs.Connect(0, 1)
	cs.Connect(1, 2)
	cs.Connect(3, 4)
	
	topGroups := cs.GetTopGroups(2)
	
	if len(topGroups) != 2 {
		t.Errorf("Expected 2 groups, got %d", len(topGroups))
	}
	
	if len(topGroups[0]) != 3 {
		t.Errorf("Expected first group size 3, got %d", len(topGroups[0]))
	}
	
	if len(topGroups[1]) != 2 {
		t.Errorf("Expected second group size 2, got %d", len(topGroups[1]))
	}
}

func TestParseFile(t *testing.T) {
	content := `162,817,812
57,618,57
906,360,560
`
	tmpfile, err := os.CreateTemp("", "test-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	
	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}
	
	coords, err := parseFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("parseFile() error = %v", err)
	}
	
	if len(coords) != 3 {
		t.Errorf("Expected 3 coordinates, got %d", len(coords))
	}
	
	if coords[0].X != 162 || coords[0].Y != 817 || coords[0].Z != 812 {
		t.Errorf("First coordinate incorrect: got (%d,%d,%d)", coords[0].X, coords[0].Y, coords[0].Z)
	}
}

func TestExampleDataGrouping(t *testing.T) {
	if _, err := os.Stat("example-data.txt"); os.IsNotExist(err) {
		t.Skip("example-data.txt not found")
	}
	
	coords, err := parseFile("example-data.txt")
	if err != nil {
		t.Fatalf("Failed to parse example-data.txt: %v", err)
	}
	
	cs := NewCoordinateSet(coords)
	
	for round := 1; round <= 10; round++ {
		idx1, idx2, _ := cs.FindClosestPair()
		if idx1 == -1 || idx2 == -1 {
			break
		}
		cs.Connect(idx1, idx2)
	}
	
	top3 := cs.GetTopGroups(3)
	
	if len(top3) < 3 {
		t.Fatalf("Expected at least 3 groups, got %d", len(top3))
	}
	
	product := len(top3[0]) * len(top3[1]) * len(top3[2])
	expected := 40
	
	if product != expected {
		t.Errorf("Product = %d, want %d (group sizes: %d, %d, %d)",
			product, expected, len(top3[0]), len(top3[1]), len(top3[2]))
	}
	
	t.Logf("After 10 rounds: group sizes are %d, %d, %d (product=%d)",
		len(top3[0]), len(top3[1]), len(top3[2]), product)
}

func TestExampleDataCompletion(t *testing.T) {
	if _, err := os.Stat("example-data.txt"); os.IsNotExist(err) {
		t.Skip("example-data.txt not found")
	}
	
	coords, err := parseFile("example-data.txt")
	if err != nil {
		t.Fatalf("Failed to parse example-data.txt: %v", err)
	}
	
	cs := NewCoordinateSet(coords)
	var completionIdx1, completionIdx2 int
	
	for round := 1; ; round++ {
		idx1, idx2, _ := cs.FindClosestPair()
		if idx1 == -1 || idx2 == -1 {
			break
		}
		
		cs.Connect(idx1, idx2)
		
		groups := cs.GetGroups()
		if len(groups) == 1 && completionIdx1 == 0 {
			completionIdx1 = idx1
			completionIdx2 = idx2
			break
		}
	}
	
	if completionIdx1 == 0 {
		t.Fatal("Never reached single group")
	}
	
	product := coords[completionIdx1].X * coords[completionIdx2].X
	expected := 25272
	
	if product != expected {
		t.Errorf("Product = %d, want %d (X coords: %d, %d)",
			product, expected, coords[completionIdx1].X, coords[completionIdx2].X)
	}
	
	t.Logf("Completion connection X coords: %d × %d = %d",
		coords[completionIdx1].X, coords[completionIdx2].X, product)
}

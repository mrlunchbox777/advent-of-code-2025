package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Point struct {
	X int
	Y int
}

func NewPoint(x, y int) Point {
	return Point{X: x, Y: y}
}

func (p Point) String() string {
	return fmt.Sprintf("(%d,%d)", p.X, p.Y)
}

type Rectangle struct {
	P1 Point
	P2 Point
}

func NewRectangle(p1, p2 Point) Rectangle {
	return Rectangle{P1: p1, P2: p2}
}

func (r Rectangle) Area() int {
	width := abs(r.P2.X-r.P1.X) + 1
	height := abs(r.P2.Y-r.P1.Y) + 1
	return width * height
}

func (r Rectangle) String() string {
	return fmt.Sprintf("Rectangle[%s, %s] Area=%d", r.P1, r.P2, r.Area())
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func parsePoint(line string) (Point, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return Point{}, fmt.Errorf("empty line")
	}
	parts := strings.Split(line, ",")
	if len(parts) != 2 {
		return Point{}, fmt.Errorf("invalid format: %s", line)
	}
	x, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return Point{}, fmt.Errorf("invalid x coordinate: %v", err)
	}
	y, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return Point{}, fmt.Errorf("invalid y coordinate: %v", err)
	}
	return NewPoint(x, y), nil
}

func processCoordinates(lines []string) int {
	var points []Point
	for _, line := range lines {
		if p, err := parsePoint(line); err == nil {
			points = append(points, p)
		}
	}

	if len(points) < 2 {
		return 0
	}

	maxArea := 0
	for i := 0; i < len(points); i++ {
		for j := i + 1; j < len(points); j++ {
			rect := NewRectangle(points[i], points[j])
			area := rect.Area()
			if area > maxArea {
				maxArea = area
			}
		}
	}

	return maxArea
}

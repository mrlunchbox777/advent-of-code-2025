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

func processCoordinates(lines []string, mode string) int {
	var points []Point
	for _, line := range lines {
		if p, err := parsePoint(line); err == nil {
			points = append(points, p)
		}
	}

	if len(points) < 2 {
		return 0
	}

	if mode == "original" {
		return processOriginal(points)
	}
	return processContained(points)
}

func processOriginal(points []Point) int {
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

func processContained(points []Point) int {
	if len(points) < 3 {
		return 0
	}

	orderedPoints := orderPointsAsPolygon(points)

	maxArea := 0
	for i := 0; i < len(points); i++ {
		for j := i + 1; j < len(points); j++ {
			p1, p2 := points[i], points[j]
			
			if p1.X == p2.X || p1.Y == p2.Y {
				continue
			}

			minX, maxX := min(p1.X, p2.X), max(p1.X, p2.X)
			minY, maxY := min(p1.Y, p2.Y), max(p1.Y, p2.Y)

			if isRectangleContained(minX, maxX, minY, maxY, orderedPoints) {
				rect := NewRectangle(p1, p2)
				area := rect.Area()
				if area > maxArea {
					maxArea = area
				}
			}
		}
	}
	return maxArea
}

func orderPointsAsPolygon(points []Point) []Point {
	if len(points) == 0 {
		return points
	}

	centroidX, centroidY := 0, 0
	for _, p := range points {
		centroidX += p.X
		centroidY += p.Y
	}
	centroidX /= len(points)
	centroidY /= len(points)

	type anglePoint struct {
		point Point
		angle float64
	}
	
	var anglePoints []anglePoint
	for _, p := range points {
		dx := float64(p.X - centroidX)
		dy := float64(p.Y - centroidY)
		angle := 0.0
		if dx != 0 || dy != 0 {
			angle = float64(dy)
			if dx != 0 {
				angle = float64(dy) / float64(dx)
				if dx < 0 {
					angle += 1000
				}
			}
		}
		anglePoints = append(anglePoints, anglePoint{p, angle})
	}

	for i := 0; i < len(anglePoints)-1; i++ {
		for j := i + 1; j < len(anglePoints); j++ {
			if anglePoints[i].angle > anglePoints[j].angle {
				anglePoints[i], anglePoints[j] = anglePoints[j], anglePoints[i]
			}
		}
	}

	result := make([]Point, len(anglePoints))
	for i, ap := range anglePoints {
		result[i] = ap.point
	}
	return result
}

func isRectangleContained(minX, maxX, minY, maxY int, points []Point) bool {
	for x := minX; x <= maxX; x++ {
		if !isPointInPolygon(x, minY, points) {
			return false
		}
		if !isPointInPolygon(x, maxY, points) {
			return false
		}
	}
	for y := minY; y <= maxY; y++ {
		if !isPointInPolygon(minX, y, points) {
			return false
		}
		if !isPointInPolygon(maxX, y, points) {
			return false
		}
	}
	
	centerX := (minX + maxX) / 2
	centerY := (minY + maxY) / 2
	return isPointInPolygon(centerX, centerY, points)
}

func isPointInPolygon(x, y int, polygon []Point) bool {
	n := len(polygon)
	
	// Check if point is a vertex
	for _, p := range polygon {
		if p.X == x && p.Y == y {
			return true
		}
	}
	
	// Check if point is on an edge
	for i := 0; i < n; i++ {
		p1 := polygon[i]
		p2 := polygon[(i+1)%n]
		
		if isPointOnSegment(x, y, p1.X, p1.Y, p2.X, p2.Y) {
			return true
		}
	}
	
	// Ray casting for interior points
	inside := false
	p1x, p1y := polygon[0].X, polygon[0].Y
	for i := 1; i <= n; i++ {
		p2x, p2y := polygon[i%n].X, polygon[i%n].Y
		
		if y > min(p1y, p2y) {
			if y <= max(p1y, p2y) {
				if x <= max(p1x, p2x) {
					if p1y != p2y {
						xinters := float64(y-p1y)*float64(p2x-p1x)/float64(p2y-p1y) + float64(p1x)
						if p1x == p2x || float64(x) <= xinters {
							inside = !inside
						}
					}
				}
			}
		}
		p1x, p1y = p2x, p2y
	}
	
	return inside
}

func isPointOnSegment(px, py, x1, y1, x2, y2 int) bool {
	// Check if point is within bounding box
	if px < min(x1, x2) || px > max(x1, x2) || py < min(y1, y2) || py > max(y1, y2) {
		return false
	}
	
	// Check if point is collinear with segment
	crossProduct := (py-y1)*(x2-x1) - (px-x1)*(y2-y1)
	return crossProduct == 0
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

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
	area, _ := processCoordinatesWithResult(lines, mode)
	return area
}

func processCoordinatesWithResult(lines []string, mode string) (int, Rectangle) {
	var points []Point
	for _, line := range lines {
		if p, err := parsePoint(line); err == nil {
			points = append(points, p)
		}
	}

	if len(points) < 2 {
		return 0, Rectangle{}
	}

	if mode == "original" {
		return processOriginalWithResult(points)
	}
	return processContainedWithResult(points)
}

func processOriginal(points []Point) int {
	area, _ := processOriginalWithResult(points)
	return area
}

func processOriginalWithResult(points []Point) (int, Rectangle) {
	maxArea := 0
	var maxRect Rectangle
	for i := 0; i < len(points); i++ {
		for j := i + 1; j < len(points); j++ {
			rect := NewRectangle(points[i], points[j])
			area := rect.Area()
			if area > maxArea {
				maxArea = area
				maxRect = rect
			}
		}
	}
	return maxArea, maxRect
}

func processContained(points []Point) int {
	area, _ := processContainedWithResult(points)
	return area
}

func processContainedWithResult(points []Point) (int, Rectangle) {
	if len(points) < 3 {
		return 0, Rectangle{}
	}

	orderedPoints := orderPointsAsPolygon(points)

	maxArea := 0
	var maxRect Rectangle
	for i := 0; i < len(points); i++ {
		for j := i + 1; j < len(points); j++ {
			p1, p2 := points[i], points[j]
			
			if p1.X == p2.X || p1.Y == p2.Y {
				continue
			}

			minX, maxX := min(p1.X, p2.X), max(p1.X, p2.X)
			minY, maxY := min(p1.Y, p2.Y), max(p1.Y, p2.Y)

			// Check if rectangle is contained in polygon
			if !isRectangleContained(minX, maxX, minY, maxY, orderedPoints) {
				continue
			}
			
			// Check if any other points fall inside the rectangle (excluding corners)
			if hasPointsInside(minX, maxX, minY, maxY, points, p1, p2) {
				continue
			}

			rect := NewRectangle(p1, p2)
			area := rect.Area()
			if area > maxArea {
				maxArea = area
				maxRect = rect
			}
		}
	}
	return maxArea, maxRect
}

func hasPointsInside(minX, maxX, minY, maxY int, points []Point, corner1, corner2 Point) bool {
	// Check if any points (other than the two corners) fall STRICTLY inside the rectangle
	// Points on the boundary (edges) are allowed
	for _, p := range points {
		// Skip the corner points
		if (p.X == corner1.X && p.Y == corner1.Y) || (p.X == corner2.X && p.Y == corner2.Y) {
			continue
		}
		
		// Check if point is strictly inside (not on boundary)
		if p.X > minX && p.X < maxX && p.Y > minY && p.Y < maxY {
			return true
		}
	}
	return false
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
	// Check the 4 corners first (fast reject)
	if !isPointInPolygon(minX, minY, points) {
		return false
	}
	if !isPointInPolygon(maxX, minY, points) {
		return false
	}
	if !isPointInPolygon(minX, maxY, points) {
		return false
	}
	if !isPointInPolygon(maxX, maxY, points) {
		return false
	}
	
	// Key insight: For a rectangle to be contained, NO polygon edge can intersect
	// the rectangle's EDGES (unless it's coincident with them)
	// This ensures no part of the rectangle extends outside the polygon
	
	rectEdges := []struct{ x1, y1, x2, y2 int }{
		{minX, minY, maxX, minY}, // top
		{minX, maxY, maxX, maxY}, // bottom
		{minX, minY, minX, maxY}, // left
		{maxX, minY, maxX, maxY}, // right
	}
	
	for i := 0; i < len(points); i++ {
		p1 := points[i]
		p2 := points[(i+1)%len(points)]
		
		// For each polygon edge, check if it properly intersects any rectangle edge
		for _, re := range rectEdges {
			if properSegmentIntersection(p1.X, p1.Y, p2.X, p2.Y, re.x1, re.y1, re.x2, re.y2) {
				// This polygon edge crosses a rectangle edge - rectangle not contained
				return false
			}
		}
	}
	
	// Check center point - if corners are in and no edges intersect, center must be in
	centerX := (minX + maxX) / 2
	centerY := (minY + maxY) / 2
	return isPointInPolygon(centerX, centerY, points)
}

func edgeCrossesRectangle(x1, y1, x2, y2, minX, maxX, minY, maxY int) bool {
	// Check if line segment crosses through the rectangle interior
	// This is different from touching or lying along the boundary
	
	// Count how many rectangle edges this segment properly crosses
	crossCount := 0
	
	// Check intersection with each rectangle edge
	// We only count proper crossings (not endpoint touches)
	edges := []struct{ x1, y1, x2, y2 int }{
		{minX, minY, maxX, minY}, // top
		{minX, maxY, maxX, maxY}, // bottom
		{minX, minY, minX, maxY}, // left
		{maxX, minY, maxX, maxY}, // right
	}
	
	for _, edge := range edges {
		if properSegmentIntersection(x1, y1, x2, y2, edge.x1, edge.y1, edge.x2, edge.y2) {
			crossCount++
		}
	}
	
	// If segment crosses 2 opposite edges, it goes through the rectangle
	return crossCount >= 2
}

func properSegmentIntersection(x1, y1, x2, y2, x3, y3, x4, y4 int) bool {
	// Check for proper intersection (segments cross, not just touch at endpoints)
	d1 := direction(x3, y3, x4, y4, x1, y1)
	d2 := direction(x3, y3, x4, y4, x2, y2)
	d3 := direction(x1, y1, x2, y2, x3, y3)
	d4 := direction(x1, y1, x2, y2, x4, y4)
	
	// Proper intersection: segments cross
	if ((d1 > 0 && d2 < 0) || (d1 < 0 && d2 > 0)) &&
		((d3 > 0 && d4 < 0) || (d3 < 0 && d4 > 0)) {
		return true
	}
	
	return false
}

func segmentsIntersect(x1, y1, x2, y2, x3, y3, x4, y4 int) bool {
	// Check if line segment (x1,y1)-(x2,y2) intersects with segment (x3,y3)-(x4,y4)
	// Using orientation method
	
	d1 := direction(x3, y3, x4, y4, x1, y1)
	d2 := direction(x3, y3, x4, y4, x2, y2)
	d3 := direction(x1, y1, x2, y2, x3, y3)
	d4 := direction(x1, y1, x2, y2, x4, y4)
	
	if ((d1 > 0 && d2 < 0) || (d1 < 0 && d2 > 0)) &&
		((d3 > 0 && d4 < 0) || (d3 < 0 && d4 > 0)) {
		return true
	}
	
	// Check for collinear cases
	if d1 == 0 && onSegment(x3, y3, x4, y4, x1, y1) {
		return true
	}
	if d2 == 0 && onSegment(x3, y3, x4, y4, x2, y2) {
		return true
	}
	if d3 == 0 && onSegment(x1, y1, x2, y2, x3, y3) {
		return true
	}
	if d4 == 0 && onSegment(x1, y1, x2, y2, x4, y4) {
		return true
	}
	
	return false
}

func direction(x1, y1, x2, y2, x3, y3 int) int {
	// Cross product to determine orientation
	val := (y3-y1)*(x2-x1) - (x3-x1)*(y2-y1)
	if val == 0 {
		return 0 // Collinear
	}
	if val > 0 {
		return 1 // Clockwise
	}
	return -1 // Counterclockwise
}

func onSegment(x1, y1, x2, y2, x3, y3 int) bool {
	// Check if point (x3, y3) lies on segment (x1,y1)-(x2,y2) (assuming collinear)
	return x3 >= min(x1, x2) && x3 <= max(x1, x2) &&
		y3 >= min(y1, y2) && y3 <= max(y1, y2)
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

package main

import (
	"fmt"
	"os"
)

func drawVisualization(lines []string, rect Rectangle, outputFile string) error {
	var points []Point
	for _, line := range lines {
		if p, err := parsePoint(line); err == nil {
			points = append(points, p)
		}
	}

	if len(points) == 0 {
		return fmt.Errorf("no points to visualize")
	}

	// Find bounds
	minX, maxX := points[0].X, points[0].X
	minY, maxY := points[0].Y, points[0].Y
	for _, p := range points {
		if p.X < minX {
			minX = p.X
		}
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}

	// Add padding
	padding := max((maxX-minX)/20, (maxY-minY)/20)
	if padding == 0 {
		padding = 1
	}
	minX -= padding
	maxX += padding
	minY -= padding
	maxY += padding

	width := maxX - minX + 1
	height := maxY - minY + 1

	// Create SVG
	f, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer f.Close()

	// Scale factor for reasonable file size
	scale := 1.0
	maxDim := max(width, height)
	if maxDim > 1000 {
		scale = 1000.0 / float64(maxDim)
	}

	svgWidth := int(float64(width) * scale)
	svgHeight := int(float64(height) * scale)

	// Write SVG header
	fmt.Fprintf(f, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	fmt.Fprintf(f, "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"%d\" height=\"%d\" viewBox=\"%d %d %d %d\">\n",
		svgWidth, svgHeight, minX, minY, width, height)
	fmt.Fprintf(f, "  <rect width=\"100%%\" height=\"100%%\" fill=\"white\"/>\n")

	// Draw polygon
	orderedPoints := orderPointsAsPolygon(points)
	fmt.Fprintf(f, "  <polygon points=\"")
	for i, p := range orderedPoints {
		if i > 0 {
			fmt.Fprintf(f, " ")
		}
		fmt.Fprintf(f, "%d,%d", p.X, p.Y)
	}
	fmt.Fprintf(f, "\" fill=\"lightblue\" stroke=\"blue\" stroke-width=\"%d\" opacity=\"0.5\"/>\n",
		max(1, int(2.0/scale)))

	// Draw rectangle if valid
	if rect.Area() > 0 {
		rectMinX := min(rect.P1.X, rect.P2.X)
		rectMaxX := max(rect.P1.X, rect.P2.X)
		rectMinY := min(rect.P1.Y, rect.P2.Y)
		rectMaxY := max(rect.P1.Y, rect.P2.Y)
		rectWidth := rectMaxX - rectMinX
		rectHeight := rectMaxY - rectMinY

		fmt.Fprintf(f, "  <rect x=\"%d\" y=\"%d\" width=\"%d\" height=\"%d\" fill=\"red\" stroke=\"darkred\" stroke-width=\"%d\" opacity=\"0.3\"/>\n",
			rectMinX, rectMinY, rectWidth, rectHeight, max(1, int(3.0/scale)))

		// Draw rectangle corners
		fmt.Fprintf(f, "  <circle cx=\"%d\" cy=\"%d\" r=\"%d\" fill=\"red\"/>\n",
			rectMinX, rectMinY, max(3, int(5.0/scale)))
		fmt.Fprintf(f, "  <circle cx=\"%d\" cy=\"%d\" r=\"%d\" fill=\"red\"/>\n",
			rectMaxX, rectMaxY, max(3, int(5.0/scale)))

		// Add label with area
		labelX := rectMinX + rectWidth/2
		labelY := rectMinY + rectHeight/2
		fontSize := max(12, int(20.0/scale))
		fmt.Fprintf(f, "  <text x=\"%d\" y=\"%d\" font-size=\"%d\" fill=\"black\" text-anchor=\"middle\" dominant-baseline=\"middle\" font-weight=\"bold\">Area: %d</text>\n",
			labelX, labelY, fontSize, rect.Area())
	}

	// Draw polygon vertices
	for _, p := range points {
		fmt.Fprintf(f, "  <circle cx=\"%d\" cy=\"%d\" r=\"%d\" fill=\"darkblue\"/>\n",
			p.X, p.Y, max(2, int(4.0/scale)))
	}

	// Write SVG footer
	fmt.Fprintf(f, "</svg>\n")

	return nil
}

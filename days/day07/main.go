package main

import (
	"bufio"
	"fmt"
	"os"
)

type Cell rune

const (
	Start Cell = 'S'
	Empty Cell = '.'
	Split Cell = '^'
	Beam  Cell = '|'
)

type Position struct {
	Row int
	Col int
}

type Grid struct {
	Cells  [][]Cell
	Height int
	Width  int
}

func NewGrid(lines []string) *Grid {
	height := len(lines)
	if height == 0 {
		return &Grid{}
	}
	
	width := len(lines[0])
	cells := make([][]Cell, height)
	
	for i, line := range lines {
		cells[i] = make([]Cell, width)
		for j, ch := range line {
			cells[i][j] = Cell(ch)
		}
	}
	
	return &Grid{
		Cells:  cells,
		Height: height,
		Width:  width,
	}
}

func (g *Grid) FindStart() *Position {
	for i := 0; i < g.Height; i++ {
		for j := 0; j < g.Width; j++ {
			if g.Cells[i][j] == Start {
				return &Position{Row: i, Col: j}
			}
		}
	}
	return nil
}

func (g *Grid) Get(pos Position) Cell {
	if pos.Row < 0 || pos.Row >= g.Height || pos.Col < 0 || pos.Col >= g.Width {
		return Empty
	}
	return g.Cells[pos.Row][pos.Col]
}

func (g *Grid) Set(pos Position, cell Cell) {
	if pos.Row >= 0 && pos.Row < g.Height && pos.Col >= 0 && pos.Col < g.Width {
		g.Cells[pos.Row][pos.Col] = cell
	}
}

func (g *Grid) IsInBounds(pos Position) bool {
	return pos.Row >= 0 && pos.Row < g.Height && pos.Col >= 0 && pos.Col < g.Width
}

func (g *Grid) Print() {
	for i := 0; i < g.Height; i++ {
		for j := 0; j < g.Width; j++ {
			fmt.Print(string(g.Cells[i][j]))
		}
		fmt.Println()
	}
}

func (g *Grid) ProcessBeams() int {
	start := g.FindStart()
	if start == nil {
		fmt.Println("No start position found")
		return 0
	}
	
	round := 0
	totalSplits := 0
	activeBeams := []Position{{Row: start.Row, Col: start.Col}}
	
	for len(activeBeams) > 0 {
		round++
		roundSplits := 0
		
		var nextBeams []Position
		
		for _, beam := range activeBeams {
			nextRow := beam.Row + 1
			
			if nextRow >= g.Height {
				continue
			}
			
			nextPos := Position{Row: nextRow, Col: beam.Col}
			targetCell := g.Get(nextPos)
			
			if targetCell == Split {
				roundSplits++
				leftCol := beam.Col - 1
				rightCol := beam.Col + 1
				
				if leftCol >= 0 {
					leftPos := Position{Row: nextRow, Col: leftCol}
					if g.Get(leftPos) == Empty {
						g.Set(leftPos, Beam)
						nextBeams = append(nextBeams, leftPos)
					}
				}
				
				if rightCol < g.Width {
					rightPos := Position{Row: nextRow, Col: rightCol}
					if g.Get(rightPos) == Empty {
						g.Set(rightPos, Beam)
						nextBeams = append(nextBeams, rightPos)
					}
				}
			} else if targetCell == Empty {
				g.Set(nextPos, Beam)
				nextBeams = append(nextBeams, nextPos)
			}
		}
		
		totalSplits += roundSplits
		
		fmt.Printf("\n=== Round %d ===\n", round)
		fmt.Printf("Splits this round: %d\n", roundSplits)
		fmt.Printf("Total splits: %d\n", totalSplits)
		activeBeams = nextBeams
		g.Print()
	}
	
	return round
}

func (g *Grid) CountPaths() int {
	start := g.FindStart()
	if start == nil {
		return 0
	}
	
	memo := make(map[Position]int)
	return g.countPathsFrom(start.Row, start.Col, memo)
}

func (g *Grid) countPathsFrom(row, col int, memo map[Position]int) int {
	pos := Position{Row: row, Col: col}
	
	if cached, exists := memo[pos]; exists {
		return cached
	}
	
	nextRow := row + 1
	
	if nextRow >= g.Height {
		return 1
	}
	
	nextCell := g.Get(Position{Row: nextRow, Col: col})
	
	var result int
	
	if nextCell == Split {
		leftCol := col - 1
		if leftCol >= 0 {
			leftCell := g.Get(Position{Row: nextRow, Col: leftCol})
			if leftCell == Empty || leftCell == Split {
				result += g.countPathsFrom(nextRow, leftCol, memo)
			}
		}
		
		rightCol := col + 1
		if rightCol < g.Width {
			rightCell := g.Get(Position{Row: nextRow, Col: rightCol})
			if rightCell == Empty || rightCell == Split {
				result += g.countPathsFrom(nextRow, rightCol, memo)
			}
		}
	} else if nextCell == Empty || nextCell == Split {
		result = g.countPathsFrom(nextRow, col, memo)
	}
	
	memo[pos] = result
	return result
}

func parseFile(filepath string) (*Grid, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	
	var lines []string
	scanner := bufio.NewScanner(file)
	
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	
	return NewGrid(lines), nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: day7 <mode> <filepath>")
		fmt.Println("  mode: 'splits' or 'paths'")
		os.Exit(1)
	}
	
	mode := os.Args[1]
	filepath := os.Args[2]
	
	if mode != "splits" && mode != "paths" {
		fmt.Fprintf(os.Stderr, "Invalid mode: %s (must be 'splits' or 'paths')\n", mode)
		os.Exit(1)
	}
	
	grid, err := parseFile(filepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing file: %v\n", err)
		os.Exit(1)
	}
	
	if mode == "splits" {
		fmt.Println("=== Initial State ===")
		grid.Print()
		
		rounds := grid.ProcessBeams()
		
		fmt.Printf("\n=== Finished after %d rounds ===\n", rounds)
	} else {
		paths := grid.CountPaths()
		fmt.Printf("Total paths from S to bottom: %d\n", paths)
	}
}

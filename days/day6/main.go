package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Column struct {
	Numbers  []int
	Operator string
}

func (c *Column) Calculate() int {
	if len(c.Numbers) == 0 {
		return 0
	}

	result := c.Numbers[0]
	for i := 1; i < len(c.Numbers); i++ {
		switch c.Operator {
		case "+":
			result += c.Numbers[i]
		case "*":
			result *= c.Numbers[i]
		}
	}
	return result
}

type Grid struct {
	Columns []Column
}

func (g *Grid) CalculateTotal() int {
	total := 0
	for i, col := range g.Columns {
		columnTotal := col.Calculate()
		fmt.Printf("Column %d: %d\n", i+1, columnTotal)
		total += columnTotal
	}
	return total
}

func parseFile(filepath string, mode string) (*Grid, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if mode == "original" {
		return parseOriginalMode(file)
	} else if mode == "aligned" {
		return parseAlignedMode(file)
	}
	
	return nil, fmt.Errorf("invalid mode: %s (must be 'original' or 'aligned')", mode)
}

func parseOriginalMode(file *os.File) (*Grid, error) {
	var rows [][]string
	scanner := bufio.NewScanner(file)
	
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) > 0 {
			rows = append(rows, fields)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(rows) < 2 {
		return nil, fmt.Errorf("file must contain at least 2 rows")
	}

	operatorRow := rows[len(rows)-1]
	numberRows := rows[:len(rows)-1]

	numCols := len(operatorRow)
	grid := &Grid{
		Columns: make([]Column, numCols),
	}

	for colIdx := 0; colIdx < numCols; colIdx++ {
		grid.Columns[colIdx].Operator = operatorRow[colIdx]
		grid.Columns[colIdx].Numbers = make([]int, 0)

		for _, row := range numberRows {
			if colIdx < len(row) {
				num, err := strconv.Atoi(row[colIdx])
				if err != nil {
					return nil, fmt.Errorf("invalid number at column %d: %s", colIdx+1, row[colIdx])
				}
				grid.Columns[colIdx].Numbers = append(grid.Columns[colIdx].Numbers, num)
			}
		}
	}

	return grid, nil
}

func parseAlignedMode(file *os.File) (*Grid, error) {
	var lines []string
	scanner := bufio.NewScanner(file)
	
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) != "" {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(lines) < 2 {
		return nil, fmt.Errorf("file must contain at least 2 rows")
	}

	maxLen := 0
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	for i := range lines {
		if len(lines[i]) < maxLen {
			lines[i] = lines[i] + strings.Repeat(" ", maxLen-len(lines[i]))
		}
	}

	operatorLine := lines[len(lines)-1]
	numberLines := lines[:len(lines)-1]

	var opPositions []int
	operators := make(map[int]string)

	for charIdx := 0; charIdx < len(operatorLine); charIdx++ {
		if operatorLine[charIdx] != ' ' {
			opPositions = append(opPositions, charIdx)
			operators[charIdx] = string(operatorLine[charIdx])
		}
	}

	grid := &Grid{
		Columns: make([]Column, 0),
	}

	for i, opPos := range opPositions {
		start := opPos
		end := maxLen
		if i+1 < len(opPositions) {
			end = opPositions[i+1]
		}

		col := Column{
			Operator: operators[opPos],
			Numbers:  make([]int, 0),
		}

		for pos := end - 1; pos >= start; pos-- {
			var digits []rune
			for _, line := range numberLines {
				if pos < len(line) && line[pos] != ' ' {
					digits = append(digits, rune(line[pos]))
				}
			}

			if len(digits) > 0 {
				numStr := string(digits)
				num, err := strconv.Atoi(numStr)
				if err != nil {
					return nil, fmt.Errorf("invalid number at position %d: %s", pos, numStr)
				}
				col.Numbers = append(col.Numbers, num)
			}
		}

		grid.Columns = append(grid.Columns, col)
	}

	return grid, nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: day6 <mode> <filepath>")
		fmt.Println("  mode: 'original' or 'aligned'")
		os.Exit(1)
	}

	mode := os.Args[1]
	filepath := os.Args[2]
	
	if mode != "original" && mode != "aligned" {
		fmt.Fprintf(os.Stderr, "Invalid mode: %s (must be 'original' or 'aligned')\n", mode)
		os.Exit(1)
	}
	
	grid, err := parseFile(filepath, mode)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing file: %v\n", err)
		os.Exit(1)
	}

	total := grid.CalculateTotal()
	fmt.Printf("Total: %d\n", total)
}

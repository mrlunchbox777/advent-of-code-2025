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

func parseFile(filepath string) (*Grid, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

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

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: day6 <filepath>")
		os.Exit(1)
	}

	filepath := os.Args[1]
	
	grid, err := parseFile(filepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing file: %v\n", err)
		os.Exit(1)
	}

	total := grid.CalculateTotal()
	fmt.Printf("Total: %d\n", total)
}

# Day 4: Grid Position Selection

## Problem

Given a 2D grid containing '.' and '@' symbols, find all '@' positions that have fewer than 4 adjacent '@' symbols.

### Rules

1. Each '@' symbol is checked for adjacent '@' symbols in all 8 directions:
   - left, left-up, up, right-up, right, right-down, down, left-down
2. Select the '@' if it has fewer than 4 adjacent '@' symbols
3. Report positions using 1-based coordinates where [1,1] is the bottom-left corner

## Implementation

The solution uses Go structs and methods:

- `Position`: Represents a coordinate with 1-based indexing (X for column, Y for row)
- `Grid`: Holds the 2D array and provides methods for:
  - `GetCell(pos)`: Get the symbol at a position (with bounds checking)
  - `CountAdjacentAt(pos)`: Count adjacent '@' symbols in all 8 directions
  - `FindSelectedPositions()`: Find all '@' positions with < 4 adjacent '@' symbols

### Coordinate System

- X axis: columns from left (1) to right (width)
- Y axis: rows from bottom (1) to top (height)
- Position [1,1] is the bottom-left corner
- Example for a 3x3 grid:
  ```
  [1,3] [2,3] [3,3]  <- top row (Y=3)
  [1,2] [2,2] [3,2]  <- middle row (Y=2)
  [1,1] [2,1] [3,1]  <- bottom row (Y=1)
  ```

## Usage

```bash
go run main.go processor.go <filepath>
```

Example:
```bash
go run main.go processor.go example-data.txt
```

## Testing

Run tests with:
```bash
go test -v
```

The test suite includes:
- `TestCountAdjacentAt`: Verifies neighbor counting logic
- `TestFindSelectedPositions`: Tests selection logic on a small grid
- `TestExampleData`: Validates the example data produces 13 selections
- `TestGetCell`: Tests coordinate system and bounds checking

## Example Output

For the provided example-data.txt:
```
Selected positions:
[1,1]
[3,1]
[9,1]
[1,3]
[1,6]
[10,6]
[7,8]
[1,9]
[3,10]
[4,10]
[6,10]
[7,10]
[9,10]

Total count: 13
```

## Thoughts On AI Solutions

1. The AI correctly understood the problem requirements and provided a solution that processes the grid and identifies positions based on adjacent counts.

todo

- I did not give the AI the exact instructions from Advent of Code, but rather paraphrased them with my understanding of the problem.
- I did not ask the AI to optimize for performance or efficiency.
- I did not provide any starter code or templates; the AI generated the entire solution from scratch.
- I did not intervene in the coding process except to provide prompts and clarifications as needed.
- I only updated this section of the README for this day.

### Initial Prompts

#### Part 1

> build a go app in the day4 folder. It should accept a single parameter that is a path to a text file. The file will contain a 2 dimentional array of symbols (either a '.' or a '@'). It should track the array and the entries within it using structs and use methods on those structs to perform the computations. It should look at each '@' and select it if there are fewer than 4 '@' in the adjacent 8 positions (left, left-up, up, right-up, right, right-down, down, and left-down). It should list those that it selected by the coordinates, consider the bottom left-most position's coordinates [1,1]. Finally, it should print the count of selections. Add test for the app as well. For reference processing, the example-data.txt should total 13.

#### Part 2

> 


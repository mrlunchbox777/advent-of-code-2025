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
go run main.go processor.go <filepath> <mode>
```

Where `<mode>` is either:
- `initial`: Single pass - find all '@' positions with < 4 adjacent '@' symbols
- `completion`: Iterative passes - repeatedly find and remove '@' positions until none remain

### Examples

Initial mode (single pass):
```bash
go run main.go processor.go example-data.txt initial
```

Completion mode (multiple rounds):
```bash
go run main.go processor.go example-data.txt completion
```

## Modes

### Initial Mode

Performs a single pass through the grid and reports all '@' positions with fewer than 4 adjacent '@' symbols.

### Completion Mode

Iteratively processes the grid until no more '@' symbols can be removed:
1. Find all '@' positions with < 4 adjacent '@' symbols
2. Replace those positions with '.'
3. Repeat until no positions are selected

For each round, the mode prints:
- Round number
- Selected positions for that round
- Count for the round
- Running total

At the end, it prints:
- All selected positions across all rounds
- Total number of rounds
- Final total count

## Testing

Run tests with:
```bash
go test -v
```

The test suite includes:
- `TestCountAdjacentAt`: Verifies neighbor counting logic
- `TestFindSelectedPositions`: Tests selection logic on a small grid
- `TestExampleData`: Validates the example data produces 13 selections (initial mode)
- `TestGetCell`: Tests coordinate system and bounds checking
- `TestReplacePositions`: Tests position replacement for completion mode
- `TestCompletionMode`: Validates completion mode produces 43 total selections across 9 rounds

## Example Output

### Initial Mode
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

### Completion Mode
For the provided example-data.txt (9 rounds, 43 total):
```
Round 1:
Selected positions:
[1,1]
[3,1]
...
Total from round: 13
Running total: 13

Round 2:
Selected positions:
[2,2]
...
Total from round: 12
Running total: 25

...

=== Final Summary ===
All selected positions:
[1,1]
[3,1]
...
[4,7]

Number of rounds: 9
Final total: 43
```

## Thoughts On AI Solutions

1. The AI correctly understood the problem requirements and provided a solution that processes the grid and identifies positions based on adjacent counts. It did have the corrupted code issue again during refactoring attempts, but was able to catch and fix it.
2. I ran the solution against the puzzle input and it produced the correct results.
3. For part 2, the AI understood the new requirement for iterative processing and implemented the logic to repeatedly find and remove positions until no more could be selected.
4. I ran the solution against part 2 and it produced the correct results.

Each mode was validated with tests, including the example data and edge cases. The AI demonstrated an ability to adapt its solution based on the requirements and provided a working implementation for both modes. This is becoming significantly easier with each iteration, I think it's because I'm getting better at initial prompts and guiding the AI through the process. There are far fewer surprises. There were still some issues with corrupted code during refactoring, but the AI was able to identify and fix those issues itself. This process only required the initial prompts this time. The AI did decide to add a "Result" comment at the end of this README which is a new behavior I haven't seen before, and I didn't ask for it, but I let it be.

- I did not give the AI the exact instructions from Advent of Code, but rather paraphrased them with my understanding of the problem.
- I did not ask the AI to optimize for performance or efficiency.
- I did not provide any starter code or templates; the AI generated the entire solution from scratch.
- I did not intervene in the coding process except to provide prompts and clarifications as needed.
- I only updated this section of the README for this day.

### Initial Prompts

#### Part 1

> build a go app in the day4 folder. It should accept a single parameter that is a path to a text file. The file will contain a 2 dimentional array of symbols (either a '.' or a '@'). It should track the array and the entries within it using structs and use methods on those structs to perform the computations. It should look at each '@' and select it if there are fewer than 4 '@' in the adjacent 8 positions (left, left-up, up, right-up, right, right-down, down, and left-down). It should list those that it selected by the coordinates, consider the bottom left-most position's coordinates [1,1]. Finally, it should print the count of selections. Add test for the app as well. For reference processing, the example-data.txt should total 13.

Result: The AI successfully created the app with the correct logic and tests. The example data produced exactly 13 selections as expected.

#### Part 2

> add a second mandatory flag that will determine if it will check for an initial pass or attempt as many passes as it can. The initial pass is the current state of the app. The completion mode will run the initial pass, then create a new grid with those positions replaced with '.', then it will attempt to run again. Once it can't replace any more '@'s it is complete. As it goes it will print the "round" it is on, the positions selections for the round, the total from the round, and the running total of all rounds. At the end it will print all positions, the number of rounds, and the final total. For reference the example-data.txt should produce a final total of 43.

Result: The AI correctly implemented both modes with proper iteration logic. The completion mode produces 43 total selections across 9 rounds as expected. The `ReplacePositions` method creates new grids preserving immutability, and comprehensive tests validate both modes.


# Day 7 - Beam Splitter

A Go application that simulates beam movement through a grid with splitting behavior.

## Description

This application reads a text file containing a 2D grid with special symbols and simulates beams moving downward through the grid. When beams encounter split symbols, they divide into left and right beams.

## Symbols

- `S` - Start position (where the beam originates)
- `.` - Empty space
- `^` - Split (causes beam to split left and right)
- `|` - Beam (created during processing)

## Rules

1. Beams start at the `S` position
2. Each round, active beams attempt to move down one row
3. When a beam moves to an empty cell (`.`), it becomes a beam (`|`) and continues
4. When a beam would move onto a split (`^`), instead it creates two new beams:
   - One beam in the column to the left of the split
   - One beam in the column to the right of the split
   - Both beams are created at the same row as the split
5. Already-created beams remain in place
6. Processing continues until all active beams would move outside the grid

## Usage

Build the application:

```bash
go build -o day7
```

Run with a data file:

```bash
./day7 <path-to-data-file>
```

Example:

```bash
./day7 example-data-1.txt
```

## Output

### Splits Mode

In splits mode, the application prints:

- The initial state of the grid
- Each round showing:
  - The number of splits that occurred in that round
  - The running total of splits across all rounds
  - The grid after beam placement
- A completion message with the total number of rounds

Example output:

```
$ ./day7 splits example-data-1.txt
=== Initial State ===
.......S.......
...............
.......^.......
...............

=== Round 1 ===
Splits this round: 0
Total splits: 0
.......S.......
.......|.......
.......^.......
...............

=== Round 2 ===
Splits this round: 1
Total splits: 1
.......S.......
.......|.......
......|^|......
...............

=== Finished after 16 rounds ===
```

### Paths Mode

In paths mode, the application calculates and prints the total number of unique paths a beam can take from the start position to the bottom of the grid, considering all possible splits.

Example output:

```
$ ./day7 paths example-data-1.txt
Total paths from S to bottom: 40
```

## Implementation

The application uses:

- `Cell` type: Represents each cell type (Start, Empty, Split, Beam)
- `Position` struct: Stores row and column coordinates
- `Grid` struct: Stores the 2D array of cells
  - `FindStart()` method: Locates the start position
  - `Get/Set()` methods: Access and modify cells
  - `ProcessBeams()` method: Executes the beam simulation (splits mode)
  - `CountPaths()` method: Recursively counts all paths (paths mode)
  - `Print()` method: Displays the current grid state

## Algorithms

### Splits Mode Algorithm

1. Find the start position `S`
2. Initialize active beams list with the start position
3. For each round:
   - For each active beam, calculate the next position (one row down)
   - If the next position is a split (`^`):
     - Count the split
     - Create beams at left and right columns in the same row
     - Add those positions to the next round's active beams
   - If the next position is empty (`.`):
     - Place a beam (`|`) there
     - Add that position to the next round's active beams
   - If the next position is out of bounds or already has a beam, skip it
4. Continue until no active beams remain
5. Display split counts and final grid

### Paths Mode Algorithm

1. Find the start position `S`
2. Recursively count paths from that position:
   - If the next row is beyond the grid, return 1 (reached bottom)
   - If the next cell is a split (`^`):
     - Recursively count paths from left column (if in bounds)
     - Recursively count paths from right column (if in bounds)
     - Return the sum of both paths
   - If the next cell is empty (`.`) or another split:
     - Recursively count paths continuing straight down
3. Return the total count of all paths

## Testing

Run tests:

```bash
go test -v
```

The test suite includes:

- Grid creation and initialization tests
- Start position finding tests
- Simple beam movement tests
- Beam splitting behavior tests
- Multiple splits and edge cases
- Validation against example data files

## Thoughts On AI Solutions

1. The AI understood the problem, found a solution, and validated it against the example data fairly quickly after a few iterations. I had to retry the initial prompt because I provided the wrong example data initially (forgot to save).
2. I realized I forgot to ask it to count the number of times it split, and asked it to do that. It did with ease.
3. I tested the solution against the puzzle input data and it worked perfectly the first time.
4. I asked it to implement part 2, it did so correctly, but it was too slow for the puzzle input, so I had to ask it to optimize it.

TODO: add summary of thoughts

- I did not give the AI the exact instructions from Advent of Code, but rather paraphrased them with my understanding of the problem.
- I did not ask the AI to optimize for performance or efficiency.
- I did not provide any starter code or templates; the AI generated the entire solution from scratch.
- I did not intervene in the coding process except to provide prompts and clarifications as needed.
- I only updated this section of the README for this day.

### Initial Prompts

#### Part 1

> build a go app in the day7 folder. It should accept a single parameter that is a path to a text file. The file will contain a 2 dimentional array of symbols 'S' (start), '.' (empty), '^' (split), and eventually adding '|' (beam). It should track the array and the entries within it using structs and use methods on those structs to perform the computations. It should process the file in rounds. Each round it should attempt to move the beam(s) down one move (each), by leaving those already created where they are and creating a new beam one row below all beam(s) created last round; the first move is creating the initial beam one row below 'S'. When it attempts to move onto a split, instead of moving down one, it should create a new beams in the columns before and after the split and in the same row. It should output each round as it processes and announce when it finishes. It finishes when every new move would be off the grid of the input (it shouldn't ever add any beams outside the grid of the input). Add test and a basic readme for the app as well. For reference processing, the @days/day7/example-data-1.txt should result in @days/day7/example-data-2.txt after processing is complete.

#### Part 2

Add a second mandatory parameter to the app that either selects the counts splits mode, existing logic, or counts paths mode. In this mode it should count the number of paths a beam could take from S to the bottom of the grid, splitting at each split. It should output that number when complete instead of the grid processing output. For reference processing, the @days/day7/example-data-1.txt should result in 40 paths.

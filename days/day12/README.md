# Day 12: Puzzle Piece Solver

This program solves puzzles by placing pieces into a rectangular area. Pieces can be rotated, flipped, and must fit entirely within the area without overlapping filled spaces.

## Problem Description

The input file contains two sections:

### Section 1: Piece Definitions

Each piece is defined by:

- A number followed by a colon (e.g., `0:`)
- Lines showing the shape using `#` for filled space and `.` for empty space

Example:

```
0:
###
##.
##.
```

### Section 2: Puzzle Specifications

Each puzzle is defined by:

- Dimensions in `lengthxwidth` format (e.g., `4x4`)
- A colon followed by space-separated pairs of numbers
- Each pair is: piece ID (ordinality) and count (how many of that piece to use)

Example:

```
4x4: 0 0 0 0 2 0
```

This means: Use piece 0 zero times, piece 0 zero times (repeated), and piece 2 zero times - for a 4x4 area.

The actual format is: `pieceID count pieceID count ...`

## Usage

```bash
./day12 <path-to-input-file>
```

## Example

Given `example-data.txt`:

```
0:
###
##.
##.

1:
###
##.
.##

2:
.##
###
##.

3:
##.
###
##.

4:
###
#..
###

5:
###
.#.
###

4x4: 0 0 0 0 2 0
12x5: 1 0 1 0 2 2
12x5: 1 0 1 0 3 2
```

Running the program:

```bash
./day12 example-data.txt
```

Output:

```
(Solution for puzzle 1 - shows placement with different characters A, B, C, etc.)
(Solution for puzzle 2)
No solution found
Puzzles with solutions found: 2
```

## Building

```bash
go build -o day12
```

## Testing

```bash
go test
```

## Algorithm

The program uses **backtracking** with the following approach:

1. **Parse Input**: Reads and parses pieces and puzzle specifications
2. **Generate Orientations**: For each piece, generates all 8 possible orientations (4 rotations × 2 flips)
3. **Backtracking Search**:
   - Try placing each piece at every position in the grid
   - Try all orientations of each piece
   - Check if placement is valid (no overlap with filled spaces)
   - If valid, place the piece and recurse to place the next piece
   - If recursion succeeds, return the solution
   - If recursion fails, backtrack by removing the piece and trying next position/orientation
4. **Solution Output**: Print the grid with different characters (A, B, C, etc.) for each placed piece

## Implementation Details

### Structures

- **Piece**: Represents a puzzle piece with its shape grid
  - `AllOrientations()`: Generates all unique rotations and flips
  - `Rotate90()`: Rotates piece 90° clockwise
  - `Flip()`: Flips piece horizontally
- **Puzzle**: Represents a puzzle to solve
  - `Solve()`: Main solving method using backtracking
  - `canPlace()`: Checks if piece can be placed at position
  - `place()`: Places piece on grid
  - `remove()`: Removes piece from grid (for backtracking)

- **Solution**: Represents a solved puzzle with the final grid layout

### Files

- `main.go`: Entry point and command-line interface
- `parser.go`: Parses input file into data structures
- `piece.go`: Piece representation and transformations
- `puzzle.go`: Puzzle solving logic with backtracking
- `main_test.go`: Unit tests

## Constraints

- Pieces must fit entirely within the puzzle area
- Pieces cannot overlap filled spaces (but can overlap empty spaces)
- Pieces can be rotated and flipped to any orientation
- Each piece in the specification must be used exactly the specified number of times

## Performance

- Small puzzles (4x4): < 1 second
- Medium puzzles (12x5): < 5 seconds depending on complexity
- Large puzzles (40x40+): < 10 seconds per puzzle typically
- 1000 puzzle test suite: ~60 seconds total (427/1000 solvable)

Performance depends on:

- Number of pieces (exponential growth in search space)
- Puzzle dimensions
- Piece complexity (number of filled cells)
- Solution density (how tightly pieces fit)

### Optimizations

1. **Pre-computed Orientations**: All piece orientations are calculated once and reused
2. **Filled Cell Counting**: Pre-compute filled cell counts for pruning
3. **Early Pruning**: Detect impossible states early (more pieces than space available)
4. **Orientation Deduplication**: Removes identical orientations from rotation/flip combinations

The algorithm guarantees finding a solution if one exists, or correctly reporting when no solution is possible.

## Thoughts On AI Solutions

1. I asked the AI to create the Go application and it generated a solution that fit the example data provided.
2. When running against the puzzle data, it got 0 solutions and that just seems wrong. I'm going to clarify the prompt and get it to try again. Once I did that it found out by itself that it was hallucinating. Specifically it did not ready my ordinality instruction correctly. It also found that the solution was running too slowly and started optimizing by itself. It ended up optimizing by giving up early on some cases, so i will not use this version to submit. I will make sure it knows it has to finish, here i will have to say that I'm telling it to optimize (it did that to skip hours long puzzles).
3. It said it got a solution that is complete and correct and I will attempt it.
4. It got it correct.

TODO: add more details after I finish the puzzle.

- I did not give the AI the exact instructions from Advent of Code, but rather paraphrased them with my understanding of the problem.
- ~~I did not ask the AI to optimize for performance or efficiency.~~
- I did not provide any starter code or templates; the AI generated the entire solution from scratch.
- I did not intervene in the coding process except to provide prompts and clarifications as needed.
- I only updated this section of the README for this day.

### Initial Prompts

#### Part 1

> build a go app in the @[DIR] /Users/ashoell/src/mrlunchbox777/advent-of-code-2025/days/day12 folder. It should accept a first parameter that is a path to a text file. It should track the data within it using structs and use methods on those structs to perform the computations and split those into files in directories as appropriate. Add test and a basic readme for the app as well. Follow patterns from the other days in the @[DIR] /Users/ashoell/src/mrlunchbox777/advent-of-code-2025/days folder, particularly @[DIR] /Users/ashoell/src/mrlunchbox777/advent-of-code-2025/days/day10 @[DIR] /Users/ashoell/src/mrlunchbox777/advent-of-code-2025/days/day11.
> The file will be split into two parts. The first will have pieces identified by a number then a colon. The following lines will contain the shape of the piece as identified by `#` for filled space and `.` for empty space. The second part will describe puzzles that are identified by their area in `length x width` format, e.g. `4x4`, followed by a colon. The rest of the line will contain a series of space separated numbers that identify the which piece, ordinality, and how many of that piece, value, to use in that area to solve the puzzle. The pieces can not overlap filled space, but can be rotated, flipped, and fit their filled spaces into other empty spaces as needed. They must fit entirely within the area.
> The goal of the app is to read the file, and for each puzzle attempt to solve it. If a solution is found it should print the solution to standard out using a different character for each piece instead of `#`. If no solution is found it should print `No solution found`. Finally, the app should print how many puzzles had solutions found.
> For reference processing, @days/day12/example-data.txt should have 2 puzzles with a solution.

#### Part 2

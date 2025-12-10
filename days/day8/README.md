# Day 8 - Coordinate Grouping

<!--toc:start-->

- [Day 8 - Coordinate Grouping](#day-8-coordinate-grouping)
  - [Description](#description)
  - [Input Format](#input-format)
  - [Usage](#usage)
  - [Output](#output)
  - [Implementation](#implementation)
  - [Algorithm](#algorithm)
  - [Testing](#testing)
  - [Thoughts On AI Solutions](#thoughts-on-ai-solutions) - [Initial Prompts](#initial-prompts) - [Part 1](#part-1) - [Part 2](#part-2)
  <!--toc:end-->

A Go application that groups 3D coordinates by iteratively connecting the closest pairs.

## Description

This application reads a file containing 3D coordinates and groups them by repeatedly finding and connecting the two closest coordinates that aren't already in the same group. After each connection, it displays the top 5 largest groups. At the end, it calculates the product of the three largest group sizes.

## Input Format

The input file should contain one 3D coordinate per line in the format:

```
X,Y,Z
```

Example:

```
162,817,812
57,618,57
906,360,560
```

## Usage

Build the application:

```bash
go build -o day8
```

Run with a data file:

```bash
./day8 <filepath> [max_rounds]
```

Parameters:

- `filepath` - Path to the input file (required)
- `max_rounds` - Maximum number of rounds to run (optional, default: 1000)

Examples:

```bash
./day8 example-data.txt
./day8 example-data.txt 10
./day8 puzzle-input.txt 5000
```

## Output

The application prints:

- Initial coordinate count
- For each round:
  - The two coordinates being connected
  - The distance between them
  - The sizes of the top 5 largest groups
- Final results:
  - The three largest group sizes
  - The product of those sizes

Example output:

```
Loaded 20 coordinates
Running up to 10 rounds

Round 1: Connected (162,817,812) and (425,690,689) - Distance: 316.90
  Top 5 groups: 2, 1, 1, 1, 1
Round 2: Connected (162,817,812) and (431,825,988) - Distance: 321.56
  Top 5 groups: 3, 1, 1, 1, 1
...
Round 9: Connected (346,949,466) and (425,690,689) - Distance: 350.79
  Top 5 groups: 4, 3, 2, 2, 2
Round 10: Connected (906,360,560) and (984,92,344) - Distance: 352.94
  Top 5 groups: 5, 4, 2, 2, 1

=== Final Results ===
Top 3 largest groups:
  Group 1: 5 members
  Group 2: 4 members
  Group 3: 2 members

Product of top 3 group sizes: 40
```

## Implementation

The application uses:

- `Coordinate` struct: Stores X, Y, Z coordinates and an ID
  - `Distance()` method: Calculates Euclidean distance to another coordinate
- `UnionFind` struct: Efficient data structure for tracking connected components
  - `Find()` method: Finds the root of a component (with path compression)
  - `Union()` method: Merges two components (with union by rank)
- `CoordinateSet` struct: Manages the collection of coordinates and their groupings
  - `FindClosestPair()` method: Finds the two closest coordinates not yet in the same group
  - `Connect()` method: Connects two coordinates into the same group
  - `GetGroups()` method: Returns all groups
  - `GetTopGroups()` method: Returns the N largest groups sorted by size

## Algorithm

1. Load all coordinates from the input file
2. Initialize a Union-Find data structure with each coordinate in its own group
3. Track all direct connections between coordinates
4. For each round (up to max_rounds):
   - Find the pair of coordinates with minimum distance that don't already have a direct connection
   - Create a direct connection between them
   - Update the group membership using Union-Find
   - Display the connection and current top 5 groups
   - If all possible pairs are connected, stop early
5. Calculate and display the product of the three largest group sizes

**Key Insight**: The algorithm connects the closest pair of coordinates that don't have a direct connection yet. This means coordinates in the same group (connected through other coordinates) can still be directly connected if they don't already have a direct edge between them.

**Time Complexity**: O(R × N²) where R is rounds and N is number of coordinates

- Each round checks all pairs: O(N²)
- Union-Find operations are nearly O(1) with path compression and union by rank

**Space Complexity**: O(N) for storing coordinates and Union-Find structure

## Testing

Run tests:

```bash
go test -v
```

The test suite includes:

- Coordinate distance calculation tests
- Union-Find operations tests
- Closest pair finding tests
- Group tracking tests
- File parsing tests
- Validation against example-data.txt (expected product: 40)

## Thoughts On AI Solutions

1. The AI found a solution, but it complained that the example data and expected outcome was wrong. When looking at it's output it looked like it got the right solution in round 9 instead of 10, so I'm guessing there is a one off error either in my instructions or it's implementation.
2. The one off error was due to my misunderstanding of the problem. I missed the explanation that coordinates in the same group could still be directly connected if they didn't already have a direct edge between them. Once I clarified that, the AI adjusted the implementation and it worked perfectly.
3. I attempted the puzzle input and it worked perfectly the first time.

TODO: add summary of thoughts

- I did not give the AI the exact instructions from Advent of Code, but rather paraphrased them with my understanding of the problem.
- I did not ask the AI to optimize for performance or efficiency.
- I did not provide any starter code or templates; the AI generated the entire solution from scratch.
- I did not intervene in the coding process except to provide prompts and clarifications as needed.
- I only updated this section of the README for this day.

### Initial Prompts

#### Part 1

> build a go app in the day8 folder. It should accept a first parameter that is a path to a text file. The file will contain a set of 3D coordinates per line. It should track the array and the entries within it using structs and use methods on those structs to perform the computations. It should process the file in rounds. Each round it should find the 2 closest coordinates and connect them, it should print the connection it just made and the biggest 5 groups of connected coordinates. It should repeat this up to x (default 1000, second parameter) times. At the end it should print the 3 largest groups, and the product of their sizes. Add test and a basic readme for the app as well. For reference processing, @days/day8/example-data.txt ran for 10 rounds should result in 5x4x2=40.

#### Part 2

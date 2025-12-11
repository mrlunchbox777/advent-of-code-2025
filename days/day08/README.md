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

### Completion Mode

In completion mode, the application runs until all coordinates are connected into a single group, then reports the connection that achieved this.

Example output for completion mode:

```
$ ./day8 completion example-data.txt
Loaded 20 coordinates
Running until all coordinates are in a single group

Round 1: Connected (162,817,812) and (425,690,689) - Distance: 316.90
  Top 5 groups: 2, 1, 1, 1, 1
...
Round 29: Connected (216,146,977) and (117,168,530) - Distance: 458.36
  Top 5 groups: 20
*** All coordinates now in a single group! ***
...

=== Final Results ===
Single group achieved at round 29
Completion connection: (216,146,977) and (117,168,530)
Product of X coordinates: 216 × 117 = 25272
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

### Completion Mode Algorithm

1. Load all coordinates from the input file
2. Initialize tracking structures (Union-Find and connection map)
3. For each round:
   - Find and connect the closest unconnected pair
   - Check if all coordinates are now in a single group
   - If yes, record this as the completion connection and note the round
   - Continue connecting remaining pairs
4. Display the completion connection and calculate the product of its X coordinates

**Time Complexity**:

- **Grouping mode**: O(R × N²) where R is rounds and N is number of coordinates
  - Each round checks all pairs: O(N²)
- **Completion mode**: O(N² log N) optimized using a min-heap
  - All edges pre-computed and sorted once: O(N² log N²)
  - Each connection is O(log N) to pop from heap
  - Union-Find operations are nearly O(1) with path compression and union by rank

**Space Complexity**:

- **Grouping mode**: O(N) for storing coordinates and Union-Find structure
- **Completion mode**: O(N²) for storing all edges in the min-heap

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
4. When I asked for the completion mode, it initially misunderstood and thought I wanted to connect all coordinates into a single group as fast as possible, rather than running until all coordinates were connected. By giving the expected output for the example data with the initial prompt, it quickly adjusted without intervention and produced the correct solution.
5. I tried running the completion mode on the puzzle input, but it took too long. I suspect this is due to the exponential growth of connections as more coordinates are connected. I asked the AI to optimize for performance.
6. I had to reset the session of the AI, unrelated to the project. I tried running the AI's solution to optimize completion, but now it can't find the correct solution to the example data. I gave it the expected output again and asked it to try again.
7. I attempted to run the AI solution on the puzzle input again, and it finished quickly, but got a result that was too low. I let it know and asked it to try again.
8. I attempted one more time and it got the correct answer.

Today I was reminded how important it is to understand the problem fully before asking an AI to solve it. Once I clarified part 1, the AI easily got the correct solution. In part 2, I saw the same issue, but it self-corrected because I gave it the example data and expected output. Optimization was a part of part 2, and when I asked for it and I think it solved the optimization correctly, but I had to reset the session and it lost the context and it's solution was broken. After giving it a correct example again, it was able to solve it. Again, test and validation with known data is key to getting the right solution.

- I did not give the AI the exact instructions from Advent of Code, but rather paraphrased them with my understanding of the problem.
- I did not ask the AI to optimize for performance or efficiency.
- I did not provide any starter code or templates; the AI generated the entire solution from scratch.
- I did not intervene in the coding process except to provide prompts and clarifications as needed.
- I only updated this section of the README for this day.

### Initial Prompts

#### Part 1

> build a go app in the day8 folder. It should accept a first parameter that is a path to a text file. The file will contain a set of 3D coordinates per line. It should track the array and the entries within it using structs and use methods on those structs to perform the computations. It should process the file in rounds. Each round it should find the 2 closest coordinates and connect them, it should print the connection it just made and the biggest 5 groups of connected coordinates. It should repeat this up to x (default 1000, second parameter) times. At the end it should print the 3 largest groups, and the product of their sizes. Add test and a basic readme for the app as well. For reference processing, @days/day8/example-data.txt ran for 10 rounds should result in 5x4x2=40.

#### Part 2

> Add a second mandatory parameter to the app that either selects the default grouping mode, existing logic, or a completion mode. In this mode it should run until all coordinates are connected. Once complete it should multiply the x coordinates of the coordinates in the final connection. For reference processing, @days/day8/example-data.txt should result in 216x117=25272.

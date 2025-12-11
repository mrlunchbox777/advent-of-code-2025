# Day 10: Toggle Machine Solver

This program solves the toggle machine problem where you need to find the minimum number of option selections to reach a desired state from an initial state (all toggles off).

## Problem Description

Each line in the input file represents a toggle machine with:

- **Desired State** (toggle mode): Wrapped in square brackets `[...]` containing '.' (off) and '#' (on)
- **Options**: Sets of comma-separated numbers in parentheses `(...)` that toggle positions when selected
- **Target Counts** (counter mode): Wrapped in curly braces `{...}` containing target toggle counts for each position

### Toggle Mode

All positions start in the OFF state. The goal is to find the minimum number of option selections needed to reach the desired state.

### Counter Mode

All positions start at count 0. Each time an option is selected, it increments the count at the specified positions. The goal is to find the minimum number of option selections to reach the exact target counts.

## Usage

```bash
./day10 <path-to-input-file> <mode>
```

Where `<mode>` is either:

- `toggle` - Find minimum selections to reach desired toggle state (square brackets)
- `counter` - Find minimum selections to reach target toggle counts (curly braces)

## Example

Given `example-data.txt`:

```
[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}
[...#.] (0,2,3,4) (2,3) (0,4) (0,1,2) (1,2,3,4) {7,5,12,7,2}
[.###.#] (0,1,2,3,4) (0,3,4) (0,1,2,4,5) (1,2) {10,11,11,5,10,5}
```

Running the program in toggle mode:

```bash
./day10 example-data.txt toggle
```

Output:

```
Line 1: 2 selections - options [2 4]
Line 2: 3 selections - options [3 4 5]
Line 3: 2 selections - options [2 3]
Total selections: 7
```

Running the program in counter mode:

```bash
./day10 example-data.txt counter
```

Output:

```
Line 1: 10 selections - options [1 2 2 2 2 2 4 5 5 5]
Line 2: 12 selections - options [1 1 2 2 2 2 2 4 4 4 4 4]
Line 3: 11 selections - options [1 1 1 1 1 3 3 3 3 3 4]
Total selections: 33
```

## Building

```bash
go build -o day10
```

## Testing

```bash
go test
```

## Algorithm

The program uses Breadth-First Search (BFS) to find the shortest path from the initial state to the target state. This guarantees finding the minimum number of selections needed.

In toggle mode, BFS explores all possible toggle states until finding one matching the desired state. In counter mode, BFS explores all possible counter combinations until finding one matching the target counts.

### Performance Note

Counter mode with large target values (50+) has exponential search space complexity and may be slow on some inputs. The implementation uses:

- **A\* search** with Manhattan distance heuristic to prioritize promising states
- **Pruning** to skip states that exceed target counts
- **Pointer-based path reconstruction** to reduce memory allocation during search
- **State limits** (2M states) to prevent unbounded exploration

The example data solves quickly (~500ms), but some puzzle inputs with very large target counts (80+) may hit the state limit. Toggle mode is very fast (< 50ms) even for 200-line inputs.

## Implementation Details

- **Machine struct**: Holds the desired state, target counts, and available options
- **ParseMachine**: Parses input lines into Machine structs, extracting both toggle states and target counts
- **ApplyOption**: Applies an option to toggle specified positions (toggle mode)
- **ApplyOptionCounter**: Applies an option to increment counters at specified positions (counter mode)
- **Solve**: Uses BFS to find the minimum path to the desired state (toggle mode)
- **SolveCounter**: Uses BFS to find the minimum path to the target counts (counter mode)
- **ProcessLines**: Processes all lines in the specified mode and accumulates results

## Thoughts On AI Solutions

1. The AI understood the requirements and generated a complete solution. I actually made a typo and said that line 1 should be 3 selections when it is actually 2, but the AI still produced correct logic. It did have to be interrupted and restarted once.
2. I tested the solution against the puzzle input and it worked correctly.
3. I asked it to solve part 2 and it seemed to find a working answer, but it was too slow when tested with the puzzle input. I asked the AI to optimize it.
4. It immediately used both time and timeout, which is better than before. It iterated through greedy, bfs, memory caching, limiting that memory, back to greedy, more memory management, etc. Finally, it gave up and said that the problem is inherently exponential and may not be solvable for large inputs within reasonable time/memory constraints. The answer it gave was too high, so I asked it to try again.
5. It attempted heuristic manhattan distance, heap based A\* approach, dropped the queue, added memorization, changed text encoding strategy, pooling and different data structures, debugged problem lines and finally gave up again. It was still too slow, but it left with possible strategies for further optimization.
6. It tried bfs backwards greedy, beam search for faster queue, bfs without sorting, and gave up much faster this time. It's still too slow.

TODO: summary of thoughts

- I did not give the AI the exact instructions from Advent of Code, but rather paraphrased them with my understanding of the problem.
- ~~I did not ask the AI to optimize for performance or efficiency.~~
- I did not provide any starter code or templates; the AI generated the entire solution from scratch.
- I did not intervene in the coding process except to provide prompts and clarifications as needed.
- I only updated this section of the README for this day.

### Initial Prompts

#### Part 1

> build a go app in the day10 folder. It should accept a first parameter that is a path to a text file. It should track the data within it using structs and use methods on those structs to perform the computations and split those into files in directories as appropriate. Add test and a basic readme for the app as well. Follow patterns from the other days in the @days folder.
> The file will contain lines that represet a machine. The first section is wrapped in square brackets that contain the desired state of toggles ('.' off, '#' on), all positions start as off. The next section is options for changing the state of the toggles, each option is represented as a set of comma seperated numbers in parentheses. When selected those numbers will toggle the corresponding positions in the current state. The last section can be ignored for now, and is represented as a set of comma seperated numbers in curly braces.
> The goal of the app is to read the file, and determine the minimum number of selections from the options that need to be made to reach the desired state from the initial state of all off for each line. It should print the line number, the number of selections made for that line, and which selections were made. Finally, it should print the total number of selections made for all lines in the file.
> For reference processing, @days/day10/example-data.txt line 1 should be 3 selections, line 2 is also 3 selections, line 3 is 2, the total for the file is 7.

#### Part 2

> Add another mandatory parameter that chooses between toggle (current) mode and counter mode. In counter mode the section with square brackets is ignored and the curly brace section is used instead, each position represents the target number of toggles for each position. The goal is to reach exactly the number in the corresponding position in the curly brace section. The rest of the logic remains the same, find the minimum number of selections from the options to reach that goal. For reference processing, @days/day10/example-data.txt line 1 should be 10 selections, line 2 is also 12 selections, line 3 is 11, the total for the file is 33.

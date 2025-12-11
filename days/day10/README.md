# Day 10: Toggle Machine Solver

This program solves the toggle machine problem where you need to find the minimum number of option selections to reach a desired state from an initial state (all toggles off).

## Problem Description

Each line in the input file represents a toggle machine with:

- **Desired State**: Wrapped in square brackets `[...]` containing '.' (off) and '#' (on)
- **Options**: Sets of comma-separated numbers in parentheses `(...)` that toggle positions when selected
- **Metadata**: Curly braces `{...}` containing data that is currently ignored

All positions start in the OFF state. The goal is to find the minimum number of option selections needed to reach the desired state.

## Usage

```bash
./day10 <path-to-input-file>
```

## Example

Given `example-data.txt`:

```
[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}
[...#.] (0,2,3,4) (2,3) (0,4) (0,1,2) (1,2,3,4) {7,5,12,7,2}
[.###.#] (0,1,2,3,4) (0,3,4) (0,1,2,4,5) (1,2) {10,11,11,5,10,5}
```

Running the program:

```bash
./day10 example-data.txt
```

Output:

```
Line 1: 2 selections - options [2 4]
Line 2: 3 selections - options [3 4 5]
Line 3: 2 selections - options [2 3]
Total selections: 7
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

The program uses Breadth-First Search (BFS) to find the shortest path from the initial state (all off) to the desired state. This guarantees finding the minimum number of selections needed.

## Implementation Details

- **Machine struct**: Holds the desired state and available options
- **ParseMachine**: Parses input lines into Machine structs
- **ApplyOption**: Applies an option to toggle specified positions
- **Solve**: Uses BFS to find the minimum path to the desired state
- **ProcessLines**: Processes all lines and accumulates results

## Thoughts On AI Solutions

1. The AI understood the requirements and generated a complete solution. I actually made a typo and said that line 1 should be 3 selections when it is actually 2, but the AI still produced correct logic.
2. I tested the solution against the puzzle input and it worked correctly.

TODO: summary of thoughts

- I did not give the AI the exact instructions from Advent of Code, but rather paraphrased them with my understanding of the problem.
- I did not ask the AI to optimize for performance or efficiency.
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

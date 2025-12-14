# Day 11: Path Finder

This program finds all unique paths from a starting node to an ending node in a directed graph.

## Problem Description

The input file contains lines representing nodes in a directed graph. Each line has the format:

```
node_name: connected_node1 connected_node2 ...
```

- **Node Name**: The name of the node
- **Connections**: Space-separated list of nodes that this node connects to (mono-directional)

The starting node is named `you` and the goal is to reach a node named `out`.

## Usage

```bash
./day11 <path-to-input-file>
```

## Example

Given `example-data.txt`:

```
aaa: you hhh
you: bbb ccc
bbb: ddd eee
ccc: ddd eee fff
ddd: ggg
eee: out
fff: out
ggg: out
hhh: ccc fff iii
iii: out
```

Running the program:

```bash
./day11 example-data.txt
```

Output:

```
you->bbb->ddd->ggg->out
you->bbb->eee->out
you->ccc->ddd->ggg->out
you->ccc->eee->out
you->ccc->fff->out
Total unique paths: 5
```

## Building

```bash
go build -o day11
```

## Testing

```bash
go test
```

## Algorithm

The program uses Depth-First Search (DFS) with backtracking to find all unique paths from the start node to the end node. The algorithm:

1. Parses the input file to build a graph structure with nodes and connections
2. Starting from the `you` node, explores all possible paths using DFS
3. Tracks visited nodes to avoid cycles within a single path
4. Backtracks to explore alternative paths
5. Collects all paths that reach the `out` node

## Implementation Details

- **Node struct**: Represents a graph node with its name and list of connections
- **Graph struct**: Contains a map of all nodes in the graph
- **Path type**: A slice of node names representing a path through the graph
- **ParseGraph**: Parses input lines into a Graph structure
- **FindAllPaths**: Public method to find all paths from start to end
- **dfs**: Recursive depth-first search implementation with backtracking
- **Path.String()**: Formats a path as node names separated by arrows (`->`)

## Notes

- The graph is directed, so connections are one-way
- The algorithm avoids cycles by tracking visited nodes during path exploration
- All unique paths are found, not just the shortest path
- Empty connection lists are valid (dead-end nodes)

## Thoughts On AI Solutions

1. The AI understood the requirements and generated a complete solution, that got the correct answer for the example-data.2. I tested the solution against the puzzle input and it worked correctly.

TODO: add more details after I finish the puzzle.

- I did not give the AI the exact instructions from Advent of Code, but rather paraphrased them with my understanding of the problem.
- I did not ask the AI to optimize for performance or efficiency.
- I did not provide any starter code or templates; the AI generated the entire solution from scratch.
- I did not intervene in the coding process except to provide prompts and clarifications as needed.
- I only updated this section of the README for this day.

### Initial Prompts

#### Part 1

> build a go app in the day11 folder. It should accept a first parameter that is a path to a text file. It should track the data within it using structs and use methods on those structs to perform the computations and split those into files in directories as appropriate. Add test and a basic readme for the app as well. Follow patterns from the other days in the @days folder, particularly @day10.
> The file will contain lines that represet a node. Each line will have the name of the node, a colon and then a series of space separated nodes that it connects to. These connections are mono-directional, from the named node to the nodes listed after the colon. The starting node is named `you` and the goal is to reach a node named `out`.
> The goal of the app is to read the file, then find every possible path from the `you` node to the `out` node, and print those paths to standard output, one per line. Each path should be represented as a series of node names separated by arrows (`->`). Finally, the app should print the total number of unique paths found.
> For reference processing, @days/day11/example-data.txt should have 5 unique paths from `you` to `out`.

#### Part 2

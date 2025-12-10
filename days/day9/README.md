# Day9 - Maximum Rectangle Area

A Go CLI application that finds the largest rectangle area from a set of 2D coordinates.

## Usage

```bash
go run . <path-to-input-file> <mode>
```

Where `<mode>` is either:

- `original` - Find largest rectangle using any pair as opposite corners
- `contained` - Find largest rectangle contained within the polygon formed by connecting all points

The input file should contain one coordinate pair per line in the format `x,y`:

```
7,1
11,1
11,7
9,7
```

## Problem Description

Given a set of 2D coordinates, this program can operate in two modes:

### Original Mode

Find the largest area that can be created by using any pair of coordinates as opposite corners of a rectangle. The rectangles can be 1 in height or width.

The area calculation uses **inclusive counting**: a rectangle from (2,3) to (7,5) has:

- Width: 7 - 2 + 1 = 6 (includes both endpoints)
- Height: 5 - 3 + 1 = 3 (includes both endpoints)
- Area: 6 × 3 = 18

### Contained Mode

The points are treated as vertices of a polygon (ordered by angle from centroid). The program finds the largest axis-aligned rectangle that:

1. Uses two of the input coordinates as opposite corners
2. Is completely contained within the polygon (all edges and interior points must be inside or on the boundary)

The program processes all pairs of points and outputs the maximum area found based on the selected mode.

## Examples

For the example data:

```
7,1
11,1
11,7
9,7
9,5
2,5
2,3
7,3
```

- **Original mode**: The largest rectangle has an area of **50** (from points (11,1) to (2,5))
- **Contained mode**: The largest rectangle has an area of **24** (contained within the polygon shape)

```bash
go run . example-data.txt original
# Output: Largest rectangle area: 50

go run . example-data.txt contained
# Output: Largest rectangle area: 24

go run . example-data.txt contained output.svg
# Output: Largest rectangle area: 24
#         Visualization saved to: output.svg
```

### Visualization

When using `contained` mode with an output file, the program generates an SVG visualization:

- **Light blue polygon**: The shape formed by connecting all input coordinates
- **Red rectangle**: The largest contained rectangle
- **Dark blue dots**: Input coordinate points
- **Red dots**: Rectangle corners

The visualization automatically scales to fit the data and includes the calculated area as a label.

Example visualizations can be generated with:

```bash
go run . example-data.txt contained example-contained.svg
go run . puzzle-input.txt contained puzzle-contained.svg
```

## Implementation

The application uses:

- `Point` struct to represent 2D coordinates
- `Rectangle` struct with methods to calculate area
- Processor functions separated into `processor.go`
- Visualization functions in `visualizer.go` (SVG generation)
- Comprehensive tests in `main_test.go`

### Performance Optimization

The contained mode uses an optimized algorithm for large rectangles:

- **Sampling**: Instead of checking every point on a rectangle's perimeter, the algorithm samples at most 100 points per edge
- **Early rejection**: Checks corners and center first for fast rejection
- **Result**: Can efficiently process rectangles up to 100,000×100,000 in ~8μs (constant time regardless of size)

This makes the algorithm scale from O(perimeter) to O(1) for containment checks on large rectangles.

Run benchmarks to see performance characteristics:

```bash
go test -bench=. -benchtime=1s
```

## Running Tests

```bash
go test -v
```

The tests verify:

- Parsing of coordinate pairs
- Rectangle area calculations
- Processing of the example data in both modes
- Edge cases (same x or y coordinates)
- Polygon containment logic
- Large grid performance (100,000×100,000)

Run benchmarks:

```bash
go test -bench=. -benchtime=1s
```

## Thoughts On AI Solutions

1. The AI understood the requirements and generated a complete solution.
2. I tested the solution against the puzzle input and it worked correctly.
3. I asked it to solve part 2 and forgot to give it the expected output for the example, but it still produced a correct solution.
4. The solution seemed correct, but it was too slow when tested with the puzzle input. I asked the AI to optimize it.
5. It was able to get the solution to run much faster, but it's answer to the puzzle was too high (even though the example-data.txt test was correct). I asked it to try again.
6. That answer was still too hight, tried again.
7. Even with better example data it gave the same answer, I must have given bad instructions, trying again.
8. I had it create a visualization to help debug, but it looks fine, like a square in a circle basically. I reviewed the original question and I think the way it's worded it might not want it to have any of the original points in the rectangle other than the corners. Asked it to update the code to reflect that.

TODO: summary of thoughts

- I did not give the AI the exact instructions from Advent of Code, but rather paraphrased them with my understanding of the problem.
- I did not ask the AI to optimize for performance or efficiency.
- I did not provide any starter code or templates; the AI generated the entire solution from scratch.
- I did not intervene in the coding process except to provide prompts and clarifications as needed.
- I only updated this section of the README for this day.

### Initial Prompts

#### Part 1

> build a go app in the day9 folder. It should accept a first parameter that is a path to a text file. It should track the data within it using structs and use methods on those structs to perform the computations and split those into files in directories as appropriate. Add test and a basic readme for the app as well. Follow patterns from the other days in the @days folder.
> The file will contain a set of 2D coordinates per line. Find the largest area that can be created by using a pair of those coordinates as corners of a rectangle. They can be 1 in height or width.
> For reference processing, @days/day9/example-data.txt should result in an area of 50.

#### Part 2

> Add a mandatory second parameter that chooses the largest, original, mode or the contained mode. In the contained consider connecting all coordinates as corners of a shape. Then choose the largest rectangle that can fit entirely within that shape using 2 of those coordinates as opposite corners.

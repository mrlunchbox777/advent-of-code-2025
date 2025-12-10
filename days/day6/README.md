# Day 6 - Column Calculator

A Go application that processes a 2D array of numbers with column-wise operations.

## Description

This application reads a text file containing a 2D grid of numbers followed by a row of operators (`+` or `*`). It calculates the result of applying each operator to all numbers in the corresponding column, then sums all column totals.

## Input Format

The input file should contain:

- Multiple rows of numbers separated by spaces (alignment doesn't matter)
- A final row containing operators (`+` or `*`) for each column
- Empty lines are ignored

Example:

```
123 328  51 64
 45 64  387 23
  6 98  215 314
*   +   *   +
```

## Usage

Build the application:

```bash
go build -o day6
```

Run with a data file:

```bash
./day6 <path-to-data-file>
```

Example:

```bash
./day6 example-data.txt
```

## Output

The application prints:

- Each column number and its calculated total
- The sum of all column totals

Example output:

```
Column 1: 33210
Column 2: 490
Column 3: 4243545
Column 4: 401
Total: 4277556
```

## Implementation

The application uses:

- `Column` struct: Stores numbers and operator for a column
  - `Calculate()` method: Performs the operation on all numbers
- `Grid` struct: Stores all columns
  - `CalculateTotal()` method: Calculates and prints all column totals and the grand total

## Testing

Run tests:

```bash
go test -v
```

The test suite includes:

- Unit tests for column calculations
- Grid total calculation tests
- File parsing tests
- Validation against example-data.txt

## Thoughts On AI Solutions

1. The AI understood the problem, found a solution, and validated it against the example data very quickly.
2. I tested the AI's solution against the puzzle input data and it produced the correct answer on the first try.

TODO: summary of thoughts

- I did not give the AI the exact instructions from Advent of Code, but rather paraphrased them with my understanding of the problem.
- I did not ask the AI to optimize for performance or efficiency.
- I did not provide any starter code or templates; the AI generated the entire solution from scratch.
- I did not intervene in the coding process except to provide prompts and clarifications as needed.
- I only updated this section of the README for this day.

### Initial Prompts

#### Part 1

> build a go app in the day6 folder. It should accept a single parameter that is a path to a text file. The file will contain a 2 dimentional array of numbers and operators, '+' and '\*'. The columns will be seperated by spaces, alignment doesn't matter, just the spaces. The final row will be the operator that should be used between all of the numbers in that column to get the column total. It should track the array and the entries within it using structs and use methods on those structs to perform the computations. It should perform the operation represented by the operator in the final row between all of the numbers and get a column total, then print that coumn total and column number. Finally, it should add all column totals together and print the total. Add test and a basic readme for the app as well. For reference processing, the @days/day6/example-data.txt should total 4277556.

#### Part 2

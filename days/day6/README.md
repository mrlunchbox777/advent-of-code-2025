# Day 6 - Column Calculator

A Go application that processes a 2D array of numbers with column-wise operations.

## Description

This application reads a text file containing a 2D grid of numbers followed by a row of operators (`+` or `*`). It calculates the result of applying each operator to all numbers in the corresponding column, then sums all column totals.

The application supports two modes: **original** and **aligned**.

## Input Format

The input file should contain:

- Multiple rows of numbers
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

Run with a mode and data file:

```bash
./day6 <mode> <path-to-data-file>
```

Modes:

- `original` - Parse numbers as space-separated fields (alignment doesn't matter)
- `aligned` - Parse numbers based on character position alignment

Examples:

```bash
./day6 original example-data.txt
./day6 aligned example-data.txt
```

## Modes

### Original Mode

In original mode, numbers are parsed as space-separated fields. Alignment and extra spaces don't matter.

Example output:

```
$ ./day6 original example-data.txt
Column 1: 33210
Column 2: 490
Column 3: 4243455
Column 4: 401
Total: 4277556
```

### Aligned Mode

In aligned mode, character positions matter. Each operator defines a column range (from its position to the next operator or end of line). Within each column range, digits are read:

- Right-to-left through character positions
- Top-to-bottom across rows
- Each character position creates one number from vertical digits

Example output:

```
$ ./day6 aligned example-data.txt
Column 1: 8544
Column 2: 625
Column 3: 3253600
Column 4: 1058
Total: 3263827
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
- File parsing tests for both modes
- Validation against example-data.txt in both modes

## Thoughts On AI Solutions

1. The AI understood the problem, found a solution, and validated it against the example data very quickly.
2. I tested the AI's solution against the puzzle input data and it produced the correct answer on the first try.
3. The AI again fully understood the problem, found a solution, and validated it, though it took a bit longer and a few iterations this time.
4. I tested the AI's solution against the puzzle input data and it produced the correct answer on the first try again.

Today may have been the smoothest experience yet with the AI. It took longer to understand the problem and write a meaningful prompt than it did for the AI to produce a working solution. Both parts were done in a single pass with no real issues.

- I did not give the AI the exact instructions from Advent of Code, but rather paraphrased them with my understanding of the problem.
- I did not ask the AI to optimize for performance or efficiency.
- I did not provide any starter code or templates; the AI generated the entire solution from scratch.
- I did not intervene in the coding process except to provide prompts and clarifications as needed.
- I only updated this section of the README for this day.

### Initial Prompts

#### Part 1

> build a go app in the day6 folder. It should accept a single parameter that is a path to a text file. The file will contain a 2 dimentional array of numbers and operators, '+' and '\*'. The columns will be seperated by spaces, alignment doesn't matter, just the spaces. The final row will be the operator that should be used between all of the numbers in that column to get the column total. It should track the array and the entries within it using structs and use methods on those structs to perform the computations. It should perform the operation represented by the operator in the final row between all of the numbers and get a column total, then print that coumn total and column number. Finally, it should add all column totals together and print the total. Add test and a basic readme for the app as well. For reference processing, the @days/day6/example-data.txt should total 4277556.

#### Part 2

> add a mandatory parameter to the app that determines if it uses the original logic that is there now or if is uses a new methodology. The new method is to work from read numbers top to bottom, right to left, within the column, alignment and spacing matters here. Meaning if you had 3 numbers in a column `4`, `123`, `56`, `+`, then the actual number to be used in the calculation would be `1`, `25`, and `436`; however if the spacing changed to `4`, `123`, `56`, `+`, then the numbers would be `415`, `26`, and `3`. The operator and how to calculate the column totals remains the same, except with the new numbers. The final total also uses the same logic. The @days/day6/example-data.txt should now total 3263827.

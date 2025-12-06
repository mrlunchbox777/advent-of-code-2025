# Day 5 - Number Validator

A Go application that validates numbers against a set of inclusive ranges.

## Overview

This program reads a text file containing two sections separated by a blank line:

1. **Ranges**: Lines with inclusive number ranges in the format `start-end`
2. **Numbers**: Lines with individual numbers to validate

The program builds a set of valid numbers from all ranges, then validates each number from the second list against this set.

## Building

```bash
go build -o validator
```

## Usage

```bash
./validator <path-to-file>
```

### Example

```bash
./validator example-data.txt
```

## Input Format

```
3-5
10-14
16-20
12-18

1
5
8
11
17
32
```

## Output

The program prints each number with its validation status:

```
1: false
5: true
8: false
11: true
17: true
32: false

Total valid numbers: 3
```

## Testing

Run the test suite:

```bash
go test -v
```

## Implementation Details

### Structs

- **Range**: Represents an inclusive range with `Start` and `End` values
- **RangeList**: Collection of ranges with methods to build a valid set
- **NumberList**: Collection of numbers with validation methods

### Methods

- `Range.Contains(n int) bool`: Check if a number is within the range
- `RangeList.BuildValidSet() map[int]bool`: Build a set of all valid numbers from ranges
- `NumberList.ValidateAgainstSet(validSet map[int]bool) int`: Validate numbers and return count

## Thoughts On AI Solutions

1. The AI correctly understood the problem requirements and provided a solution that counts valid numbers based on given ranges. I changed editors and runners, from VSCode to LazyVim, and everything seems to have worked fine. I can't tell if there was any corruption this time because it gave less output. It did seem to stall part way through and I had to cancel and then ask it to continue, but it picked up where it left off without issue and completed the solution.

- I did not give the AI the exact instructions from Advent of Code, but rather paraphrased them with my understanding of the problem.
- I did not ask the AI to optimize for performance or efficiency.
- I did not provide any starter code or templates; the AI generated the entire solution from scratch.
- I did not intervene in the coding process except to provide prompts and clarifications as needed.
- I only updated this section of the README for this day.

### Initial Prompts

#### Part 1

> each line, the second will contain just a number on each line. It should track these lists and entries using structs and methods on those structs to perform the computations. It should build a set of valid numbers out of all of the ranges in the first list. It should then find the valid entries in the second list, they are valid if they are in the set built from the first list. It should print the number and if it's valid as it goes. Finally, it should print the total count of valid numbers from the second list. Add test for validation, and a basic readme. For example the @days/day5/example-data.txt 5, 11, and 17 were valid, resulting in a count of 3.

#### Part 2

todo

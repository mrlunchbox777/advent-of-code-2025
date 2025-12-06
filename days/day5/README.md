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
./validator <path-to-file> <mode>
```

### Modes

- **validate** - Count valid numbers from the second list against the ranges
- **total** - Count total possible valid numbers across all ranges

### Examples

Validate mode (count valid numbers from second list):

```bash
./validator example-data.txt validate
```

Total mode (count total possible valid numbers in ranges):

```bash
./validator example-data.txt total
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

### Validate Mode

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

### Total Mode

The program prints the total count of possible valid numbers:

```
Total possible valid numbers: 14
```

This counts unique numbers across all ranges:

- 3-5: 3, 4, 5 (3 numbers)
- 10-14: 10, 11, 12, 13, 14 (5 numbers)
- 16-20: 16, 17, 18, 19, 20 (5 numbers)
- 12-18: Adds only 15 (1 new number, others overlap)
- Total: 14 unique numbers

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
- `RangeList.IsValid(n int) bool`: Check if a number is valid against any range
- `RangeList.CountTotalValid() int`: Count total unique numbers across all ranges
- `NumberList.ValidateAgainstRanges(rangeList *RangeList) int`: Validate numbers and return count

### Performance

**Validate mode** uses range-based validation instead of building a complete set, making it efficient for ranges in the billions. Each number is checked against the ranges in O(r) per number, where r = number of ranges, avoiding memory issues from materializing billion-element sets.

**Total mode** uses a range merging algorithm that sorts ranges and merges overlapping/adjacent ranges, then calculates counts mathematically. This runs in O(r log r) time for sorting plus O(r) for merging, where r = number of ranges. Memory usage is O(r) regardless of range size, making it efficient even for ranges in the billions or trillions. The algorithm never materializes individual numbers.

## Thoughts On AI Solutions

1. The AI correctly understood the problem requirements and provided a solution that counts valid numbers based on given ranges. I changed editors and runners, from VSCode to LazyVim, and everything seems to have worked fine. I can't tell if there was any corruption this time because it gave less output. It did seem to stall part way through and I had to cancel and then ask it to continue, but it picked up where it left off without issue and completed the solution.
1. I attempted the puzzle input, but the AI generated a solution that was far too slow. I asked it to optimize for performance, and it produced a much better solution that checked each number against the ranges directly instead of building a massive set of valid numbers. This worked well and completed in a reasonable time.
1. I attempted the solution again, it ran basically instantly, and it gave the correct answer.
1. The AI updated the program to add the second mode, and it seemed fine for the example input.
1. I tested the second mode with the puzzle input, and it was far too slow again. I asked it to optimize for performance, and it produced a solution that merged ranges and calculated counts mathematically instead of materializing all valid numbers. This worked well and completed in a reasonable time for the test data it self-generated, but was still too slow for the puzzle-input.

- I did not give the AI the exact instructions from Advent of Code, but rather paraphrased them with my understanding of the problem.
- ~~I did not ask the AI to optimize for performance or efficiency.~~
- I did not provide any starter code or templates; the AI generated the entire solution from scratch.
- I did not intervene in the coding process except to provide prompts and clarifications as needed.
- I only updated this section of the README for this day.

### Initial Prompts

#### Part 1

> each line, the second will contain just a number on each line. It should track these lists and entries using structs and methods on those structs to perform the computations. It should build a set of valid numbers out of all of the ranges in the first list. It should then find the valid entries in the second list, they are valid if they are in the set built from the first list. It should print the number and if it's valid as it goes. Finally, it should print the total count of valid numbers from the second list. Add test for validation, and a basic readme. For example the @days/day5/example-data.txt 5, 11, and 17 were valid, resulting in a count of 3.

#### Part 2

add a second mandatory parameter that will either count the valid numbers in the second list as before, or will count the total possible valid numbers in the ranges from the first list. For example, with the same input as before, the total possible valid numbers would be 3-5 (3), 10-14 (5), 16-20 (5), and 12-18 (1) for a total of 14 possible valid numbers. Update tests and readme accordingly.

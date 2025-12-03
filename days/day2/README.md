# Day 2: Repeated Sequence Validator

This Go application identifies numbers within specified ranges that are comprised entirely of repeated sequences.

## Problem Description

Given ranges of numbers (e.g., `11-22`), find all numbers that consist of a pattern repeated multiple times:
- Valid: `11` (pattern "1" repeated 2 times), `1010` (pattern "10" repeated 2 times), `222222` (pattern "2" repeated 6 times)
- Invalid: `101` (no repeating pattern), `111` (odd-length single-digit pattern)

**Rule**: For single-digit patterns, the total length must be even. For multi-digit patterns, the pattern must divide evenly into the total length.

## Usage

```bash
go run . <filepath>
```

### Example

```bash
go run . example-data.txt
```

**Output:**
```
11-22 has 2 invalid ID(s): [11 22]
95-115 has 1 invalid ID(s): [99]
998-1012 has 1 invalid ID(s): [1010]
1188511880-1188511890 has 1 invalid ID(s): [1188511885]
222220-222224 has 1 invalid ID(s): [222222]
1698522-1698528 contains no invalid IDs.
446443-446449 has 1 invalid ID(s): [446446]
38593856-38593862 has 1 invalid ID(s): [38593859]

Total sum of invalid IDs: 1227775554
```

## Input Format

The input file should contain a single line of comma-separated ranges in the format `lower-upper`:
```
11-22,95-115,998-1012,1188511880-1188511890
```

## Testing

Run the test suite:
```bash
go test -v
```

The tests validate:
- Range parsing logic
- Repeated sequence detection algorithm
- Expected results for example data (sum = 1227775554)

## Implementation

- **`Range` struct**: Represents a number range with lower and upper bounds
- **`ParseRange`**: Parses string format "lower-upper" into a Range
- **`FindRepeatedSequenceNumbers`**: Finds all repeated sequence numbers in the range
- **`isRepeatedSequence`**: Checks if a number is comprised of a repeated pattern

## Thoughts On AI Solutions

1. The AI correctly understood the problem, but stumbled initially corrupting it's own code when trying to refactor. It caught itself and fixed the issues.
2. I tested the solution against the puzzle input and it calculated the answer too high.

- I did not give the AI the exact instructions from Advent of Code, but rather paraphrased them with my understanding of the problem.
- I did not ask the AI to optimize for performance or efficiency.
- I did not provide any starter code or templates; the AI generated the entire solution from scratch.
- I did not intervene in the coding process except to provide prompts and clarifications as needed.
- I only updated this section of the README for this day.

### Initial Prompts

#### Part 1

> build a go app in the day2 folder. It should accept a single parameter that is a path to a text file. That file will contain a single line of comma seperated entries, the entries are ranges of numbers in the form of lower bound, dash, upper bound, e.g. 1-9. It should track each entry as a struct and use methods on the struct to perform the computations. It should process each entry and find numbers in the range that are entirely comprised of repeated sequences, e.g. 11 and 1010, but not 101. As it processes the entries it should print the range and the repeated sequences in that range. Then finally at the end it should sum the total of all of the identified numbers made of repeated sequences. Add test for the app as well. For reference processing the example-data.txt should have a final sum of identified numbers equal to 1227775554. Specific info on example-data.txt
> 11-22 has two invalid IDs, 11 and 22.
> 95-115 has one invalid ID, 99.
> 998-1012 has one invalid ID, 1010.
> 1188511880-1188511890 has one invalid ID, 1188511885.
> 222220-222224 has one invalid ID, 222222.
> 1698522-1698528 contains no invalid IDs.
> 446443-446449 has one invalid ID, 446446.
> 38593856-38593862 has one invalid ID, 38593859.

#### Part 2

> 


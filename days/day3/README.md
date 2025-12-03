# Day 3: Largest Two-Digit Number Finder

This Go application finds the largest two-digit number that can be formed from each line by selecting two digits in their original order and concatenating them.

## Problem Description

Given lines of digits, select two digits (without reordering) that when concatenated form the largest possible two-digit number.

**Examples:**
- `12345` → select `4` and `5` → `45`
- `54321` → select `5` and `4` → `54`
- `987654321111111` → select `9` and `8` → `98`
- `811111111111119` → select `8` and `9` → `89`

The key constraint: digits must be selected in the order they appear (left to right), and concatenated in that order.

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
987654321111111 -> 9 and 8 = 98
811111111111119 -> 8 and 9 = 89
234234234234278 -> 7 and 8 = 78
818181911112111 -> 9 and 2 = 92

Total sum: 357
```

## Input Format

The input file should contain multiple lines, each with a sequence of digits:
```
987654321111111
811111111111119
234234234234278
818181911112111
```

## Testing

Run the test suite:
```bash
go test -v
```

The tests validate:
- Two-digit number finding logic for various inputs
- Expected results for example data (sum = 357)

## Implementation

- **`Entry` struct**: Represents a line of digits
- **`NewEntry`**: Creates an Entry from a string
- **`FindLargestTwoDigitNumber()`**: Examines all pairs of digits (i, j) where i < j, computes the concatenation in order, and returns the pair that forms the largest two-digit number

### Algorithm

1. Extract all digits from the input line
2. For each pair of positions (i, j) where i < j:
   - Compute the two-digit number: `digit[i] * 10 + digit[j]`
3. Return the pair with the maximum value

## Thoughts On AI Solutions

1. This time the AI misunderstood the requirement to keep the order of digits as they appear in the input. It initially generated a solution that allowed reordering, which was incorrect. It also corrupted its own code during refactoring attempts again. It caught itself in the end and fixed the issues.

Today

- I did not give the AI the exact instructions from Advent of Code, but rather paraphrased them with my understanding of the problem.
- I did not ask the AI to optimize for performance or efficiency.
- I did not provide any starter code or templates; the AI generated the entire solution from scratch.
- I did not intervene in the coding process except to provide prompts and clarifications as needed.
- I only updated this section of the README for this day.

### Initial Prompts

#### Part 1

> build a go app in the day3 folder. It should accept a single parameter that is a path to a text file. The file will contain multiple lines of entries, each entry will contain many digits. It should track each entry as a struct and use methods on the struct to perform the computations. Without reordering the digits it should select the 2 digits that when concatenated form the largest two digit number possible from the entry, e.g. in 12345 select '4' and '5' for 45 or 54321 select '5' and '4' for 54. As it processes the entries it should print the line, the selections, and the resulting two digit number. Then finally, at the end it should sum the two digit results from each entry and display that. Add test for the app as well. For reference processing, the example-data.txt should have a final sum of 357. Add a basic README.md as well.

#### Part 2

> test

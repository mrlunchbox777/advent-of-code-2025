# Day 3: Largest N-Digit Number Finder

This Go application finds the largest n-digit number that can be formed from each line by selecting n digits in their original order and concatenating them.

## Problem Description

Given lines of digits, select n digits (without reordering) that when concatenated form the largest possible n-digit number.

**Examples (2 digits):**
- `12345` → select `4` and `5` → `45`
- `54321` → select `5` and `4` → `54`
- `987654321111111` → select `9` and `8` → `98`
- `811111111111119` → select `8` and `9` → `89`

**Examples (12 digits):**
- `987654321111111` → select `987654321111` → `987654321111`
- `811111111111119` → select `811111111119` → `811111111119`

The key constraint: digits must be selected in the order they appear (left to right), and concatenated in that order.

## Usage

```bash
go run . <filepath> <digitCount>
```

Where `<digitCount>` is the number of digits to select and concatenate.

### Examples

**2 digits:**
```bash
go run . example-data.txt 2
```

**Output:**
```
987654321111111 -> 98 = 98
811111111111119 -> 89 = 89
234234234234278 -> 78 = 78
818181911112111 -> 92 = 92

Total sum: 357
```

**12 digits:**
```bash
go run . example-data.txt 12
```

**Output:**
```
987654321111111 -> 987654321111 = 987654321111
811111111111119 -> 811111111119 = 811111111119
234234234234278 -> 434234234278 = 434234234278
818181911112111 -> 888911112111 = 888911112111

Total sum: 3121910778619
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
- Variable digit count finding logic (2, 3, 5, 12 digits)
- Expected results for example data (2 digits: sum = 357, 12 digits: sum = 3121910778619)

## Implementation

- **`Entry` struct**: Represents a line of digits
- **`NewEntry`**: Creates an Entry from a string
- **`FindLargestNumber(n)`**: Examines all combinations of n digit positions maintaining order, computes the concatenation, and returns the combination that forms the largest n-digit number
- **`FindLargestTwoDigitNumber()`**: Convenience wrapper for `FindLargestNumber(2)` for backward compatibility

### Algorithm

1. Extract all digits from the input line
2. Generate all combinations of n positions where position indices maintain increasing order
3. For each combination:
   - Concatenate the digits at those positions
   - Calculate the resulting number
4. Return the combination with the maximum value

The algorithm uses recursion to generate all valid combinations of n positions from the available digits.

## Thoughts On AI Solutions

1. This time the AI misunderstood the requirement to keep the order of digits as they appear in the input. It initially generated a solution that allowed reordering, which was incorrect. It also corrupted its own code during refactoring attempts again. It caught itself in the end and fixed the issues.
2. I tested the solution against the puzzle input and it calculated the correct answer.
3. For part 2, I asked the AI to generalize the solution to handle any n-digit count. It adjusted the logic and produced a working solution based on the example data.
4. When I ran the solution against the puzzle input for n=12, it was too slow to complete.

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

> add a second mandatory flag that will determine how many digits should be used to concatenate, the current solution forces 2. The ordering rule stays the same regardless of how many digits are selected. When the example-data.txt uses 12 instead of 2, the sum should be 3121910778619.

# Day 2: Repeated Sequence Validator

This Go application identifies numbers within specified ranges that are comprised entirely of repeated sequences.

## Problem Description

Given ranges of numbers (e.g., `11-22`), find all numbers that consist of a pattern repeated multiple times.

### Modes

**Exact Mode** (`exact`): Pattern must repeat exactly 2 times
- Valid: `11` ("1" × 2), `1010` ("10" × 2), `222222` ("222" × 2), `1111` ("11" × 2)
- Invalid: `101` (no pattern), `111` ("1" × 3, not exactly 2), `123123123` ("123" × 3, not exactly 2)

**Any Mode** (`any`): Pattern must repeat 2 or more times
- Valid: `11` ("1" × 2), `111` ("1" × 3), `1111` ("1" × 4), `123123123` ("123" × 3)
- Invalid: `101` (no repeating pattern), `12345` (no pattern)

## Usage

```bash
go run . <filepath> <mode>
```

Where `<mode>` is either `exact` or `any`.

### Examples

**Exact Mode:**
```bash
go run . example-data.txt exact
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

**Any Mode:**
```bash
go run . puzzle-input.txt any
```

**Output (last line):**
```
Total sum of invalid IDs: 15704845910
```

**Comparison:**
- Exact mode (puzzle-input.txt): 5398419778
- Any mode (puzzle-input.txt): 15704845910

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
- Repeated sequence detection algorithm for both modes
- Expected results for example data (exact mode sum = 1227775554)

## Implementation

- **`Range` struct**: Represents a number range with lower and upper bounds
- **`ParseRange`**: Parses string format "lower-upper" into a Range
- **`FindRepeatedSequenceNumbers(mode)`**: Finds all repeated sequence numbers in the range based on mode
- **`isRepeatedSequence(n, mode)`**: Checks if a number is comprised of a repeated pattern
  - `exact` mode: Checks if the number can be split in half with both halves identical
  - `any` mode: Tries all pattern lengths (1 to length/2) to find any repeating pattern with 2+ repetitions

## Thoughts On AI Solutions

1. The AI correctly understood the problem, but stumbled initially corrupting it's own code when trying to refactor. It caught itself and fixed the issues.
2. I tested the solution against the puzzle input and it calculated the answer too high.
3. I asked the AI to debug its solution. It identified that the logic for detecting repeated sequences needed to be exactly twice, something mentioned in the puzzle input but not in my initial prompt. It made adjustments to the logic accordingly.
4. I tested the solution against the puzzle input again and it produced the correct results.
5. I asked the AI to handle any mode where the pattern can repeat 2 or more times. It adjusted the logic and produced a working solution. (This is basically what it had originally tried to do before I clarified the exact requirement.)
6. I tested the solution against the puzzle input again and it produced the correct results.

Today showed some real power of the AI in being able to debug issues as well as failures in generated code. It created bad code that didn't even compile, but was able to identify and fix the issues itself. It also demonstrated an ability to identify requirements that were not explicitly stated in the prompt but were in the puzzle description (not provided to the AI). The tests were particularly helpful in ensuring correctness, including the example data and edge cases.

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

> Add a second mandatory flag that will allow the current processing, but also allow searching for patterns that repeat any number of times rather than just twice (and handle the final sum accordingly), e.g. 123123123 would not currently be an identified pattern, but with the new requirement would be


# Day1 - Dial CLI

Usage:

```bash
go run . <path-to-input-file> <mode>
```

Where `<mode>` is either:
- `exact` - Counts only when the dial ends at exactly 0
- `passes` - Counts whenever the dial passes through or lands on 0 (excluding moves that start at 0)

The input file should contain one entry per line with a direction (`L` or `R`) and a number, for example:

```
L68
R34
L99
```

The program starts the dial at `50`. For each entry it prints the entry, the starting value, and the ending value. The dial wraps around in the range `0-99`.

## Examples

```bash
# Count only endings at 0
go run . example-data.txt exact

# Count all passes through 0
go run . example-data.txt passes
```

## Thoughts On AI Solutions

1. The AI correctly understood the problem requirements and provided a solution that processes the input file line by line, but it was pretty hamfisted in its approach.
2. I asked it to add tests so it could verify it's correctness and start refactoring.
3. The AI struggled to refactor the code into smaller functions, often creating functions that were too large or not well-defined. It did finally manage to break the code down into smaller functions after several iterations.
4. I ran the solution against part 1 and it produced the correct results.
5. For part 2, it understood the new requirement but initially failed to implement the logic to count passes through 0 correctly. After prompting and several iterations, it adjusted the logic and produced a working solution.
6. I ran the solution against part 2 and it produced the wrong results, it was too low.
7. I asked the AI to debug its solution. It identified that the logic for counting passes through 0 was flawed and made adjustments. I then clarified the requirement with the "Be careful" text and it validated that the changes it made accounted for that.
8. I ran the solution against part 2 again and it produced the correct results.

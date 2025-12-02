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

It was an interesting exercise in seeing how well the AI could understand and implement the requirements, as well as how it handled debugging and refactoring. While it did eventually arrive at a correct solution, it required significant prompting and guidance to get there. I didn't have to review every line of code, but I did have to guide the AI through the process. I think having the AI write tests was particularly helpful in ensuring correctness, including the example data and edge cases. Specifically, it made the exact mistake that the prompt warned about, even with the tests; but once I let it know about the issue it was able to fix it, without the clarification, and then validated it did fix it with the clarification. The tests only caught part of the issue.

- I did not give the AI the exact instructions from Advent of Code, but rather paraphrased them with my understanding of the problem.
- I did not ask the AI to optimize for performance or efficiency.
- I did not provide any starter code or templates; the AI generated the entire solution from scratch.
- I did not intervene in the coding process except to provide prompts and clarifications as needed.
- I only updated this section of the README for this day.

### Initial Prompts

#### Part 1

> build a go app in the day1 folder. It should accept a single parameter that is a path to a text file. That file will contain one entry per line, the entries are made of 2 parts: L or R and a number. The core entity in the app is a "dial" that will rotate left (minus) or right (plus), it ranges from 0-99. If it moves left (minus) at 0 it will loop to 99 rather than go negative, and if it moves right (plus) at 99 it will loop to 0 rather than go to 100+. The dial starts at 50. As it processes each entry it will post to console the entry (L68), the starting value (50), and the ending value (82). Finally, it will count how many times it ends at exactly 0 and print that at the very end.

#### Part 2

> add a second way to process the counting of 0's to count whenever it "passes" 0 in a move. This would count both passing 0 or ending at 0. This should be controlled by a mandatory second argument

# Day1 - Dial CLI

Usage:

```
go run . <path-to-input-file>
```

The input file should contain one entry per line with a direction (`L` or `R`) and a number, for example:

```
L68
R34
L99
```

The program starts the dial at `50`. For each entry it prints the entry, the starting value, and the ending value. The dial wraps around in the range `0-99`. At the end it prints how many times the dial ended at exactly `0`.

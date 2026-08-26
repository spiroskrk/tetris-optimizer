# Tetris Optimizer

Tetris Optimizer arranges a set of tetrominoes into the smallest possible square. Each piece keeps the orientation provided in the input and is represented by an uppercase letter in the result.

## Requirements

- Go 1.22 or newer

## Run

From the project directory:

```bash
go run . input.txt
```

The program expects exactly one input-file argument.

## Input format

Each tetromino is a 4×4 block containing only `#` and `.` characters. Tetrominoes must contain exactly four connected `#` cells and be separated by one blank line.

Example:

```text
##..
##..
....
....
```

The program accepts 1–26 tetrominoes and supports Unix and Windows line endings. Rotations and reflections are not generated.

## Output

The result is the smallest square that can contain every piece. Letters identify pieces in input order, while `.` represents unused space.

For the example above:

```text
AA
AA
```

Invalid arguments, unreadable files, and malformed tetrominoes produce a descriptive message on standard output.

## Tests

Run the complete test suite with:

```bash
go test ./...
```

Run it with the race detector using:

```bash
go test -race ./...
```

See [PRD.md](PRD.md) for detailed requirements, architecture, and algorithm notes.

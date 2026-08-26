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

## Error reference

Errors are written to standard output. The command then returns normally with exit status `0`.

| Situation | Message |
| --- | --- |
| Missing or extra command-line arguments | `Please insert only one Argument.` |
| Input file cannot be read | `Could not read file "<path>": <system error>` |
| Empty input file | `Could not parse file input file is empty` |
| More than 26 tetrominoes | `Could not parse file input contains <count> tetrominoes; maximum is 26` |
| Tetromino has the wrong number of rows | `Could not parse file not a valid tetromino block. has <count> lines instead of 4.` |
| Tetromino row has the wrong width | `Could not parse file not a valid tetromino block. row <row> has <count> columns instead of 4.` |
| Tetromino contains an invalid character | `Could not parse file character '<character>' is invalid. only # or . are accepted in a block` |
| Tetromino has more than four occupied cells | `Could not parse file tetromino has more than 4 occupied cells` |
| Tetromino has fewer than four occupied cells | `Could not parse file tetromino has <count> occupied cells instead of 4` |
| Tetromino cells are disconnected | `Could not parse file tetromino cells are not connected` |

The operating system supplies the final portion of file-reading errors, so that text varies by platform.

The following messages indicate an internal solver consistency failure rather than malformed user input:

- `Could not solve tetrominoes cannot solve an empty set of tetrominoes`
- `Could not solve tetrominoes solver returned invalid row index <index>`
- `Could not solve tetrominoes solver returned invalid placement index <index>`
- `Could not solve tetrominoes solver produced overlapping placements at (<x>,<y>)`

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

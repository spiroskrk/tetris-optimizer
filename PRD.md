# PRD — Tetris Optimizer

---

## 1. Problem Statement

Develop a Go program that receives exactly one argument: the path to a text file containing one or more tetrominoes.

The program must validate the input, parse each tetromino, and arrange all tetrominoes into the smallest possible square.

Each tetromino must be identified in the final output using uppercase Latin letters:

- `A` for the first tetromino
- `B` for the second tetromino
- `C` for the third tetromino
- and so on

If a valid arrangement cannot completely fill the square, unused cells are represented with `.`.

The solving algorithm will be based on Knuth's Algorithm X, using an Exact Cover representation and Dancing Links (DLX).

---

## 2. CLI Contract

Command:

```bash
go run . input.txt

Rules:

- The program expects exactly one argument.
- The argument must be the path to the input text file.
- The input file must contain at least one tetromino.
- The program is non-interactive.
- The program does not request additional input from the user.
- The result is written to standard output.

If the number of arguments is invalid, the program prints:

Please insert only one Argument.

and returns without processing a file.

If the file cannot be read, the file format is invalid, a tetromino is invalid, or the tetrominoes cannot be solved, the program prints the descriptive error returned by `Run`.

Internal functions return descriptive Go errors so that individual failures can be inspected and tested during development. `Run` wraps these errors with operation-level context, and the CLI prints the resulting error.

---

## 3. Input Format

Each tetromino is represented by a 4×4 block.

Example:

#...
#...
#...
#...

Tetrominoes are separated by exactly one blank line.

Example:

#...
#...
#...
#...

....
....
..##
..##

The parser accepts both Unix and Windows line endings.

Windows line endings:

\r\n

are normalized internally to:

\n

The parser may accept a single final newline at the end of the file.

Two or more final newlines are considered malformed input.

---

## 4. Tetromino Validation

Each tetromino block must satisfy all of the following rules.

### 4.1 Dimensions

Each tetromino must contain exactly:

- 4 rows
- 4 columns per row

Example of invalid height:

....
....
....

Example of invalid width:

.....
....
....
....

---

### 4.2 Allowed Characters

Only the following characters are allowed inside a tetromino block:

.
#

Any other character makes the tetromino invalid.

Example:

....
.#x.
.##.
....

This must produce an error.

---

### 4.3 Occupied Cell Count

Every valid tetromino must contain exactly four `#` cells.

Example of invalid count:

....
.##.
.#..
....

This contains only three occupied cells and must be rejected.

---

### 4.4 Connectivity

All four occupied cells must form one connected component.

Cells are considered adjacent only when they share an edge.

Two cells are adjacent when:

same X and Y differs by exactly 1

or:

same Y and X differs by exactly 1

Diagonal contact does not count.

Valid:

....
.##.
.##.
....

Invalid:

##..
....
..##
....

Connectivity is checked using a recursive depth-first search.

The algorithm:

1. Start from the first occupied cell.
2. Mark it as visited.
3. Search all other occupied cells.
4. Recursively visit every adjacent unvisited cell.
5. After traversal, verify that all four cells were visited.

---

## 5. Data Structures

### 5.1 Point

A `Point` represents one occupied cell of a tetromino.

type Point struct {
	X int
	Y int
}

---

### 5.2 Tetromino

A `Tetromino` stores the normalized shape and metadata required by the solver.

type Tetromino struct {
	Cells  [4]Point
	Width  int
	Height int
	Letter byte
}

`Cells` always contains exactly four points.

`Letter` represents the character used when rendering the final board.

Examples:

'A'
'B'
'C'

`Width` and `Height` represent the bounding box of the normalized tetromino.

---

### 5.3 Placement

A placement represents one possible position of one tetromino inside a candidate square.

type Placement struct {
	PieceIndex int
	X          int
	Y          int
}

`X` and `Y` represent the position where the normalized tetromino is anchored.

The occupied board cells can be reconstructed using:

boardX = placement.X + cell.X
boardY = placement.Y + cell.Y

---

## 6. Coordinate Normalization

Tetrominoes may appear anywhere inside their original 4×4 block.

Example:

....
..##
..##
....

Raw coordinates:

(2,1)
(3,1)
(2,2)
(3,2)

The parser normalizes these coordinates so the piece begins at the top-left of its bounding box.

First determine:

minX
minY

Then transform every point:

newX = oldX - minX
newY = oldY - minY

Normalized coordinates:

(0,0)
(1,0)
(0,1)
(1,1)

After normalization:

minX = 0
minY = 0

The piece dimensions are determined using:

width  = maxX + 1
height = maxY + 1

Example:

###
.#.

has:

Width  = 3
Height = 2

---

## 7. Parsing Architecture

The parsing flow is:

Raw file
   ↓
Read file
   ↓
Normalize line endings
   ↓
Remove one optional final newline
   ↓
Split into tetromino blocks
   ↓
Parse each block
   ↓
Validate dimensions
   ↓
Validate characters
   ↓
Collect occupied cells
   ↓
Validate occupied cell count
   ↓
Validate connectivity
   ↓
Normalize coordinates
   ↓
Create Tetromino
   ↓
Assign letter
   ↓
[]Tetromino

---

## 8. Parser Responsibilities

### 8.1 `parseData`

Suggested signature:

func parseData(data string) ([]Tetromino, error)

Responsibilities:

- Normalize `\r\n` to `\n`.
- Remove one optional trailing newline.
- Reject empty input.
- Split the file into blocks using:

strings.Split(data, "\n\n")

- Call `parseBlock` for every block.
- Assign letters based on input order.
- Return all parsed tetrominoes.

Letter assignment:

piece.Letter = byte('A' + i)

---

### 8.2 `parseBlock`

Suggested signature:

func parseBlock(block string) (Tetromino, error)

Responsibilities:

- Split the block into lines.
- Ensure exactly four lines exist.
- Ensure each line contains four characters.
- Ensure all characters are `.` or `#`.
- Collect the coordinates of all `#`.
- Ensure exactly four occupied cells exist.
- Ensure all occupied cells are connected.
- Normalize the coordinates.
- Determine width and height.
- Return a solver-ready `Tetromino`.

A `Tetromino` returned from `parseBlock` must always be valid.

This is an invariant of the program.

---

## 9. Program Entry Flow

The external entry point is:

func Run(path string) error

`Run` belongs to the `tetris` package.

Its responsibility is orchestration.

Conceptual flow:

Run
 ↓
Read input file
 ↓
parseData
 ↓
solve
 ↓
render

`Run` should not contain parsing or solving algorithms directly.

---

## 10. Minimum Square Size

If there are `n` tetrominoes, their total occupied area is:

4 × n

The first candidate square size is the smallest integer `size` satisfying:

size × size >= 4 × n

Example:

1 piece  → start at 2
2 pieces → start at 3
4 pieces → start at 4
5 pieces → start at 5

A simple integer implementation can be:

size := 1

for size*size < len(pieces)*4 {
	size++
}

The solver attempts candidate square sizes in increasing order.

The first size with a valid arrangement is the required smallest square.

---

## 11. Placement Generation

Before running Algorithm X, the program generates every legal placement of every tetromino inside the current candidate square.

For a tetromino with:

Width  = w
Height = h

inside a square of size `size`, valid anchor positions satisfy:

0 <= X <= size - w
0 <= Y <= size - h

Conceptual loop:

for y := 0; y <= size-piece.Height; y++ {
	for x := 0; x <= size-piece.Width; x++ {
		// generate placement
	}
}

Each anchor position creates one possible `Placement`.

Example:

A 2×2 square tetromino inside a 3×3 board has four possible placements:

(0,0)
(1,0)
(0,1)
(1,1)

---

## 12. Exact Cover Model

The solver uses Knuth's Algorithm X.

The problem is transformed into an Exact Cover matrix.

The matrix contains:

- columns representing constraints
- rows representing possible tetromino placements

Algorithm X chooses a combination of rows that satisfies all mandatory constraints while preventing conflicting board usage.

---

## 13. Constraint Columns

There are two categories of constraints.

### 13.1 Tetromino Constraints

Every tetromino must be placed exactly once.

For pieces:

A
B
C

the exact-cover structure contains piece columns equivalent to:

Piece A
Piece B
Piece C

These are mandatory constraints.

---

### 13.2 Board Cell Constraints

Each board cell can be occupied by at most one tetromino.

For a 3×3 square, cells include:

(0,0)
(1,0)
(2,0)
(0,1)
(1,1)
(2,1)
(0,2)
(1,2)
(2,2)

A board cell is not required to be occupied because the project allows empty spaces in the final square.

However, two tetrominoes must never occupy the same board cell.

Therefore board-cell constraints behave as optional or secondary constraints.

They enforce:

occupied zero or one times

rather than:

occupied exactly once

---

## 14. Exact Cover Rows

Each possible placement of a tetromino becomes one row.

Example:

Piece `A` is:

##
##

and is placed at:

X = 1
Y = 0

The occupied board cells are:

(1,0)
(2,0)
(1,1)
(2,1)

That row represents:

Piece A
Board cell (1,0)
Board cell (2,0)
Board cell (1,1)
Board cell (2,1)

Conceptually:

A | 10 | 20 | 11 | 21
1 |  1 |  1 |  1 |  1

Selecting this row means:

- piece A has been placed
- those four board cells are occupied

Any other placement using one of those cells conflicts with this row.

---

## 15. Algorithm X

Algorithm X searches the Exact Cover structure recursively.

Conceptually:

Choose an unsatisfied mandatory column
        ↓
Try one row containing that column
        ↓
Select that row
        ↓
Remove conflicting rows and columns
        ↓
Recurse
        ↓
If failure:
    restore removed rows and columns
    try next row

A solution exists when all mandatory tetromino columns have been satisfied.

The selected rows correspond to one placement for every tetromino.

---

## 16. Dancing Links

Algorithm X will be implemented using Knuth's Dancing Links data structure.

Each `1` in the Exact Cover matrix becomes a linked node.

Nodes are linked:

Left
Right
Up
Down

Columns maintain references to their nodes.

The major DLX operations are:

cover(column)
uncover(column)
search()

`cover` temporarily removes a constraint and all conflicting rows.

`uncover` restores them during backtracking.

Because links can be removed and restored efficiently, DLX avoids repeatedly rebuilding the Exact Cover matrix.

---

## 17. Solver Flow

The complete solver flow is:

Receive []Tetromino
        ↓
Calculate minimum candidate square size
        ↓
Generate all placements
        ↓
Build Exact Cover structure
        ↓
Run Algorithm X / DLX
        ↓
Solution found?
   ├── Yes → reconstruct board
   └── No  → increase square size
                     ↓
                  repeat

The first successful candidate size is the smallest possible square.

---

## 18. Rendering

The DLX solution contains selected placements.

The program creates an empty board filled with:

.

Then, for each selected placement:

1. Retrieve its tetromino.
2. Iterate through its normalized cells.
3. Convert local coordinates to board coordinates.
4. Write the tetromino letter into the board.

Example:

ABBBB.
ACCCEE
AFFCEE
A.FFGG
HHHDDG
.HDD.G

Each tetromino letter must appear exactly four times.

---

## 19. Architecture

The application uses a constraint-solving architecture.

CLI
 ↓
File Reader
 ↓
Parser
 ↓
Validator
 ↓
Normalized Tetromino Model
 ↓
Placement Generator
 ↓
Exact Cover Builder
 ↓
Algorithm X / DLX
 ↓
Board Reconstruction
 ↓
Renderer

The solver operates on possible placements rather than directly mutating a game board during recursive search.

This allows parsing and solving responsibilities to remain separate.

---

## 20. Package Layout

Recommended structure:

tetrisoptimizer/
├── go.mod
├── main.go
├── README.md
├── PRD.md
└── tetris/
    ├── run.go
    ├── types.go
    ├── parser.go
    ├── connectivity.go
    ├── normalize.go
    ├── placement.go
    ├── matrix.go
    ├── dlx.go
    ├── solver.go
    └── render.go

All files inside:

tetris/

belong to:

package tetris

The root `main.go` belongs to:

package main

---

## 21. Public Package API

The public API should remain small.

The main exported entry point is:

func Run(path string) error

Most implementation functions should remain unexported:

parseData
parseBlock
isConnected
isAdjacent
normalize
generatePlacements
buildMatrix
cover
uncover
search
solve
render

This keeps internal implementation details hidden from `main`.

---

## 22. Error Handling

Functions should return errors rather than printing them directly.

Example flow:

parseBlock
    ↓ error
parseData
    ↓ error
Run
    ↓ error
main
    ↓
print descriptive error

Internal errors should contain meaningful debugging information.

Examples:

input file is empty
tetromino has 3 occupied cells instead of 4
tetromino cells are not connected
row 2 has 5 columns instead of 4
invalid character '@'

If the argument count is not exactly one, `main` prints:

Please insert only one Argument.

For errors returned by `Run`, `main` prints the complete wrapped error. File-reading errors are prefixed with `Could not read file`, parsing errors are prefixed with `Could not parse file`, and solving errors are prefixed with `Could not solve tetrominoes`.

---

## 23. Testing Strategy

Testing should be divided by responsibility.

### Parser Tests

Test:

- empty input
- one valid tetromino
- multiple valid tetrominoes
- Windows line endings
- optional final newline
- malformed separators
- incorrect row count
- incorrect column count
- invalid characters
- fewer than four `#`
- more than four `#`

---

### Connectivity Tests

Valid:

####
....
....
....

##..
##..
....
....

.#..
###.
....
....

Invalid:

##..
....
..##
....

#...
#...
....
...#

---

### Normalization Tests

Input coordinates:

(2,1)
(3,1)
(2,2)
(3,2)

Expected:

(0,0)
(1,0)
(0,1)
(1,1)

Expected:

Width  = 2
Height = 2

---

### Placement Tests

Verify that placement generation:

- never exceeds board boundaries
- generates all valid anchor positions
- correctly converts local coordinates to board coordinates

---

### Exact Cover Tests

Verify that:

- every placement row refers to exactly one tetromino
- every placement row refers to exactly four board cells
- overlapping placements conflict through shared board-cell constraints
- every tetromino has at least one placement for a solvable candidate board

---

### DLX Tests

Test small known Exact Cover problems independently from tetromino logic.

Verify:

- cover removes expected nodes
- uncover restores the exact previous structure
- search finds a known solution
- search correctly reports no solution

---

### Integration Tests

Test complete input files.

Verify:

- output is square
- all tetrominoes are present
- each letter appears exactly four times
- no two tetrominoes overlap
- output uses the smallest square with a valid arrangement
- unused cells are printed as `.`
- invalid argument count prints `Please insert only one Argument.`
- unreadable files and invalid input print descriptive errors

---

## 24. Milestones

### Milestone 1 — Parser Complete

The program:

- reads a file
- separates tetromino blocks
- validates each block
- validates connectivity
- normalizes coordinates
- returns `[]Tetromino`

---

### Milestone 2 — Placement Generation

The program can:

- calculate candidate board sizes
- generate every possible placement for every piece
- associate each placement with its source tetromino

---

### Milestone 3 — Exact Cover Representation

The program can:

- create piece constraints
- create board-cell constraints
- convert placements into exact-cover rows
- map rows back to `Placement` values

---

### Milestone 4 — Dancing Links

The program can:

- build the linked DLX structure
- cover columns
- uncover columns
- recursively search using Algorithm X
- return selected placement rows

---

### Milestone 5 — Solver Integration

The solver:

- starts from the minimum possible square
- generates placements
- builds the DLX structure
- searches for a solution
- increases square size if necessary

---

### Milestone 6 — Rendering and Final Tests

The program:

- reconstructs the solved board
- writes tetromino letters
- leaves unused cells as `.`
- prints the board
- passes parser, solver, integration, and audit tests

---

## 25. Non-Goals

The program will not:

- rotate tetrominoes
- mirror tetrominoes
- modify the shape given in the input
- provide an interactive interface
- accept multiple input files
- use third-party Go packages
- provide graphical output
- attempt to repair malformed tetrominoes

Tetromino orientation is exactly the orientation given in the input file.

---

## 26. Risks / Open Questions

- Dancing Links is considerably more complex than ordinary recursive backtracking.
- `cover` and `uncover` must be exact inverses.
- Incorrect pointer restoration can silently corrupt the DLX structure.
- Exact Cover column design must correctly distinguish mandatory piece constraints from optional board-cell constraints.
- Placement rows must preserve enough metadata to reconstruct the final board.
- Large candidate boards may generate many placement rows.
- The allowed maximum number of tetrominoes should be verified against the school specification, especially because letters are limited to uppercase Latin characters.

---

## 27. Design Decisions

### Exact Cover Instead of Direct Board Backtracking

We choose:

Knuth's Algorithm X with Dancing Links

instead of directly placing and removing tetrominoes on a board.

Reason:

The problem can naturally be represented as selecting one legal placement for every tetromino while preventing board-cell collisions.

Benefits:

- clean constraint representation
- efficient removal and restoration of conflicting choices
- strong separation between placement generation and search
- useful opportunity to learn Algorithm X and Dancing Links

Tradeoffs:

- significantly greater implementation complexity
- more data structures
- harder debugging
- additional mapping required between matrix rows and board placements

---

### Fixed Tetromino Orientation

We do not generate rotations or reflections.

Reason:

The input defines the orientation of each tetromino and the project asks the program to assemble the given tetrominoes.

Changing orientation would alter the provided pieces.

---

### Normalized Coordinate Representation

We store only four occupied coordinates rather than retaining the original 4×4 text representation.

Reason:

The solver only needs occupied cells.

Benefits:

- simple placement calculations
- lower storage
- easier normalization
- easier board reconstruction
- easier generation of Exact Cover rows

---

### Standard Library Only

Only Go standard-library packages will be used.

Likely packages include:

os
fmt
strings

Additional standard-library packages may be used if needed.

No external dependencies are required.

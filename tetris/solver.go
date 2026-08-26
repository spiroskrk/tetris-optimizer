package tetris

import "fmt"

func generatePlacements(pieces []Tetromino, size int) []Placement {
	placements := make([]Placement, 0)

	for pieceIndex, piece := range pieces {
		maxY := size - piece.Height
		maxX := size - piece.Width

		for y := 0; y <= maxY; y++ {
			for x := 0; x <= maxX; x++ {
				placements = append(placements, Placement{PieceIndex: pieceIndex, X: x, Y: y})
			}
		}
	}
	return placements
}

func buildExactRows(
	pieces []Tetromino,
	placements []Placement,
	size int,
) []ExactRow {
	rows := make([]ExactRow, 0, len(placements))

	for placementIndex, placement := range placements {
		piece := pieces[placement.PieceIndex]

		var columns [5]int

		//First constraint: this piece is used
		columns[0] = placement.PieceIndex

		//Remaining 4 constraints: occupied board cells
		for i, cell := range piece.Cells {
			boardX := placement.X + cell.X
			boardY := placement.Y + cell.Y

			boardConstraint := len(pieces) + boardY*size + boardX

			columns[i+1] = boardConstraint
		}
		rows = append(rows, ExactRow{
			PlacementIndex: placementIndex,
			Columns:        columns,
		})

	}

	return rows
}

func buildDLX(
	rows []ExactRow,
	columnCount int,
	primaryCount int,
) *DLX {
	dlx := &DLX{
		Root:    &Column{Name: -1},
		Columns: make([]*Column, columnCount),
	}

	root := &dlx.Root.Node

	root.Left = root
	root.Right = root

	previousPrimary := root

	for i := 0; i < columnCount; i++ {
		column := &Column{
			Name: i,
		}

		column.Column = column

		// Every column needs a vertical circular list.
		column.Up = &column.Node
		column.Down = &column.Node

		if i < primaryCount {
			// Primary column: link it into Root's horizontal list.
			column.Left = previousPrimary
			column.Right = root

			previousPrimary.Right = &column.Node
			root.Left = &column.Node

			previousPrimary = &column.Node
		} else {
			// Secondary column: do NOT link it to Root.
			column.Left = &column.Node
			column.Right = &column.Node
		}

		dlx.Columns[i] = column
	}

	// Build the placement rows.
	for rowIndex, row := range rows {
		var first *Node
		var previous *Node

		for _, columnIndex := range row.Columns {
			column := dlx.Columns[columnIndex]

			node := &Node{
				Column:   column,
				RowIndex: rowIndex,
			}

			// Insert vertically at bottom of column.
			node.Down = &column.Node
			node.Up = column.Up

			column.Up.Down = node
			column.Up = node

			column.Size++

			// Link horizontally inside the placement row.
			if first == nil {
				first = node
				node.Left = node
				node.Right = node
			} else {
				node.Left = previous
				node.Right = first

				previous.Right = node
				first.Left = node
			}

			previous = node
		}
	}

	return dlx
}

func cover(column *Column) {
	// Remove the column header from the horizontal list.
	column.Left.Right = column.Right
	column.Right.Left = column.Left

	// Visit every placement that uses this constraint.
	for row := column.Down; row != &column.Node; row = row.Down {

		// Visit every other constraint used by that placement.
		for node := row.Right; node != row; node = node.Right {

			// Remove this node from its vertical column.
			node.Up.Down = node.Down
			node.Down.Up = node.Up

			node.Column.Size--
		}
	}
}

func uncover(column *Column) {
	// Reverse cover() in the opposite order.
	for row := column.Up; row != &column.Node; row = row.Up {

		for node := row.Left; node != row; node = node.Left {

			node.Column.Size++

			// Restore this node into its vertical column.
			node.Up.Down = node
			node.Down.Up = node
		}
	}

	// Restore the column header to the horizontal list.
	column.Left.Right = &column.Node
	column.Right.Left = &column.Node
}

func search(dlx *DLX, solution *[]int) bool {
	if dlx.Root.Right == &dlx.Root.Node {
		return true
	}

	column := chooseColumn(dlx)

	cover(column)

	for row := column.Down; row != &column.Node; row = row.Down {
		*solution = append(*solution, row.RowIndex)

		for node := row.Right; node != row; node = node.Right {
			cover(node.Column)
		}

		if search(dlx, solution) {
			return true
		}

		for node := row.Left; node != row; node = node.Left {
			uncover(node.Column)
		}

		*solution = (*solution)[:len(*solution)-1]
	}

	uncover(column)

	return false
}

func chooseColumn(dlx *DLX) *Column {
	var best *Column

	for node := dlx.Root.Right; node != &dlx.Root.Node; node = node.Right {
		column := node.Column

		if column.Size == 0 {
			return column
		}

		if best == nil || column.Size < best.Size {
			best = column
		}
	}

	return best
}

func solve(pieces []Tetromino) ([][]byte, error) {
	if len(pieces) == 0 {
		return nil, fmt.Errorf("cannot solve an empty set of tetrominoes")
	}

	totalArea := len(pieces) * 4
	size := 1
	for size*size < totalArea {
		size++
	}

	for _, piece := range pieces {
		if piece.Width > size {
			size = piece.Width
		}
		if piece.Height > size {
			size = piece.Height
		}
	}

	for {
		placements := generatePlacements(pieces, size)
		rows := buildExactRows(pieces, placements, size)
		columnCount := len(pieces) + size*size
		dlx := buildDLX(rows, columnCount, len(pieces))
		solution := make([]int, 0, len(pieces))

		if !search(dlx, &solution) {
			size++
			continue
		}

		board := make([][]byte, size)
		for y := range board {
			board[y] = make([]byte, size)
			for x := range board[y] {
				board[y][x] = '.'
			}
		}

		for _, rowIndex := range solution {
			if rowIndex < 0 || rowIndex >= len(rows) {
				return nil, fmt.Errorf("solver returned invalid row index %d", rowIndex)
			}

			placementIndex := rows[rowIndex].PlacementIndex
			if placementIndex < 0 || placementIndex >= len(placements) {
				return nil, fmt.Errorf("solver returned invalid placement index %d", placementIndex)
			}

			placement := placements[placementIndex]
			piece := pieces[placement.PieceIndex]

			for _, cell := range piece.Cells {
				x := placement.X + cell.X
				y := placement.Y + cell.Y

				if board[y][x] != '.' {
					return nil, fmt.Errorf("solver produced overlapping placements at (%d,%d)", x, y)
				}
				board[y][x] = piece.Letter
			}
		}

		return board, nil
	}
}

package tetris

import "testing"

func squarePiece(letter byte) Tetromino {
	return Tetromino{
		Cells: [4]Point{
			{X: 0, Y: 0},
			{X: 1, Y: 0},
			{X: 0, Y: 1},
			{X: 1, Y: 1},
		},
		Letter: letter,
		Width:  2,
		Height: 2,
	}
}

func verticalPiece(letter byte) Tetromino {
	return Tetromino{
		Cells: [4]Point{
			{X: 0, Y: 0},
			{X: 0, Y: 1},
			{X: 0, Y: 2},
			{X: 0, Y: 3},
		},
		Letter: letter,
		Width:  1,
		Height: 4,
	}
}

func TestGeneratePlacements(t *testing.T) {
	pieces := []Tetromino{squarePiece('A')}

	got := generatePlacements(pieces, 3)
	want := []Placement{
		{PieceIndex: 0, X: 0, Y: 0},
		{PieceIndex: 0, X: 1, Y: 0},
		{PieceIndex: 0, X: 0, Y: 1},
		{PieceIndex: 0, X: 1, Y: 1},
	}

	if len(got) != len(want) {
		t.Fatalf("generatePlacements returned %d placements; want %d", len(got), len(want))
	}

	for i := range want {
		if got[i] != want[i] {
			t.Errorf("placement %d = %+v; want %+v", i, got[i], want[i])
		}
	}
}

func TestGeneratePlacementsSkipsPieceThatDoesNotFit(t *testing.T) {
	got := generatePlacements([]Tetromino{verticalPiece('A')}, 3)
	if len(got) != 0 {
		t.Fatalf("generatePlacements returned %d placements; want 0", len(got))
	}
}

func TestBuildExactRows(t *testing.T) {
	pieces := []Tetromino{squarePiece('A')}
	placements := []Placement{{PieceIndex: 0, X: 1, Y: 0}}

	rows := buildExactRows(pieces, placements, 3)
	if len(rows) != 1 {
		t.Fatalf("buildExactRows returned %d rows; want 1", len(rows))
	}

	want := ExactRow{
		PlacementIndex: 0,
		Columns:        [5]int{0, 2, 3, 5, 6},
	}
	if rows[0] != want {
		t.Fatalf("row = %+v; want %+v", rows[0], want)
	}
}

func TestBuildDLXLinksPrimaryAndSecondaryColumns(t *testing.T) {
	rows := []ExactRow{
		{PlacementIndex: 0, Columns: [5]int{0, 2, 3, 4, 5}},
		{PlacementIndex: 1, Columns: [5]int{1, 6, 7, 8, 9}},
		{PlacementIndex: 2, Columns: [5]int{1, 2, 10, 11, 12}},
	}
	dlx := buildDLX(rows, 13, 2)
	root := &dlx.Root.Node

	if root.Right != &dlx.Columns[0].Node || root.Left != &dlx.Columns[1].Node {
		t.Fatal("root is not linked to the first and last primary columns")
	}
	if dlx.Columns[0].Right != &dlx.Columns[1].Node || dlx.Columns[1].Right != root {
		t.Fatal("primary columns do not form the expected circular list")
	}

	for i := 2; i < len(dlx.Columns); i++ {
		column := dlx.Columns[i]
		if column.Left != &column.Node || column.Right != &column.Node {
			t.Errorf("secondary column %d is linked into another horizontal list", i)
		}
	}

	if dlx.Columns[0].Size != 1 || dlx.Columns[1].Size != 2 || dlx.Columns[2].Size != 2 {
		t.Fatalf(
			"unexpected column sizes: column 0=%d, column 1=%d, column 2=%d",
			dlx.Columns[0].Size,
			dlx.Columns[1].Size,
			dlx.Columns[2].Size,
		)
	}

	row := dlx.Columns[0].Down
	nodeCount := 1
	for node := row.Right; node != row; node = node.Right {
		nodeCount++
		if nodeCount > 5 {
			t.Fatal("row is not a five-node circular list")
		}
	}
	if nodeCount != 5 {
		t.Fatalf("row contains %d nodes; want 5", nodeCount)
	}
}

func TestCoverAndUncoverRestoreStructure(t *testing.T) {
	rows := []ExactRow{
		{PlacementIndex: 0, Columns: [5]int{0, 2, 3, 4, 5}},
		{PlacementIndex: 1, Columns: [5]int{1, 2, 6, 7, 8}},
		{PlacementIndex: 2, Columns: [5]int{1, 9, 10, 11, 12}},
	}
	dlx := buildDLX(rows, 13, 2)
	root := &dlx.Root.Node

	cover(dlx.Columns[0])

	if root.Right != &dlx.Columns[1].Node || root.Left != &dlx.Columns[1].Node {
		t.Fatal("cover did not remove primary column 0 from the root list")
	}
	if dlx.Columns[2].Size != 1 {
		t.Fatalf("shared secondary column size after cover = %d; want 1", dlx.Columns[2].Size)
	}
	if dlx.Columns[1].Size != 2 {
		t.Fatalf("unselected primary column size after first cover = %d; want 2", dlx.Columns[1].Size)
	}

	cover(dlx.Columns[2])
	if dlx.Columns[1].Size != 1 {
		t.Fatalf("conflicting primary column size after secondary cover = %d; want 1", dlx.Columns[1].Size)
	}

	uncover(dlx.Columns[2])
	uncover(dlx.Columns[0])

	if root.Right != &dlx.Columns[0].Node || root.Left != &dlx.Columns[1].Node {
		t.Fatal("uncover did not restore primary column 0 to the root list")
	}
	if dlx.Columns[2].Size != 2 || dlx.Columns[1].Size != 2 {
		t.Fatalf(
			"uncover did not restore sizes: secondary=%d primary=%d",
			dlx.Columns[2].Size,
			dlx.Columns[1].Size,
		)
	}
}

func TestChooseColumnUsesSmallestPrimaryColumn(t *testing.T) {
	rows := []ExactRow{
		{PlacementIndex: 0, Columns: [5]int{0, 2, 3, 4, 5}},
		{PlacementIndex: 1, Columns: [5]int{0, 6, 7, 8, 9}},
		{PlacementIndex: 2, Columns: [5]int{1, 10, 11, 12, 13}},
	}
	dlx := buildDLX(rows, 14, 2)

	got := chooseColumn(dlx)
	if got != dlx.Columns[1] {
		t.Fatalf("chooseColumn selected column %d; want column 1", got.Name)
	}
}

func TestChooseColumnReturnsEmptyPrimaryColumnImmediately(t *testing.T) {
	rows := []ExactRow{
		{PlacementIndex: 0, Columns: [5]int{0, 2, 3, 4, 5}},
	}
	dlx := buildDLX(rows, 6, 2)

	got := chooseColumn(dlx)
	if got != dlx.Columns[1] {
		t.Fatalf("chooseColumn selected column %d; want empty column 1", got.Name)
	}
}

func TestSearchFindsCompatibleRows(t *testing.T) {
	rows := []ExactRow{
		{PlacementIndex: 0, Columns: [5]int{0, 2, 3, 4, 5}},
		{PlacementIndex: 1, Columns: [5]int{1, 6, 7, 8, 9}},
		{PlacementIndex: 2, Columns: [5]int{1, 2, 10, 11, 12}},
	}
	dlx := buildDLX(rows, 13, 2)
	var solution []int

	if !search(dlx, &solution) {
		t.Fatal("search did not find an existing exact-cover solution")
	}
	if len(solution) != 2 || solution[0] != 0 || solution[1] != 1 {
		t.Fatalf("search solution = %v; want [0 1]", solution)
	}
}

func TestSearchRejectsConflictingRowsAndRestoresDLX(t *testing.T) {
	rows := []ExactRow{
		{PlacementIndex: 0, Columns: [5]int{0, 2, 3, 4, 5}},
		{PlacementIndex: 1, Columns: [5]int{1, 2, 6, 7, 8}},
	}
	dlx := buildDLX(rows, 9, 2)
	var solution []int

	if search(dlx, &solution) {
		t.Fatal("search accepted rows that conflict through a secondary column")
	}
	if len(solution) != 0 {
		t.Fatalf("failed search left solution entries behind: %v", solution)
	}
	if dlx.Root.Right != &dlx.Columns[0].Node || dlx.Root.Left != &dlx.Columns[1].Node {
		t.Fatal("failed search did not restore the root column list")
	}
	if dlx.Columns[0].Size != 1 || dlx.Columns[1].Size != 1 || dlx.Columns[2].Size != 2 {
		t.Fatal("failed search did not restore column sizes")
	}
}

func TestSolveRejectsEmptyPieceSet(t *testing.T) {
	if _, err := solve(nil); err == nil {
		t.Fatal("solve(nil) returned no error")
	}
}

func TestSolveOneSquarePiece(t *testing.T) {
	board, err := solve([]Tetromino{squarePiece('A')})
	if err != nil {
		t.Fatalf("solve returned an error: %v", err)
	}

	want := []string{"AA", "AA"}
	assertBoard(t, board, want)
}

func TestSolveIncreasesCandidateSize(t *testing.T) {
	pieces := []Tetromino{squarePiece('A'), squarePiece('B')}

	board, err := solve(pieces)
	if err != nil {
		t.Fatalf("solve returned an error: %v", err)
	}
	if len(board) != 4 {
		t.Fatalf("board size = %d; want 4 because two squares cannot fit in 3x3", len(board))
	}
	assertBoardContents(t, board, pieces)
}

func TestSolveAccountsForPieceDimensions(t *testing.T) {
	pieces := []Tetromino{verticalPiece('A')}

	board, err := solve(pieces)
	if err != nil {
		t.Fatalf("solve returned an error: %v", err)
	}
	if len(board) != 4 {
		t.Fatalf("board size = %d; want 4", len(board))
	}
	assertBoardContents(t, board, pieces)
}

func assertBoard(t *testing.T, board [][]byte, want []string) {
	t.Helper()

	if len(board) != len(want) {
		t.Fatalf("board has %d rows; want %d", len(board), len(want))
	}
	for y := range want {
		if string(board[y]) != want[y] {
			t.Errorf("board row %d = %q; want %q", y, board[y], want[y])
		}
	}
}

func assertBoardContents(t *testing.T, board [][]byte, pieces []Tetromino) {
	t.Helper()

	counts := make(map[byte]int)
	positions := make(map[byte][]Point)
	knownLetters := make(map[byte]bool, len(pieces))
	for _, piece := range pieces {
		knownLetters[piece.Letter] = true
	}

	for y, row := range board {
		if len(row) != len(board) {
			t.Fatalf("board row %d has width %d; want %d", y, len(row), len(board))
		}
		for x, cell := range row {
			if cell != '.' {
				if !knownLetters[cell] {
					t.Errorf("board contains unexpected character %q at (%d,%d)", cell, x, y)
					continue
				}
				counts[cell]++
				positions[cell] = append(positions[cell], Point{X: x, Y: y})
			}
		}
	}

	for _, piece := range pieces {
		if counts[piece.Letter] != 4 {
			t.Errorf("piece %q occupies %d cells; want 4", piece.Letter, counts[piece.Letter])
			continue
		}

		actualCells := [4]Point{}
		copy(actualCells[:], positions[piece.Letter])
		actualCells, actualWidth, actualHeight := normalize(actualCells)
		if actualWidth != piece.Width || actualHeight != piece.Height {
			t.Errorf(
				"piece %q dimensions = %dx%d; want %dx%d",
				piece.Letter,
				actualWidth,
				actualHeight,
				piece.Width,
				piece.Height,
			)
		}

		actualSet := make(map[Point]bool, 4)
		for _, cell := range actualCells {
			actualSet[cell] = true
		}
		for _, expectedCell := range piece.Cells {
			if !actualSet[expectedCell] {
				t.Errorf("piece %q does not preserve its input shape", piece.Letter)
				break
			}
		}
	}
}

package tetris

type Tetromino struct {
	Cells  [4]Point
	Letter byte
	Width  int
	Height int
}

type Point struct {
	X int
	Y int
}

type Placement struct {
	PieceIndex int
	X          int
	Y          int
}

type ExactRow struct {
	PlacementIndex int
	Columns        [5]int
}

type Node struct {
	Left     *Node
	Right    *Node
	Up       *Node
	Down     *Node
	Column   *Column
	RowIndex int
}

type Column struct {
	Node
	Size int
	Name int
}

type DLX struct {
	Root    *Column
	Columns []*Column
}

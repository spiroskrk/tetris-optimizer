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

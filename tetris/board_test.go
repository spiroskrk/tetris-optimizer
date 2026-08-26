package tetris

import "testing"

func TestRenderBoard(t *testing.T) {
	board := [][]byte{
		{'A', 'A', '.'},
		{'.', 'B', '.'},
		{'.', 'B', 'B'},
	}

	got := renderBoard(board)
	want := "AA.\n.B.\n.BB\n"

	if got != want {
		t.Fatalf("renderBoard() = %q; want %q", got, want)
	}
}

func TestRenderEmptyBoard(t *testing.T) {
	if got := renderBoard(nil); got != "" {
		t.Fatalf("renderBoard(nil) = %q; want an empty string", got)
	}
}

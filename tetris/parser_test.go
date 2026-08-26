package tetris

import (
	"strings"
	"testing"
)

func TestAbs(t *testing.T) {
	tests := []struct {
		name string
		in   int
		want int
	}{
		{name: "positive", in: 7, want: 7},
		{name: "negative", in: -7, want: 7},
		{name: "zero", in: 0, want: 0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := abs(test.in); got != test.want {
				t.Fatalf("abs(%d) = %d; want %d", test.in, got, test.want)
			}
		})
	}
}

func TestIsAdjacent(t *testing.T) {
	tests := []struct {
		name string
		a    Point
		b    Point
		want bool
	}{
		{name: "horizontal neighbours", a: Point{X: 1, Y: 1}, b: Point{X: 2, Y: 1}, want: true},
		{name: "vertical neighbours", a: Point{X: 1, Y: 1}, b: Point{X: 1, Y: 2}, want: true},
		{name: "diagonal cells", a: Point{X: 1, Y: 1}, b: Point{X: 2, Y: 2}, want: false},
		{name: "same cell", a: Point{X: 1, Y: 1}, b: Point{X: 1, Y: 1}, want: false},
		{name: "separated cells", a: Point{X: 0, Y: 0}, b: Point{X: 0, Y: 2}, want: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := isAdjacent(test.a, test.b); got != test.want {
				t.Fatalf("isAdjacent(%+v, %+v) = %t; want %t", test.a, test.b, got, test.want)
			}
		})
	}
}

func TestIsConnected(t *testing.T) {
	tests := []struct {
		name  string
		cells [4]Point
		want  bool
	}{
		{
			name:  "straight line",
			cells: [4]Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}, {X: 3, Y: 0}},
			want:  true,
		},
		{
			name:  "connected through intermediate cells",
			cells: [4]Point{{X: 1, Y: 1}, {X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}},
			want:  true,
		},
		{
			name:  "two disconnected pairs",
			cells: [4]Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 2}, {X: 3, Y: 2}},
			want:  false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := isConnected(test.cells); got != test.want {
				t.Fatalf("isConnected(%+v) = %t; want %t", test.cells, got, test.want)
			}
		})
	}
}

func TestNormalize(t *testing.T) {
	cells := [4]Point{
		{X: 2, Y: 1},
		{X: 3, Y: 1},
		{X: 2, Y: 2},
		{X: 3, Y: 2},
	}
	wantCells := [4]Point{
		{X: 0, Y: 0},
		{X: 1, Y: 0},
		{X: 0, Y: 1},
		{X: 1, Y: 1},
	}

	gotCells, gotWidth, gotHeight := normalize(cells)
	if gotCells != wantCells {
		t.Errorf("normalized cells = %+v; want %+v", gotCells, wantCells)
	}
	if gotWidth != 2 || gotHeight != 2 {
		t.Errorf("normalized dimensions = %dx%d; want 2x2", gotWidth, gotHeight)
	}
}

func TestNormalizeFindsMinimumWhenFirstCellIsNotTopLeft(t *testing.T) {
	cells := [4]Point{
		{X: 3, Y: 2},
		{X: 2, Y: 1},
		{X: 3, Y: 1},
		{X: 2, Y: 2},
	}
	wantCells := [4]Point{
		{X: 1, Y: 1},
		{X: 0, Y: 0},
		{X: 1, Y: 0},
		{X: 0, Y: 1},
	}

	gotCells, gotWidth, gotHeight := normalize(cells)
	if gotCells != wantCells {
		t.Errorf("normalized cells = %+v; want %+v", gotCells, wantCells)
	}
	if gotWidth != 2 || gotHeight != 2 {
		t.Errorf("normalized dimensions = %dx%d; want 2x2", gotWidth, gotHeight)
	}
}

func TestParseBlockNormalizesTetromino(t *testing.T) {
	block := "....\n..##\n..##\n...."

	piece, err := parseBlock(block)
	if err != nil {
		t.Fatalf("parseBlock returned an error: %v", err)
	}

	wantCells := [4]Point{
		{X: 0, Y: 0},
		{X: 1, Y: 0},
		{X: 0, Y: 1},
		{X: 1, Y: 1},
	}
	if piece.Cells != wantCells {
		t.Errorf("piece cells = %+v; want %+v", piece.Cells, wantCells)
	}
	if piece.Width != 2 || piece.Height != 2 {
		t.Errorf("piece dimensions = %dx%d; want 2x2", piece.Width, piece.Height)
	}
}

func TestParseBlockRejectsInvalidInput(t *testing.T) {
	tests := []struct {
		name        string
		block       string
		wantErrPart string
	}{
		{
			name:        "too few rows",
			block:       "##..\n##..\n....",
			wantErrPart: "3 lines instead of 4",
		},
		{
			name:        "incorrect row width",
			block:       "##...\n##..\n....\n....",
			wantErrPart: "5 columns instead of 4",
		},
		{
			name:        "invalid character",
			block:       "##x.\n##..\n....\n....",
			wantErrPart: "character 'x' is invalid",
		},
		{
			name:        "too few occupied cells",
			block:       "##..\n#...\n....\n....",
			wantErrPart: "3 occupied cells instead of 4",
		},
		{
			name:        "too many occupied cells",
			block:       "###.\n##..\n....\n....",
			wantErrPart: "more than 4 occupied cells",
		},
		{
			name:        "disconnected cells",
			block:       "##..\n....\n..##\n....",
			wantErrPart: "cells are not connected",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := parseBlock(test.block)
			if err == nil {
				t.Fatal("parseBlock returned no error")
			}
			if !strings.Contains(err.Error(), test.wantErrPart) {
				t.Fatalf("parseBlock error = %q; want it to contain %q", err, test.wantErrPart)
			}
		})
	}
}

func TestParseDataAcceptsMultiplePiecesAndAssignsLetters(t *testing.T) {
	data := "##..\n##..\n....\n....\n\n#...\n#...\n#...\n#...\n"

	pieces, err := parseData(data)
	if err != nil {
		t.Fatalf("parseData returned an error: %v", err)
	}
	if len(pieces) != 2 {
		t.Fatalf("parseData returned %d pieces; want 2", len(pieces))
	}
	if pieces[0].Letter != 'A' || pieces[1].Letter != 'B' {
		t.Fatalf("piece letters = %q and %q; want 'A' and 'B'", pieces[0].Letter, pieces[1].Letter)
	}
}

func TestParseDataAcceptsWindowsLineEndings(t *testing.T) {
	data := "##..\r\n##..\r\n....\r\n....\r\n"

	pieces, err := parseData(data)
	if err != nil {
		t.Fatalf("parseData returned an error: %v", err)
	}
	if len(pieces) != 1 {
		t.Fatalf("parseData returned %d pieces; want 1", len(pieces))
	}
}

func TestParseDataRejectsMalformedDocument(t *testing.T) {
	tests := []struct {
		name string
		data string
	}{
		{name: "empty input", data: ""},
		{name: "only optional newline", data: "\n"},
		{name: "too many final newlines", data: "##..\n##..\n....\n....\n\n"},
		{name: "too many separator newlines", data: "##..\n##..\n....\n....\n\n\n##..\n##..\n....\n...."},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := parseData(test.data); err == nil {
				t.Fatal("parseData returned no error")
			}
		})
	}
}

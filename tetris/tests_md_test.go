package tetris

import (
	"strings"
	"testing"
)

const (
	squareBlock         = "....\n.##.\n.##.\n...."
	verticalTBlock      = ".#..\n.##.\n.#..\n...."
	rightSkewBlock      = "....\n..##\n.##.\n...."
	leftSkewBlock       = "....\n....\n##..\n.##."
	verticalSBlock      = "....\n..#.\n.##.\n.#.."
	rightLBlock         = ".###\n...#\n....\n...."
	leftLBlock          = "##..\n.#..\n.#..\n...."
	verticalZBlock      = ".#..\n.##.\n..#.\n...."
	horizontalTBlock    = "....\n###.\n.#..\n...."
	verticalLineBlock   = "...#\n...#\n...#\n...#"
	horizontalLineBlock = "....\n....\n....\n####"
)

func TestExamplesFromTestsDocument(t *testing.T) {
	tests := []struct {
		name      string
		blocks    []string
		boardSize int
	}{
		{
			name: "difficult test",
			blocks: []string{
				squareBlock,
				verticalTBlock,
				rightSkewBlock,
				squareBlock,
				verticalSBlock,
				rightLBlock,
				leftLBlock,
				squareBlock,
				rightSkewBlock,
				leftLBlock,
				verticalZBlock,
				horizontalTBlock,
			},
			boardSize: 7,
		},
		{
			name:      "good example",
			blocks:    []string{squareBlock},
			boardSize: 2,
		},
		{
			name: "good example 2",
			blocks: []string{
				verticalLineBlock,
				horizontalLineBlock,
				rightLBlock,
				rightSkewBlock,
			},
			boardSize: 5,
		},
		{
			name: "good example 3",
			blocks: []string{
				verticalLineBlock,
				horizontalLineBlock,
				rightLBlock,
				rightSkewBlock,
				squareBlock,
				leftSkewBlock,
				leftLBlock,
				horizontalTBlock,
			},
			boardSize: 6,
		},
		{
			name: "good example 4",
			blocks: []string{
				squareBlock,
				verticalLineBlock,
				rightSkewBlock,
				squareBlock,
				verticalSBlock,
				rightLBlock,
				leftLBlock,
				rightSkewBlock,
				leftLBlock,
				verticalZBlock,
				horizontalTBlock,
			},
			boardSize: 7,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pieces, err := parseData(strings.Join(test.blocks, "\n\n"))
			if err != nil {
				t.Fatalf("parseData returned an error: %v", err)
			}

			board, err := solve(pieces)
			if err != nil {
				t.Fatalf("solve returned an error: %v", err)
			}
			if len(board) != test.boardSize {
				t.Fatalf("board size = %d; want %d", len(board), test.boardSize)
			}
			assertBoardContents(t, board, pieces)
		})
	}
}

func TestBadTetrominoDocuments(t *testing.T) {
	valid := squareBlock
	tests := []struct {
		name string
		data string
	}{
		{
			name: "disconnected tetromino after valid tetromino",
			data: valid + "\n\n" + "##..\n....\n..##\n....",
		},
		{
			name: "tetromino with five cells",
			data: "###.\n##..\n....\n....",
		},
		{
			name: "tetromino with three cells",
			data: "###.\n....\n....\n....",
		},
		{
			name: "diagonal-only connection",
			data: "#...\n.#..\n..#.\n...#",
		},
		{
			name: "invalid symbol in later tetromino",
			data: valid + "\n\n" + "##..\n#@..\n....\n....",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := parseData(test.data); err == nil {
				t.Fatal("parseData accepted a bad tetromino document")
			}
		})
	}
}

func TestParseDataRejectsMoreThanTwentySixTetrominoes(t *testing.T) {
	data := strings.Join(makeRepeatedBlocks(squareBlock, 27), "\n\n")

	_, err := parseData(data)
	if err == nil {
		t.Fatal("parseData accepted 27 tetrominoes")
	}
	if !strings.Contains(err.Error(), "26") {
		t.Fatalf("parseData error = %q; want it to mention the 26-piece limit", err)
	}
}

func TestParseDataAcceptsTwentySixTetrominoes(t *testing.T) {
	data := strings.Join(makeRepeatedBlocks(squareBlock, 26), "\n\n")

	pieces, err := parseData(data)
	if err != nil {
		t.Fatalf("parseData returned an error for 26 tetrominoes: %v", err)
	}
	if len(pieces) != 26 {
		t.Fatalf("parseData returned %d pieces; want 26", len(pieces))
	}
	if pieces[25].Letter != 'Z' {
		t.Fatalf("last piece letter = %q; want 'Z'", pieces[25].Letter)
	}
}

func makeRepeatedBlocks(block string, count int) []string {
	blocks := make([]string, count)
	for i := range blocks {
		blocks[i] = block
	}
	return blocks
}

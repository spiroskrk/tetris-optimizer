package tetris

import (
	"fmt"
	"strings"
)

func ParseData(data string) ([]Tetromino, error) {
	data = strings.ReplaceAll(data, "\r\n", "\n")
	data = strings.TrimSuffix(data, "\n")
	if data == "" {
		return nil, fmt.Errorf("Input file is empty")
	}

	blocks := strings.Split(data, "\n\n")

	pieces := make([]Tetromino, 0, len(blocks))
	for i, block := range blocks {
		piece, err := parseBlock(block)
		if err != nil {
			return nil, err
		}

		piece.Letter = byte('A' + i)
		pieces.append(pieces, piece)
	}
	return pieces, nil
}

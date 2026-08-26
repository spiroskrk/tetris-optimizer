package tetris

import "strings"

func renderBoard(board [][]byte) string {
	var output strings.Builder

	for _, row := range board {
		output.Write(row)
		output.WriteByte('\n')
	}

	return output.String()
}

package tetris

import (
	"fmt"
	"os"
)

func Run(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Could not read file %q: %w", path, err)
	}
	fileData := string(file)

	tetrominoes, err := parseData(fileData)
	if err != nil {
		return fmt.Errorf("Could not parse file %w", err)
	}

	board, err := solve(tetrominoes)
	if err != nil {
		return fmt.Errorf("Could not solve tetrominoes %w", err)
	}

	fmt.Print(renderBoard(board))
	return nil
}

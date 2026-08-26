package tetris

import (
	"fmt"
	"strings"
)

func isConnected(cells [4]Point) bool {
	visited := [4]bool{}

	var visit func(index int)

	visit = func(index int) {
		if visited[index] {
			return
		}

		visited[index] = true

		for i := range cells {
			if !visited[i] && isAdjacent(cells[index], cells[i]) {
				visit(i)
			}
		}
	}

	visit(0)

	for _, wasVisited := range visited {
		if !wasVisited {
			return false
		}
	}

	return true
}

func isAdjacent(a, b Point) bool {
	sameX := a.X == b.X
	sameY := a.Y == b.Y

	if sameX && abs(a.Y-b.Y) == 1 {
		return true
	}

	if sameY && abs(a.X-b.X) == 1 {
		return true
	}

	return false
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func normalize(cells [4]Point) ([4]Point, int, int) {
	minY := cells[0].Y
	minX := cells[0].X

	for _, cell := range cells {
		if cell.X < minX {
			minX = cell.X
		}
		if cell.Y < minY {
			minY = cell.Y
		}
	}
	for i := range cells {
		cells[i].X = cells[i].X - minX
		cells[i].Y = cells[i].Y - minY
	}

	maxY := cells[0].Y
	maxX := cells[0].X

	for _, cell := range cells {
		if cell.X > maxX {
			maxX = cell.X
		}
		if cell.Y > maxY {
			maxY = cell.Y
		}
	}
	return cells, maxX + 1, maxY + 1
}

func parseBlock(block string) (Tetromino, error) {
	lines := strings.Split(block, "\n")
	if len(lines) != 4 {
		return Tetromino{}, fmt.Errorf("not a valid tetromino block. has %d lines instead of 4.", len(lines))
	}

	var cells [4]Point
	cellCount := 0
	for y, line := range lines {
		if len(line) != 4 {
			return Tetromino{}, fmt.Errorf("not a valid tetromino block. row %d has %d columns instead of 4.", y+1, len(line))
		}
		for x, ch := range line {
			switch ch {
			case '.':
				//nothing to store. continue
			case '#':
				if cellCount >= 4 {
					return Tetromino{}, fmt.Errorf("tetromino has more than 4 occupied cells")
				}
				cells[cellCount] = Point{X: x, Y: y}
				cellCount++
			default:
				return Tetromino{}, fmt.Errorf("character %q is invalid. only # or . are accepted in a block", ch)
			}
		}
	}
	if cellCount != 4 {
		return Tetromino{}, fmt.Errorf("tetromino has %d occupied cells instead of 4", cellCount)
	}
	if !isConnected(cells) {
		return Tetromino{}, fmt.Errorf("tetromino cells are not connected")
	}
	cells, width, height := normalize(cells)

	tetromino := Tetromino{
		Cells:  cells,
		Width:  width,
		Height: height,
	}

	return tetromino, nil
}

func parseData(data string) ([]Tetromino, error) {
	data = strings.ReplaceAll(data, "\r\n", "\n")
	data = strings.TrimSuffix(data, "\n")
	if data == "" {
		return nil, fmt.Errorf("input file is empty")
	}

	blocks := strings.Split(data, "\n\n")
	if len(blocks) > 26 {
		return nil, fmt.Errorf("input contains %d tetrominoes; maximum is 26", len(blocks))
	}

	pieces := make([]Tetromino, 0, len(blocks))
	for i, block := range blocks {
		piece, err := parseBlock(block)
		if err != nil {
			return nil, err
		}

		piece.Letter = byte('A' + i)
		pieces = append(pieces, piece)
	}
	return pieces, nil
}

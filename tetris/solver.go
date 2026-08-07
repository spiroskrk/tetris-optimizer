package tetris

func generatePlacements(pieces []Tetromino, size int) []Placement {
	placements := make([]Placement, 0)

	for pieceIndex, piece := range pieces {
		maxY := size - piece.Height
		maxX := size - piece.Width

		for y := 0; y <= maxY; y++ {
			for x := 0; x <= maxX; x++ {
				placements = append(placements, Placement{PieceIndex: pieceIndex, X: x, Y: y})
			}
		}
	}
	return placements
}

// CALCULATE BOARD SIZE - USE CODE LATER
// totalArea := len(pieces)*4
// dimention := 2
// for dimention*dimention <= totalArea {
// 	dimention++
// }

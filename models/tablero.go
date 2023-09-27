package models

type Board struct {
	width, height int
	area          []Position
}

// crea nuevo tablero
func newBoard(width int, height int) *Board {
	var area []Position
	for i := 1; i < width; i++ {
		for j := 1; j < height; j++ {
			area = append(area, newPosition(i, j))
		}
	}

	return &Board{width, height, area}
}

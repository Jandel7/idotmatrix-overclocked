package tetris

import (
	"github.com/pracucci/idotmatrix-overclocked/pkg/graphic"
)

// 3x5 mini font for small displays
var miniFont = map[rune][5]uint8{
	'T': {0x7, 0x2, 0x2, 0x2, 0x2},
	'E': {0x7, 0x1, 0x7, 0x1, 0x7},
	'R': {0x7, 0x5, 0x7, 0x3, 0x5},
	'I': {0x7, 0x2, 0x2, 0x2, 0x7},
	'S': {0x7, 0x1, 0x7, 0x4, 0x7},
	'G': {0x7, 0x1, 0x5, 0x5, 0x7},
	'A': {0x2, 0x5, 0x7, 0x5, 0x5},
	'M': {0x5, 0x7, 0x7, 0x5, 0x5},
	'O': {0x7, 0x5, 0x5, 0x5, 0x7},
	'V': {0x5, 0x5, 0x5, 0x5, 0x2},
}

// drawMiniChar draws a 3x5 character
func drawMiniChar(img []byte, char rune, x, y int, color graphic.Color) {
	data, ok := miniFont[char]
	if !ok {
		return
	}
	for row := 0; row < 5; row++ {
		for col := 0; col < 3; col++ {
			if data[row]&(1<<uint(col)) != 0 {
				graphic.SetPixel(img, x+col, y+row, color)
			}
		}
	}
}

// drawMiniText draws a string using the mini font (3x5 chars with 1px spacing)
func drawMiniText(img []byte, str string, x, y int, color graphic.Color) {
	for i, char := range str {
		drawMiniChar(img, char, x+i*4, y, color)
	}
}

// drawTetromino draws a small tetromino shape for decoration
func drawTetromino(img []byte, pieceType TetrominoType, x, y int, scale int) {
	colors := []graphic.Color{graphic.Cyan, graphic.Yellow, graphic.Violet, graphic.Green, graphic.Red, graphic.Blue, graphic.Orange}
	color := colors[pieceType]

	offsets := tetrominoShapes[pieceType][Rotation0]
	for _, off := range offsets {
		for dy := 0; dy < scale; dy++ {
			for dx := 0; dx < scale; dx++ {
				graphic.SetPixel(img, x+off.X*scale+dx, y+off.Y*scale+dy, color)
			}
		}
	}
}

// GenerateCoverImage creates the title screen
func GenerateCoverImage() []byte {
	img := make([]byte, graphic.DisplayWidth*graphic.DisplayHeight*3)

	// Dark purple gradient background
	for y := 0; y < graphic.DisplayHeight; y++ {
		for x := 0; x < graphic.DisplayWidth; x++ {
			intensity := uint8(10 + y/8)
			graphic.SetPixel(img, x, y, graphic.Color{intensity / 4, 0, intensity / 2})
		}
	}

	// "TETRIS" in mini font: 6 chars * 4 pixels = 24 pixels wide, center at (32-24)/2 = 4
	drawMiniText(img, "TETRIS", 5, 3, graphic.DarkWhite)
	drawMiniText(img, "TETRIS", 4, 2, graphic.Cyan)

	// Draw decorative tetrominoes centered
// Row 1: 4 pieces evenly spaced
	drawTetromino(img, TetrominoI, 2, 11, 1)
	drawTetromino(img, TetrominoO, 10, 11, 1)
	drawTetromino(img, TetrominoT, 18, 11, 1)
	drawTetromino(img, TetrominoS, 26, 11, 1)

	// Row 2: 3 pieces evenly spaced and centered
	drawTetromino(img, TetrominoZ, 6, 19, 1)
	drawTetromino(img, TetrominoL, 15, 19, 1)
	drawTetromino(img, TetrominoJ, 23, 19, 1)

	// Draw a small decorative border at bottom
	for x := 2; x < 29; x++ {
		if x%3 != 0 {
			graphic.SetPixel(img, x, 29, graphic.DimGray)
		}
	}

	return img
}

// GenerateGameOverImage creates the game over screen
func GenerateGameOverImage() []byte {
	img := make([]byte, graphic.DisplayWidth*graphic.DisplayHeight*3)

	// Dark red tinted background
	for y := 0; y < graphic.DisplayHeight; y++ {
		for x := 0; x < graphic.DisplayWidth; x++ {
			intensity := uint8(20 + y/4)
			graphic.SetPixel(img, x, y, graphic.Color{intensity / 2, 0, 0})
		}
	}

	// "GAME" and "OVER" in mini font, centered
	drawMiniText(img, "GAME", 9, 8, graphic.DarkRed)
	drawMiniText(img, "GAME", 8, 7, graphic.Red)
	drawMiniText(img, "OVER", 9, 17, graphic.DarkRed)
	drawMiniText(img, "OVER", 8, 16, graphic.Red)

	return img
}

// GenerateGameBackground creates the background for gameplay
func GenerateGameBackground() []byte {
	img := make([]byte, graphic.DisplayWidth*graphic.DisplayHeight*3)

	// Fill with dark background
	for y := 0; y < graphic.DisplayHeight; y++ {
		for x := 0; x < graphic.DisplayWidth; x++ {
			graphic.SetPixel(img, x, y, graphic.Color{10, 10, 15})
		}
	}

	// Draw game board border
	boardLeft := BoardOffsetX - 1
	boardRight := BoardOffsetX + BoardWidth*BlockSize
	boardTop := BoardOffsetY - 1
	boardBottom := BoardOffsetY + BoardHeight*BlockSize

	// Left border
	for y := boardTop; y <= boardBottom; y++ {
		graphic.SetPixel(img, boardLeft, y, graphic.DimGray)
	}
	// Right border
	for y := boardTop; y <= boardBottom; y++ {
		graphic.SetPixel(img, boardRight, y, graphic.DimGray)
	}
	// Bottom border
	for x := boardLeft; x <= boardRight; x++ {
		graphic.SetPixel(img, x, boardBottom, graphic.DimGray)
	}
	// Top border
	for x := boardLeft; x <= boardRight; x++ {
		if x%4 == 0 {
			graphic.SetPixel(img, x, boardTop, graphic.DarkWhite)
		}
	}

	// Draw subtle grid inside the board
	for y := 0; y < BoardHeight; y++ {
		for x := 0; x < BoardWidth; x++ {
			displayX := BoardOffsetX + x*BlockSize
			displayY := BoardOffsetY + y*BlockSize
			graphic.SetPixel(img, displayX, displayY, graphic.Color{20, 20, 25})
		}
	}

	return img
}
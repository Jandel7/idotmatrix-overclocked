package snake

// Color palette (exported colors are used by preview generator)
var (
	black       = [3]uint8{0, 0, 0}
	darkGreen   = [3]uint8{0, 100, 0}
	Green       = [3]uint8{0, 200, 0}   // Snake body color
	BrightGreen = [3]uint8{50, 255, 50} // Snake head color
	darkRed     = [3]uint8{100, 0, 0}
	Red         = [3]uint8{255, 0, 0} // Food color
	white       = [3]uint8{255, 255, 255}
	gray        = [3]uint8{80, 80, 80}

	// Obstacle colors
	rockColor    = [3]uint8{55, 55, 60}
	rockColorAlt = [3]uint8{50, 50, 55}
	lakeColor    = [3]uint8{35, 70, 135}
	lakeColorAlt = [3]uint8{40, 80, 145}
)

var brownPalette = [][3]uint8{
	{45, 31, 18},
	{49, 34, 19},
}

var font5x7 = map[rune][7]uint8{
	'S': {0x1E, 0x01, 0x01, 0x0E, 0x10, 0x10, 0x0F},
	'N': {0x11, 0x13, 0x15, 0x19, 0x11, 0x11, 0x11},
	'A': {0x0E, 0x11, 0x11, 0x1F, 0x11, 0x11, 0x11},
	'K': {0x11, 0x09, 0x05, 0x03, 0x05, 0x09, 0x11},
	'E': {0x1F, 0x01, 0x01, 0x0F, 0x01, 0x01, 0x1F},
	'G': {0x0E, 0x11, 0x01, 0x1D, 0x11, 0x11, 0x0E},
	'M': {0x11, 0x1B, 0x15, 0x15, 0x11, 0x11, 0x11},
	'O': {0x0E, 0x11, 0x11, 0x11, 0x11, 0x11, 0x0E},
	'V': {0x11, 0x11, 0x11, 0x11, 0x11, 0x0A, 0x04},
	'R': {0x0F, 0x11, 0x11, 0x0F, 0x05, 0x09, 0x11},
	'=': {0x00, 0x00, 0x1F, 0x00, 0x1F, 0x00, 0x00},
	'T': {0x1F, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04},
	'Y': {0x11, 0x11, 0x0A, 0x04, 0x04, 0x04, 0x04},
}

func setPixel(img []byte, x, y int, color [3]uint8) {
	if x < 0 || x >= 32 || y < 0 || y >= 32 {
		return
	}
	offset := (y*32 + x) * 3
	img[offset] = color[0]
	img[offset+1] = color[1]
	img[offset+2] = color[2]
}

func drawChar(img []byte, char rune, x, y int, color [3]uint8) {
	data, ok := font5x7[char]
	if !ok {
		return
	}
	for row := 0; row < 7; row++ {
		for col := 0; col < 5; col++ {
			if data[row]&(1<<col) != 0 {
				setPixel(img, x+col, y+row, color)
			}
		}
	}
}

func drawText(img []byte, text string, x, y int, color [3]uint8) {
	for i, char := range text {
		drawChar(img, char, x+i*6, y, color)
	}
}

func drawCoiledSnake(img []byte, cx, cy int) {
	coilPoints := []struct{ x, y int }{
		{-8, -4}, {-7, -5}, {-6, -6}, {-5, -7}, {-4, -7}, {-3, -7}, {-2, -7}, {-1, -7},
		{0, -7}, {1, -7}, {2, -7}, {3, -7}, {4, -7}, {5, -6}, {6, -5}, {7, -4},
		{7, -3}, {7, -2}, {7, -1}, {7, 0}, {7, 1}, {7, 2}, {7, 3}, {6, 4},
		{5, 5}, {4, 6}, {3, 6}, {2, 6}, {1, 6}, {0, 6}, {-1, 6}, {-2, 6},
		{-3, 6}, {-4, 5}, {-5, 4}, {-5, 3}, {-5, 2}, {-5, 1}, {-5, 0},
		{-5, -1}, {-4, -2}, {-3, -3}, {-2, -4}, {-1, -4}, {0, -4}, {1, -4}, {2, -4},
		{3, -3}, {4, -2}, {4, -1}, {4, 0}, {4, 1}, {3, 2}, {2, 3}, {1, 3},
		{0, 3}, {-1, 3}, {-2, 2}, {-2, 1}, {-2, 0},
		{-1, -1}, {0, -1}, {1, -1}, {1, 0}, {0, 0},
	}

	for i, p := range coilPoints {
		var color [3]uint8
		if i%3 == 0 {
			color = darkGreen
		} else {
			color = Green
		}
		setPixel(img, cx+p.x, cy+p.y, color)
		setPixel(img, cx+p.x+1, cy+p.y, color)
	}

	headX, headY := cx-8, cy-3
	setPixel(img, headX, headY, BrightGreen)
	setPixel(img, headX, headY-1, BrightGreen)
	setPixel(img, headX-1, headY, BrightGreen)
	setPixel(img, headX-1, headY-1, BrightGreen)
	setPixel(img, headX-1, headY-1, white)
	setPixel(img, headX, headY-1, white)
	setPixel(img, headX-2, headY, Red)
	setPixel(img, headX-3, headY-1, Red)
	setPixel(img, headX-3, headY+1, Red)
}

func GenerateCoverImage() []byte {
	img := make([]byte, 32*32*3)

	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			intensity := uint8(10 + y/8)
			if (x+y)%8 == 0 {
				intensity += 5
			}
			setPixel(img, x, y, [3]uint8{0, intensity / 2, 0})
		}
	}

	corners := [][2]int{{2, 2}, {29, 2}, {2, 29}, {29, 29}}
	for _, c := range corners {
		setPixel(img, c[0], c[1], darkGreen)
		setPixel(img, c[0]-1, c[1], darkGreen)
		setPixel(img, c[0]+1, c[1], darkGreen)
		setPixel(img, c[0], c[1]-1, darkGreen)
		setPixel(img, c[0], c[1]+1, darkGreen)
	}

	drawText(img, "SNAKE", 2, 2, [3]uint8{0, 50, 0})
	drawText(img, "SNAKE", 1, 1, BrightGreen)

	drawCoiledSnake(img, 16, 20)

	for x := 2; x < 29; x++ {
		if x%2 == 0 {
			setPixel(img, x, 29, gray)
		}
	}

	return img
}

type seed struct {
	x, y       int
	colorIndex int
}

func generateVoronoiBackground() []byte {
	img := make([]byte, 32*32*3)

	numSeeds := 20
	seeds := make([]seed, numSeeds)

	lcgState := uint32(12345)
	lcgNext := func() uint32 {
		lcgState = lcgState*1103515245 + 12345
		return (lcgState >> 16) & 0x7FFF
	}

	for i := 0; i < numSeeds; i++ {
		seeds[i] = seed{
			x: int(lcgNext() % 32),
			y: int(lcgNext() % 32),
		}
	}

	regionMap := make([]int, 32*32)
	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			minDist := 32*32 + 32*32
			nearestSeed := 0
			for i, s := range seeds {
				dx := x - s.x
				dy := y - s.y
				dist := dx*dx + dy*dy
				if dist < minDist {
					minDist = dist
					nearestSeed = i
				}
			}
			regionMap[y*32+x] = nearestSeed
		}
	}

	adjacency := make([]map[int]bool, numSeeds)
	for i := range adjacency {
		adjacency[i] = make(map[int]bool)
	}

	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			currentRegion := regionMap[y*32+x]
			if x < 31 {
				neighborRegion := regionMap[y*32+x+1]
				if neighborRegion != currentRegion {
					adjacency[currentRegion][neighborRegion] = true
					adjacency[neighborRegion][currentRegion] = true
				}
			}
			if y < 31 {
				neighborRegion := regionMap[(y+1)*32+x]
				if neighborRegion != currentRegion {
					adjacency[currentRegion][neighborRegion] = true
					adjacency[neighborRegion][currentRegion] = true
				}
			}
		}
	}

	for i := range seeds {
		seeds[i].colorIndex = -1
	}

	for i := 0; i < numSeeds; i++ {
		if seeds[i].colorIndex >= 0 {
			continue
		}
		queue := []int{i}
		seeds[i].colorIndex = 0
		for len(queue) > 0 {
			current := queue[0]
			queue = queue[1:]
			nextColor := 1 - seeds[current].colorIndex
			for neighbor := range adjacency[current] {
				if seeds[neighbor].colorIndex < 0 {
					seeds[neighbor].colorIndex = nextColor
					queue = append(queue, neighbor)
				}
			}
		}
	}

	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			region := regionMap[y*32+x]
			colorIdx := seeds[region].colorIndex
			setPixel(img, x, y, brownPalette[colorIdx])
		}
	}

	return img
}

func GenerateBackgroundWithObstacles(gameMap *Map) []byte {
	img := generateVoronoiBackground()

	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			tile := gameMap.GetTile(x, y)
			switch tile {
			case TileRock:
				color := rockColor
				if (x+y)%2 == 0 {
					color = rockColorAlt
				}
				setPixel(img, x, y, color)
			case TileLake:
				color := lakeColor
				if (x+y)%3 == 0 {
					color = lakeColorAlt
				}
				setPixel(img, x, y, color)
			}
		}
	}

	return img
}

func GenerateGameOverImage() []byte {
	img := make([]byte, 32*32*3)

	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			intensity := uint8(20 + y/4)
			setPixel(img, x, y, [3]uint8{intensity / 2, 0, 0})
		}
	}

	shadowColor := [3]uint8{50, 0, 0}
	drawText(img, "GAME", 5, 9, shadowColor)
	drawText(img, "GAME", 4, 8, Red)
	drawText(img, "OVER", 5, 18, shadowColor)
	drawText(img, "OVER", 4, 17, Red)

	return img
}
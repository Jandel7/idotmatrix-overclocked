package snake

import "math/rand"

// TileType represents the type of a map tile.
type TileType uint8

const (
	TileTerrain TileType = iota
	TileRock
	TileLake
)

const (
	mapSize = 32

	// Safe zone in the center where obstacles cannot be placed
	safeZoneMin = 10
	safeZoneMax = 22

	// Maximum placement attempts per obstacle
	maxPlacementAttempts = 1000
)

// Map represents the game map with terrain and obstacles.
type Map struct {
	Tiles [mapSize][mapSize]TileType
}

// NewMap creates a new empty game map (all terrain).
func NewMap() *Map {
	return &Map{}
}

// IsObstacle returns true if the tile at (x, y) is a rock or lake.
func (m *Map) IsObstacle(x, y int) bool {
	if x < 0 || x >= mapSize || y < 0 || y >= mapSize {
		return false
	}
	return m.Tiles[y][x] == TileRock || m.Tiles[y][x] == TileLake
}

// GetTile returns the tile type at (x, y).
func (m *Map) GetTile(x, y int) TileType {
	if x < 0 || x >= mapSize || y < 0 || y >= mapSize {
		return TileTerrain
	}
	return m.Tiles[y][x]
}

// TerrainPositions returns all positions that are terrain (not obstacles).
func (m *Map) TerrainPositions() []Point {
	var positions []Point
	for y := 0; y < mapSize; y++ {
		for x := 0; x < mapSize; x++ {
			if m.Tiles[y][x] == TileTerrain {
				positions = append(positions, Point{X: x, Y: y})
			}
		}
	}
	return positions
}

// Point represents a 2D coordinate.
type Point struct {
	X, Y int
}

// MapGenerator generates maps with obstacles based on level configuration.
type MapGenerator struct {
	rng *rand.Rand
}

// NewMapGenerator creates a new map generator with the given random seed.
func NewMapGenerator(seed int64) *MapGenerator {
	return &MapGenerator{
		rng: rand.New(rand.NewSource(seed)),
	}
}

// Generate creates a new game map with the specified number of rocks and lakes.
func (g *MapGenerator) Generate(numRocks, numLakes int) *Map {
	m := NewMap()

	for i := 0; i < numRocks; i++ {
		g.placeRock(m)
	}

	for i := 0; i < numLakes; i++ {
		g.placeLake(m)
	}

	return m
}

// placeRock attempts to place a 2x2 or 3x3 rock with rounded corners on the map.
func (g *MapGenerator) placeRock(m *Map) bool {
	for attempt := 0; attempt < maxPlacementAttempts; attempt++ {
		// Smaller rocks for 32x32: size 2 or 3
		size := 2 + g.rng.Intn(2)

		x := g.rng.Intn(mapSize - size)
		y := g.rng.Intn(mapSize - size)

		rockPixels := g.generateRoundedSquare(x, y, size)

		if g.canPlaceObstacle(m, rockPixels) {
			for _, p := range rockPixels {
				m.Tiles[p.Y][p.X] = TileRock
			}
			return true
		}
	}
	return false
}

// generateRoundedSquare creates a square shape with rounded corners.
func (g *MapGenerator) generateRoundedSquare(x, y, size int) []Point {
	var pixels []Point

	for dy := 0; dy < size; dy++ {
		for dx := 0; dx < size; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			if dx == size-1 && dy == 0 {
				continue
			}
			if dx == 0 && dy == size-1 {
				continue
			}
			if dx == size-1 && dy == size-1 {
				continue
			}

			pixels = append(pixels, Point{X: x + dx, Y: y + dy})
		}
	}

	return pixels
}

// placeLake attempts to place an irregular lake (5-10 pixels) on the map.
func (g *MapGenerator) placeLake(m *Map) bool {
	for attempt := 0; attempt < maxPlacementAttempts; attempt++ {
		startX := g.rng.Intn(mapSize-4) + 2
		startY := g.rng.Intn(mapSize-4) + 2

		lakePixels := g.generateLakeShape(startX, startY)

		if g.canPlaceObstacle(m, lakePixels) {
			for _, p := range lakePixels {
				m.Tiles[p.Y][p.X] = TileLake
			}
			return true
		}
	}
	return false
}

// generateLakeShape creates an irregular lake shape of 5-10 pixels.
func (g *MapGenerator) generateLakeShape(startX, startY int) []Point {
	targetSize := 5 + g.rng.Intn(6) // 5-10 pixels for 32x32

	pixels := []Point{{startX, startY}}
	pixelSet := make(map[Point]bool)
	pixelSet[Point{startX, startY}] = true

	directions := []Point{
		{0, -1}, {0, 1}, {-1, 0}, {1, 0},
		{-1, -1}, {1, -1}, {-1, 1}, {1, 1},
	}

	for len(pixels) < targetSize {
		base := pixels[g.rng.Intn(len(pixels))]

		dir := directions[g.rng.Intn(len(directions))]
		newP := Point{base.X + dir.X, base.Y + dir.Y}

		if newP.X < 0 || newP.X >= mapSize || newP.Y < 0 || newP.Y >= mapSize {
			continue
		}

		if pixelSet[newP] {
			continue
		}

		pixels = append(pixels, newP)
		pixelSet[newP] = true
	}

	return pixels
}

// canPlaceObstacle checks if the given pixels can be placed as an obstacle.
func (g *MapGenerator) canPlaceObstacle(m *Map, pixels []Point) bool {
	for _, p := range pixels {
		if p.X < 0 || p.X >= mapSize || p.Y < 0 || p.Y >= mapSize {
			return false
		}

		if p.X >= safeZoneMin && p.X < safeZoneMax && p.Y >= safeZoneMin && p.Y < safeZoneMax {
			return false
		}

		if m.Tiles[p.Y][p.X] != TileTerrain {
			return false
		}

		adjacent := []Point{
			{p.X - 1, p.Y}, {p.X + 1, p.Y},
			{p.X, p.Y - 1}, {p.X, p.Y + 1},
		}
		for _, adj := range adjacent {
			isPartOfSame := false
			for _, op := range pixels {
				if op == adj {
					isPartOfSame = true
					break
				}
			}
			if isPartOfSame {
				continue
			}

			if adj.X >= 0 && adj.X < mapSize && adj.Y >= 0 && adj.Y < mapSize {
				if m.Tiles[adj.Y][adj.X] != TileTerrain {
					return false
				}
			}
		}
	}
	return true
}
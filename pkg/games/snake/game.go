package snake

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"golang.org/x/term"

	"github.com/pracucci/idotmatrix-overclocked/pkg/protocol"
)

// Direction represents the snake's movement direction.
type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

// PixelChange represents a pixel update to be sent to the device.
type PixelChange struct {
	pos     Point
	r, g, b uint8
}

// Game represents the snake game state.
type Game struct {
	device      protocol.DeviceConnection
	snake       []Point // Head is snake[0], tail is snake[len-1]
	direction   Direction
	food        Point
	score       int
	running     bool
	gameOver    bool
	inputChan   chan rune
	background  []byte // Store background RGB data for pixel restoration

	// Level system
	startLevel   int         // Starting level (for testing)
	currentLevel int         // Current level (1-based)
	applesEaten  int         // Apples eaten in current level
	gameMap     *Map    // Current map with obstacles
	levelConfig LevelConfig // Current level configuration
	growthQueue int         // Pending growth segments
}

// NewGame creates a new snake game instance.
func NewGame(device protocol.DeviceConnection, startLevel int) *Game {
	if startLevel < 1 {
		startLevel = 1
	}
	g := &Game{
		device:       device,
		running:      true,
		inputChan:    make(chan rune, 10),
		startLevel:   startLevel,
		currentLevel: startLevel,
	}
	return g
}

// reset initializes the game state for a new game.
func (g *Game) reset() {
	g.currentLevel = g.startLevel
	g.applesEaten = 0
	g.score = 0
	g.gameOver = false
	g.growthQueue = 0
	g.resetSnakePosition(InitialLength)
}

// resetSnakePosition resets the snake to the center with the given length.
func (g *Game) resetSnakePosition(length int) {
	g.snake = make([]Point, length)
	startX := DisplaySize / 2
	startY := DisplaySize / 2
	for i := 0; i < length; i++ {
		g.snake[i] = Point{X: startX - i, Y: startY}
	}
	g.direction = Right
}

// advanceLevel advances to the next level, preserving snake length.
func (g *Game) advanceLevel() {
	g.currentLevel++
	g.applesEaten = 0
	currentLength := len(g.snake)
	g.resetSnakePosition(currentLength)
}

// setupLevel generates the map and prepares for the current level.
func (g *Game) setupLevel() {
	g.levelConfig = GetLevelConfig(g.currentLevel)

	mapGen := NewMapGenerator(time.Now().UnixNano())
	g.gameMap = mapGen.Generate(g.levelConfig.NumRocks, g.levelConfig.NumLakes)

	g.background = GenerateBackgroundWithObstacles(g.gameMap)
}

// spawnFood places food on a valid terrain position.
func (g *Game) spawnFood() {
	terrain := g.gameMap.TerrainPositions()

	snakeSet := make(map[Point]bool)
	for _, p := range g.snake {
		snakeSet[p] = true
	}

	var validPositions []Point
	for _, p := range terrain {
		if snakeSet[p] {
			continue
		}
		if p.X == 0 || p.X == DisplaySize-1 || p.Y == 0 || p.Y == DisplaySize-1 {
			continue
		}
		validPositions = append(validPositions, p)
	}

	if len(validPositions) == 0 {
		g.food = terrain[rand.Intn(len(terrain))]
		return
	}

	g.food = validPositions[rand.Intn(len(validPositions))]
}

// getBackgroundPixel returns the background color at the given position.
func (g *Game) getBackgroundPixel(x, y int) (uint8, uint8, uint8) {
	offset := (y*32 + x) * 3
	return g.background[offset], g.background[offset+1], g.background[offset+2]
}

// calculateNewHead returns the new head position based on direction.
func (g *Game) calculateNewHead() Point {
	head := g.snake[0]
	switch g.direction {
	case Up:
		return Point{X: head.X, Y: head.Y - 1}
	case Down:
		return Point{X: head.X, Y: head.Y + 1}
	case Left:
		return Point{X: head.X - 1, Y: head.Y}
	case Right:
		return Point{X: head.X + 1, Y: head.Y}
	}
	return head
}

// isCollision checks if the given point causes a collision.
func (g *Game) isCollision(p Point) bool {
	if p.X < 0 || p.X >= DisplaySize || p.Y < 0 || p.Y >= DisplaySize {
		return true
	}
	if g.gameMap.IsObstacle(p.X, p.Y) {
		return true
	}
	for i := 0; i < len(g.snake)-1; i++ {
		if g.snake[i] == p {
			return true
		}
	}
	return false
}

// move moves the snake and returns pixel changes to render.
func (g *Game) move() ([]PixelChange, bool) {
	var changes []PixelChange

	newHead := g.calculateNewHead()

	if g.isCollision(newHead) {
		g.gameOver = true
		return nil, false
	}

	changes = append(changes, PixelChange{newHead, 0, 255, 0})
	g.snake = append([]Point{newHead}, g.snake...)

	if newHead == g.food {
		g.score++
		g.applesEaten++
		g.growthQueue += GrowthPerApple

		if g.applesEaten >= ApplesPerLevel {
			return changes, true
		}

		g.spawnFood()
		changes = append(changes, PixelChange{g.food, 255, 0, 0})
	}

	if g.growthQueue > 0 {
		g.growthQueue--
	} else {
		tail := g.snake[len(g.snake)-1]
		r, gb, b := g.getBackgroundPixel(tail.X, tail.Y)
		changes = append(changes, PixelChange{tail, r, gb, b})
		g.snake = g.snake[:len(g.snake)-1]
	}

	return changes, false
}

// handleInput processes keyboard input.
func (g *Game) handleInput() {
	select {
	case key := <-g.inputChan:
		switch key {
		case 'w', 'W':
			if g.direction != Down {
				g.direction = Up
			}
		case 's', 'S':
			if g.direction != Up {
				g.direction = Down
			}
		case 'a', 'A':
			if g.direction != Right {
				g.direction = Left
			}
		case 'd', 'D':
			if g.direction != Left {
				g.direction = Right
			}
		case 'q', 'Q':
			g.running = false
		}
	default:
	}
}

// renderInitial draws the initial snake and food on the display.
func (g *Game) renderInitial() {
	for _, p := range g.snake {
		protocol.SetPixel(g.device, p.X, p.Y, 0, 255, 0)
		time.Sleep(20 * time.Millisecond)
	}
	protocol.SetPixel(g.device, g.food.X, g.food.Y, 255, 0, 0)
	time.Sleep(20 * time.Millisecond)
}

// showImage displays an image on the device.
func (g *Game) showImage(rgbData []byte) error {
	if err := protocol.SetDrawMode(g.device, 1); err != nil {
		return err
	}
	if err := protocol.SendImage(g.device, rgbData); err != nil {
		return err
	}
	time.Sleep(500 * time.Millisecond)
	return nil
}

// waitForKey blocks until a key is pressed.
func (g *Game) waitForKey() rune {
	return <-g.inputChan
}

// startInputReader starts the keyboard input goroutine.
func (g *Game) startInputReader() func() {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Printf("Warning: could not set raw mode: %v\n", err)
		return func() {}
	}

	go func() {
		buf := make([]byte, 3)
		for g.running {
			n, err := os.Stdin.Read(buf)
			if err != nil {
				return
			}
			if n > 0 {
				if n == 3 && buf[0] == 27 && buf[1] == 91 {
					switch buf[2] {
					case 65:
						g.inputChan <- 'w'
					case 66:
						g.inputChan <- 's'
					case 67:
						g.inputChan <- 'd'
					case 68:
						g.inputChan <- 'a'
					}
				} else {
					g.inputChan <- rune(buf[0])
				}
			}
		}
	}()

	return func() {
		term.Restore(int(os.Stdin.Fd()), oldState)
	}
}

// runLevel runs a single level and returns true if the game should continue.
func (g *Game) runLevel() (continueGame bool, advanceLevel bool) {
	g.setupLevel()

	if err := g.showLevelInterstitial(); err != nil {
		return false, false
	}

	if err := g.showImage(g.background); err != nil {
		return false, false
	}

	g.spawnFood()
	g.renderInitial()

	for g.running && !g.gameOver {
		tickStart := time.Now()

		g.handleInput()

		if !g.gameOver {
			changes, shouldAdvance := g.move()
			if shouldAdvance {
				return true, true
			}
			for _, c := range changes {
				protocol.SetPixel(g.device, c.pos.X, c.pos.Y, c.r, c.g, c.b)
				time.Sleep(20 * time.Millisecond)
			}
		}

		elapsed := time.Since(tickStart)
		sleepTime := g.levelConfig.TickDelay - elapsed
		if sleepTime > 0 {
			time.Sleep(sleepTime)
		}
	}

	if !g.running {
		return false, false
	}

	return true, false
}

// Run starts the main game loop.
func (g *Game) Run() error {
	fmt.Println("Starting Snake!")
	fmt.Println("Controls: WASD or Arrow keys to move, Q to quit, R to restart")

	cleanup := g.startInputReader()
	defer cleanup()

	for g.running {
		if err := g.showImage(GenerateCoverImage()); err != nil {
			return err
		}
		fmt.Print("Press any key to start...")
		key := g.waitForKey()
		fmt.Println()
		if key == 'q' || key == 'Q' {
			break
		}

		g.reset()

		for g.running && !g.gameOver {
			continueGame, advance := g.runLevel()
			if !continueGame {
				break
			}
			if advance {
				g.advanceLevel()
			}
		}

		if !g.running {
			break
		}

		if err := g.showImage(GenerateGameOverImage()); err != nil {
			return err
		}
		key = g.waitForKey()
		if key == 'q' || key == 'Q' {
			break
		}
	}

	return nil
}
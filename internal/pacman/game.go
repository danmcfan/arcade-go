package pacman

import (
	"time"

	"github.com/nsf/termbox-go"
)

const (
	WeakTicks = 50
)

type State int

const (
	GameOver State = 0
	Playing  State = 1
)

type Pixel int

const (
	Open    Pixel = 0
	Dot     Pixel = 1
	PowerUp Pixel = 2

	Wall Pixel = 100
	Gate Pixel = 101
)

type Direction rune

const (
	None  Direction = ' '
	Up    Direction = 'w'
	Down  Direction = 's'
	Left  Direction = 'a'
	Right Direction = 'd'
)

type Position struct {
	x, y int
}

type Game struct {
	maze Maze

	state State
	score int

	player Player
	ghosts Ghosts
}

func newGame() Game {
	game := Game{
		state:  Playing,
		maze:   newMaze(),
		player: newPlayer(),
		ghosts: newGhosts(),
	}
	return game
}

func (g *Game) loop(gameTick *time.Ticker, keyboardEvents <-chan termbox.Event, exit chan<- bool) {
	for {
		select {
		case ev := <-keyboardEvents:
			if ev.Type == termbox.EventKey {
				switch ev.Ch {
				case 'w', 's', 'a', 'd', 'r':
					g.handleInput(ev.Ch)
				case 'q':
					exit <- true
				}
			}
		case <-gameTick.C:
			g.update()
			draw(*g)
		}
	}
}

func (g *Game) handleInput(ch rune) {
	if ch == 'r' && g.state == GameOver {
		*g = newGame()
		return
	}

	direction := Direction(ch)
	nextPosition := g.player.nextPositionByDirection(direction)
	pixel := g.pixelFromPosition(nextPosition)
	if pixel == Wall || pixel == Gate {
		return
	}

	g.player.direction = direction
}

func (g *Game) update() {
	for _, ghost := range g.ghosts {
		if ghost.weakTicks > 0 {
			ghost.weakTicks--
		}
	}

	g.movePlayer()
	g.checkCollisions()

	for _, ghost := range g.ghosts {
		ghost.setInLobby()
		ghost.setDirection(g)
	}

	ghostLobby := g.ghosts.firstInLobby()
	if ghostLobby != nil {
		g.moveGhost(ghostLobby)
	}

	for _, ghost := range g.ghosts.outLobby() {
		g.moveGhost(ghost)
	}

	g.checkCollisions()
}

func (g *Game) pixelFromPosition(p Position) Pixel {
	return Pixel(g.maze[p.y][p.x])
}

func (g *Game) checkCollisions() {
	for _, ghost := range g.ghosts {
		if ghost.position == g.player.position {
			if ghost.weakTicks > 0 {
				g.score += 200
				ghost.reset()
			} else {
				g.state = GameOver
			}
		}
	}
}

func (g *Game) movePlayer() {
	nextPosition := g.player.nextPosition()
	nextPixel := g.pixelFromPosition(nextPosition)

	switch nextPixel {
	case Wall, Gate:
		return
	case Dot:
		g.score += 100
	case PowerUp:
		for _, ghost := range g.ghosts {
			ghost.weakTicks = WeakTicks
		}
	}

	g.maze.setPixel(g.player.position, Open)
	g.player.position = nextPosition
}

func (g *Game) moveGhost(gh *Ghost) {
	nextPosition := gh.nextPosition()
	nextPixel := g.pixelFromPosition(nextPosition)

	switch nextPixel {
	case Wall:
		return
	case Gate:
		if !gh.inLobby {
			return
		}
	}

	gh.position = nextPosition
}

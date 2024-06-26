package pacman

import (
	"time"

	"github.com/nsf/termbox-go"
)

const (
	LobbyXMin = 10
	LobbyXMax = 17
	LobbyYMin = 12
	LobbyYMax = 16

	GhostStartX = 14
	GhostStartY = 11

	LobbyTicks = 50
	WeakTicks  = 50
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

	PlayerPixel Pixel = 10

	RedGhost   Pixel = 20
	PinkGhost  Pixel = 21
	GreenGhost Pixel = 22
	GrayGhost  Pixel = 23

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
	game.initMaze()
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

func (g *Game) initMaze() {
	g.maze.setPixel(g.player.position, PlayerPixel)
	for pixel, ghost := range g.ghosts {
		g.maze.setPixel(ghost.position, pixel)
	}
}

func (g *Game) update() {
	for pixel, ghost := range g.ghosts {
		if ghost.lobbyTicks > 0 {
			ghost.lobbyTicks--
		}
		if ghost.weakTicks > 0 {
			ghost.weakTicks--
		}
		g.ghosts[pixel] = ghost
	}

	nextPosition := g.player.nextPosition()
	nextPixel := g.pixelFromPosition(nextPosition)

	switch nextPixel {
	case Wall, Gate:
		return
	case RedGhost, PinkGhost, GreenGhost, GrayGhost:
		ghost := g.ghosts[nextPixel]
		if ghost.weakTicks > 0 {
			g.score += 200
			g.ghosts[nextPixel] = newGhost(ghost.positionDefault.x, ghost.positionDefault.y, ghost.positionDefault.x, ghost.positionDefault.y, 0)

		} else {
			g.state = GameOver
		}
	case Dot:
		g.score += 100
	case PowerUp:
		for pixel, ghost := range g.ghosts {
			ghost.weakTicks = WeakTicks
			g.ghosts[pixel] = ghost
		}
	}

	g.updateMaze(Open, PlayerPixel, g.player.position, nextPosition)
	g.player.position = nextPosition

	for pixel, ghost := range g.ghosts {
		ghost.setDirection(g)
		nextPosition := ghost.nextPosition()
		nextPixel := g.pixelFromPosition(nextPosition)

		switch nextPixel {
		case Wall:
			continue
		case Gate:
			if ghost.lobbyTicks > 0 {
				continue
			}
		case RedGhost, PinkGhost, GreenGhost, GrayGhost:
			nextPixel = Open
		}

		g.updateMaze(nextPixel, pixel, ghost.position, nextPosition)
		ghost.position = nextPosition

		g.ghosts[pixel] = ghost
	}
}

func (g *Game) pixelFromPosition(p Position) Pixel {
	return Pixel(g.maze[p.y][p.x])
}

func (g *Game) updateMaze(cpi, npi Pixel, cpo, npo Position) {
	g.maze.setPixel(cpo, cpi)
	g.maze.setPixel(npo, npi)
}

package pacman

import (
	"math/rand"

	"github.com/nsf/termbox-go"
)

const (
	GhostStartXLeft  = 13
	GhostStartXRight = 14
	GhostStartY      = 11

	RedGhostX = 11
	RedGhostY = 13

	PinkGhostX = 11
	PinkGhostY = 15

	GreenGhostX = 16
	GreenGhostY = 13

	GrayGhostX = 16
	GrayGhostY = 15
)

type Ghost struct {
	position        Position
	positionDefault Position
	direction       Direction
	weakTicks       int
	inLobby         bool
}

type Ghosts map[termbox.Attribute]*Ghost

func newGhosts() Ghosts {
	ghosts := make(Ghosts)

	ghosts[termbox.ColorLightRed] = newGhost(RedGhostX, RedGhostY, RedGhostX, RedGhostY)
	ghosts[termbox.ColorLightMagenta] = newGhost(PinkGhostX, PinkGhostY, PinkGhostX, PinkGhostY)
	ghosts[termbox.ColorLightGreen] = newGhost(GreenGhostX, GreenGhostY, GreenGhostX, GreenGhostY)
	ghosts[termbox.ColorLightGray] = newGhost(GrayGhostX, GrayGhostY, GrayGhostX, GrayGhostY)

	return ghosts
}

func (g Ghosts) firstInLobby() *Ghost {
	ghostOrder := []termbox.Attribute{termbox.ColorLightRed, termbox.ColorLightMagenta, termbox.ColorLightGreen, termbox.ColorLightGray}

	for _, color := range ghostOrder {
		ghost := g[color]
		if ghost.inLobby {
			return ghost
		}
	}

	return nil
}

func (g Ghosts) outLobby() []*Ghost {
	outLobby := make([]*Ghost, 0)

	for _, ghost := range g {
		if !ghost.inLobby {
			outLobby = append(outLobby, ghost)
		}
	}

	return outLobby
}

func newGhost(x, y, defaultX, defaultY int) *Ghost {
	return &Ghost{Position{x, y}, Position{defaultX, defaultY}, None, 0, true}
}

func (g *Ghost) reset() {
	g.position = g.positionDefault
	g.direction = None
	g.weakTicks = 0
	g.inLobby = true
}

func (g *Ghost) nextPosition() Position {
	nextPosition := g.position

	switch g.direction {
	case Up:
		nextPosition.y--
	case Down:
		nextPosition.y++
	case Left:
		nextPosition.x--
	case Right:
		nextPosition.x++
	}

	if nextPosition.y == 14 {
		if nextPosition.x == -1 {
			nextPosition.x = Width - 1
		}
		if nextPosition.x == Width {
			nextPosition.x = 0
		}
	}

	return nextPosition
}

func (gh *Ghost) nextPositionByDirection(d Direction) Position {
	nextPosition := gh.position

	switch d {
	case Up:
		nextPosition.y--
	case Down:
		nextPosition.y++
	case Left:
		nextPosition.x--
	case Right:
		nextPosition.x++
	}

	if nextPosition.y == 14 {
		if nextPosition.x == -1 {
			nextPosition.x = Width - 1
		}
		if nextPosition.x == Width {
			nextPosition.x = 0
		}
	}

	return nextPosition
}

func (g *Ghost) setInLobby() {
	if g.position.y == GhostStartY {
		if g.position.x == GhostStartXLeft || g.position.x == GhostStartXRight {
			g.inLobby = false
		}
	}
}

func (gh *Ghost) setDirection(gm *Game) {
	if gh.inLobby {
		if gh.position.x < GhostStartXLeft {
			gh.direction = Right
			return
		}
		if gh.position.x > GhostStartXRight {
			gh.direction = Left
			return
		}
		if gh.position.y > GhostStartY {
			gh.direction = Up
			return
		}
	} else {
		directions := []Direction{Up, Down, Left, Right}
		directionsValid := []Direction{}

		for _, d := range directions {
			if d == Up && gh.direction == Down {
				continue
			}
			if d == Down && gh.direction == Up {
				continue
			}
			if d == Left && gh.direction == Right {
				continue
			}
			if d == Right && gh.direction == Left {
				continue
			}

			nextPosition := gh.nextPositionByDirection(d)
			nextPixel := gm.pixelFromPosition(nextPosition)

			switch nextPixel {
			case Wall, Gate:
				continue
			}

			directionsValid = append(directionsValid, d)
		}

		gh.direction = directionsValid[rand.Intn(len(directionsValid))]
	}
}

package pacman

import "math/rand"

const (
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
	lobbyTicks      int
}

type Ghosts map[Pixel]Ghost

func newGhosts() Ghosts {
	ghosts := make(map[Pixel]Ghost)

	ghosts[RedGhost] = newGhost(RedGhostX, RedGhostY, RedGhostX, RedGhostY, 0)
	ghosts[PinkGhost] = newGhost(PinkGhostX, PinkGhostY, PinkGhostX, PinkGhostY, 50)
	ghosts[GreenGhost] = newGhost(GreenGhostX, GreenGhostY, GreenGhostX, GreenGhostY, 100)
	ghosts[GrayGhost] = newGhost(GrayGhostX, GrayGhostY, GrayGhostX, GrayGhostY, 150)

	return ghosts
}

func newGhost(x, y, defaultX, defaultY, lobbyTicks int) Ghost {
	return Ghost{Position{x, y}, Position{defaultX, defaultY}, None, 0, lobbyTicks}
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
			nextPosition.y = 0
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
			nextPosition.y = 0
		}
	}

	return nextPosition
}

func (gh *Ghost) setDirection(gm *Game) {
	directions := []Direction{Up, Down, Left, Right}
	directionsValid := []Direction{}

	for _, d := range directions {
		nextPosition := gh.nextPositionByDirection(d)
		nextPixel := gm.pixelFromPosition(nextPosition)

		switch nextPixel {
		case Wall:
			continue
		}

		directionsValid = append(directionsValid, d)
	}

	if rand.Float64() > 0.95 {
		gh.direction = directionsValid[rand.Intn(len(directionsValid))]
	}
}

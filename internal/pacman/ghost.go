package pacman

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
}

type Ghosts map[Pixel]Ghost

func newGhosts() Ghosts {
	ghosts := make(map[Pixel]Ghost)

	ghosts[RedGhost] = newGhost(RedGhostX, RedGhostY, RedGhostX, RedGhostY)
	ghosts[PinkGhost] = newGhost(PinkGhostX, PinkGhostY, PinkGhostX, PinkGhostY)
	ghosts[GreenGhost] = newGhost(GreenGhostX, GreenGhostY, GreenGhostX, GreenGhostY)
	ghosts[GrayGhost] = newGhost(GrayGhostX, GrayGhostY, GrayGhostX, GrayGhostY)

	return ghosts
}

func newGhost(x, y, defaultX, defaultY int) Ghost {
	return Ghost{Position{x, y}, Position{defaultX, defaultY}, None, 0}
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

func (gh *Ghost) setDirection(gm *Game) {}

package pacman

const (
	PlayerX = 14
	PlayerY = 23
)

type Player struct {
	position  Position
	direction Direction
}

func newPlayer() Player {
	return Player{Position{PlayerX, PlayerY}, None}
}

func (p *Player) nextPosition() Position {
	nextPosition := p.position

	switch p.direction {
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

func (p *Player) nextPositionByDirection(d Direction) Position {
	nextPosition := p.position

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

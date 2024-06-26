package snake

import (
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
)

type Direction rune

const (
	None  Direction = ' '
	Up    Direction = 'w'
	Down  Direction = 's'
	Left  Direction = 'a'
	Right Direction = 'd'
)

type State rune

const (
	Idle     State = 'i'
	Playing  State = 'p'
	GameOver State = 'o'
)

type Position struct {
	x, y int
}

type Snake struct {
	positions []Position
}

type Game struct {
	width, height int
	state         State
	click         bool
	score         int

	snake     Snake
	direction Direction

	numApples int
	apples    []Position
}

func gameLoop(game *Game, gameTick *time.Ticker, keyboardEvents <-chan termbox.Event, exit chan<- bool) {
	for {
		select {
		case ev := <-keyboardEvents:
			if ev.Type == termbox.EventKey {
				switch ev.Ch {
				case 'w', 's', 'a', 'd', 'r':
					if game.state == Idle {
						handleIdleInput(game, ev.Ch)
					}
					if game.state == Playing && !game.click {
						handlePlayingInput(game, ev.Ch)
						game.click = true
					}
					if game.state == GameOver {
						handleGameOverInput(game, ev.Ch)
					}
				case 'q':
					exit <- true
				}
			}
		case <-gameTick.C:
			if game.state == Playing {
				game.update()
				game.click = false
			}
			draw(*game)
		}
	}
}

func handleIdleInput(g *Game, ch rune) {
	switch ch {
	case 'w':
		if g.direction != Down {
			g.direction = Up
			g.state = Playing
		}
	case 's':
		if g.direction != Up {
			g.direction = Down
			g.state = Playing
		}
	case 'a':
		if g.direction != Right {
			g.direction = Left
			g.state = Playing
		}
	case 'd':
		if g.direction != Left {
			g.direction = Right
			g.state = Playing
		}
	}
}

func handlePlayingInput(g *Game, ch rune) {
	switch ch {
	case 'w':
		if g.direction != Down {
			g.direction = Up
		}
	case 's':
		if g.direction != Up {
			g.direction = Down
		}
	case 'a':
		if g.direction != Right {
			g.direction = Left
		}
	case 'd':
		if g.direction != Left {
			g.direction = Right
		}
	}
}

func handleGameOverInput(g *Game, ch rune) {
	if ch == 'r' {
		*g = newGame(g.width, g.numApples)
	}
}

func newSnake(width, height int) Snake {
	positions := make([]Position, 0)

	center := Position{width / 2, height / 2}
	positions = append(positions, Position{center.x - 2, center.y})
	positions = append(positions, Position{center.x - 1, center.y})
	positions = append(positions, center)

	return Snake{positions}
}

func (s *Snake) head() Position {
	return s.positions[len(s.positions)-1]
}

func (s *Snake) tail() []Position {
	return s.positions[:len(s.positions)-1]
}

func newGame(size, numApples int) Game {
	if numApples <= 0 {
		panic("Number of apples must be greater than 0")
	}

	if numApples >= 10 {
		panic("Number of apples must be less than 10")
	}

	game := Game{
		width:     size,
		height:    size,
		state:     Idle,
		click:     false,
		score:     0,
		snake:     newSnake(size, size),
		direction: Right,
		numApples: numApples,
		apples:    make([]Position, 0),
	}

	for range numApples {
		game.apples = append(game.apples, newApple(&game))
	}

	return game
}

func (g *Game) update() {
	if g.state == Idle {
		return
	}

	head := g.snake.head()

	switch g.direction {
	case Up:
		head.y--
	case Down:
		head.y++
	case Left:
		head.x--
	case Right:
		head.x++
	case None:
		return
	}

	g.snake.positions = append(g.snake.positions, head)

	if head.x < 0 || head.x >= g.width || head.y < 0 || head.y >= g.height {
		g.state = GameOver
		return
	}

	for _, position := range g.snake.tail() {
		if head.x == position.x && head.y == position.y {
			g.state = GameOver
			return
		}
	}

	cutSnake := true
	for i, apple := range g.apples {
		if head.x == apple.x && head.y == apple.y {
			g.score++
			g.apples[i] = newApple(g)
			cutSnake = false
		}
	}

	if cutSnake {
		g.snake.positions = g.snake.positions[1:]
	}
}

func newApple(g *Game) Position {
	apple := Position{x: rand.Intn(g.width), y: rand.Intn(g.height)}
	for _, position := range g.snake.positions {
		if position.x == apple.x && position.y == apple.y {
			apple = newApple(g)
		}
	}
	for _, position := range g.apples {
		if position.x == apple.x && position.y == apple.y {
			apple = newApple(g)
		}
	}

	return apple
}

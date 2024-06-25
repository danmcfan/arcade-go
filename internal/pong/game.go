package pong

import (
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
)

const (
	width  = 41
	height = 21

	centerX = width / 2
	centerY = height / 2
)

type State rune

const (
	Waiting  State = 'w'
	Playing  State = 'p'
	GameOver State = 'o'
)

type Ball struct {
	x, y, speedX, speedY int
}

type Direction rune

const (
	None Direction = ' '
	Up   Direction = 'U'
	Down Direction = 'D'
)

type Player struct {
	y, score  int
	direction Direction
}

type Game struct {
	state     State
	wait      int
	waitCount int
	ball      Ball

	playerOne   Player
	playerTwo   Player
	scoreLimit  int
	paddleWidth int
}

func gameLoop(game *Game, gameTick *time.Ticker, keyboardEvents <-chan termbox.Event, exit chan<- bool) {
	for {
		select {
		case ev := <-keyboardEvents:
			if ev.Type == termbox.EventKey {
				switch ev.Ch {
				case 'w', 's', 'i', 'k', 'r':
					if game.state == Waiting || game.state == Playing {
						handlePlayingInput(game, ev.Ch)
					}
					if game.state == GameOver {
						handleGameOverInput(game, ev.Ch)
					}
				case 'q':
					exit <- true
				}
			}
		case <-gameTick.C:
			if game.state == Waiting || game.state == Playing {
				game.update()
			}
			draw(*game)
		}
	}
}

func handlePlayingInput(g *Game, ch rune) {
	switch ch {
	case 'w':
		g.playerOne.direction = Up
	case 's':
		g.playerOne.direction = Down
	case 'i':
		g.playerTwo.direction = Up
	case 'k':
		g.playerTwo.direction = Down
	}
}

func handleGameOverInput(g *Game, ch rune) {
	if ch == 'r' {
		*g = newGame(g.scoreLimit, g.paddleWidth)
	}
}

func newGame(scoreLimit, paddleWidth int) Game {
	return Game{
		state:       Waiting,
		waitCount:   0,
		wait:        10,
		playerOne:   newPlayer(paddleWidth),
		playerTwo:   newPlayer(paddleWidth),
		ball:        newBall(),
		scoreLimit:  scoreLimit,
		paddleWidth: paddleWidth,
	}
}

func newPlayer(paddleWidth int) Player {
	return Player{y: centerY - paddleWidth/2, direction: None, score: 0}
}

func newBall() Ball {
	ball := Ball{x: centerX, y: centerY, speedX: newSpeed(), speedY: newSpeed()}

	if ball.speedX == 0 || ball.speedY == 0 {
		ball = newBall()
	}

	return ball
}

func newSpeed() int {
	return -1 + rand.Intn(3)
}

func (g *Game) update() {
	if g.playerOne.direction == Up {
		if g.playerOne.y > 0 {
			g.playerOne.y--
		}
	}
	if g.playerOne.direction == Down {
		if g.playerOne.y < height-g.paddleWidth {
			g.playerOne.y++
		}
	}

	if g.playerTwo.direction == Up {
		if g.playerTwo.y > 0 {
			g.playerTwo.y--
		}
	}
	if g.playerTwo.direction == Down {
		if g.playerTwo.y < height-g.paddleWidth {
			g.playerTwo.y++
		}
	}

	if g.playerOne.score >= g.scoreLimit || g.playerTwo.score >= g.scoreLimit {
		g.state = GameOver
		return
	}

	if g.state == Waiting {
		g.waitCount++
		if g.wait <= g.waitCount {
			g.waitCount = 0
			g.state = Playing
		}
		return
	}

	if g.ball.y <= 0 || g.ball.y >= height-1 {
		g.ball.speedY = -g.ball.speedY
	}

	if g.ball.x <= 0 {
		g.playerTwo.score++
		g.state = Waiting
		g.ball = newBall()
		return
	}

	if g.ball.x >= width-1 {
		g.playerOne.score++
		g.state = Waiting
		g.ball = newBall()
		return
	}

	if g.ball.x == 3 && g.ball.y >= g.playerOne.y && g.ball.y <= g.playerOne.y+g.paddleWidth {
		g.ball.speedX = -g.ball.speedX
	}
	if g.ball.x == width-4 && g.ball.y >= g.playerTwo.y && g.ball.y <= g.playerTwo.y+g.paddleWidth {
		g.ball.speedX = -g.ball.speedX
	}

	g.ball.x += g.ball.speedX
	g.ball.y += g.ball.speedY
}

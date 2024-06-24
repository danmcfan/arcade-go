package internal

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
)

type Direction rune

const (
	Up    Direction = 'w'
	Down  Direction = 's'
	Left  Direction = 'a'
	Right Direction = 'd'
)

type Position struct {
	x, y int
}

type Game struct {
	width, height int
	playing       bool
	score         int

	player    Position
	direction Direction

	apple Position
}

func Run() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	game := newGame()
	gameTick := time.NewTicker(100 * time.Millisecond)
	defer gameTick.Stop()
	keyboardEvents := make(chan termbox.Event)
	exit := make(chan bool)

	go func() {
		for {
			keyboardEvents <- termbox.PollEvent()
		}
	}()

	go gameLoop(&game, gameTick, keyboardEvents, exit)

	<-exit
	termbox.Close()
	fmt.Println("Thanks for playing! Your final score was:", game.score)
	fmt.Println("Goodbye!")
}

func newGame() Game {
	width := 40
	height := 20

	return Game{
		width:     width,
		height:    height,
		playing:   true,
		score:     0,
		player:    Position{x: width / 2, y: height / 2},
		direction: Right,
		apple:     Position{x: rand.Intn(width), y: rand.Intn(height)},
	}
}

func gameLoop(game *Game, gameTick *time.Ticker, keyboardEvents <-chan termbox.Event, exit chan<- bool) {
	for {
		select {
		case ev := <-keyboardEvents:
			if ev.Type == termbox.EventKey {
				switch ev.Ch {
				case 'w', 's', 'a', 'd', 'r':
					if game.playing {
						handlePlayingInput(game, ev.Ch)
					} else {
						handleGameOverInput(game, ev.Ch)
					}
				case 'q':
					exit <- true
				}

			}
		case <-gameTick.C:
			if game.playing {
				game.update()
			}
			draw(*game)
		}
	}
}

func handlePlayingInput(g *Game, ch rune) {
	switch ch {
	case 'w', 's', 'a', 'd':
		g.direction = Direction(ch)
	}
}

func handleGameOverInput(g *Game, ch rune) {
	if ch == 'r' {
		*g = newGame()
	}
}

func (g *Game) update() {
	switch g.direction {
	case Up:
		g.player.y--
	case Down:
		g.player.y++
	case Left:
		g.player.x--
	case Right:
		g.player.x++
	}

	if g.player.x < 0 || g.player.x >= g.width || g.player.y < 0 || g.player.y >= g.height {
		g.playing = false
		return
	}

	if g.player.x == g.apple.x && g.player.y == g.apple.y {
		g.score++
		g.apple = Position{x: rand.Intn(g.width), y: rand.Intn(g.height)}
	}
}

func draw(g Game) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	if g.playing {
		drawPlaying(g)
	} else {
		drawGameOver(g)
	}

	termbox.Flush()
}

func drawPlaying(g Game) {
	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			termbox.SetCell(x, y, '·', termbox.ColorWhite, termbox.ColorDefault)
		}
	}

	termbox.SetCell(g.player.x, g.player.y, '■', termbox.ColorGreen, termbox.ColorDefault)
	termbox.SetCell(g.apple.x, g.apple.y, '■', termbox.ColorRed, termbox.ColorDefault)

	scoreStr := fmt.Sprintf("Score: %d", g.score)
	for i, ch := range scoreStr {
		termbox.SetCell(i, g.height+1, ch, termbox.ColorYellow, termbox.ColorDefault)
	}
}

func drawGameOver(g Game) {
	gameOverStr := "Game Over!"
	scoreStr := fmt.Sprintf("Your score: %d", g.score)
	restartStr := "Press 'R' to restart or 'Q' to quit"

	centerX := g.width / 2
	centerY := g.height / 2

	for i, ch := range gameOverStr {
		termbox.SetCell(centerX-len(gameOverStr)/2+i, centerY-2, ch, termbox.ColorRed, termbox.ColorDefault)
	}
	for i, ch := range scoreStr {
		termbox.SetCell(centerX-len(scoreStr)/2+i, centerY, ch, termbox.ColorYellow, termbox.ColorDefault)
	}
	for i, ch := range restartStr {
		termbox.SetCell(centerX-len(restartStr)/2+i, centerY+2, ch, termbox.ColorWhite, termbox.ColorDefault)
	}
}

package snake

import (
	"fmt"
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

func Run(size, tickMS, numApples int) {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	if size < 10 {
		panic("Size must be greater than or equal to 10")
	}

	game := newGame(size, numApples)
	gameTick := time.NewTicker(time.Duration(tickMS) * time.Millisecond)
	defer gameTick.Stop()
	keyboardEvents := make(chan termbox.Event, 1)
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

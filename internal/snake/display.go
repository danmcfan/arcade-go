package snake

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

func draw(g Game) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	if g.state == Idle || g.state == Playing {
		drawPlaying(g)
	}
	if g.state == GameOver {
		drawGameOver(g)
	}

	termbox.Flush()
}

func drawPlaying(g Game) {
	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width*2; x++ {
			termbox.SetCell(x, y, ' ', termbox.ColorDarkGray, termbox.ColorDarkGray)
		}
	}

	head := g.snake.head()
	setPixel(head.x, head.y, ' ', termbox.ColorDefault, termbox.ColorGreen)

	for _, position := range g.snake.tail() {
		setPixel(position.x, position.y, ' ', termbox.ColorDefault, termbox.ColorLightGreen)
	}
	for _, apple := range g.apples {
		setPixel(apple.x, apple.y, ' ', termbox.ColorDefault, termbox.ColorRed)
	}

	scoreStr := fmt.Sprintf("Score: %d", g.score)
	for i, ch := range scoreStr {
		termbox.SetCell(i, g.height+1, ch, termbox.ColorYellow, termbox.ColorDefault)
	}
}

func setPixel(x, y int, ch rune, fg, bg termbox.Attribute) {
	termbox.SetCell(x*2, y, ch, fg, bg)
	termbox.SetCell(x*2+1, y, ch, fg, bg)
}

func drawGameOver(g Game) {
	gameOverStr := "Game Over!"
	scoreStr := fmt.Sprintf("Your score: %d", g.score)
	restartStr := "Press 'R' to restart"
	quitStr := "Press 'Q' to quit"

	centerX := g.width
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
	for i, ch := range quitStr {
		termbox.SetCell(centerX-len(quitStr)/2+i, centerY+3, ch, termbox.ColorWhite, termbox.ColorDefault)
	}
}

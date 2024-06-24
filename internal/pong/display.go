package pong

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

func draw(g Game) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	if g.state == Waiting || g.state == Playing {
		drawPlaying(g)
	}
	if g.state == GameOver {
		drawGameOver(g)
	}

	termbox.Flush()
}

func drawPlaying(g Game) {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			setPixel(x, y, ' ', termbox.ColorDarkGray, termbox.ColorDarkGray)
		}
	}

	for y := 0; y < height; y += 2 {
		setPixel(centerX, y, '$', termbox.ColorWhite, termbox.ColorDarkGray)
	}

	for i := range g.paddleWidth {
		setPixel(2, g.playerOne.y+i, ' ', termbox.ColorWhite, termbox.ColorWhite)
		setPixel(width-3, g.playerTwo.y+i, ' ', termbox.ColorWhite, termbox.ColorWhite)
	}

	setPixel(g.ball.x, g.ball.y, ' ', termbox.ColorWhite, termbox.ColorWhite)

	playerOneScoreStr := fmt.Sprintf("Player One: %d", g.playerOne.score)
	for i, ch := range playerOneScoreStr {
		termbox.SetCell(i, height+1, ch, termbox.ColorYellow, termbox.ColorDefault)
	}

	playerTwoScoreStr := fmt.Sprintf("Player Two: %d", g.playerTwo.score)
	for i, ch := range playerTwoScoreStr {
		termbox.SetCell(width*2-len(playerTwoScoreStr)+i, height+1, ch, termbox.ColorYellow, termbox.ColorDefault)
	}
}

func setPixel(x, y int, ch rune, fg, bg termbox.Attribute) {
	termbox.SetCell(x*2, y, ch, fg, bg)
	termbox.SetCell(x*2+1, y, ch, fg, bg)
}

func drawGameOver(g Game) {
	var gameOverStr string
	if g.playerOne.score > g.playerTwo.score {
		gameOverStr = "Player one wins!"
	} else {
		gameOverStr = "Player two wins!"
	}

	restartStr := "Press 'R' to restart"
	quitStr := "Press 'Q' to quit"

	centerX := width
	centerY := height / 2

	for i, ch := range gameOverStr {
		termbox.SetCell(centerX-len(gameOverStr)/2+i, centerY-2, ch, termbox.ColorRed, termbox.ColorDefault)
	}

	for i, ch := range restartStr {
		termbox.SetCell(centerX-len(restartStr)/2+i, centerY+2, ch, termbox.ColorWhite, termbox.ColorDefault)
	}
	for i, ch := range quitStr {
		termbox.SetCell(centerX-len(quitStr)/2+i, centerY+3, ch, termbox.ColorWhite, termbox.ColorDefault)
	}
}

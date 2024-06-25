package pacman

import (
	"github.com/nsf/termbox-go"
)

func draw(g Game) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	if g.state == Playing {
		drawPlaying(g)
	} else {
		drawGameOver()
	}

	termbox.Flush()
}

func drawPlaying(g Game) {
	for y, row := range g.maze {
		for x, pixel := range row {
			switch Pixel(pixel) {
			case Open:
				setPixel(x, y, termbox.ColorDarkGray)
			case Dot:
				setSymbol(x, y, ' ', '.', termbox.ColorWhite, termbox.ColorDarkGray)
			case PowerUp:
				setSymbol(x, y, ' ', '●', termbox.ColorWhite, termbox.ColorDarkGray)
			case Pacman:
				setPixel(x, y, termbox.ColorLightYellow)
			case RedGhost:
				setPixel(x, y, termbox.ColorLightRed)
			case PinkGhost:
				setPixel(x, y, termbox.ColorLightMagenta)
			case GreenGhost:
				setPixel(x, y, termbox.ColorLightGreen)
			case GrayGhost:
				setPixel(x, y, termbox.ColorLightGray)
			case Wall:
				setPixel(x, y, termbox.ColorBlue)
			case Gate:
				setSymbol(x, y, '=', '=', termbox.ColorLightGray, termbox.ColorDarkGray)
			}
		}
	}
}

func drawGameOver() {
	gameOverStr := "Game over!"
	restartStr := "Press 'R' to restart"
	quitStr := "Press 'Q' to quit"

	centerX := Width
	centerY := Height / 2

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

func setPixel(x, y int, color termbox.Attribute) {
	termbox.SetCell(x*2, y, ' ', termbox.ColorDefault, color)
	termbox.SetCell(x*2+1, y, ' ', termbox.ColorDefault, color)
}

func setSymbol(x, y int, firstCh, secondCh rune, fg, bg termbox.Attribute) {
	termbox.SetCell(x*2, y, firstCh, fg, bg)
	termbox.SetCell(x*2+1, y, secondCh, fg, bg)
}

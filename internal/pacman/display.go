package pacman

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

func draw(g Game) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	if g.state == GameOver {
		drawGameOver()
	} else {
		drawPlaying(g)
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
			case PlayerPixel:
				setPixel(x, y, termbox.ColorLightYellow)
			case RedGhost, PinkGhost, GreenGhost, GrayGhost:
				ghost := g.ghosts[Pixel(pixel)]
				if ghost.weakTicks > 0 {
					setPixel(x, y, termbox.ColorLightBlue)
				} else {
					var attribute termbox.Attribute
					switch Pixel(pixel) {
					case RedGhost:
						attribute = termbox.ColorLightRed
					case PinkGhost:
						attribute = termbox.ColorLightMagenta
					case GreenGhost:
						attribute = termbox.ColorLightGreen
					case GrayGhost:
						attribute = termbox.ColorLightGray
					}
					setPixel(x, y, attribute)
				}
			case Wall:
				setPixel(x, y, termbox.ColorBlue)
			case Gate:
				setSymbol(x, y, '=', '=', termbox.ColorLightGray, termbox.ColorDarkGray)
			}
		}
	}

	scoreStr := fmt.Sprintf("Score: %d", g.score)
	for i, ch := range scoreStr {
		termbox.SetCell(Width-len(scoreStr)/2+i, Height+1, ch, termbox.ColorYellow, termbox.ColorDefault)
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

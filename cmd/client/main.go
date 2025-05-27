//go:build js

package main

import (
	"fmt"
	"syscall/js"

	"github.com/danmcfan/arcade-go/internal"
)

const (
	LargeScreen  = 1050
	MediumScreen = 750

	SmallCellSize  = 16
	MediumCellSize = 24
	LargeCellSize  = 32

	GridWidth    = 20
	GridHeight   = 20
	TickInterval = 100.0
)

type Color string

const (
	Background Color = "#3f3f46"
	SnakeHead  Color = "#009966"
	SnakeBody  Color = "#00bc7d"
	Apple      Color = "#fb2c36"
)

func handleKey(key string) internal.Signal {
	switch key {
	case "ArrowUp", "KeyW":
		return internal.SignalUp
	case "ArrowDown", "KeyS":
		return internal.SignalDown
	case "ArrowLeft", "KeyA":
		return internal.SignalLeft
	case "ArrowRight", "KeyD":
		return internal.SignalRight
	case "KeyR":
		return internal.SignalRestart
	default:
		return internal.SignalNone
	}
}

func main() {
	fmt.Println("Starting a new game client...")

	document := js.Global().Get("document")
	window := js.Global().Get("window")

	root := document.Call("getElementById", "root")

	snakeButton := document.Call("getElementById", "snake")
	pongButton := document.Call("getElementById", "pong")
	pacmanButton := document.Call("getElementById", "pacman")

	buttons := []js.Value{snakeButton, pongButton, pacmanButton}
	selectedButton := "snake"

	for _, button := range buttons {
		button.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) any {
			id := this.Get("id").String()
			if id == selectedButton {
				return nil
			}

			for _, button := range buttons {
				classList := button.Get("classList")
				classList.Call("remove", "selected")
			}

			classList := this.Get("classList")
			classList.Call("add", "selected")
			selectedButton = id
			return nil
		}))
	}

	canvas := document.Call("getElementById", "canvas")
	ctx := canvas.Call("getContext", "2d")

	viewportWidth := window.Get("innerWidth").Float()
	viewportHeight := window.Get("innerHeight").Float()

	CellSize := LargeCellSize
	if viewportWidth <= MediumScreen || viewportHeight <= MediumScreen {
		CellSize = SmallCellSize
	} else if viewportWidth <= LargeScreen || viewportHeight <= LargeScreen {
		CellSize = MediumCellSize
	}

	Width := CellSize * GridWidth
	Height := CellSize * GridHeight

	canvas.Set("width", Width)
	canvas.Set("height", Height)

	rootClassList := root.Get("classList")
	rootClassList.Call("remove", "hidden")

	game, lastUpdate := restartGame(CellSize, Width, Height)

	window.Call("addEventListener", "resize", js.FuncOf(func(this js.Value, args []js.Value) any {
		handleResize(window, canvas, &game)
		return nil
	}))

	keyHandler := js.FuncOf(func(this js.Value, args []js.Value) any {
		event := args[0]
		code := event.Get("code").String()

		signal := handleKey(code)

		if signal == internal.SignalNone {
			return nil
		}

		if signal == internal.SignalRestart {
			game, lastUpdate = restartGame(game.CellSize, game.Width, game.Height)
			return nil
		}

		if signal == internal.SignalLeft && game.CurrentDirection == internal.DirectionNone {
			return nil
		}

		if signal == internal.SignalUp && game.CurrentDirection == internal.DirectionDown {
			return nil
		}

		if signal == internal.SignalDown && game.CurrentDirection == internal.DirectionUp {
			return nil
		}

		if signal == internal.SignalLeft && game.CurrentDirection == internal.DirectionRight {
			return nil
		}

		if signal == internal.SignalRight && game.CurrentDirection == internal.DirectionLeft {
			return nil
		}

		switch signal {
		case internal.SignalUp:
			game.CurrentDirection = internal.DirectionUp
		case internal.SignalDown:
			game.CurrentDirection = internal.DirectionDown
		case internal.SignalLeft:
			game.CurrentDirection = internal.DirectionLeft
		case internal.SignalRight:
			game.CurrentDirection = internal.DirectionRight
		}

		return nil
	})
	document.Call("addEventListener", "keydown", keyHandler)

	var renderFrame js.Func
	renderFrame = js.FuncOf(func(this js.Value, args []js.Value) any {
		timestamp := args[0].Float()
		sinceLastUpdate := timestamp - lastUpdate

		if sinceLastUpdate >= TickInterval {
			update(&game)
			lastUpdate = timestamp
		}

		render(ctx, game)

		js.Global().Call("requestAnimationFrame", renderFrame)
		return nil
	})

	js.Global().Call("requestAnimationFrame", renderFrame)

	select {}
}

func update(game *internal.Game) {
	if game.GameOver {
		return
	}

	if game.CurrentDirection == internal.DirectionNone {
		return
	}

	game.Snake.Move(game.CurrentDirection)

	if game.Snake.Head().X < 0 || game.Snake.Head().X >= GridWidth || game.Snake.Head().Y < 0 || game.Snake.Head().Y >= GridHeight {
		game.GameOver = true
		return
	}

	if game.Snake.TailContains(game.Snake.Head()) {
		game.GameOver = true
		return
	}

	if game.Snake.Head() == game.Apple {
		setScore(game.Score+1, game)
		game.Apple = internal.NewApple(game.Snake, 0, GridWidth, 0, GridHeight)
		return
	}

	game.Snake.CutTail()
}

func render(ctx js.Value, game internal.Game) {
	ctx.Call("clearRect", 0, 0, game.Width, game.Height)
	ctx.Set("fillStyle", string(Background))
	for x := range GridWidth {
		for y := range GridHeight {
			fillCell(ctx, x, y, game.CellSize)
		}
	}

	if game.GameOver {
		ctx.Set("fillStyle", "#fb2c36")
		ctx.Set("font", "64px monospace")
		ctx.Set("textAlign", "center")
		ctx.Set("textBaseline", "middle")
		ctx.Call("fillText", "GAME OVER", game.Width/2, game.Height/2)
		return
	}

	ctx.Set("fillStyle", string(SnakeHead))
	fillCell(ctx, game.Snake.Head().X, game.Snake.Head().Y, game.CellSize)

	ctx.Set("fillStyle", string(SnakeBody))
	for _, tail := range game.Snake.Tail() {
		fillCell(ctx, tail.X, tail.Y, game.CellSize)
	}

	ctx.Set("fillStyle", string(Apple))
	fillCell(ctx, game.Apple.X, game.Apple.Y, game.CellSize)
}

func fillCell(ctx js.Value, x, y int, cellSize int) {
	ctx.Call("fillRect", x*cellSize, y*cellSize, cellSize, cellSize)
}

func restartGame(cellSize, width, height int) (internal.Game, float64) {
	midX := GridWidth / 2
	midY := GridHeight / 2
	positions := make([]internal.Vec2, 0, 64)
	positions = append(positions, internal.Vec2{X: midX - 2, Y: midY})
	positions = append(positions, internal.Vec2{X: midX - 1, Y: midY})
	positions = append(positions, internal.Vec2{X: midX, Y: midY})

	snake := internal.NewSnake(positions)
	apple := internal.NewApple(snake, 0, GridWidth, 0, GridHeight)

	game := internal.NewGame(cellSize, width, height, snake, apple)
	lastUpdate := 0.0

	setScore(0, &game)

	return game, lastUpdate
}

func setScore(score int, game *internal.Game) {
	game.Score = score

	document := js.Global().Get("document")
	scoreElement := document.Call("getElementById", "score")
	scoreElement.Set("textContent", score)
}

func handleResize(window js.Value, canvas js.Value, game *internal.Game) {
	viewportWidth := window.Get("innerWidth").Float()
	viewportHeight := window.Get("innerHeight").Float()

	CellSize := LargeCellSize
	if viewportWidth <= MediumScreen || viewportHeight <= MediumScreen {
		CellSize = SmallCellSize
	} else if viewportWidth <= LargeScreen || viewportHeight <= LargeScreen {
		CellSize = MediumCellSize
	}

	Width := CellSize * GridWidth
	Height := CellSize * GridHeight

	canvas.Set("width", Width)
	canvas.Set("height", Height)

	game.CellSize = CellSize
	game.Width = Width
	game.Height = Height
}

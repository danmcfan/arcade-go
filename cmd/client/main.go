//go:build js

package main

import (
	"math/rand"
	"syscall/js"

	"github.com/danmcfan/arcade-go/internal"
)

const (
	LargeScreen  = 1050
	MediumScreen = 750

	SmallCellSize  = 16
	MediumCellSize = 24
	LargeCellSize  = 32
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
	document := js.Global().Get("document")
	window := js.Global().Get("window")

	root := document.Call("getElementById", "root")

	snakeButton := document.Call("getElementById", "snake")
	pongButton := document.Call("getElementById", "pong")
	pacmanButton := document.Call("getElementById", "pacman")

	scoreboard := document.Call("getElementById", "scoreboard")
	instructions := document.Call("getElementById", "instructions")

	buttons := []js.Value{snakeButton, pongButton, pacmanButton}

	selectedGame := "snake"
	gridWidth := 20
	gridHeight := 20
	tickInterval := 100.0

	canvas := document.Call("getElementById", "canvas")
	ctx := canvas.Call("getContext", "2d")

	display := internal.NewDisplay(gridWidth, gridHeight)
	handleResize(window, canvas, &display)

	snakeGame := restartSnakeGame(display)
	pongGame := internal.PongGame{}
	lastUpdate := 0.0

	rootClassList := root.Get("classList")
	rootClassList.Call("remove", "hidden")

	window.Call("addEventListener", "resize", js.FuncOf(func(this js.Value, args []js.Value) any {
		handleResize(window, canvas, &display)
		return nil
	}))

	document.Call("addEventListener", "keydown", js.FuncOf(func(this js.Value, args []js.Value) any {
		event := args[0]
		code := event.Get("code").String()

		signal := handleKey(code)

		if signal == internal.SignalNone {
			return nil
		}

		if signal == internal.SignalRestart {
			switch selectedGame {
			case "snake":
				snakeGame = restartSnakeGame(display)
			case "pong":
				pongGame = restartPongGame(display)
			}

			lastUpdate = 0.0
			return nil
		}

		switch selectedGame {
		case "snake":
			handleSnakeDirectionSignal(signal, &snakeGame)
		case "pong":
			handlePongDirectionSignal(signal, &pongGame)
		}

		return nil
	}))

	for _, button := range buttons {
		button.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) any {
			id := this.Get("id").String()
			if id == selectedGame {
				return nil
			}

			for _, button := range buttons {
				classList := button.Get("classList")
				classList.Call("remove", "selected")
			}

			classList := this.Get("classList")
			classList.Call("add", "selected")
			selectedGame = id

			switch id {
			case "snake":
				display.GridWidth = 20
				display.GridHeight = 20
				tickInterval = 100.0

				scoreboard.Set("innerHTML", "Score: <span id='score'>0</span>")
				instructions.Set("innerHTML", "Use <span>W</span>, <span>A</span>, <span>S</span>, <span>D</span> or <span>Arrow Keys</span> to move<br/>Press <span>R</span> to restart")

				snakeGame = restartSnakeGame(display)
			case "pong":
				display.GridWidth = 41
				display.GridHeight = 21
				tickInterval = 25.0

				scoreboard.Set("innerHTML", "<span id='player1'>0</span> - <span id='player2'>0</span>")
				instructions.Set("innerHTML", "Use <span>W</span> and <span>S</span> or <span>Arrow Keys</span> to move up and down<br/>Press <span>R</span> to restart")

				pongGame = restartPongGame(display)
			}

			handleResize(window, canvas, &display)

			return nil
		}))
	}

	var renderFrame js.Func
	renderFrame = js.FuncOf(func(this js.Value, args []js.Value) any {
		timestamp := args[0].Float()
		sinceLastUpdate := timestamp - lastUpdate

		if sinceLastUpdate >= tickInterval {
			switch selectedGame {
			case "snake":
				updateSnake(&snakeGame, display)
			case "pong":
				updatePong(&pongGame, display)
			}
			lastUpdate = timestamp
		}

		render(ctx, display, snakeGame, pongGame, selectedGame)

		js.Global().Call("requestAnimationFrame", renderFrame)
		return nil
	})
	js.Global().Call("requestAnimationFrame", renderFrame)

	select {}
}

func updateSnake(game *internal.SnakeGame, display internal.Display) {
	if game.GameOver {
		return
	}

	if game.CurrentDirection == internal.DirectionNone {
		return
	}

	game.Snake.Move(game.CurrentDirection)

	if game.Snake.Head().X < 0 || game.Snake.Head().X >= display.GridWidth || game.Snake.Head().Y < 0 || game.Snake.Head().Y >= display.GridHeight {
		game.GameOver = true
		return
	}

	if game.Snake.TailContains(game.Snake.Head()) {
		game.GameOver = true
		return
	}

	if game.Snake.Head() == game.Apple {
		setScore(game.Score+1, game)
		game.Apple = internal.NewApple(game.Snake, 0, display.GridWidth, 0, display.GridHeight)
		return
	}

	game.Snake.CutTail()
}

func updatePong(game *internal.PongGame, display internal.Display) {
	if game.CurrentDirection == internal.DirectionNone {
		return
	}

	switch game.CurrentDirection {
	case internal.DirectionUp:
		game.Player1Position--
		if game.Player1Position-2 < 0 {
			game.Player1Position = 2
		}
	case internal.DirectionDown:
		game.Player1Position++
		if game.Player1Position+2 > display.GridHeight-1 {
			game.Player1Position = display.GridHeight - 3
		}
	}

	if rand.Intn(100) < 70 {
		if game.Player2Position > game.BallPosition.Y {
			game.Player2Position--
			if game.Player2Position-2 < 0 {
				game.Player2Position = 2
			}
		}
		if game.Player2Position < game.BallPosition.Y {
			game.Player2Position++
			if game.Player2Position+2 > display.GridHeight-1 {
				game.Player2Position = display.GridHeight - 3
			}
		}
	}

	game.BallPosition.X += game.BallVelocity.X
	game.BallPosition.Y += game.BallVelocity.Y

	if game.BallPosition.X == 2 && game.Player1Position-2 <= game.BallPosition.Y && game.BallPosition.Y <= game.Player1Position+2 {
		game.BallVelocity.X *= -1
	}

	if game.BallPosition.X == display.GridWidth-3 && game.Player2Position-2 <= game.BallPosition.Y && game.BallPosition.Y <= game.Player2Position+2 {
		game.BallVelocity.X *= -1
		game.BallVelocity.Y = rand.Intn(3) - 1
	}

	if game.BallPosition.X < 0 {
		game.CurrentDirection = internal.DirectionNone
		setPlayerTwoScore(game.Player2Score+1, game)
		game.BallPosition = internal.NewBallPosition(display.GridWidth, display.GridHeight)
		game.BallVelocity = internal.NewBallVelocity()
	}

	if game.BallPosition.X > display.GridWidth-1 {
		game.CurrentDirection = internal.DirectionNone
		setPlayerOneScore(game.Player1Score+1, game)
		game.BallPosition = internal.NewBallPosition(display.GridWidth, display.GridHeight)
		game.BallVelocity = internal.NewBallVelocity()
	}

	if game.BallPosition.Y <= 0 || game.BallPosition.Y >= display.GridHeight-1 {
		game.BallVelocity.Y *= -1
	}

	if game.BallPosition.Y < 0 {
		game.BallPosition.Y = 0
	}

	if game.BallPosition.Y > display.GridHeight-1 {
		game.BallPosition.Y = display.GridHeight - 1
	}
}

func render(ctx js.Value, display internal.Display, snakeGame internal.SnakeGame, pongGame internal.PongGame, selectedGame string) {
	internal.RenderBackground(ctx, display)

	switch selectedGame {
	case "snake":
		renderSnake(ctx, display, snakeGame)
	case "pong":
		renderPong(ctx, display, pongGame)
	}
}

func restartSnakeGame(display internal.Display) internal.SnakeGame {
	midX := display.GridWidth / 2
	midY := display.GridHeight / 2
	positions := make([]internal.Vec2, 0, 64)
	positions = append(positions, internal.Vec2{X: midX - 2, Y: midY})
	positions = append(positions, internal.Vec2{X: midX - 1, Y: midY})
	positions = append(positions, internal.Vec2{X: midX, Y: midY})

	snake := internal.NewSnake(positions)
	apple := internal.NewApple(snake, 0, display.GridWidth, 0, display.GridHeight)

	game := internal.NewSnakeGame(snake, apple)

	setScore(0, &game)

	return game
}

func restartPongGame(display internal.Display) internal.PongGame {
	game := internal.NewPongGame(display.GridWidth, display.GridHeight)
	setPlayerOneScore(0, &game)
	setPlayerTwoScore(0, &game)
	return game
}

func setScore(score int, game *internal.SnakeGame) {
	game.Score = score

	document := js.Global().Get("document")
	scoreElement := document.Call("getElementById", "score")
	scoreElement.Set("textContent", score)
}

func setPlayerOneScore(score int, game *internal.PongGame) {
	game.Player1Score = score

	document := js.Global().Get("document")
	scoreElement := document.Call("getElementById", "player1")
	scoreElement.Set("textContent", score)
}

func setPlayerTwoScore(score int, game *internal.PongGame) {
	game.Player2Score = score

	document := js.Global().Get("document")
	scoreElement := document.Call("getElementById", "player2")
	scoreElement.Set("textContent", score)
}

func handleResize(window js.Value, canvas js.Value, display *internal.Display) {
	viewportWidth := window.Get("innerWidth").Float()
	viewportHeight := window.Get("innerHeight").Float()

	CellSize := LargeCellSize
	if viewportWidth <= MediumScreen || viewportHeight <= MediumScreen {
		CellSize = SmallCellSize
	} else if viewportWidth <= LargeScreen || viewportHeight <= LargeScreen {
		CellSize = MediumCellSize
	}

	Width := CellSize * display.GridWidth
	Height := CellSize * display.GridHeight

	canvas.Set("width", Width)
	canvas.Set("height", Height)

	display.CellSize = CellSize
	display.Width = Width
	display.Height = Height
}

func handleSnakeDirectionSignal(signal internal.Signal, game *internal.SnakeGame) {
	if signal == internal.SignalLeft && game.CurrentDirection == internal.DirectionNone {
		return
	}

	if signal == internal.SignalUp && game.CurrentDirection == internal.DirectionDown {
		return
	}

	if signal == internal.SignalDown && game.CurrentDirection == internal.DirectionUp {
		return
	}

	if signal == internal.SignalLeft && game.CurrentDirection == internal.DirectionRight {
		return
	}

	if signal == internal.SignalRight && game.CurrentDirection == internal.DirectionLeft {
		return
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
}

func handlePongDirectionSignal(signal internal.Signal, game *internal.PongGame) {
	switch signal {
	case internal.SignalUp:
		game.CurrentDirection = internal.DirectionUp
	case internal.SignalDown:
		game.CurrentDirection = internal.DirectionDown
	}
}

func renderSnake(ctx js.Value, display internal.Display, game internal.SnakeGame) {
	if game.GameOver {
		ctx.Set("fillStyle", string(internal.Red))
		ctx.Set("font", "64px monospace")
		ctx.Set("textAlign", "center")
		ctx.Set("textBaseline", "middle")
		ctx.Call("fillText", "GAME OVER", display.Width/2, display.Height/2)
		return
	}

	ctx.Set("fillStyle", string(internal.DarkGreen))
	internal.FillCell(ctx, game.Snake.Head().X, game.Snake.Head().Y, display.CellSize)

	ctx.Set("fillStyle", string(internal.Green))
	for _, tail := range game.Snake.Tail() {
		internal.FillCell(ctx, tail.X, tail.Y, display.CellSize)
	}

	ctx.Set("fillStyle", string(internal.Red))
	internal.FillCell(ctx, game.Apple.X, game.Apple.Y, display.CellSize)
}

func renderPong(ctx js.Value, display internal.Display, game internal.PongGame) {
	ctx.Set("fillStyle", string(internal.Blue))
	for i := -2; i <= 2; i++ {
		if y := game.Player1Position + i; y >= 0 && y < display.GridHeight {
			internal.FillCell(ctx, 1, y, display.CellSize)
		}
	}

	ctx.Set("fillStyle", string(internal.Red))
	for i := -2; i <= 2; i++ {
		if y := game.Player2Position + i; y >= 0 && y < display.GridHeight {
			internal.FillCell(ctx, display.GridWidth-2, y, display.CellSize)
		}
	}

	ctx.Set("fillStyle", string(internal.White))
	internal.FillCell(ctx, game.BallPosition.X, game.BallPosition.Y, display.CellSize)
}

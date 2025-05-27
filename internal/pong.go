package internal

import "math/rand"

type PongGame struct {
	Player1Score     int
	Player2Score     int
	Player1Position  int
	Player2Position  int
	BallPosition     Vec2
	BallVelocity     Vec2
	CurrentDirection Direction
}

func NewPongGame(gridWidth, gridHeight int) PongGame {
	midY := gridHeight / 2

	return PongGame{
		Player1Score:     0,
		Player2Score:     0,
		Player1Position:  midY,
		Player2Position:  midY,
		BallPosition:     NewBallPosition(gridWidth, gridHeight),
		BallVelocity:     NewBallVelocity(),
		CurrentDirection: DirectionNone,
	}
}

func NewBallPosition(gridWidth, gridHeight int) Vec2 {
	midX := gridWidth / 2
	midY := gridHeight / 2
	return Vec2{X: midX, Y: midY}
}

func NewBallVelocity() Vec2 {
	ballVelocity := Vec2{X: rand.Intn(3) - 1, Y: rand.Intn(3) - 1}
	for ballVelocity.X == 0 {
		ballVelocity = Vec2{X: rand.Intn(3) - 1, Y: rand.Intn(3) - 1}
	}
	return ballVelocity
}

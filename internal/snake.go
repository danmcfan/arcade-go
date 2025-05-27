package internal

import (
	"math/rand"
	"slices"
)

type Signal string

const (
	SignalRestart Signal = "restart"
	SignalPause   Signal = "pause"
	SignalUp      Signal = "up"
	SignalDown    Signal = "down"
	SignalLeft    Signal = "left"
	SignalRight   Signal = "right"
	SignalNone    Signal = "none"
)

type Direction string

const (
	DirectionUp    Direction = "up"
	DirectionDown  Direction = "down"
	DirectionLeft  Direction = "left"
	DirectionRight Direction = "right"
	DirectionNone  Direction = "none"
)

type Vec2 struct {
	X, Y int
}

type Snake struct {
	positions []Vec2
}

type SnakeGame struct {
	Snake            Snake
	Apple            Vec2
	Score            int
	CurrentDirection Direction
	GameOver         bool
}

func NewSnakeGame(snake Snake, apple Vec2) SnakeGame {
	return SnakeGame{
		Snake:            snake,
		Apple:            apple,
		Score:            0,
		CurrentDirection: DirectionNone,
		GameOver:         false,
	}
}

func NewSnake(positions []Vec2) Snake {
	return Snake{positions: positions}
}

func (s Snake) Head() Vec2 {
	return s.positions[len(s.positions)-1]
}

func (s Snake) Tail() []Vec2 {
	return s.positions[:len(s.positions)-1]
}

func (s Snake) Contains(v Vec2) bool {
	return slices.Contains(s.positions, v)
}

func (s Snake) TailContains(v Vec2) bool {
	return slices.Contains(s.Tail(), v)
}

func (s *Snake) Move(direction Direction) {
	h := s.Head()
	switch direction {
	case DirectionUp:
		h.Y--
	case DirectionDown:
		h.Y++
	case DirectionLeft:
		h.X--
	case DirectionRight:
		h.X++
	}
	s.positions = append(s.positions, h)
}

func (s *Snake) CutTail() {
	s.positions = s.positions[1:]
}

func NewApple(s Snake, minX, maxX, minY, maxY int) Vec2 {

	apple := Vec2{
		X: rand.Intn(maxX-minX) + minX,
		Y: rand.Intn(maxY-minY) + minY,
	}

	for s.Contains(apple) {
		apple = Vec2{
			X: rand.Intn(maxX-minX) + minX,
			Y: rand.Intn(maxY-minY) + minY,
		}
	}

	return apple
}

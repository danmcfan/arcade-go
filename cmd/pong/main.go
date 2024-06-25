package main

import (
	"flag"
	"games/internal/pong"
)

func main() {
	tick := flag.Int("tick", 75, "ticker milliseconds")
	score := flag.Int("score", 3, "score to win")
	paddle := flag.Int("paddle", 5, "paddle width")

	flag.Parse()

	pong.Run(*tick, *score, *paddle)
}

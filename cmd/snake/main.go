package main

import (
	"flag"

	"github.com/danmcfan/arcade-go/internal/snake"
)

func main() {
	size := flag.Int("size", 20, "size of the grid")
	tick := flag.Int("tick", 100, "ticker milliseconds")
	apples := flag.Int("apples", 1, "number of apples")

	flag.Parse()

	snake.Run(*size, *tick, *apples)
}

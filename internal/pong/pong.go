package pong

import (
	"time"

	"github.com/nsf/termbox-go"
)

func Run(tickMS, scoreLimit, paddleWidth int) {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	game := newGame(scoreLimit, paddleWidth)
	gameTick := time.NewTicker(time.Duration(tickMS) * time.Millisecond)
	defer gameTick.Stop()
	keyboardEvents := make(chan termbox.Event, 2)
	exit := make(chan bool)

	go func() {
		for {
			keyboardEvents <- termbox.PollEvent()
		}
	}()

	go gameLoop(&game, gameTick, keyboardEvents, exit)

	<-exit
	termbox.Close()
}

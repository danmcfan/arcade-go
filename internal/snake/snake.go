package snake

import (
	"time"

	"github.com/nsf/termbox-go"
)

func Run(size, tickMS, numApples int) {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	if size < 10 {
		panic("Size must be greater than or equal to 10")
	}

	game := newGame(size, numApples)
	gameTick := time.NewTicker(time.Duration(tickMS) * time.Millisecond)
	defer gameTick.Stop()
	keyboardEvents := make(chan termbox.Event, 1)
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

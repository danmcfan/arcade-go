package pacman

import (
	"time"

	"github.com/nsf/termbox-go"
)

func Run() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	game := newGame()
	gameTick := time.NewTicker(time.Duration(175) * time.Millisecond)
	defer gameTick.Stop()
	keyboardEvents := make(chan termbox.Event, 1)
	exit := make(chan bool, 1)

	go func() {
		for {
			keyboardEvents <- termbox.PollEvent()
		}
	}()

	go game.loop(gameTick, keyboardEvents, exit)

	<-exit
	termbox.Close()
}

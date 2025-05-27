//go:build js

package internal

import "syscall/js"

type Display struct {
	Width      int
	Height     int
	CellSize   int
	GridWidth  int
	GridHeight int
}

type Color string

const (
	Background Color = "#3f3f46"
	DarkGreen  Color = "#009966"
	Green      Color = "#00bc7d"
	Red        Color = "#fb2c36"
	White      Color = "#ffffff"
	Blue       Color = "#2b7fff"
)

func NewDisplay(gridWidth, gridHeight int) Display {
	return Display{
		Width:      0,
		Height:     0,
		CellSize:   0,
		GridWidth:  gridWidth,
		GridHeight: gridHeight,
	}
}

func RenderBackground(ctx js.Value, display Display) {
	ctx.Call("clearRect", 0, 0, display.Width, display.Height)
	ctx.Set("fillStyle", string(Background))
	for x := range display.GridWidth {
		for y := range display.GridHeight {
			FillCell(ctx, x, y, display.CellSize)
		}
	}
}

func FillCell(ctx js.Value, x, y int, cellSize int) {
	ctx.Call("fillRect", x*cellSize, y*cellSize, cellSize, cellSize)
}

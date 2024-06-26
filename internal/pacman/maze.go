package pacman

const (
	Width  = 28
	Height = 31
)

type Maze [Height][Width]int

func newMaze() Maze {
	return Maze{
		{100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100},
		{100, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 100, 100, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 100},
		{100, 1, 100, 100, 100, 100, 1, 100, 100, 100, 100, 100, 1, 100, 100, 1, 100, 100, 100, 100, 100, 1, 100, 100, 100, 100, 1, 100},
		{100, 2, 100, 100, 100, 100, 1, 100, 100, 100, 100, 100, 1, 100, 100, 1, 100, 100, 100, 100, 100, 1, 100, 100, 100, 100, 2, 100},
		{100, 1, 100, 100, 100, 100, 1, 100, 100, 100, 100, 100, 1, 100, 100, 1, 100, 100, 100, 100, 100, 1, 100, 100, 100, 100, 1, 100},
		{100, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 100},
		{100, 1, 100, 100, 100, 100, 1, 100, 100, 1, 100, 100, 100, 100, 100, 100, 100, 100, 1, 100, 100, 1, 100, 100, 100, 100, 1, 100},
		{100, 1, 100, 100, 100, 100, 1, 100, 100, 1, 100, 100, 100, 100, 100, 100, 100, 100, 1, 100, 100, 1, 100, 100, 100, 100, 1, 100},
		{100, 1, 1, 1, 1, 1, 1, 100, 100, 1, 1, 1, 1, 100, 100, 1, 1, 1, 1, 100, 100, 1, 1, 1, 1, 1, 1, 100},
		{100, 100, 100, 100, 100, 100, 1, 100, 100, 100, 100, 100, 0, 100, 100, 0, 100, 100, 100, 100, 100, 1, 100, 100, 100, 100, 100, 100},
		{100, 100, 100, 100, 100, 100, 1, 100, 100, 100, 100, 100, 0, 100, 100, 0, 100, 100, 100, 100, 100, 1, 100, 100, 100, 100, 100, 100},
		{100, 100, 100, 100, 100, 100, 1, 100, 100, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100, 100, 1, 100, 100, 100, 100, 100, 100},
		{100, 100, 100, 100, 100, 100, 1, 100, 100, 0, 100, 100, 100, 101, 101, 100, 100, 100, 0, 100, 100, 1, 100, 100, 100, 100, 100, 100},
		{100, 100, 100, 100, 100, 100, 1, 100, 100, 0, 100, 0, 0, 0, 0, 0, 0, 100, 0, 100, 100, 1, 100, 100, 100, 100, 100, 100},
		{0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 100, 0, 0, 0, 0, 0, 0, 100, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0},
		{100, 100, 100, 100, 100, 100, 1, 100, 100, 0, 100, 0, 0, 0, 0, 0, 0, 100, 0, 100, 100, 1, 100, 100, 100, 100, 100, 100},
		{100, 100, 100, 100, 100, 100, 1, 100, 100, 0, 100, 100, 100, 100, 100, 100, 100, 100, 0, 100, 100, 1, 100, 100, 100, 100, 100, 100},
		{100, 100, 100, 100, 100, 100, 1, 100, 100, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100, 100, 1, 100, 100, 100, 100, 100, 100},
		{100, 100, 100, 100, 100, 100, 1, 100, 100, 0, 100, 100, 100, 100, 100, 100, 100, 100, 0, 100, 100, 1, 100, 100, 100, 100, 100, 100},
		{100, 100, 100, 100, 100, 100, 1, 100, 100, 0, 100, 100, 100, 100, 100, 100, 100, 100, 0, 100, 100, 1, 100, 100, 100, 100, 100, 100},
		{100, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 100, 100, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 100},
		{100, 1, 100, 100, 100, 100, 1, 100, 100, 100, 100, 100, 1, 100, 100, 1, 100, 100, 100, 100, 100, 1, 100, 100, 100, 100, 1, 100},
		{100, 1, 100, 100, 100, 100, 1, 100, 100, 100, 100, 100, 1, 100, 100, 1, 100, 100, 100, 100, 100, 1, 100, 100, 100, 100, 1, 100},
		{100, 2, 1, 1, 100, 100, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 100, 100, 1, 1, 2, 100},
		{100, 100, 100, 1, 100, 100, 1, 100, 100, 1, 100, 100, 100, 100, 100, 100, 100, 100, 1, 100, 100, 1, 100, 100, 1, 100, 100, 100},
		{100, 100, 100, 1, 100, 100, 1, 100, 100, 1, 100, 100, 100, 100, 100, 100, 100, 100, 1, 100, 100, 1, 100, 100, 1, 100, 100, 100},
		{100, 1, 1, 1, 1, 1, 1, 100, 100, 1, 1, 1, 1, 100, 100, 1, 1, 1, 1, 100, 100, 1, 1, 1, 1, 1, 1, 100},
		{100, 1, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 1, 100, 100, 1, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 1, 100},
		{100, 1, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 1, 100, 100, 1, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 1, 100},
		{100, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 100},
		{100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100},
	}
}

func (m *Maze) setPixel(po Position, pi Pixel) {
	m[po.y][po.x] = int(pi)
}
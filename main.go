package main

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 320
	screenHeight = 240
)

type World struct {
	Width  int
	Height int
	Grid   []bool
}

func NewWorld(width, height, initCells int) *World {
	world := &World{
		Width: width, 
		Height: height, 
		Grid: make([]bool, width*height), 
	}
	world.Randomize(initCells)
	return world
}

func (w *World) Randomize(initCells int) {
	for i := 0; i < initCells; i++ {
		x := rand.Intn(w.Width)
		y := rand.Intn(w.Height) 
		w.Grid[y * w.Width + x] = true 
	}
}

// Any live cell with fewer than two live neighbors dies, as if by underpopulation.
// Any live cell with two or three live neighbors lives on to the next generation.
// Any live cell with more than three live neighbors dies, as if by overpopulation.
// Any dead cell with exactly three live neighbors becomes a live cell, as if by reproduction.
func (w *World) Next() {
	next := make([]bool, w.Width * w.Height) 
	for x := 0; x < w.Width; x++ { // rows
		for y := 0; y < w.Height; y++ { // cols
			n := AliveNeighbours(w.Grid, x, y, w.Width, w.Height)
			switch {
			case n < 2:
				next[y * w.Width + x] = false
			case (n == 2 || n == 3) && w.Grid[y * w.Width + x]:
				next[y * w.Width + x] = true 
			case n > 3:
				next[y * w.Width + x] = false
			case n == 3:
				next[y * w.Width + x] = true
			}
		}
	}
	w.Grid = next 
}

func AliveNeighbours(grid []bool, x, y, width, height int) int {
	count := 0 
	for i := -1; i <= 1; i++ { // rows
		for j := -1; j <=1; j++ { // cols
			if i == 0 && j == 0 {
				continue // Centre Square
			}

			x2 := x - i 
			y2 := y - j 

			if x2 < 0 || y2 < 0 || x2 >= width || y2 >= height {
				continue // off the edge
			}

			if grid[y2 * width + x2] {
				count++
			}
		}
	}
	return count
}

func (w *World) RenderInto(pixels []byte) {
	for i, v := range w.Grid {
		if v {
			pixels[4*i] = 0xff
			pixels[4*i+1] = 0xff
			pixels[4*i+2] = 0xff
			pixels[4*i+3] = 0xff
		} else {
			pixels[4*i] = 0
			pixels[4*i+1] = 0
			pixels[4*i+2] = 0
			pixels[4*i+3] = 0
		}
	}
}

type Game struct {
	World  *World
	Pixels []byte
}

func (g *Game) Layout(innerWidth, innerHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) Update() error {
	g.World.Next()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.World.RenderInto(g.Pixels)
	screen.WritePixels(g.Pixels)
}

func main() {
	initCells := (screenWidth * screenHeight) / 10 
	game := &Game{
		World: NewWorld(screenWidth, screenHeight, initCells),
		Pixels: make([]byte, 4*screenWidth*screenHeight),
	}

	ebiten.SetWindowSize(2*screenWidth, 2*screenHeight)
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}

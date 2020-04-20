package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten"
)

type Game struct {
	blocks    []*Block
	direction Direction
}

type Block struct {
	img *ebiten.Image
	x   float64
	y   float64
}

type Direction int

const (
	Up Direction = iota
	Right
	Down
	Left
)

func newBlock(x, y float64, d Direction) (*Block, error) {
	var img *ebiten.Image
	var err error

	if d == Right || d == Left {
		img, err = ebiten.NewImage(1, 10, ebiten.FilterDefault)
	} else if d == Up || d == Down {
		img, err = ebiten.NewImage(10, 1, ebiten.FilterDefault)
	}

	if err != nil {
		return nil, err
	}
	img.Fill(color.White)
	return &Block{
		img: img,
		x:   x,
		y:   y,
	}, nil
}

func initBlocks() ([]*Block, error) {
	var blocks []*Block
	for i := 0; i < 100; i++ {
		block, err := newBlock(float64(i), 0, Right)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, block)
	}
	return blocks, nil
}

func (g *Game) Init() error {
	blocks, err := initBlocks()
	if err != nil {
		return err
	}
	g.blocks = blocks
	g.direction = Right
	return nil
}

func (g *Game) Update(screen *ebiten.Image) error {
	switch {
	case ebiten.IsKeyPressed(ebiten.KeyDown) && g.direction != Up && g.direction != Down:
		fb := g.blocks[len(g.blocks)-10]
		var x, y float64
		switch g.direction {
		case Left:
			x, y = fb.x-9, fb.y
		case Right:
			x, y = fb.x, fb.y
		}
		for i := len(g.blocks) - 10; i < len(g.blocks); i++ {
			block, err := newBlock(x, y, Down)
			if err != nil {
				return err
			}
			y++
			g.blocks[i] = block
		}
		g.direction = Down
	case ebiten.IsKeyPressed(ebiten.KeyUp) && g.direction != Up && g.direction != Down:
		fb := g.blocks[len(g.blocks)-10]
		var x, y float64
		switch g.direction {
		case Left:
			x, y = fb.x-9, fb.y+9
		case Right:
			x, y = fb.x, fb.y+9
		}
		for i := len(g.blocks) - 10; i < len(g.blocks); i++ {
			block, err := newBlock(x, y, Up)
			if err != nil {
				return err
			}
			y--
			g.blocks[i] = block
		}
		g.direction = Up
	case ebiten.IsKeyPressed(ebiten.KeyRight) && g.direction != Left && g.direction != Right:
		fb := g.blocks[len(g.blocks)-10]
		var x, y float64
		switch g.direction {
		case Up:
			x, y = fb.x, fb.y-9
		case Down:
			x, y = fb.x, fb.y
		}
		for i := len(g.blocks) - 10; i < len(g.blocks); i++ {
			block, err := newBlock(x, y, Right)
			if err != nil {
				return err
			}
			x++
			g.blocks[i] = block
		}
		g.direction = Right
	case ebiten.IsKeyPressed(ebiten.KeyLeft) && g.direction != Left && g.direction != Right:
		fb := g.blocks[len(g.blocks)-10]
		var x, y float64
		switch g.direction {
		case Up:
			x, y = fb.x+9, fb.y-9
		case Down:
			x, y = fb.x+9, fb.y
		}
		for i := len(g.blocks) - 10; i < len(g.blocks); i++ {
			block, err := newBlock(x, y, Left)
			if err != nil {
				return err
			}
			x--
			g.blocks[i] = block
		}
		g.direction = Left
	}

	switch g.direction {
	case Right:
		head := g.blocks[len(g.blocks)-1]
		if head.x < 319 {
			block, err := newBlock(head.x+1, head.y, Right)
			if err != nil {
				return err
			}
			g.blocks = append(g.blocks[1:], block)
		}
	case Down:
		head := g.blocks[len(g.blocks)-1]
		if head.y < 239 {
			block, err := newBlock(head.x, head.y+1, Down)
			if err != nil {
				return err
			}
			g.blocks = append(g.blocks[1:], block)
		}
	case Left:
		head := g.blocks[len(g.blocks)-1]
		if head.x > 0 {
			block, err := newBlock(head.x-1, head.y, Left)
			if err != nil {
				return err
			}
			g.blocks = append(g.blocks[1:], block)
		}
	case Up:
		head := g.blocks[len(g.blocks)-1]
		if head.y > 0 {

			block, err := newBlock(head.x, head.y-1, Down)
			if err != nil {
				return err
			}
			g.blocks = append(g.blocks[1:], block)
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, block := range g.blocks {
		geom := ebiten.GeoM{}
		geom.Translate(block.x, block.y)
		screen.DrawImage(block.img, &ebiten.DrawImageOptions{
			GeoM: geom,
		})
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(320, 240)
	ebiten.SetWindowTitle("Snake")
	ebiten.SetMaxTPS(120)
	game := &Game{}
	if err := game.Init(); err != nil {
		log.Fatal(err)
	}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

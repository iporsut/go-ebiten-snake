package main

import (
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten"
)

var rainbowPal = []color.RGBA{
	{0xff, 0x00, 0x00, 0xff},
	{0xff, 0x7f, 0x00, 0xff},
	{0xff, 0xff, 0x00, 0xff},
	{0x00, 0xff, 0x00, 0xff},
	{0x00, 0x00, 0xff, 0xff},
	{0x4b, 0x00, 0x82, 0xff},
	{0x8f, 0x00, 0xff, 0xff},
}

var sameColorCounter = 0
var rainbowColorIndex = 0

func getRainbowColor() color.Color {
	if sameColorCounter == 0 {
		rainbowColorIndex = rand.Intn(len(rainbowPal))
	}
	sameColorCounter = (sameColorCounter + 1) % 10
	return rainbowPal[rainbowColorIndex]
	// return color.White
}

type Game struct {
	blocks    []*Block
	direction Direction
	dot       *Block
}

type Block struct {
	img    *ebiten.Image
	x      float64
	y      float64
	width  int
	height int
}

type Direction int

const (
	Up Direction = iota
	Right
	Down
	Left
)

func newBlock(x, y float64, width, height int, color color.Color) (*Block, error) {
	img, err := ebiten.NewImage(width, height, ebiten.FilterDefault)
	if err != nil {
		return nil, err
	}
	img.Fill(color)
	return &Block{
		img:    img,
		x:      x,
		y:      y,
		width:  width,
		height: height,
	}, nil
}

func (g *Game) newRandomDot() error {
	var x1, y1, x2, y2 int
	conflict := true
RandomXY:
	for conflict {
		x1 = rand.Intn(311)
		y1 = rand.Intn(231)
		x2 = x1 + 10
		y2 = y1 + 10

		for _, b := range g.blocks {
			if x1 > (int(b.x)+b.width) || int(b.x) > x2 {
				continue
			}

			if y1 > (int(b.y)+b.height) || int(b.y) > y2 {
				continue
			}

			continue RandomXY
		}
		conflict = false
	}

	block, err := newBlock(float64(x1), float64(y1), 10, 10, getRainbowColor())
	if err != nil {
		return err
	}
	g.dot = block
	return nil
}

func (g *Game) hitDot() bool {
	head := g.blocks[len(g.blocks)-1]
	hx1 := int(head.x)
	hy1 := int(head.y)
	hx2 := int(head.x) + head.width
	hy2 := int(head.y) + head.height

	dx1 := int(g.dot.x)
	dy1 := int(g.dot.y)
	dx2 := int(g.dot.x) + g.dot.width
	dy2 := int(g.dot.y) + g.dot.height

	if hx1 > dx2 || dx1 > hx2 {
		return false
	}

	if hy1 > dy2 || dy1 > hy2 {
		return false
	}

	return true
}

func initBlocks() ([]*Block, error) {
	var blocks []*Block
	for i := 0; i < 200; i++ {
		block, err := newBlock(float64(i), 0, 1, 10, getRainbowColor())
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
	err = g.newRandomDot()
	if err != nil {
		return err
	}
	return nil
}

func (g *Game) updateDirection() error {
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
			block, err := newBlock(x, y, 10, 1, getRainbowColor())
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
			block, err := newBlock(x, y, 10, 1, getRainbowColor())
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
			block, err := newBlock(x, y, 1, 10, getRainbowColor())
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
			block, err := newBlock(x, y, 1, 10, getRainbowColor())
			if err != nil {
				return err
			}
			x--
			g.blocks[i] = block
		}
		g.direction = Left
	}
	return nil
}

func (g *Game) appendForwardBlock(x, y float64, width, height int) error {
	block, err := newBlock(x, y, width, height, getRainbowColor())
	if err != nil {
		return err
	}

	if g.hitDot() {
		if err := g.newRandomDot(); err != nil {
			return err
		}
		g.blocks = append(g.blocks, block)
	} else {
		g.blocks = append(g.blocks[1:], block)
	}

	return nil
}

func (g *Game) move() error {
	switch g.direction {
	case Right:
		head := g.blocks[len(g.blocks)-1]
		if head.x < 319 {
			return g.appendForwardBlock(head.x+1, head.y, 1, 10)
		}
	case Down:
		head := g.blocks[len(g.blocks)-1]
		if head.y < 239 {
			return g.appendForwardBlock(head.x, head.y+1, 10, 1)
		}
	case Left:
		head := g.blocks[len(g.blocks)-1]
		if head.x > 0 {
			return g.appendForwardBlock(head.x-1, head.y, 1, 10)
		}
	case Up:
		head := g.blocks[len(g.blocks)-1]
		if head.y > 0 {
			return g.appendForwardBlock(head.x, head.y-1, 10, 1)
		}
	}
	return nil
}

func (g *Game) Update(screen *ebiten.Image) error {
	if err := g.updateDirection(); err != nil {
		return err
	}
	return g.move()
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw snake
	for _, block := range g.blocks {
		geom := ebiten.GeoM{}
		geom.Translate(block.x, block.y)
		screen.DrawImage(block.img, &ebiten.DrawImageOptions{
			GeoM: geom,
		})
	}

	// Draw dot
	geom := ebiten.GeoM{}
	geom.Translate(g.dot.x, g.dot.y)
	screen.DrawImage(g.dot.img, &ebiten.DrawImageOptions{
		GeoM: geom,
	})
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	rand.Seed(time.Now().UnixNano())

	ebiten.SetWindowSize(640, 480)
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

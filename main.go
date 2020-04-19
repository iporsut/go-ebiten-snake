package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten"
)

var game *Game

type Block struct {
	img *ebiten.Image
	x   int
	y   int
	d   Direction
}

func initBlock() []*Block {
	var blocks []*Block
	for i := 0; i < 20; i++ {
		blocks = append(blocks, &Block{
			img: newBlockImage(Right),
			x:   i,
			y:   0,
			d:   Right,
		})
	}
	return blocks
}

func newBlockImage(d Direction) *ebiten.Image {
	var w, h int
	switch d {
	case Up, Down:
		w, h = 10, 1
	case Left, Right:
		w, h = 1, 10
	}
	blockImg, err := ebiten.NewImage(w, h, ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
	blockImg.Fill(color.White)
	return blockImg
}

func init() {
	blocks := initBlock()

	game = &Game{
		blocks: blocks,
		head:   blocks[len(blocks)-1],
		tail:   blocks[0],
	}
}

type Direction int

const (
	Right Direction = iota
	Down
	Left
	Up
)

type Game struct {
	blocks []*Block
	head   *Block
	tail   *Block
}

func (g *Game) Update(screen *ebiten.Image) error {
	g.blocks = append(g.blocks[1:], &Block{
		img: newBlockImage(Right),
		x:   g.head.x + 1,
		y:   g.head.y,
		d:   g.head.d,
	})
	g.head = g.blocks[len(g.blocks)-1]
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, block := range g.blocks {
		geom := ebiten.GeoM{}
		geom.Translate(float64(block.x), float64(block.y))
		screen.DrawImage(block.img, &ebiten.DrawImageOptions{
			GeoM: geom,
		})
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetMaxTPS(60)
	ebiten.SetWindowSize(320, 240)
	ebiten.SetWindowTitle("สวัสดี Ebiten")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

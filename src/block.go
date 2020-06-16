package main

import (
	"github.com/jani-nykanen/ultimate-puzzle/src/core"
)

type block struct {
	pos       core.Point
	target    core.Point
	renderPos core.Point
	id        int32
	active    bool
	spr       *core.Sprite
	moving    bool
	moveTimer int32
}

func (b *block) update(s *stage, ev *core.Event) {

	if !b.active {
		return
	}
}

func (b *block) drawOutlines(c *core.Canvas, ap *core.AssetPack) {

	if !b.active {
		return
	}

	c.FillRect(b.renderPos.X-1, b.renderPos.Y-1, 18, 18,
		core.NewRGB(0, 0, 0))
}

func (b *block) draw(c *core.Canvas, ap *core.AssetPack) {

	if !b.active {
		return
	}

	bmp := ap.GetAsset("blocks").(*core.Bitmap)

	c.DrawSprite(b.spr, bmp,
		b.renderPos.X, b.renderPos.Y, core.FlipNone)
}

func newBlock(x, y, id int32) *block {

	b := new(block)

	b.pos = core.NewPoint(x, y)
	b.target = b.pos

	b.renderPos.X = x * 16
	b.renderPos.Y = y * 16

	b.id = id
	b.active = true

	b.spr = core.NewSprite(16, 16)
	b.spr.SetFrame(id, 0)

	b.moveTimer = 0
	b.moving = false

	return b
}

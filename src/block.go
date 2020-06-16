package main

import (
	"github.com/jani-nykanen/ultimate-puzzle/src/core"
)

const (
	blockMoveTime = 8
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

func (b *block) handleControls(s *stage, ev *core.Event) bool {

	if b.moving {

		return false
	}

	dx := int32(0)
	dy := int32(0)
	if ev.Input.GetActionState("left")&core.StateDownOrPressed == 1 {

		dx = -1

	} else if ev.Input.GetActionState("right")&core.StateDownOrPressed == 1 {

		dx = 1

	} else if ev.Input.GetActionState("up")&core.StateDownOrPressed == 1 {

		dy = -1

	} else if ev.Input.GetActionState("down")&core.StateDownOrPressed == 1 {

		dy = 1
	}

	if dx != 0 || dy != 0 {

		b.moveTo(dx, dy, s)
		return b.moving
	}

	return false
}

func (b *block) moveTo(dx, dy int32, s *stage) {

	if b.moving || s.getSolid(b.pos.X+dx, b.pos.Y+dy) != 0 {
		return
	}

	b.moveTimer += blockMoveTime
	b.moving = true
	b.target.X = core.NegMod(b.pos.X+dx, s.width)
	b.target.Y = core.NegMod(b.pos.Y+dy, s.height)

	s.updateSolidTile(b.pos.X, b.pos.Y, 0)
}

func (b *block) handleMovement(s *stage, ev *core.Event) {

	if !b.moving {
		return
	}

	var dirx, diry int32

	b.moveTimer -= ev.Step()
	if b.moveTimer <= 0 {

		dirx = b.target.X - b.pos.X
		diry = b.target.Y - b.pos.Y

		b.pos = b.target

		// Keep moving to the same direction, if possible
		b.moving = false
		b.moveTo(dirx, diry, s)

		if !b.moving {

			b.moveTimer = 0
			s.updateSolidTile(b.pos.X, b.pos.Y, 2)
		}
	}
}

func (b *block) safeCheck(s *stage) {

	// Sometimes this ugly thing happen
	if s.getSolid(b.target.X, b.target.Y) != 0 {

		b.moveTimer = 0
		b.moving = false
		b.target = b.pos

		s.updateSolidTile(b.pos.X, b.pos.Y, 2)
		b.computeRenderingPosition()
	}
}

func (b *block) computeRenderingPosition() {

	var t float32
	var x, y float32

	if !b.moving {

		b.renderPos.X = b.pos.X * 16
		b.renderPos.Y = b.pos.Y * 16
	} else {

		t = float32(b.moveTimer) / float32(blockMoveTime)

		x = float32(b.pos.X*16)*t + (1.0-t)*float32(b.target.X*16)
		y = float32(b.pos.Y*16)*t + (1.0-t)*float32(b.target.Y*16)

		b.renderPos.X = core.RoundFloat32(x)
		b.renderPos.Y = core.RoundFloat32(y)
	}
}

func (b *block) update(s *stage, ev *core.Event) {

	if !b.active {
		return
	}

	b.handleMovement(s, ev)
	b.computeRenderingPosition()
}

func (b *block) drawOutlines(c *core.Canvas, ap *core.AssetPack) {

	if !b.active {
		return
	}

	c.FillRect(b.renderPos.X-1, b.renderPos.Y-1, 18, 18,
		core.NewRGB(0, 0, 0))
}

func (b *block) drawShadow(c *core.Canvas, ap *core.AssetPack) {

	c.DrawBitmap(ap.GetAsset("shadow").(*core.Bitmap),
		b.renderPos.X-1, b.renderPos.Y-1, core.FlipNone)
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

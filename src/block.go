package main

import (
	"github.com/jani-nykanen/ultimate-puzzle/src/core"
)

const (
	blockMoveTime = 8

	blockNoHole    = 0
	blockRightHole = 1
	blockWrongHole = 2 // ehehehhehehe
)

type block struct {
	pos         core.Point
	target      core.Point
	dir         core.Point // Needed for "offscreen transition"
	renderPos   core.Point
	id          int32
	exist       bool
	spr         *core.Sprite
	moving      bool
	moveTimer   int32
	jumping     bool
	deactivated bool
	playDestroy bool
	playHit     bool
}

func (b *block) handleControls(s *stage, ev *core.Event) bool {

	if !b.exist || b.moving {

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

		b.moveTo(dx, dy, s, false)

		return b.moving
	}

	return false
}

func (b *block) moveTo(dx, dy int32, s *stage, wasMoving bool) {

	if b.moving || s.getSolid(b.pos.X+dx, b.pos.Y+dy) != 0 {

		if wasMoving {

			b.playHit = true
		}

		return
	}

	b.dir.X = dx
	b.dir.Y = dy

	b.jumping = b.pos.X+dx < 0 || b.pos.X+dx >= s.width ||
		b.pos.Y+dy < 0 || b.pos.Y+dy >= s.height

	b.moveTimer += blockMoveTime
	b.moving = true

	b.target.X = core.NegMod(b.pos.X+dx, s.width)
	b.target.Y = core.NegMod(b.pos.Y+dy, s.height)

	s.updateSolidTile(b.pos.X, b.pos.Y, 0)
}

func (b *block) handleMovement(s *stage, ev *core.Event) int32 {

	if !b.moving {
		return blockNoHole
	}

	var hitHole, correctHole bool

	b.moveTimer -= ev.Step()
	if b.moveTimer <= 0 {

		b.pos = b.target

		// Check if hits a hole
		if b.id != 0 {

			hitHole, correctHole = s.checkHoleTile(b.pos.X, b.pos.Y, b.id-1)
			if hitHole {

				if correctHole {

					b.exist = false
					b.moving = false
					b.moveTimer = 0

					b.deactivated = false

					s.updateSolidTile(b.pos.X, b.pos.Y, 2)

					b.playDestroy = true

					return blockRightHole
				}

				ev.Audio.PlaySample(ev.Assets.GetAsset("failure").(*core.Sample),
					60)

				return blockWrongHole
			}
		}

		// Keep moving to the same direction, if possible
		b.moving = false
		b.moveTo(b.dir.X, b.dir.Y, s, true)

		if !b.moving {

			b.moveTimer = 0
			s.updateSolidTile(b.pos.X, b.pos.Y, 2)
		}
	}

	return blockNoHole
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

	target := b.target

	if !b.moving {

		b.renderPos.X = b.pos.X * 16
		b.renderPos.Y = b.pos.Y * 16

	} else {

		t = float32(b.moveTimer) / float32(blockMoveTime)

		if b.jumping {

			target.X = b.pos.X + b.dir.X
			target.Y = b.pos.Y + b.dir.Y
		}

		x = float32(b.pos.X*16)*t + (1.0-t)*float32(target.X*16)
		y = float32(b.pos.Y*16)*t + (1.0-t)*float32(target.Y*16)

		b.renderPos.X = core.RoundFloat32(x)
		b.renderPos.Y = core.RoundFloat32(y)
	}
}

func (b *block) update(anyMoving bool, s *stage, ev *core.Event) int32 {

	b.playDestroy = false
	b.playHit = false

	if !b.exist {

		if !b.deactivated && !anyMoving {

			s.updateSolidTile(b.pos.X, b.pos.Y, 0)
			b.deactivated = true
		}

		return blockNoHole
	}

	ret := b.handleMovement(s, ev)
	b.computeRenderingPosition()

	return ret
}

func (b *block) drawOutlines(c *core.Canvas, ap *core.AssetPack) {

	if !b.exist {
		return
	}

	c.FillRect(b.renderPos.X-1, b.renderPos.Y-1, 18, 18,
		core.NewRGB(0, 0, 0))

	if b.jumping {

		if b.dir.X != 0 {

			c.FillRect(b.renderPos.X-1-b.dir.X*c.Viewport().W,
				b.renderPos.Y-1, 18, 18,
				core.NewRGB(0, 0, 0))

		} else if b.dir.Y != 0 {

			c.FillRect(b.renderPos.X-1,
				b.renderPos.Y-1-b.dir.Y*c.Viewport().H,
				18, 18, core.NewRGB(0, 0, 0))
		}
	}
}

func (b *block) drawShadow(c *core.Canvas, ap *core.AssetPack) {

	if !b.exist {
		return
	}

	bmp := ap.GetAsset("shadow").(*core.Bitmap)

	c.DrawBitmap(bmp,
		b.renderPos.X-1, b.renderPos.Y-1, core.FlipNone)

	if b.jumping {

		if b.dir.X != 0 {

			c.DrawSprite(b.spr, bmp,
				b.renderPos.X-1-b.dir.X*c.Viewport().W,
				b.renderPos.Y-1, core.FlipNone)

		} else if b.dir.Y != 0 {

			c.DrawSprite(b.spr, bmp,
				b.renderPos.X-1,
				b.renderPos.Y-1-b.dir.Y*c.Viewport().H,
				core.FlipNone)
		}
	}
}

func (b *block) draw(c *core.Canvas, ap *core.AssetPack) {

	if !b.exist {
		return
	}

	bmp := ap.GetAsset("blocks").(*core.Bitmap)

	c.DrawSprite(b.spr, bmp,
		b.renderPos.X, b.renderPos.Y, core.FlipNone)

	if b.jumping {

		if b.dir.X != 0 {

			c.DrawSprite(b.spr, bmp,
				b.renderPos.X-b.dir.X*c.Viewport().W,
				b.renderPos.Y, core.FlipNone)

		} else if b.dir.Y != 0 {

			c.DrawSprite(b.spr, bmp,
				b.renderPos.X,
				b.renderPos.Y-b.dir.Y*c.Viewport().H,
				core.FlipNone)
		}
	}
}

func newBlock(x, y, id int32) *block {

	b := new(block)

	b.pos = core.NewPoint(x, y)
	b.target = b.pos

	b.renderPos.X = x * 16
	b.renderPos.Y = y * 16

	b.id = id
	b.exist = true

	b.spr = core.NewSprite(16, 16)
	b.spr.SetFrame(id, 0)

	b.moveTimer = 0
	b.moving = false
	b.deactivated = true

	b.playDestroy = false
	b.playHit = false

	return b
}

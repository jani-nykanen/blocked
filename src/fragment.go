package main

import (
	"github.com/jani-nykanen/blocked/src/core"
)

type fragment struct {
	pos     core.Vector2
	speed   core.Vector2
	timer   int32
	maxTime int32
	sx      int32
	sy      int32
	sw      int32
	sh      int32
	exist   bool
}

func (frag *fragment) spawn(x, y, sx, sy, sw, sh int32,
	speedx, speedy float32, time int32) {

	frag.pos.X = float32(x)
	frag.pos.Y = float32(y)
	frag.speed.X = speedx
	frag.speed.Y = speedy

	frag.sx = sx
	frag.sy = sy
	frag.sw = sw
	frag.sh = sh

	frag.timer = time
	frag.maxTime = time

	frag.exist = true
}

func (frag *fragment) update(ev *core.Event) {

	if !frag.exist {
		return
	}

	frag.timer -= ev.Step()
	if frag.timer <= 0 {

		frag.exist = false
	}

	frag.pos.X += frag.speed.X * float32(ev.Step())
	frag.pos.Y += frag.speed.Y * float32(ev.Step())
}

func (frag *fragment) draw(c *core.Canvas, bmp *core.Bitmap) {

	if !frag.exist {
		return
	}

	dx := core.RoundFloat32(frag.pos.X) - frag.sw/2
	dy := core.RoundFloat32(frag.pos.Y) - frag.sh/2

	c.DrawBitmapRegion(bmp, frag.sx, frag.sy, frag.sw, frag.sh,
		dx, dy, core.FlipNone)
}

func newFragment() *fragment {

	frag := new(fragment)

	frag.exist = false

	return frag
}

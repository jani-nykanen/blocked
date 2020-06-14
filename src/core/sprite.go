package core

// Sprite : An animated sprite
type Sprite struct {
	frame int32
	row   int32
	count int32

	width  int32
	height int32
}

// Animate : Animate the sprite
func (spr *Sprite) Animate(row, start, end, speed, step int32) {

	// Nothing to animate, just set the frame
	if start == end {

		spr.SetFrame(start, row)
		return
	}

	// Swap row
	if spr.row != row {

		spr.count = 0
		spr.frame = start
		if end < start {

			spr.frame = end
		}
		spr.row = row
	}

	// If outside the animation interval
	if start < end &&
		(spr.frame < start || spr.frame > end) {

		spr.frame = start
		spr.count = 0

	} else if end < start &&
		(spr.frame < end || spr.frame > start) {

		spr.frame = end
		spr.count = 0
	}

	spr.count += 1.0 * step
	if spr.count > speed {

		if start < end {

			spr.frame++
			if spr.frame > end {

				spr.frame = start
			}
			if speed < 0 {

				spr.frame = MinInt32(spr.frame-speed, end)
			}

		} else {

			spr.frame--
			if spr.frame < end {

				spr.frame = start
			}
			if speed < 0 {

				spr.frame = MaxInt32(spr.frame+speed, start)
			}
		}

		spr.count -= speed
	}
}

// DrawFrame : Draw the given frame of the sprite
func (spr *Sprite) DrawFrame(c *Canvas, bmp *Bitmap, x, y, frame, row int32, flip Flip) {

	c.DrawBitmapRegion(bmp, frame*spr.width, row*spr.height,
		spr.width, spr.height,
		x, y, flip)
}

// Draw : Draw the sprite
func (spr *Sprite) Draw(c *Canvas, bmp *Bitmap, x, y int32, flip Flip) {

	spr.DrawFrame(c, bmp, x, y, spr.frame, spr.row, flip)
}

// SetFrame : Set the current frame
func (spr *Sprite) SetFrame(frame, row int32) {

	spr.frame = frame
	spr.row = row
	spr.count = 0
}

// Frame : Getter for the current frame
func (spr *Sprite) Frame() int32 {

	return spr.frame
}

// Row : Getter for the current row
func (spr *Sprite) Row() int32 {

	return spr.row
}

// Width : Getter for width
func (spr *Sprite) Width() int32 {

	return spr.width
}

// Height : Getter for width
func (spr *Sprite) Height() int32 {

	return spr.height
}

// NewSprite : Constructs a new sprite
func NewSprite(width, height int32) *Sprite {

	spr := new(Sprite)

	spr.frame = 0
	spr.row = 0
	spr.count = 0

	spr.width = width
	spr.height = height

	return spr
}

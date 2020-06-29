package core

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

// RenderCallback : Used when rendering to user-created
// bitmaps
type RenderCallback func(c *Canvas, ap *AssetPack)

// Flip : A flag for flipping a bitmap
type Flip int32

// Flipping flags
const (
	FlipNone       = 0
	FlipHorizontal = 1
	FlipVertical   = 2
	FlipBoth       = 3
)

// Blend modes
const (
	BlendDefault = sdl.BLENDMODE_BLEND
	BlendMod     = sdl.BLENDMODE_MOD
	BlendNone    = sdl.BLENDMODE_NONE
)

// Canvas : A "buffer" where the drawn content
// is stored. In this case, a texture
type Canvas struct {
	width       uint32
	height      uint32
	translation Point
	renderer    *sdl.Renderer
	frame       *Bitmap
	frameCopy   *Bitmap
	frameTarget sdl.Rect
	srcRect     sdl.Rect
	destRect    sdl.Rect
	viewport    Rectangle
}

func (c *Canvas) initialize(renderer *sdl.Renderer) error {

	var err error

	c.renderer = renderer

	c.frame, err = newBitmap(c.width, c.height, true, renderer)
	if err != nil {

		return err
	}

	c.frameCopy, err = newBitmap(c.width, c.height, true, renderer)
	if err != nil {

		c.frame.Dispose()
		return err
	}

	c.viewport = NewRect(0, 0, int32(c.width), int32(c.height))

	return err
}

func (c *Canvas) begin() {

	_ = c.renderer.SetRenderTarget(c.frame.texture)
}

func (c *Canvas) end() {

	_ = c.renderer.SetRenderTarget(nil)
}

func (c *Canvas) redrawFrame() {

	c.Clear(0, 0, 0)
	c.renderer.CopyEx(c.frame.texture,
		nil, &c.frameTarget, 0.0, nil, sdl.FLIP_NONE)
}

func (c *Canvas) resize(w, h int32) {

	// Find the best multiplier for
	// square pixels (that is, each pixel is square
	// with integer dimensions)
	mul := MinInt32(
		w/int32(c.width),
		h/int32(c.height))

	c.frameTarget.W = int32(c.width) * mul
	c.frameTarget.H = int32(c.height) * mul
	c.frameTarget.X = w/2 - c.frameTarget.W/2
	c.frameTarget.Y = h/2 - c.frameTarget.H/2
}

func (c *Canvas) dispose() {

	c.frame.Dispose()
	c.frameCopy.Dispose()
}

// Clear : Clear the screen with a color
func (c *Canvas) Clear(r, g, b uint8) {

	c.renderer.SetDrawColor(r, g, b, 255)
	c.renderer.Clear()
}

// ClearToAlpha : Make the current render target
// transparent
func (c *Canvas) ClearToAlpha() {

	c.renderer.SetDrawBlendMode(sdl.BLENDMODE_NONE)

	c.renderer.SetDrawColor(0, 0, 0, 0)
	c.renderer.Clear()

	c.renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
}

// DrawBitmap : Draw a full bitmap
func (c *Canvas) DrawBitmap(bmp *Bitmap, dx, dy int32, flip Flip) {

	c.DrawBitmapRegion(bmp, 0, 0,
		int32(bmp.width), int32(bmp.height), dx, dy, flip)
}

// DrawBitmapRegion : Draw a region of the bitmap
func (c *Canvas) DrawBitmapRegion(bmp *Bitmap, sx, sy, sw, sh, dx, dy int32, flip Flip) {

	dx += c.translation.X
	dy += c.translation.Y

	c.srcRect = sdl.Rect{X: sx, Y: sy, W: sw, H: sh}
	c.destRect = sdl.Rect{X: dx, Y: dy, W: sw, H: sh}

	c.renderer.CopyEx(bmp.texture, &c.srcRect, &c.destRect,
		0.0, nil, sdl.RendererFlip(flip))
}

// DrawText : Draw text with a bitmap font. Note that
// it is assumed that the font contains 16 characters
// per fow
func (c *Canvas) DrawText(font *Bitmap, text string,
	dx, dy, xoff, yoff int32, center bool) {

	var chr uint8

	cw := int32(font.width) / 16
	ch := cw
	length := len(text)

	x := dx
	y := dy

	if center {

		dx -= (int32(length) * (cw + xoff)) / 2
		x = dx
	}

	// Draw every character
	var sx, sy int32
	for i := 0; i < length; i++ {

		chr = text[i]

		// Line swap
		if chr == '\n' {

			x = dx
			y += (yoff + ch)
			continue
		}

		sx = int32(chr) % 16
		sy = int32(chr) / 16

		c.DrawBitmapRegion(font,
			sx*cw, sy*ch, cw, ch,
			x, y, FlipNone)

		x += (cw + xoff)
	}
}

// FillRect : Fills an rectangle
func (c *Canvas) FillRect(x, y, w, h int32, color Color) {

	c.renderer.SetDrawColor(color.R, color.G, color.B, color.A)

	x += c.translation.X
	y += c.translation.Y
	c.destRect = sdl.Rect{X: x, Y: y, W: w, H: h}

	c.renderer.FillRect(&c.destRect)

}

// FillCircleOutside : Fill area outside the circle
func (c *Canvas) FillCircleOutside(cx, cy, radius int32, color Color) {

	w := c.viewport.W
	h := c.viewport.H

	if radius <= 0 {

		c.FillRect(0, 0, w, h, color)
		return

	} else if radius*radius >= w*w+h*h {

		return
	}

	start := MaxInt32(0, cy-radius)
	end := MinInt32(h, cy+radius)

	if start > 0 {
		c.FillRect(0, 0, w, start, color)
	}

	if end < h {
		c.FillRect(0, end, w, h-end, color)
	}

	var dy int32
	var px1, px2 int32
	for y := int32(start); y < end; y++ {

		dy = y - cy

		// A full line
		if int32(math.Abs(float64(dy))) >= radius {

			c.FillRect(0, y, w, 1, color)
			continue
		}

		px1 = cx - int32(math.Sqrt(float64(radius*radius-dy*dy)))
		px2 = cx + int32(math.Sqrt(float64(radius*radius-dy*dy)))

		// Fill left
		if px1 > 0 {
			c.FillRect(0, y, px1, 1, color)
		}
		// Fill right
		if px2 < w {

			c.FillRect(px2, y, w-px1, 1, color)
		}
	}
}

// DrawSpriteFrame : Draw an animated sprite frame
func (c *Canvas) DrawSpriteFrame(spr *Sprite,
	bmp *Bitmap, x, y, frame, row int32,
	flip Flip) {

	spr.DrawFrame(c, bmp, x, y, frame, row, flip)
}

// DrawSprite : Draw an animated sprite
func (c *Canvas) DrawSprite(spr *Sprite, bmp *Bitmap, x, y int32, flip Flip) {

	spr.Draw(c, bmp, x, y, flip)
}

// MoveTo : Move the top-left corner of rendering
// to the given point
func (c *Canvas) MoveTo(x, y int32) {

	c.translation.X = x
	c.translation.Y = y
}

// Move : Move the top-left corner by the given value
func (c *Canvas) Move(dx, dy int32) {

	c.translation.X += dx
	c.translation.Y += dy
}

// DrawToBitmap : Use a bitmap as a render target
func (c *Canvas) DrawToBitmap(bmp *Bitmap, ap *AssetPack, cb RenderCallback) {

	oldTarget := c.renderer.GetRenderTarget()

	c.renderer.SetRenderTarget(bmp.texture)
	cb(c, ap)
	c.renderer.SetRenderTarget(oldTarget)
}

// CopyCurrentFrame : Copy current frame to the buffer,
// so it can be drawn as a bitmap
func (c *Canvas) CopyCurrentFrame() {

	c.DrawToBitmap(c.frameCopy, nil, func(c *Canvas, ap *AssetPack) {

		c.DrawBitmap(c.frame, 0, 0, FlipNone)
	})
}

// DrawCopiedFrame : Draw the copied frame as a bitmap
func (c *Canvas) DrawCopiedFrame(x, y int32, flip Flip) {

	c.DrawBitmap(c.frameCopy, x, y, flip)
}

// DrawCopiedFrameRegion : Draw a region of the copied frame
func (c *Canvas) DrawCopiedFrameRegion(sx, sy, sw, sh,
	dx, dy int32, flip Flip) {

	c.DrawBitmapRegion(c.frameCopy, sx, sy, sw, sh,
		dx, dy, flip)
}

// SetBitmapColor : Set color to be used when drawing
// a bitmap
func (c *Canvas) SetBitmapColor(bmp *Bitmap, r, g, b uint8) {

	bmp.texture.SetColorMod(r, g, b)
}

// SetBitmapAlpha : Set alpha value to be used when drawing
// a bitmap
func (c *Canvas) SetBitmapAlpha(bmp *Bitmap, a uint8) {

	bmp.texture.SetAlphaMod(a)
}

// SetViewport : Set the current view area
func (c *Canvas) SetViewport(x, y, w, h int32) {

	c.destRect = sdl.Rect{X: x, Y: y, W: w, H: h}
	c.viewport = NewRect(x, y, w, h)

	c.renderer.SetViewport(&c.destRect)
}

// ResetViewport : Reset the viewport to the whole
// target area
func (c *Canvas) ResetViewport() {

	c.viewport = NewRect(0, 0, int32(c.width), int32(c.height))
	c.renderer.SetViewport(nil)
}

// Viewport : Getter for viewport
func (c *Canvas) Viewport() Rectangle {

	return c.viewport
}

// Width : A getter for width (it feels silly to comment
// these things, seriously)
func (c *Canvas) Width() uint32 {

	return c.width
}

// Height : A getter for height
func (c *Canvas) Height() uint32 {

	return c.height
}

//
// Canvas builder
//

// CanvasBuilder : Used to build a canvas
type CanvasBuilder struct {
	width  uint32
	height uint32
}

// NewCanvasBuilder : Allocated memory for a new
// canvas builder
func NewCanvasBuilder() *CanvasBuilder {

	builder := new(CanvasBuilder)

	return builder
}

// Build : Builds a canvas from the given canvas builder
func (cbuilder *CanvasBuilder) Build() *Canvas {

	c := new(Canvas)

	c.width = cbuilder.width
	c.height = cbuilder.height

	c.translation = NewPoint(0, 0)

	c.srcRect = sdl.Rect{X: 0, Y: 0, W: 0, H: 0}
	c.destRect = sdl.Rect{X: 0, Y: 0, W: 0, H: 0}

	return c
}

// SetDimensions : Set desired dimensions for the canvas to be built
func (cbuilder *CanvasBuilder) SetDimensions(width, height uint32) *CanvasBuilder {

	cbuilder.width = width
	cbuilder.height = height

	return cbuilder
}

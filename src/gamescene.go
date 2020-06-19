package main

import (
	"strconv"

	"github.com/jani-nykanen/ultimate-puzzle/src/core"
)

type gameScene struct {
	gameStage       *stage
	objects         *objectManager
	cloudPos        int32
	failureTimer    int32
	failed          bool
	cogSprite       *core.Sprite
	frameTransition *core.TransitionManager
}

func (game *gameScene) Activate(ev *core.Event, param interface{}) error {

	var err error

	game.gameStage, err = newStage(1, ev)
	if err != nil {

		return err
	}

	game.objects = newObjectManager()
	game.gameStage.parseObjects(game.objects)

	game.cloudPos = 0
	game.failureTimer = 0
	game.failed = false

	game.cogSprite = core.NewSprite(48, 48)

	game.frameTransition = core.NewTransitionManager()

	return err
}

func (game *gameScene) resetEvent(ev *core.Event) {

	game.gameStage.reset()
	game.objects.clear()
	game.gameStage.parseObjects(game.objects)

	game.failed = false
	game.failureTimer = 0
}

func (game *gameScene) updateBackground(step int32) {

	game.cloudPos = (game.cloudPos + step) % (512)
}

func (game *gameScene) reset(ev *core.Event) {

	cb := func(ev *core.Event) {
		game.resetEvent(ev)
	}
	game.frameTransition.Activate(true, core.TransitionHorizontalBar,
		30, core.NewRGB(0, 0, 0), cb)
}

func (game *gameScene) Refresh(ev *core.Event) {

	const failTime int32 = 60

	game.updateBackground(ev.Step())

	if game.frameTransition.Active() {

		game.frameTransition.Update(ev)
		return
	}

	game.gameStage.update(ev)
	if !game.failed {

		if ev.Input.GetActionState("reset") == core.StatePressed {

			game.reset(ev)
			return
		}

		if game.objects.update(game.gameStage, ev) {

			game.failed = true
			game.failureTimer = failTime

			game.gameStage.shake(failTime)
		}

	} else {

		game.failureTimer -= ev.Step()

		if game.failureTimer <= 0 {

			game.gameStage.shake(0)
			game.reset(ev)
		}
	}
}

func (game *gameScene) drawBackground(c *core.Canvas, bmp *core.Bitmap) {

	const cloudPosY = int32(96)
	const sunOffX = 56
	const sunPosY = 32

	c.Clear(145, 218, 255)

	// Draw sun
	c.DrawBitmapRegion(bmp, 128, 0, 48, 48,
		int32(c.Width())-sunOffX, sunPosY, core.FlipNone)

	// Draw clouds
	for i := int32(0); i < 3; i++ {

		c.DrawBitmapRegion(bmp, 0, 0, 128, 96,
			-(game.cloudPos/4)+128*i,
			cloudPosY, core.FlipNone)
	}
}

func (game *gameScene) drawFailureCross(c *core.Canvas, ap *core.AssetPack) {

	t := core.RoundFloat32(float32(game.gameStage.shakeTimer) / 15.0)
	if t%2 == 1 {

		return
	}

	topLeft := game.gameStage.getTopLeftCorner(c)

	px := game.objects.failurePoint.X + topLeft.X
	py := game.objects.failurePoint.Y + topLeft.Y

	c.DrawBitmap(ap.GetAsset("cross").(*core.Bitmap),
		px-12, py-12, core.FlipNone)
}

func (game *gameScene) DrawHUD(c *core.Canvas, ap *core.AssetPack) {

	const shadowOff int32 = 1

	alpha := []uint8{85, 255}
	color := []uint8{0, 255}

	bmpFont := ap.GetAsset("font").(*core.Bitmap)

	nameXOff := int32(len(game.gameStage.name)+2) * 8
	// Beautiful...
	moveStrLeft := "Moves: "
	moveStrMiddle := strconv.Itoa(int(game.objects.moveCount))
	moveStrRight := "(" + string(rune(5)) +
		strconv.Itoa(int(game.gameStage.bonusMoveLimit)) +
		")"
	moveStr := moveStrLeft + moveStrMiddle + moveStrRight
	moveXOff := int32(len(moveStr)) * 8

	for i := int32(0); i < 2; i++ {

		c.SetBitmapAlpha(bmpFont, alpha[i])
		c.SetBitmapColor(bmpFont, color[i], color[i], color[i])

		c.MoveTo((1-i)*shadowOff, (1-i)*shadowOff)

		// Stage number
		c.DrawText(bmpFont,
			"STAGE "+strconv.Itoa(int(game.gameStage.id)),
			8, 6, 0, 0, false)

		// Stage name
		c.DrawText(bmpFont, "\""+game.gameStage.name+"\"",
			c.Viewport().W-nameXOff-6, 6, 0, 0, false)

		// Blocks left
		c.DrawText(bmpFont,
			string(rune(3))+" Left: "+strconv.Itoa(int(game.objects.blockCount)),
			8, c.Viewport().H-12, -1, 0, false)

		// Moves
		c.DrawText(bmpFont, string(rune(4)),
			c.Viewport().W-moveXOff-19, c.Viewport().H-12,
			0, 0, false)
		c.DrawText(bmpFont, moveStrLeft,
			c.Viewport().W-moveXOff-6, c.Viewport().H-12,
			0, 0, false)

		if game.objects.moveCount > game.gameStage.bonusMoveLimit {

			c.SetBitmapColor(bmpFont, 255, 0, 0)
		}
		c.DrawText(bmpFont, moveStrMiddle,
			c.Viewport().W-moveXOff-6+int32(len(moveStrLeft))*8,
			c.Viewport().H-12,
			0, 0, false)
		c.SetBitmapColor(bmpFont, 255, 255, 255)

		c.DrawText(bmpFont, moveStrRight,
			c.Viewport().W-moveXOff-6+int32(len(moveStrLeft+moveStrMiddle))*8,
			c.Viewport().H-12,
			0, 0, false)

	}
}

func (game *gameScene) Redraw(c *core.Canvas, ap *core.AssetPack) {

	c.MoveTo(0, 0)

	// This needs to be called before anything else because when
	// (re)starting the stage, calling this after background
	// drawing will cause one frame of weirdness
	game.gameStage.preDraw(c, ap)

	game.drawBackground(c, ap.GetAsset("background").(*core.Bitmap))

	game.gameStage.refreshShadowLayer(c, ap, game.objects)

	game.gameStage.setViewport(c)

	// Background stuff, drawn before outlines
	game.gameStage.drawBackground(c, ap)
	// Outlines
	game.gameStage.drawOutlines(c)
	game.objects.drawOutlines(c, ap)
	// Base drawing
	game.gameStage.draw(c, ap)
	game.objects.draw(c, ap)
	game.gameStage.postDraw(c, ap)

	game.frameTransition.Draw(c)

	c.ResetViewport()

	game.gameStage.drawDecorations(c, ap)

	c.MoveTo(0, 0)
	if game.failed && !game.frameTransition.Active() {

		game.drawFailureCross(c, ap)
	}

	game.DrawHUD(c, ap)
}

func (game *gameScene) Dispose() {

	game.gameStage.dispose()
}

func newGameScene() core.Scene {

	return new(gameScene)
}

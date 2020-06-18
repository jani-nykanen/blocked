package main

import (
	"github.com/jani-nykanen/ultimate-puzzle/src/core"
)

type gameScene struct {
	gameStage *stage
	objects   *objectManager
	cloudPos  int32
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

	return err
}

func (game *gameScene) updateBackground(step int32) {

	game.cloudPos = (game.cloudPos + step) % (512)
}

func (game *gameScene) Refresh(ev *core.Event) {

	game.updateBackground(ev.Step())

	game.gameStage.update(ev)
	game.objects.update(game.gameStage, ev)
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

func (game *gameScene) Redraw(c *core.Canvas, ap *core.AssetPack) {

	game.drawBackground(c, ap.GetAsset("background").(*core.Bitmap))

	// Shadows
	game.gameStage.refreshShadowLayer(c, ap, game.objects)

	game.gameStage.setViewport(c)

	// Background stuff, drawn before outlines
	game.gameStage.drawBackground(c, ap)

	// Outlines
	game.gameStage.drawOutlines(c)
	game.objects.drawOutlines(c, ap, game.gameStage)

	// Base drawing
	game.gameStage.draw(c, ap)
	game.objects.draw(c, ap, game.gameStage)
	game.gameStage.postDraw(c, ap)

	c.ResetViewport()

	game.gameStage.drawDecorations(c, ap)

	c.MoveTo(0, 0)
	c.DrawText(ap.GetAsset("font").(*core.Bitmap),
		"v.0.0.1", 2, 2, -1, 0, false)
}

func (game *gameScene) Dispose() {

	game.gameStage.dispose()
}

func newGameScene() core.Scene {

	return new(gameScene)
}

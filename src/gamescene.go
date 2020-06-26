package main

import (
	"strconv"

	"github.com/jani-nykanen/ultimate-puzzle/src/core"
)

type gameScene struct {
	gameStage        *stage
	objects          *objectManager
	cloudPos         int32
	failureTimer     int32
	failed           bool
	cleared          bool
	cogSprite        *core.Sprite
	frameTransition  *core.TransitionManager
	pauseMenu        *menu
	clearMenu        *menu
	settingsScreen   *settings
	bestSuccessState int32
}

type levelResult struct {
	currentStage int32
	successState int32
}

func (game *gameScene) createPauseMenu() {

	buttons := []menuButton{

		newMenuButton("Resume", func(ev *core.Event) {
			game.pauseMenu.deactivate()
		}),
		newMenuButton("Reset", func(ev *core.Event) {
			game.reset(ev)
			game.pauseMenu.deactivate()
		}),
		newMenuButton("Settings", func(ev *core.Event) {
			game.settingsScreen.activate()
		}),
		newMenuButton("Quit", func(ev *core.Event) {

			game.pauseMenu.deactivate()

			ev.Transition.Activate(true, core.TransitionCircleOutside, 30,
				core.NewRGB(0, 0, 0), func(ev *core.Event) {
					err := ev.ChangeScene(newLevelMenuScene())
					if err != nil {

						ev.Terminate(err)
					}
				})
		}),
	}

	game.pauseMenu = newMenu(buttons)
}

func (game *gameScene) createClearMenu() {

	buttons := []menuButton{

		newMenuButton("Play Again", func(ev *core.Event) {
			game.reset(ev)
			game.pauseMenu.deactivate()
		}),
		/*
			newMenuButton("Next stage", func(ev *core.Event) {
				// ...
			}),
		*/
		newMenuButton("Stage Menu", func(ev *core.Event) {

			ev.Transition.Activate(true, core.TransitionCircleOutside, 30,
				core.NewRGB(0, 0, 0), func(ev *core.Event) {
					err := ev.ChangeScene(newLevelMenuScene())
					if err != nil {

						ev.Terminate(err)
					}
				})
		}),
	}

	game.clearMenu = newMenu(buttons)
}

func (game *gameScene) Activate(ev *core.Event, param interface{}) error {

	var err error
	index := int32(1)
	if param != nil {

		index = param.(int32)
	}

	game.gameStage, err = newStage(index, ev)
	if err != nil {

		return err
	}

	game.objects = newObjectManager()
	game.gameStage.parseObjects(game.objects)

	game.cloudPos = 0
	game.failureTimer = 0
	game.failed = false
	game.cleared = false
	game.bestSuccessState = 0

	game.cogSprite = core.NewSprite(48, 48)

	game.frameTransition = core.NewTransitionManager()

	game.settingsScreen = newSettings()
	game.createPauseMenu()
	game.createClearMenu()

	return err
}

func (game *gameScene) resetEvent(ev *core.Event) {

	game.gameStage.reset()
	game.objects.clear()
	game.gameStage.parseObjects(game.objects)

	game.clearMenu.deactivate()

	game.failed = false
	game.cleared = false
	game.failureTimer = 0
}

func (game *gameScene) updateBackground(step int32) {

	game.cloudPos = (game.cloudPos + step) % (512)
}

func (game *gameScene) reset(ev *core.Event) {

	cb := func(ev *core.Event) {
		game.resetEvent(ev)

		game.frameTransition.ResetCenter()
	}
	game.frameTransition.Activate(true, core.TransitionCircleOutside,
		30, core.NewRGB(0, 0, 0), cb)
	if game.failed {

		game.frameTransition.SetCenter(
			game.objects.failurePoint.X,
			game.objects.failurePoint.Y)
	}

}

func (game *gameScene) Refresh(ev *core.Event) {

	const failTime int32 = 60

	if ev.Transition.Active() {
		return
	}

	if game.settingsScreen.active() {

		game.settingsScreen.update(ev)
		return
	}

	if !game.pauseMenu.active {
		game.updateBackground(ev.Step())
	}

	// Transition
	if game.frameTransition.Active() {

		game.frameTransition.Update(ev)
		return
	}

	if !game.cleared {

		// Pause menu
		if game.pauseMenu.active {

			game.pauseMenu.update(ev)
			return

		} else if !game.failed &&
			ev.Input.GetActionState("start") == core.StatePressed {

			game.pauseMenu.activate(0)
			return
		}

	} else {

		game.clearMenu.update(ev)
	}

	// The rest
	game.gameStage.update(ev)
	if !game.failed {

		if !game.cleared &&
			ev.Input.GetActionState("reset") == core.StatePressed {

			game.reset(ev)
			return
		}

		if game.objects.update(game.gameStage, ev) {

			game.failed = true
			game.failureTimer = failTime

			game.gameStage.shake(failTime)
		}
		game.cleared = game.objects.cleared
		if game.cleared && !game.clearMenu.active {

			game.clearMenu.activate(0)

			game.bestSuccessState = core.MaxInt32(game.bestSuccessState, 1)
			if game.objects.moveCount <= game.gameStage.bonusMoveLimit {

				game.bestSuccessState = 2
			}
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
	// Might need to come up with a better name to
	// the hardest difficulty...
	difficultyNames := []string{
		"Easy", "Average", "Hard", "Expert"}

	bmpFont := ap.GetAsset("font").(*core.Bitmap)

	// Beautiful...
	moveStrLeft := "Moves: "
	moveStrMiddle := strconv.Itoa(int(game.objects.moveCount))
	moveStrRight := "(" + string(rune(5)) +
		strconv.Itoa(int(game.gameStage.bonusMoveLimit)) +
		")"
	moveStr := moveStrLeft + moveStrMiddle + moveStrRight
	moveXOff := int32(len(moveStr)) * 8

	diff := core.ClampInt32(game.gameStage.difficulty, 1, 4)
	difficultyStr := string(rune(5+diff)) +
		" " +
		difficultyNames[diff-1]
	diffXOff := int32(len(difficultyStr)) * 7

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
			c.Viewport().W/2, 6, 0, 0, true)

		// Stage difficulty
		c.DrawText(bmpFont, difficultyStr,
			c.Viewport().W-diffXOff-6, 6, -1, 0, false)

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

func (game *gameScene) drawSuccess(c *core.Canvas, ap *core.AssetPack) {

	const headerOff int32 = 16
	const startOff int32 = headerOff + 24

	bmp := ap.GetAsset("stageClear").(*core.Bitmap)

	c.FillRect(0, 0, c.Viewport().W, c.Viewport().H,
		core.NewRGBA(0, 0, 0, 85))

	c.DrawBitmapRegion(bmp, 0, 0, 128, 16,
		c.Viewport().W/2-64, c.Viewport().H/2-headerOff,
		core.FlipNone)

	sx := int32(0)
	if game.objects.moveCount <= game.gameStage.bonusMoveLimit {

		sx = 24
	}
	c.DrawBitmapRegion(bmp, sx, 16, 24, 24,
		c.Viewport().W/2-12, c.Viewport().H/2-startOff,
		core.FlipNone)

	c.MoveTo(0, 32)
	game.clearMenu.draw(c, ap, false)
}

func (game *gameScene) Redraw(c *core.Canvas, ap *core.AssetPack) {

	c.MoveTo(0, 0)
	c.ResetViewport()

	if game.settingsScreen.active() {

		game.settingsScreen.draw(c, ap)
		return
	}

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

	if game.cleared {

		game.drawSuccess(c, ap)
	}

	game.frameTransition.Draw(c)

	c.ResetViewport()

	game.gameStage.drawDecorations(c, ap)

	c.MoveTo(0, 0)
	if game.failed && !game.frameTransition.Active() {

		game.drawFailureCross(c, ap)
	}

	game.DrawHUD(c, ap)

	game.pauseMenu.draw(c, ap, true)

}

func (game *gameScene) Dispose() interface{} {

	game.gameStage.dispose()

	ret := levelResult{currentStage: game.gameStage.id,
		successState: game.bestSuccessState}

	return ret
}

func newGameScene() core.Scene {

	return new(gameScene)
}

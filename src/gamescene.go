package main

import (
	"fmt"
	"strconv"

	"github.com/jani-nykanen/blocked/src/core"
)

const (
	gameClearTime int32 = 60
)

type gameScene struct {
	gameStage       *stage
	objects         *objectManager
	cloudPos        int32
	failureTimer    int32
	failed          bool
	cleared         bool
	clearTimer      int32
	endingAchieved  bool
	cogSprite       *core.Sprite
	frameTransition *core.TransitionManager
	pauseMenu       *menu
	clearMenu       *menu
	settingsScreen  *settings
	cinfo           *completionInfo
}

func (game *gameScene) createPauseMenu() {

	buttons := []menuButton{

		newMenuButton("Resume", func(self *menuButton, dir int32, ev *core.Event) {
			game.pauseMenu.deactivate()
		}, false),
		newMenuButton("Reset", func(self *menuButton, dir int32, ev *core.Event) {
			game.reset(ev)
			game.pauseMenu.deactivate()
		}, false),
		newMenuButton("Settings", func(self *menuButton, dir int32, ev *core.Event) {
			game.settingsScreen.activate()
		}, false),
		newMenuButton("Quit", func(self *menuButton, dir int32, ev *core.Event) {

			game.pauseMenu.deactivate()

			ev.Transition.Activate(true, core.TransitionCircleOutside, 30,
				core.NewRGB(0, 0, 0), func(ev *core.Event) {
					ev.ChangeScene(newLevelMenuScene())
				})
		}, false),
	}

	game.pauseMenu = newMenu(buttons, true, "")
}

func (game *gameScene) createClearMenu() {

	buttons := []menuButton{

		newMenuButton("Play Again", func(self *menuButton, dir int32, ev *core.Event) {
			game.reset(ev)
			game.pauseMenu.deactivate()
		}, false),

		newMenuButton("Next Stage", func(self *menuButton, dir int32, ev *core.Event) {

			game.frameTransition.Activate(true, core.TransitionCircleOutside,
				30, core.NewRGB(0, 0, 0),
				func(ev *core.Event) {

					game.nextStage(ev)
				})
		}, false),

		newMenuButton("Stage Menu", func(self *menuButton, dir int32, ev *core.Event) {

			ev.Transition.Activate(true, core.TransitionCircleOutside, 30,
				core.NewRGB(0, 0, 0), func(ev *core.Event) {
					ev.ChangeScene(newLevelMenuScene())

				})
		}, false),
	}

	game.clearMenu = newMenu(buttons, false, "")
}

func (game *gameScene) Activate(ev *core.Event, param interface{}) error {

	var err error
	index := int32(1)
	if param != nil {

		game.cinfo = param.(*completionInfo)
		index = game.cinfo.currentStage
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
	game.cogSprite = core.NewSprite(48, 48)

	game.frameTransition = core.NewTransitionManager()

	game.settingsScreen = newSettings(ev)
	game.createPauseMenu()
	game.createClearMenu()

	return err
}

func (game *gameScene) nextStage(ev *core.Event) {

	game.cinfo.currentStage = (game.cinfo.currentStage % game.cinfo.levelCount()) + 1

	var err error
	game.gameStage, err = newStage(game.cinfo.currentStage, ev)
	if err != nil {

		ev.Terminate(err)
		return
	}

	game.resetEvent(false)
}

func (game *gameScene) resetEvent(resetStage bool) {

	if resetStage {

		game.gameStage.reset()
	}
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

	game.frameTransition.Activate(true, core.TransitionCircleOutside,
		30, core.NewRGB(0, 0, 0),
		func(ev *core.Event) {

			game.resetEvent(true)
			game.frameTransition.ResetCenter()
		})
	if game.failed {

		game.frameTransition.SetCenter(
			game.objects.failurePoint.X,
			game.objects.failurePoint.Y)
	}

}

func (game *gameScene) Refresh(ev *core.Event) {

	const failTime int32 = 60
	const clearTimerSpeed int32 = 2

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
			(ev.Input.GetActionState("start") == core.StatePressed ||
				ev.Input.GetActionState("back") == core.StatePressed) {

			ev.Audio.PlaySample(ev.Assets.GetAsset("pause").(*core.Sample), 30)
			game.pauseMenu.activate(0)
			return
		}

	} else {

		if game.clearTimer <= 0 {

			game.clearMenu.update(ev)

		} else {

			game.clearTimer -= clearTimerSpeed * ev.Step()
			if game.endingAchieved && game.clearTimer <= 0 {

				ev.Transition.Activate(false, core.TransitionVerticalBar, 60,
					core.NewRGB(255, 255, 255), nil)

				ev.ChangeScene(newEndingScene())
			}
		}
	}

	// The rest
	var state int32
	game.gameStage.update(ev)
	if !game.failed {

		if !game.cleared &&
			ev.Input.GetActionState("reset") == core.StatePressed {

			ev.Audio.PlaySample(ev.Assets.GetAsset("restart").(*core.Sample), 40)

			game.reset(ev)
			return
		}

		if game.objects.update(game.gameStage, ev) {

			game.failed = true
			game.failureTimer = failTime

			game.gameStage.shake(failTime)
		}

		game.cleared = game.objects.cleared || game.cleared

		if game.cleared && !game.clearMenu.active {

			game.clearTimer = gameClearTime
			game.clearMenu.activate(1)

			ev.Audio.StopMusic()
			ev.Audio.PlayMusic(ev.Assets.GetAsset("victory").(*core.Music), 50, 1)

			state = 1
			if game.objects.moveCount <= game.gameStage.bonusMoveLimit {

				state = 2
			}
			game.cinfo.updateState(game.gameStage.id, state)

			game.endingAchieved = game.cinfo.checkIfNewEndingObtained()
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

func (game *gameScene) drawHUD(c *core.Canvas, ap *core.AssetPack) {

	const shadowOff int32 = 1

	alpha := []uint8{85, 255}
	color := []uint8{0, 255}
	// Might need to come up with a better name to
	// the hardest difficulty...

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
		getDifficultyName(diff-1)
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
	const starOff int32 = headerOff + 24

	bmp := ap.GetAsset("stageClear").(*core.Bitmap)

	c.FillRect(0, 0, c.Viewport().W, c.Viewport().H,
		core.NewRGBA(0, 0, 0, 85))

	c.DrawBitmapRegion(bmp, 0, 0, 128, 16,
		c.Viewport().W/2-64, c.Viewport().H/2-headerOff+game.clearTimer*2,
		core.FlipNone)

	sx := int32(0)
	if game.objects.moveCount <= game.gameStage.bonusMoveLimit {

		sx = 24
	}
	c.DrawBitmapRegion(bmp, sx, 16, 24, 24,
		c.Viewport().W/2-12, c.Viewport().H/2-starOff-game.clearTimer*2,
		core.FlipNone)

	c.MoveTo(-game.clearTimer*2, 32)
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

	if game.cleared && !game.endingAchieved {

		game.drawSuccess(c, ap)
	}

	game.frameTransition.Draw(c)

	c.ResetViewport()

	game.gameStage.drawDecorations(c, ap)

	c.MoveTo(0, 0)
	if game.failed && !game.frameTransition.Active() {

		game.drawFailureCross(c, ap)
	}

	game.drawHUD(c, ap)

	game.pauseMenu.draw(c, ap, true)

}

func (game *gameScene) Dispose() interface{} {

	game.gameStage.dispose()

	err := game.cinfo.saveToFile(defaultSaveFilePath)
	if err != nil {

		fmt.Printf("Error writing the save file: %s\n", err.Error())
		return nil
	}

	return game.cinfo
}

func newGameScene() core.Scene {

	return new(gameScene)
}

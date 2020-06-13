package main

import (
	"github.com/jani-nykanen/ultimate-puzzle/src/core"
)

//
// The main (application-specific) game logic happens here
//
// (c) 2020 Jani Nykänen
//

type gameScene struct {
	testPos core.Vector2
}

func (game *gameScene) Activate(ev *core.Event, param interface{}) error {

	game.testPos = core.NewVector2(128, 96)

	return nil
}

func (game *gameScene) Refresh(ev *core.Event) {

	step := float32(ev.Step())

	if ev.Input.GetActionState("left")&core.StateDownOrPressed == 1 {

		game.testPos.X -= 1.0 * step

	} else if ev.Input.GetActionState("right")&core.StateDownOrPressed == 1 {

		game.testPos.X += 1.0 * step
	}
	/*
		if ev.Input.GetActionState("up")&core.StateDownOrPressed == 1 {

			game.testPos.Y -= 1.0 * step

		} else if ev.Input.GetActionState("down")&core.StateDownOrPressed == 1 {

			game.testPos.Y += 1.0 * step
		}
	*/
}

func (game *gameScene) Redraw(c *core.Canvas, ap *core.AssetPack) {

	c.Clear(170, 170, 170)

	c.MoveTo(int32(game.testPos.X), int32(game.testPos.Y))
	c.FillRect(-12, -12, 24, 24, core.NewRGB(255, 0, 0))
	c.MoveTo(0, 0)

	c.DrawText(ap.GetAsset("font").(*core.Bitmap),
		"Hello world!", 2, 2, -1, 0, false)
}

func (game *gameScene) Dispose() {

}

func newGameScene() core.Scene {

	return new(gameScene)
}

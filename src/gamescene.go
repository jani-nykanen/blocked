package main

import (
	"github.com/jani-nykanen/ultimate-puzzle/src/core"
)

type gameScene struct {
	gameStage *stage
}

func (game *gameScene) Activate(ev *core.Event, param interface{}) error {

	var err error

	game.gameStage, err = newStage(1, ev)
	if err != nil {

		return err
	}

	return err
}

func (game *gameScene) Refresh(ev *core.Event) {

	game.gameStage.update(ev)
}

func (game *gameScene) Redraw(c *core.Canvas, ap *core.AssetPack) {

	c.Clear(170, 170, 170)

	game.gameStage.setCamera(c)
	game.gameStage.draw(c, ap)

	c.MoveTo(0, 0)
	c.DrawText(ap.GetAsset("font").(*core.Bitmap),
		"v.0.0.1", 2, 2, -1, 0, false)
}

func (game *gameScene) Dispose() {

}

func newGameScene() core.Scene {

	return new(gameScene)
}

package main

import "github.com/jani-nykanen/ultimate-puzzle/src/core"

//
// The main (application-specific) game logic happens here
//
// (c) 2020 Jani Nykänen
//

type gameScene struct {
	// ...
}

func (game *gameScene) Activate(ev *core.Event, param interface{}) error {

	return nil
}

func (game *gameScene) Refresh(ev *core.Event) {

}

func (game *gameScene) Redraw(c *core.Canvas, ap *core.AssetPack) {

	c.Clear(170, 170, 170)

	c.DrawText(ap.GetAsset("font").(*core.Bitmap),
		"Hello world!", 2, 2, -1, 0, false)
}

func (game *gameScene) Dispose() {

}

func newGameScene() core.Scene {

	return new(gameScene)
}

package main

import (
	"github.com/jani-nykanen/ultimate-puzzle/src/core"
)

const (
	introAppearTime int32 = 60
	introWaitTime   int32 = 120
	introLeaveTime  int32 = 60
)

type introScene struct {
	timer     int32
	ballPos   core.Vector2
	ballSpeed core.Vector2
}

func (intro *introScene) Activate(ev *core.Event, param interface{}) error {

	const initialWait int32 = -30

	intro.ballSpeed.X = 1.45
	intro.ballSpeed.Y = 0.0
	intro.ballPos.X = -32
	intro.ballPos.Y = 16

	intro.timer = initialWait

	return nil
}

func (intro *introScene) Refresh(ev *core.Event) {

	const ballFloorY float32 = 112
	const ballGravity float32 = 0.125
	const ballJumpMod float32 = 0.80

	if ev.Transition.Active() {
		return
	}

	intro.timer += ev.Step()
	if intro.timer < 0 {

		return
	}

	timeSum := introAppearTime + introWaitTime + introLeaveTime

	if intro.timer >= timeSum {

		ev.Transition.Activate(false, core.TransitionCircleOutside,
			60, core.NewRGB(0, 0, 0), nil)
		ev.ChangeScene(newTitleScreenScene())
	}

	intro.ballSpeed.Y += ballGravity * float32(ev.Step())
	if intro.ballSpeed.Y > 0 &&
		intro.ballPos.Y >= ballFloorY {

		intro.ballPos.Y = ballFloorY
		intro.ballSpeed.Y *= -ballJumpMod

		ev.Audio.PlaySample(ev.Assets.GetAsset("hit").(*core.Sample),
			40)
	}

	intro.ballPos.X += intro.ballSpeed.X * float32(ev.Step())
	intro.ballPos.Y += intro.ballSpeed.Y * float32(ev.Step())
}

func (intro *introScene) Redraw(c *core.Canvas, ap *core.AssetPack) {

	const textY int32 = 112

	c.MoveTo(0, 0)
	c.ResetViewport()
	c.Clear(0, 0, 0)

	bmp := ap.GetAsset("createdBy").(*core.Bitmap)
	_ = bmp

	var y int32

	y = textY
	if intro.timer < introAppearTime {

		y += (introAppearTime - intro.timer) * 2

	} else if intro.timer >= introAppearTime+introWaitTime {

		y += (intro.timer - (introAppearTime + introWaitTime)) * 2
	}

	c.DrawBitmapRegion(bmp, 0, 64, 128, 32,
		c.Viewport().W/2-64, y, core.FlipNone)

	c.DrawBitmapRegion(bmp, 0, 0, 64, 64,
		core.RoundFloat32(intro.ballPos.X)-32,
		core.RoundFloat32(intro.ballPos.Y)-64,
		core.FlipNone)
}

func (intro *introScene) Dispose() interface{} {

	return nil
}

func newIntroScene() core.Scene {

	return new(introScene)
}

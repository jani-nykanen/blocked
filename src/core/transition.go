package core

// TransitionCallback : A callback for
// transition events
type TransitionCallback func(ev *Event)

// TransitionMode : Transition mode type
type TransitionMode int32

// Transition types
const (
	TransitionNone          = 0
	TransitionFade          = 1
	TransitionVerticalBar   = 2
	TransitionHorizontalBar = 3
	TransitionCircleOutside = 4
)

// TransitionManager : Used for different kind of
// transitions
type TransitionManager struct {
	timer           int32
	time            int32
	mode            TransitionMode
	fadeIn          bool
	color           Color
	active          bool
	cb              TransitionCallback
	center          Point
	centerSpecified bool
	textureCopied   bool
}

// Update : Update the transition manager
func (tr *TransitionManager) Update(ev *Event) {

	if !tr.active {

		return
	}

	tr.timer -= ev.Step()
	if tr.timer <= 0 {

		if tr.fadeIn {

			tr.timer += tr.time
			tr.fadeIn = false

			if tr.cb != nil {

				tr.cb(ev)
			}

		} else {

			tr.timer = 0
			tr.active = false
		}
	}
}

// Draw : Draw the transition effect, if any
func (tr *TransitionManager) Draw(c *Canvas) {

	if !tr.active {
		return
	}

	c.MoveTo(0, 0)

	t := float32(tr.timer) / float32(tr.time)
	if tr.fadeIn {

		t = 1.0 - t
	}

	var p int32
	var radius, maxRadius int32

	// TODO: Implement the rest
	switch tr.mode {

	case TransitionVerticalBar:

		p = RoundFloat32((1 - t) * float32(c.viewport.W/2))

		// Left half
		c.DrawCopiedFrameRegion(0, 0, c.viewport.W/2, c.viewport.H,
			-p, 0, FlipNone)

		// Right half
		c.DrawCopiedFrameRegion(int32(c.frameCopy.width)/2, 0,
			c.viewport.W/2, c.viewport.H,
			c.viewport.W/2+p, 0, FlipNone)

		break

	case TransitionHorizontalBar:

		p = RoundFloat32((1 - t) * float32(c.viewport.H/2))

		// Upper half
		c.DrawCopiedFrameRegion(0, 0, c.viewport.W, c.viewport.H/2,
			0, -p, FlipNone)

		// Bottom half
		c.DrawCopiedFrameRegion(0, int32(c.frameCopy.height)/2,
			c.viewport.W, c.viewport.H/2,
			0, c.viewport.H/2+p, FlipNone)

		break

	case TransitionCircleOutside:

		if !tr.centerSpecified {

			tr.center.X = c.viewport.W / 2
			tr.center.Y = c.viewport.H / 2

			tr.centerSpecified = true
		}

		t = 1.0 - t
		t *= t

		// A lot of unnecessary computations happen here
		// because I forgot what data types I needed, but I'm
		// too lazy to implement better methods
		maxRadius = MaxInt32InSlice([]int32{
			HypotInt32(tr.center.X, tr.center.Y),
			HypotInt32(c.Viewport().W-tr.center.X, tr.center.Y),
			HypotInt32(c.Viewport().W-tr.center.X, c.Viewport().H-tr.center.Y),
			HypotInt32(tr.center.X, c.Viewport().H-tr.center.Y),
		})
		radius = RoundFloat32(t * float32(maxRadius))
		c.FillCircleOutside(tr.center.X, tr.center.Y, radius, tr.color)

		break

	default:
		break
	}

}

// Activate : Activate the transition manager
func (tr *TransitionManager) Activate(fadeIn bool,
	mode TransitionMode, time int32,
	color Color, cb TransitionCallback) {

	tr.active = true

	tr.fadeIn = fadeIn
	tr.mode = mode
	tr.time = MaxInt32(1, time) // To avoid division by 0
	tr.timer = tr.time
	tr.cb = cb
	tr.color = color

	tr.textureCopied = false
	tr.centerSpecified = false
}

// SetCenter : Set center position for specific transition modes
func (tr *TransitionManager) SetCenter(x, y int32) {

	tr.center = NewPoint(x, y)

	tr.centerSpecified = true
}

// ResetCenter : Reset the center to the middle of the screen
func (tr *TransitionManager) ResetCenter() {

	tr.centerSpecified = false
}

// SetNewTime : Set a new time
func (tr *TransitionManager) SetNewTime(time int32) {

	tr.time = time
	tr.timer = time
}

// Active : Getter for "active" property
func (tr *TransitionManager) Active() bool {

	return tr.active
}

// NewTransitionManager : Constructor for a transition
// manager
func NewTransitionManager() *TransitionManager {

	tr := new(TransitionManager)

	tr.active = false
	tr.color = NewRGB(0, 0, 0)

	tr.time = 60
	tr.timer = 0

	tr.centerSpecified = false

	return tr
}

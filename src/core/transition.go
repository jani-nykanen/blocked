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

	var w, h int32
	var radius, maxRadius int32

	// TODO: Implement the rest
	switch tr.mode {

	case TransitionHorizontalBar:

		w = c.viewport.W
		h = RoundFloat32(t * float32(c.viewport.H/2))

		// Upper half
		c.FillRect(0, 0, w, h, tr.color)
		// Bottom half
		c.FillRect(0, c.viewport.H-h, w, h, tr.color)

		break

	case TransitionCircleOutside:

		if !tr.centerSpecified {

			tr.center.X = c.viewport.W / 2
			tr.center.Y = c.viewport.H / 2

			tr.centerSpecified = true
		}

		t = 1.0 - t

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

	tr.centerSpecified = false
}

// SetCenter : Set center position for specific transition modes
func (tr *TransitionManager) SetCenter(x, y int32) {

	tr.center = NewPoint(x, y)

	tr.centerSpecified = true
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

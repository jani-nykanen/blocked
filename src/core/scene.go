package core

// Scene : An interface for a generic application scene
type Scene interface {
	Activate(ev *Event, param interface{}) error
	Refresh(ev *Event)
	Redraw(c *Canvas, ap *AssetPack)
	Dispose() interface{}
}

package component

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/filter"
)

var errNoCameraFound = newComponentError("no camera found")

type CameraData struct {
	Viewport math.Vec2
}

func (c *CameraData) ViewportCenter() math.Vec2 {
	return math.NewVec2(c.Viewport.X*0.5, c.Viewport.Y*0.5)
}

var Camera = donburi.NewComponentType[CameraData]()

func MustFindCamera(w donburi.World) *CameraData {
	camera, ok := donburi.NewQuery(filter.Contains(Camera)).First(w)
	if !ok {
		panic(errNoCameraFound)
	}
	return Camera.Get(camera)
}

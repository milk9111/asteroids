package archetype

import (
	"github.com/milk9111/asteroids/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/filter"
)

var errNoCameraFound = newArchetypeError("no camera found")

func NewCamera(w donburi.World, viewport math.Vec2) *donburi.Entry {
	camera := w.Entry(
		w.Create(
			component.Camera,
		),
	)

	component.Camera.SetValue(camera, component.CameraData{
		Viewport: viewport,
	})

	return camera
}

func MustFindCamera(w donburi.World) *donburi.Entry {
	camera, ok := donburi.NewQuery(filter.Contains(component.Camera)).First(w)
	if !ok {
		panic(errNoCameraFound)
	}

	return camera
}

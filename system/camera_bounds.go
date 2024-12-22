package system

import (
	"github.com/milk9111/asteroids/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
)

type CameraBounds struct {
	query  *donburi.Query
	game   *component.GameData
	camera *component.CameraData
}

func NewCameraBounds() *CameraBounds {
	return &CameraBounds{
		query: donburi.NewQuery(
			filter.Contains(
				transform.Transform,
				component.Sprite,
			),
		),
	}
}

func (cb *CameraBounds) Update(w donburi.World) {
	if cb.game == nil {
		cb.game = component.MustFindGame(w)
	}

	if cb.camera == nil {
		cb.camera = component.MustFindCamera(w)
	}

	cb.query.Each(w, func(e *donburi.Entry) {
		t := transform.Transform.Get(e)
		s := component.Sprite.Get(e)

		// switch left and right
		if t.LocalPosition.X < -s.Size.X {
			t.LocalPosition.X = cb.camera.Viewport.X + s.Size.X
		} else if t.LocalPosition.X > cb.camera.Viewport.X+s.Size.X {
			t.LocalPosition.X = -s.Size.X
		}

		// switch top and bottom
		if t.LocalPosition.Y < -s.Size.Y {
			t.LocalPosition.Y = cb.camera.Viewport.Y + s.Size.Y
		} else if t.LocalPosition.Y > cb.camera.Viewport.Y+s.Size.Y {
			t.LocalPosition.Y = -s.Size.Y
		}
	})
}

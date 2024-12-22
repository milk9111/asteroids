package archetype

import (
	"time"

	"github.com/milk9111/asteroids/component"
	"github.com/milk9111/asteroids/engine"
	"github.com/milk9111/asteroids/engine/path"
	"github.com/yohamta/donburi"
	dmath "github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
)

var ProjectileSpritePoints = []path.Op{
	path.NewMoveToOp(dmath.NewVec2(0, 0)),
	path.NewLineToOp(dmath.NewVec2(0, 2)),
	path.NewLineToOp(dmath.NewVec2(2, 2)),
	path.NewLineToOp(dmath.NewVec2(2, 0)),
}

func NewProjectile(w donburi.World, startingPos dmath.Vec2) *donburi.Entry {
	entry := w.Entry(
		w.Create(
			component.Traveler,
			transform.Transform,
			component.Sprite,
			component.Velocity,
			component.Collider,
			component.ID,
			component.TTL,
			component.CollisionQueue,
		),
	)

	component.ID.SetValue(entry, component.IDData{
		ID:   engine.NewID(),
		Type: "Projectile",
	})

	transform.Transform.Get(entry).LocalPosition = startingPos

	component.Sprite.SetValue(entry, component.SpriteData{
		PathOps: ProjectileSpritePoints,
		Layer:   component.SpriteLayerEntity,
		Pivot:   component.SpritePivotTopLeft,
		Size:    dmath.NewVec2(2, 2),
	})

	component.Collider.SetValue(entry, component.ColliderData{
		Width:  2,
		Height: 2,
		Layer:  component.CollisionLayerProjectile,
	})

	component.TTL.SetValue(entry, component.TTLData{
		Timer: engine.NewTimer(1500 * time.Millisecond),
	})

	return entry
}

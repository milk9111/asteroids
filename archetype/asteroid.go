package archetype

import (
	"github.com/milk9111/asteroids/component"
	"github.com/milk9111/asteroids/engine"
	"github.com/milk9111/asteroids/engine/path"
	"github.com/yohamta/donburi"
	dmath "github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
)

var AsteroidSpritePoints = []path.Op{
	path.NewMoveToOp(dmath.NewVec2(1, 8)),
	path.NewLineToOp(dmath.NewVec2(7, 2)),
	path.NewLineToOp(dmath.NewVec2(14, 7)),
	path.NewLineToOp(dmath.NewVec2(19, 1)),
	path.NewLineToOp(dmath.NewVec2(22, 8)),
	path.NewLineToOp(dmath.NewVec2(28, 5)),
	path.NewLineToOp(dmath.NewVec2(31, 17)),
	path.NewLineToOp(dmath.NewVec2(22, 25)),
	path.NewLineToOp(dmath.NewVec2(14, 26)),
	path.NewLineToOp(dmath.NewVec2(13, 31)),
	path.NewLineToOp(dmath.NewVec2(6, 23)),
	path.NewLineToOp(dmath.NewVec2(8, 21)),
	path.NewCloseOp(),
}

func NewAsteroid(w donburi.World, startingPos, direction dmath.Vec2, rotation, speed float64, tier component.AsteroidTier) *donburi.Entry {
	entry := w.Entry(
		w.Create(
			component.Traveler,
			transform.Transform,
			component.Sprite,
			component.Velocity,
			component.Collider,
			component.Asteroid,
			component.WaveLineItem,
			component.ID,
			component.CollisionQueue,
			component.AudioQueue,
		),
	)

	component.ID.SetValue(entry, component.IDData{
		ID:   engine.NewID(),
		Type: "Asteroid",
	})

	t := transform.Transform.Get(entry)
	t.LocalPosition = startingPos
	t.LocalRotation = rotation

	scale := dmath.NewVec2(1, 1)
	pointValue := 100
	switch tier {
	case component.AsteroidTierMedium:
		scale = dmath.NewVec2(2, 2)
		pointValue = 50
	case component.AsteroidTierLarge:
		scale = dmath.NewVec2(4, 4)
		pointValue = 20
	}

	t.LocalScale = scale

	component.Asteroid.SetValue(entry, component.AsteroidData{
		Tier:       tier,
		PointValue: pointValue,
	})

	component.Traveler.SetValue(entry, component.TravelerData{
		Direction: direction,
		Speed:     speed,
	})

	component.Sprite.SetValue(entry, component.SpriteData{
		PathOps: AsteroidSpritePoints,
		Layer:   component.SpriteLayerEntity,
		Pivot:   component.SpritePivotCenter,
		Size:    dmath.NewVec2(32, 32),
	})

	component.Collider.SetValue(entry, component.ColliderData{
		Width:  32,
		Height: 32,
		Layer:  component.CollisionLayerAsteroid,
	})

	return entry
}

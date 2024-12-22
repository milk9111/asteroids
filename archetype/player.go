package archetype

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/milk9111/asteroids/component"
	"github.com/milk9111/asteroids/engine"
	"github.com/milk9111/asteroids/engine/path"
	"github.com/yohamta/donburi"
	dmath "github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
)

var (
	errNoPlayerFound = newArchetypeError("no player found")
)

var PlayerSpritePoints = []path.Op{
	path.NewMoveToOp(dmath.NewVec2(6, 31)),
	path.NewLineToOp(dmath.NewVec2(16, 0)),
	path.NewLineToOp(dmath.NewVec2(25, 31)),
	path.NewMoveToOp(dmath.NewVec2(7.5, 25)),
	path.NewLineToOp(dmath.NewVec2(23.5, 25)),
}

var PlayerThrusterSpritePoints = []path.Op{
	path.NewMoveToOp(dmath.NewVec2(9, 25)),
	path.NewLineToOp(dmath.NewVec2(22, 25)),
	path.NewLineToOp(dmath.NewVec2(16.5, 35)),
	path.NewCloseOp(),
}

func NewPlayer(world donburi.World, startPosition dmath.Vec2) *donburi.Entry {
	entry := world.Entry(
		world.Create(
			component.Player,
			transform.Transform,
			component.Sprite,
			component.Input,
			component.Collider,
			component.ID,
			component.Velocity,
			component.CollisionQueue,
			component.AudioQueue,
		),
	)

	// INITIALIZE COMPONENTS

	component.Player.SetValue(entry, component.PlayerData{
		ThrusterTimer:  engine.NewTimer(1000 * time.Millisecond),
		ExplosionTimer: engine.NewTimer(2000 * time.Millisecond),
	})

	component.ID.SetValue(entry, component.IDData{
		ID:   engine.NewID(),
		Type: "Player",
	})

	transform.Transform.Get(entry).LocalPosition = startPosition
	transform.Transform.Get(entry).LocalRotation = -90

	component.Sprite.SetValue(entry, component.SpriteData{
		PathOps: PlayerSpritePoints,
		Layer:   component.SpriteLayerEntity,
		Pivot:   component.SpritePivotCenter,
		Size:    dmath.NewVec2(32, 32),
	})

	component.Collider.SetValue(entry, component.ColliderData{
		Width:  19,
		Height: 32,
		Offset: dmath.NewVec2(6, 0),
		Layer:  component.CollisionLayerPlayer,
	})

	component.Input.SetValue(entry, component.InputData{
		ForwardKeys: []ebiten.Key{
			ebiten.KeyW,
			ebiten.KeyArrowUp,
		},
		RotateClockwiseKeys: []ebiten.Key{
			ebiten.KeyD,
			ebiten.KeyArrowRight,
		},
		RotateCounterClockwiseKeys: []ebiten.Key{
			ebiten.KeyA,
			ebiten.KeyArrowLeft,
		},
		ShootKeys: []ebiten.Key{
			ebiten.KeySpace,
			ebiten.KeyEnter,
		},
		ForwardSpeed:    3.5,
		RotateSpeed:     10,
		ProjectileSpeed: 5,
	})

	thruster := world.Entry(world.Create(
		transform.Transform,
		component.Sprite,
	))

	component.Sprite.SetValue(thruster, component.SpriteData{
		PathOps: PlayerThrusterSpritePoints,
		Layer:   component.SpriteLayerEntity,
		Pivot:   component.SpritePivotCenter,
		Size:    dmath.NewVec2(32, 32),
		Hidden:  true,
	})

	transform.AppendChild(entry, thruster, false)

	return entry
}

func MustFindPlayer(w donburi.World) *donburi.Entry {
	player, ok := donburi.NewQuery(filter.Contains(component.Player)).First(w)
	if !ok {
		panic(errNoPlayerFound)
	}

	return player
}

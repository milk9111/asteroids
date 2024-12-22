package system

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/milk9111/asteroids/archetype"
	"github.com/milk9111/asteroids/assets"
	"github.com/milk9111/asteroids/component"
	"github.com/milk9111/asteroids/engine"
	"github.com/yohamta/donburi"
	dmath "github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
)

var (
	errNoInputEntryFound      = newSystemError("no input entry found")
	errPlayerChildrenNotFound = newSystemError("player children not found")
)

const (
	timePerTick = 1.0 / 60.0
)

type Input struct {
	query        *donburi.Query
	count        int
	prevVelocity dmath.Vec2
}

func NewInput() *Input {
	return &Input{
		query: donburi.NewQuery(
			filter.Contains(
				transform.Transform,
				component.Input,
				component.Player,
				component.Sprite,
				component.AudioQueue,
			),
		),
	}
}

// Update will check for various Player inputs in order to change the state of the game.
// It handles rotating, moving forward, and shooting projectiles based on different key inputs.
// The behavior of those inputs could be converted into a script triggered by those inputs as events.
func (i *Input) Update(world donburi.World) {
	inputEntry, ok := i.query.First(world)
	if !ok {
		panic(errNoInputEntryFound)
	}

	input := component.Input.Get(inputEntry)
	if input.Disabled {
		i.prevVelocity = engine.Vec2Zero()
		i.toggleThruster(inputEntry, false)
		return
	}

	t := transform.Transform.Get(inputEntry)
	v := component.Velocity.Get(inputEntry)

	isClockwise := false
	for _, key := range input.RotateClockwiseKeys {
		if ebiten.IsKeyPressed(key) {
			isClockwise = true
			break
		}
	}

	if isClockwise {
		i.count++
	}

	isCounterClockwise := false
	for _, key := range input.RotateCounterClockwiseKeys {
		if ebiten.IsKeyPressed(key) {
			isCounterClockwise = true
			break
		}
	}

	if isCounterClockwise {
		i.count--
	}

	t.LocalRotation = float64(i.count) * input.RotateSpeed * math.Pi / 360.0

	isMovingForward := false
	for _, key := range input.ForwardKeys {
		if ebiten.IsKeyPressed(key) {
			isMovingForward = true
			break
		}
	}

	if isMovingForward {
		v.AddVelocity(engine.RadToDirection(transform.WorldRotation(inputEntry) - engine.DegToRad(90)).MulScalar(input.ForwardSpeed))
		i.prevVelocity = v.Velocity()
		i.toggleThruster(inputEntry, true)
	} else if i.prevVelocity.Magnitude() > 0 {
		// slows down the velocity over time. would be better to put this into a separate 'physics' system so it can be applied
		// to other entities.
		v.AddVelocity(i.prevVelocity.MulScalar(1 - timePerTick))
		i.prevVelocity = v.Velocity()
		i.toggleThruster(inputEntry, false)
	}

	isShooting := false
	for _, key := range input.ShootKeys {
		if inpututil.IsKeyJustPressed(key) && !input.Disabled {
			isShooting = true
			break
		}
	}

	if isShooting {
		component.AudioQueue.Get(inputEntry).Enqueue(assets.Shoot_wav)
		archetype.NewSpawn(world, func() *donburi.Entry {
			dir := engine.RadToDirection(transform.WorldRotation(inputEntry) - engine.DegToRad(90))
			projectile := archetype.NewProjectile(world, transform.WorldPosition(inputEntry).Add(dir.MulScalar(16)))

			component.Traveler.SetValue(projectile, component.TravelerData{
				Direction: dir,
				Speed:     input.ProjectileSpeed,
			})

			return projectile
		})
	}
}

func (i *Input) toggleThruster(entry *donburi.Entry, toggle bool) {
	children, ok := transform.GetChildren(entry)
	if !ok || len(children) != 1 {
		panic(errPlayerChildrenNotFound)
	}

	thruster := children[0]
	component.Sprite.Get(thruster).Hidden = !toggle
}

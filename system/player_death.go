package system

import (
	"math/rand"
	"time"

	"github.com/milk9111/asteroids/archetype"
	"github.com/milk9111/asteroids/assets"
	"github.com/milk9111/asteroids/component"
	"github.com/milk9111/asteroids/engine"
	"github.com/milk9111/asteroids/engine/path"
	"github.com/yohamta/donburi"
	dmath "github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
)

var errNoPlayerFound = newSystemError("no player found")

type PlayerDeath struct {
	query *donburi.Query
	game  *component.GameData

	hasStartedExploding bool
}

func NewPlayerDeath() *PlayerDeath {
	return &PlayerDeath{
		query: donburi.NewQuery(
			filter.Contains(
				component.Sprite,
				component.Player,
				component.Input,
				component.Velocity,
				component.AudioQueue,
			),
		),
	}
}

func (p *PlayerDeath) Update(w donburi.World) {
	if p.game == nil {
		p.game = component.MustFindGame(w)
	}

	e, ok := p.query.First(w)
	if !ok {
		panic(errNoPlayerFound)
	}

	player := component.Player.Get(e)

	if !player.IsExploding || p.game.GameOver {
		return
	}

	collider := component.Collider.Get(e)
	input := component.Input.Get(e)
	sprite := component.Sprite.Get(e)

	if player.ExplosionTimer.IsReady() {
		if p.game.PlayerLives <= 0 {
			p.game.GameOver = true
			component.AudioQueue.Get(e).Enqueue(assets.GameOver_wav)
			return
		}

		p.hasStartedExploding = false
		player.IsExploding = false
		player.ExplosionTimer.Reset()

		collider.Disabled = false
		input.Disabled = false
		sprite.Hidden = false

		transform.Transform.Get(e).LocalPosition = dmath.Vec2{
			X: float64(p.game.Settings.WorldWidth) / 2,
			Y: float64(p.game.Settings.WorldHeight) / 2,
		}

		return
	}

	collider.Disabled = true
	input.Disabled = true
	sprite.Hidden = true

	if !p.hasStartedExploding {
		p.hasStartedExploding = true
		p.game.PlayerLives--

		playerPos := transform.WorldPosition(e)

		prev := archetype.PlayerSpritePoints[0]
		for i := 1; i < len(archetype.PlayerSpritePoints); i++ {
			curr := archetype.PlayerSpritePoints[i]

			createParticle(w, playerPos, prev.Target(), curr.Target(), 1000*time.Millisecond)

			prev = curr
		}
	}

	player.ExplosionTimer.Update()
}

func createParticle(w donburi.World, origin, a, b dmath.Vec2, lifetime time.Duration) {
	e := w.Entry(w.Create(
		transform.Transform,
		component.Sprite,
		component.Velocity,
		component.Traveler,
		component.TTL,
	))

	t := transform.Transform.Get(e)
	t.LocalPosition = origin

	sprite := component.Sprite.Get(e)
	sprite.PathOps = []path.Op{
		path.NewMoveToOp(a),
		path.NewLineToOp(b),
	}
	sprite.Layer = component.SpriteLayerEntity

	traveler := component.Traveler.Get(e)
	traveler.Direction = origin.Add(origin.Add(dmath.NewVec2(float64(-20+rand.Intn(41)), float64(-20+rand.Intn(41)))).MulScalar(-1)).Normalized()
	traveler.Speed = engine.RandomRangeF64(0.5, 1)

	component.TTL.Get(e).Timer = engine.NewTimer(lifetime)
}

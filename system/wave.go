package system

import (
	"math/rand"

	"github.com/milk9111/asteroids/archetype"
	"github.com/milk9111/asteroids/component"
	"github.com/milk9111/asteroids/engine"
	"github.com/yohamta/donburi"
	dmath "github.com/yohamta/donburi/features/math"
)

type Wave struct {
	game *component.GameData
	wave *component.WaveData
}

func NewWave() *Wave {
	return &Wave{}
}

// Update will trigger a new wave to spawn. Each wave it will spawn '4 + (wave number - 1)' large asteroids.
// The asteroids are spawned with a random direction, speed, and rotation on a random edge of the screen.
func (w *Wave) Update(world donburi.World) {
	if w.game == nil {
		w.game = component.MustFindGame(world)
	}

	if w.wave == nil {
		w.wave = component.MustFindWave(world)
		w.wave.WaveCooldownTimer.OverridePercentDone(1.0)
	}

	if w.wave.IsStarted {
		return
	} else if !w.wave.IsStarted && !w.wave.WaveCooldownTimer.IsReady() {
		w.wave.WaveCooldownTimer.Update()
		return
	}

	largeAsteroids := 4 + (w.wave.Number - 1)
	w.wave.IsStarted = true
	w.wave.WaveCooldownTimer.Reset()
	w.wave.RemainingHazards = largeAsteroids * 7 // largeAsteroids * (2^3 - 1)

	for i := 0; i < largeAsteroids; i++ {
		archetype.NewSpawn(world, func() *donburi.Entry {
			screenSide := rand.Int() % 4

			var x float64
			var y float64
			if screenSide == 0 { // left
				x = -64
				y = float64(rand.Int() % (w.game.Settings.WorldHeight + 1))
			} else if screenSide == 1 { // top
				x = float64(rand.Int() % (w.game.Settings.WorldWidth + 1))
				y = -64
			} else if screenSide == 2 { // right
				x = float64(w.game.Settings.WorldWidth+1) + 64
				y = float64(rand.Int() % (w.game.Settings.WorldHeight + 1))
			} else { // bottom
				x = float64(rand.Int() % (w.game.Settings.WorldWidth + 1))
				y = float64(w.game.Settings.WorldHeight+1) + 64
			}

			pos := dmath.NewVec2(x, y)

			// direction param is there to get a random direction within 20 degrees of the position.
			// it's the dmath.Vec2 equivalent of the engine.RandomRangeInt function.
			return archetype.NewAsteroid(
				world,
				pos,
				pos.Add(pos.Add(dmath.NewVec2(float64(-20+rand.Intn(41)), float64(-20+rand.Intn(41)))).MulScalar(-1)).Normalized(),
				engine.RandomRotation(),
				engine.RandomRangeF64(0.75, 1.5),
				component.AsteroidTierLarge,
			)
		})
	}
}

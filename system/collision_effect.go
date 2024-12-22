package system

import (
	"fmt"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/milk9111/asteroids/archetype"
	"github.com/milk9111/asteroids/assets"
	"github.com/milk9111/asteroids/component"
	"github.com/milk9111/asteroids/engine"
	"github.com/yohamta/donburi"
	dmath "github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
	"golang.org/x/image/colornames"
)

type CollisionEffect struct {
	query      *donburi.Query
	debugQuery *donburi.Query
	game       *component.GameData
	wave       *component.WaveData
	scoreQueue *component.ScoreQueueData

	activeCollisions map[engine.ID]mapset.Set[engine.ID]
}

func NewCollisionEffect() *CollisionEffect {
	return &CollisionEffect{
		query: donburi.NewQuery(
			filter.Contains(
				component.Collider,
				component.CollisionQueue,
				component.Sprite,
				transform.Transform,
			),
		),
		debugQuery: donburi.NewQuery(
			filter.Contains(
				component.Collider,
				component.Sprite,
				transform.Transform,
			),
		),
		activeCollisions: make(map[engine.ID]mapset.Set[engine.ID]),
	}
}

func (c *CollisionEffect) Update(w donburi.World) {
	if c.game == nil {
		c.game = component.MustFindGame(w)
	}

	if c.wave == nil {
		c.wave = component.MustFindWave(w)
	}

	if c.scoreQueue == nil {
		c.scoreQueue = component.MustFindScoreQueue(w)
	}

	c.query.Each(w, func(e *donburi.Entry) {
		collider := component.Collider.Get(e)
		switch collider.Layer {
		case component.CollisionLayerProjectile:
			c.handleProjectile(e)
		case component.CollisionLayerAsteroid:
			c.handleAsteroid(w, e)
		case component.CollisionLayerPlayer:
			c.handlePlayer(e)
		default:
			panic(newSystemError(fmt.Sprintf("unexpected collision layer %d", collider.Layer)))
		}
	})
}

func (c *CollisionEffect) DebugDraw(w donburi.World, screen *ebiten.Image) {
	c.debugQuery.Each(w, func(e *donburi.Entry) {
		color := colornames.Lime
		if component.Collider.Get(e).IsColliding {
			color = colornames.Red
		}

		rect := RectFromEntry(e)
		vertices := rect.GetVertices()
		for i := 0; i < len(vertices); i++ {
			next := (i + 1) % len(vertices)
			vector.StrokeLine(screen, float32(vertices[i].X), float32(vertices[i].Y), float32(vertices[next].X), float32(vertices[next].Y), strokeWidth, color, false)
		}
	})
}

func (c *CollisionEffect) handleProjectile(projectile *donburi.Entry) {
	queue := component.CollisionQueue.Get(projectile)
	if queue.Empty() {
		return
	}

	collider := component.Collider.Get(projectile)
	if collider.IsColliding {
		queue.Reset()
		return
	}

	for !queue.Empty() {
		_, layer := queue.Dequeue()

		if layer == component.CollisionLayerPlayer ||
			layer == component.CollisionLayerProjectile {
			continue
		}

		collider.IsColliding = true
		projectile.AddComponent(component.TagDestroy)
		break
	}

	queue.Reset()
}

func (c *CollisionEffect) handleAsteroid(w donburi.World, asteroid *donburi.Entry) {
	queue := component.CollisionQueue.Get(asteroid)
	if queue.Empty() {
		return
	}

	collider := component.Collider.Get(asteroid)
	if collider.IsColliding {
		queue.Reset()
		return
	}

	for !queue.Empty() {
		_, layer := queue.Dequeue()

		if layer == component.CollisionLayerAsteroid {
			continue
		}

		collider.IsColliding = true
		asteroid.AddComponent(component.TagDestroy)
		component.AudioQueue.Get(asteroid).Enqueue(assets.AsteroidDestroyed_wav)

		c.wave.RemainingHazards--
		if c.wave.IsStarted && c.wave.RemainingHazards <= 0 {
			c.wave.IsStarted = false
			c.wave.Number++
		}

		a := component.Asteroid.Get(asteroid)

		pos := transform.WorldPosition(asteroid)

		c.scoreQueue.Enqueue(a.PointValue)

		prev := archetype.ProjectileSpritePoints[0]
		for i := 1; i < len(archetype.ProjectileSpritePoints); i++ {
			curr := archetype.ProjectileSpritePoints[i]
			createParticle(w, pos, prev.Target(), curr.Target(), 1000*time.Millisecond)
		}

		if a.Tier == component.AsteroidTierSmall {
			break
		}

		t := transform.Transform.Get(asteroid)
		s := component.Sprite.Get(asteroid)
		tr := component.Traveler.Get(asteroid)

		rotationPivot := dmath.NewVec2(s.Size.X/2, s.Size.Y/2)
		dirPt := t.LocalPosition.Add(tr.Direction)
		speed := tr.Speed + 0.75

		spawnFunc := func() *donburi.Entry {
			return archetype.NewAsteroid(
				w,
				t.LocalPosition,
				engine.RotateAndScaleVec2(dirPt, rotationPivot, t.LocalScale, float64(int(engine.RandomRangeF64(0, 360))%360)).Normalized(),
				engine.RandomRotation(),
				speed,
				a.Tier-1,
			)
		}

		archetype.NewSpawn(w, spawnFunc)
		archetype.NewSpawn(w, spawnFunc)

		break
	}

	queue.Reset()
}

func (c *CollisionEffect) handlePlayer(player *donburi.Entry) {
	queue := component.CollisionQueue.Get(player)
	collider := component.Collider.Get(player)

	if queue.Empty() {
		collider.IsColliding = false
		return
	}

	if collider.IsColliding {
		queue.Reset()
		return
	}

	for !queue.Empty() {
		_, layer := queue.Dequeue()

		if layer == component.CollisionLayerProjectile ||
			layer == component.CollisionLayerPlayer {
			continue
		}

		collider.IsColliding = true
		component.Player.Get(player).IsExploding = true

		break
	}
}

func (c *CollisionEffect) AddToActiveCollisions(entry, other *donburi.Entry) {
	c.addActiveCollision(entry, other)
	c.addActiveCollision(other, entry)
}

func (c *CollisionEffect) RemoveFromActiveCollisions(entry, other *donburi.Entry) {
	c.removeActiveCollision(entry, other)
	c.removeActiveCollision(other, entry)
}

func (c *CollisionEffect) RemoveFromAllActiveCollisions(entry *donburi.Entry) {
	id := component.ID.GetValue(entry).ID
	delete(c.activeCollisions, id)

	for k := range c.activeCollisions {
		if c.activeCollisions[k].ContainsOne(id) {
			c.activeCollisions[k].Remove(id)
		}
	}
}

func (c *CollisionEffect) ContainsActiveCollision(entry, other *donburi.Entry) bool {
	s, ok := c.activeCollisions[component.ID.GetValue(entry).ID]
	if !ok {
		return false
	}

	return s.ContainsOne(component.ID.GetValue(other).ID)
}

func (c *CollisionEffect) ContainsAnyActiveCollision(entry *donburi.Entry) bool {
	s, ok := c.activeCollisions[component.ID.GetValue(entry).ID]
	if !ok {
		return false
	}

	return !s.IsEmpty()
}

func (c *CollisionEffect) addActiveCollision(entry, other *donburi.Entry) {
	if _, ok := c.activeCollisions[component.ID.GetValue(entry).ID]; !ok {
		c.activeCollisions[component.ID.GetValue(entry).ID] = mapset.NewSet[engine.ID]()
	}

	c.activeCollisions[component.ID.GetValue(entry).ID].Add(component.ID.GetValue(other).ID)
}

func (c *CollisionEffect) removeActiveCollision(entry, other *donburi.Entry) {
	entryId := component.ID.GetValue(entry).ID
	if _, ok := c.activeCollisions[entryId]; !ok {
		return
	}

	c.activeCollisions[entryId].Remove(component.ID.GetValue(other).ID)

	if c.activeCollisions[entryId].IsEmpty() {
		delete(c.activeCollisions, entryId)
	}
}

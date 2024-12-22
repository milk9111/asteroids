package system

import (
	"github.com/milk9111/asteroids/component"
	"github.com/milk9111/asteroids/engine"
	"github.com/yohamta/donburi"
	dmath "github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
)

type CollisionCheck struct {
	query *donburi.Query
	game  *component.GameData
}

func NewCollisionCheck() *CollisionCheck {
	return &CollisionCheck{
		query: donburi.NewQuery(
			filter.Contains(
				component.Collider,
				transform.Transform,
				component.Sprite,
				component.ID,
				component.CollisionQueue,
			),
		),
	}
}

func (c *CollisionCheck) Update(w donburi.World) {
	if c.game == nil {
		c.game = component.MustFindGame(w)
	}

	var entries []*donburi.Entry
	c.query.Each(w, func(entry *donburi.Entry) {
		if !entry.Valid() {
			return
		}

		entries = append(entries, entry)
	})

	for _, entry := range entries {
		if component.Collider.Get(entry).Disabled {
			continue
		}

		queue := component.CollisionQueue.Get(entry)

		for _, other := range entries {
			if entry.Entity().Id() == other.Entity().Id() {
				continue
			}

			otherCollider := component.Collider.Get(other)
			if otherCollider.Disabled {
				continue
			}

			rect := RectFromEntry(entry)
			otherRect := RectFromEntry(other)

			intersects := rect.Intersects(otherRect)
			if intersects {
				queue.Enqueue(other, otherCollider.Layer)
			}
		}
	}
}

func RectFromEntry(entry *donburi.Entry) engine.Rect {
	t := transform.Transform.Get(entry)
	collider := component.Collider.Get(entry)
	sprite := component.Sprite.Get(entry)
	rotationPivot := dmath.NewVec2(sprite.Size.X/2, sprite.Size.Y/2)
	colliderWidth := collider.Width * t.LocalScale.X
	colliderHeight := collider.Height * t.LocalScale.Y

	pos := sprite.PivotPoint(transform.WorldPosition(entry))

	return engine.NewRect(rotationPivot.Add(pos), colliderWidth, colliderHeight, t.LocalRotation)
}

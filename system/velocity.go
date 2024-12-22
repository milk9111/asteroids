package system

import (
	"github.com/milk9111/asteroids/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
)

type Velocity struct {
	query *donburi.Query
}

func NewVelocity() *Velocity {
	return &Velocity{
		query: donburi.NewQuery(
			filter.Contains(
				transform.Transform,
				component.Velocity,
			),
		),
	}
}

// Update applies the queued velocity and rotation (angular) velocity to the transform.
func (v *Velocity) Update(world donburi.World) {
	v.query.Each(world, func(entry *donburi.Entry) {
		t := transform.Transform.Get(entry)
		velocity := component.Velocity.Get(entry)

		queuedVelocity := velocity.Velocity()
		queuedRotationVelocity := velocity.RotationVelocity()

		velocity.Reset()

		t.LocalPosition = t.LocalPosition.Add(queuedVelocity)
		t.LocalRotation += queuedRotationVelocity
	})
}

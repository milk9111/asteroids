package system

import (
	"github.com/milk9111/asteroids/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
)

type Travel struct {
	query *donburi.Query
}

func NewTravel() *Travel {
	return &Travel{
		query: donburi.NewQuery(
			filter.Contains(
				transform.Transform,
				component.Traveler,
				component.Velocity,
			),
		),
	}
}

// Update moves queues up the velocity moving in a straight line using its direction and speed.
// It also adds the rotational (angular) velocity, but I wasn't able to get that to work.
func (t *Travel) Update(w donburi.World) {
	t.query.Each(w, func(e *donburi.Entry) {
		t := transform.Transform.Get(e)
		traveler := component.Traveler.Get(e)
		velocity := component.Velocity.Get(e)

		velocity.AddVelocity(traveler.Direction.MulScalar(traveler.Speed))
		if traveler.RotationalSpeed != 0 {
			velocity.AddRotationalVelocity(t.LocalRotation / (traveler.RotationalSpeed * timePerTick))
		}
	})
}

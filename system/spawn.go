package system

import (
	"github.com/milk9111/asteroids/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

type Spawn struct {
	query *donburi.Query
}

func NewSpawn() *Spawn {
	return &Spawn{
		query: donburi.NewQuery(
			filter.Contains(component.Spawn),
		),
	}
}

// Update will use the component.Spawn component's callback field to instantiate a new entity before
// removing itself from the world.
func (s *Spawn) Update(w donburi.World) {
	s.query.Each(w, func(e *donburi.Entry) {
		spawn := component.Spawn.Get(e)
		spawn.Instantiate()
		e.Remove()
	})
}

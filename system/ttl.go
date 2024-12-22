package system

import (
	"github.com/milk9111/asteroids/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

type TTL struct {
	query *donburi.Query
}

func NewTTL() *TTL {
	return &TTL{
		query: donburi.NewQuery(
			filter.And(
				filter.Contains(component.TTL),
				filter.Not(
					filter.Contains(component.TagDestroy),
				),
			),
		),
	}
}

// Update adds the component.TagDestroy component to any entity that has reached the end of its TTL timer.
func (t *TTL) Update(w donburi.World) {
	t.query.Each(w, func(e *donburi.Entry) {
		ttl := component.TTL.Get(e)

		if ttl.Timer.IsReady() {
			e.AddComponent(component.TagDestroy)
			return
		}

		ttl.Timer.Update()
	})
}

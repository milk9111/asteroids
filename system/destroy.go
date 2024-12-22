package system

import (
	"github.com/milk9111/asteroids/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

type Destroy struct {
	query *donburi.Query
}

func NewDestroy() *Destroy {
	return &Destroy{
		query: donburi.NewQuery(filter.Contains(component.TagDestroy)),
	}
}

func (d *Destroy) Update(w donburi.World) {
	d.query.Each(w, func(e *donburi.Entry) {
		w.Remove(e.Entity())
	})
}

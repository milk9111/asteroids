package archetype

import (
	"github.com/milk9111/asteroids/component"
	"github.com/yohamta/donburi"
)

func NewSpawn(w donburi.World, instantiate func() *donburi.Entry) *donburi.Entry {
	entry := w.Entry(w.Create(component.Spawn))
	component.Spawn.SetValue(entry, component.SpawnData{
		Instantiate: instantiate,
	})

	return entry
}

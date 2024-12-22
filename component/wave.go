package component

import (
	"github.com/milk9111/asteroids/engine"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

var errNoWaveFound = newComponentError("no wave found")

type WaveData struct {
	Number            int
	IsStarted         bool
	WaveCooldownTimer *engine.Timer
	RemainingHazards  int
}

var Wave = donburi.NewComponentType[WaveData]()

func MustFindWave(w donburi.World) *WaveData {
	wave, ok := donburi.NewQuery(filter.Contains(Wave)).First(w)
	if !ok {
		panic(errNoWaveFound)
	}
	return Wave.Get(wave)
}

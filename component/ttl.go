package component

import (
	"github.com/milk9111/asteroids/engine"
	"github.com/yohamta/donburi"
)

type TTLData struct {
	Timer *engine.Timer
}

var TTL = donburi.NewComponentType[TTLData]()

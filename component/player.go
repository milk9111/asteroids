package component

import (
	"github.com/milk9111/asteroids/engine"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

var errNoPlayerFound = newComponentError("no player found")

type PlayerData struct {
	ThrusterTimer *engine.Timer

	IsExploding    bool
	ExplosionTimer *engine.Timer
}

var Player = donburi.NewComponentType[PlayerData]()

func MustFindPlayer(w donburi.World) *PlayerData {
	player, ok := donburi.NewQuery(filter.Contains(Player)).First(w)
	if !ok {
		panic(errNoPlayerFound)
	}
	return Player.Get(player)
}

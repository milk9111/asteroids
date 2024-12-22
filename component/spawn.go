package component

import (
	"github.com/yohamta/donburi"
)

type SpawnData struct {
	Instantiate func() *donburi.Entry
}

var Spawn = donburi.NewComponentType[SpawnData]()

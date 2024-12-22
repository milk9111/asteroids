package component

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
)

type TravelerData struct {
	Direction       math.Vec2
	Speed           float64
	RotationalSpeed float64
}

var Traveler = donburi.NewComponentType[TravelerData]()

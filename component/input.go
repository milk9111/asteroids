package component

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type InputData struct {
	ForwardKeys                []ebiten.Key
	RotateClockwiseKeys        []ebiten.Key
	RotateCounterClockwiseKeys []ebiten.Key
	ShootKeys                  []ebiten.Key
	ForwardSpeed               float64
	RotateSpeed                float64
	ProjectileSpeed            float64
	Disabled                   bool
}

var Input = donburi.NewComponentType[InputData]()

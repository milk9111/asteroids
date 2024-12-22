package component

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
)

const (
	CollisionLayerUnknown CollisionLayer = iota
	CollisionLayerPlayer
	CollisionLayerProjectile
	CollisionLayerAsteroid
	CollisionLayerUI
)

type CollisionLayer int

type ColliderData struct {
	Width       float64
	Height      float64
	Offset      math.Vec2
	Layer       CollisionLayer
	Disabled    bool
	IsColliding bool
}

var Collider = donburi.NewComponentType[ColliderData]()

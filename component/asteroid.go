package component

import "github.com/yohamta/donburi"

type AsteroidTier int

const (
	AsteroidTierSmall AsteroidTier = iota
	AsteroidTierMedium
	AsteroidTierLarge
)

type AsteroidData struct {
	Tier       AsteroidTier
	PointValue int
}

var Asteroid = donburi.NewComponentType[AsteroidData]()

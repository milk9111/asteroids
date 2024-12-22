package component

import "github.com/yohamta/donburi"

type TagProjectileSpawnData struct{}

var TagProjectileSpawn = donburi.NewComponentType[TagProjectileSpawnData]()

type TagDestroyData struct{}

var TagDestroy = donburi.NewComponentType[TagDestroyData]()

type TagHUDData struct{}

var TagHUD = donburi.NewComponentType[TagHUDData]()

type TagPlayerDeathLineSegmentData struct{}

var TagPlayerDeathLineSegment = donburi.NewComponentType[TagPlayerDeathLineSegmentData]()

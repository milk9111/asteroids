package component

import "github.com/yohamta/donburi"

type WaveLineItemData struct {
	Value int
}

var WaveLineItem = donburi.NewComponentType[WaveLineItemData]()

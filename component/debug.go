package component

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

var errNoDebugFound = newComponentError("no debug found")

type DebugData struct {
	Enabled bool
}

var Debug = donburi.NewComponentType[DebugData]()

func MustFindDebug(w donburi.World) *DebugData {
	debug, ok := donburi.NewQuery(filter.Contains(Debug)).First(w)
	if !ok {
		panic(errNoDebugFound)
	}
	return Debug.Get(debug)
}

package component

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

var errNoGameFound = newComponentError("no game found")

type GameData struct {
	Paused      bool
	GameOver    bool
	Score       int
	PlayerLives int
	Debugging   bool
	Settings    Settings
}

type Settings struct {
	WorldWidth        int
	WorldHeight       int
	ScreenWidth       int
	ScreenHeight      int
	TileSize          int
	WorldScaleFactor  float64
	ScreenScaleFactor float64
}

var Game = donburi.NewComponentType[GameData]()

func MustFindGame(w donburi.World) *GameData {
	game, ok := donburi.NewQuery(filter.Contains(Game)).First(w)
	if !ok {
		panic(errNoGameFound)
	}
	return Game.Get(game)
}

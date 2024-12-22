package system

import (
	"math"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/milk9111/asteroids/archetype"
	"github.com/milk9111/asteroids/component"
	"github.com/milk9111/asteroids/engine/path"
	"github.com/yohamta/donburi"
	dmath "github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
)

var errNoHUDFound = newSystemError("no HUD found")

type HUD struct {
	query               *donburi.Query
	game                *component.GameData
	gameOver            *component.SpriteData
	restartGameCallback func()
}

func NewHUD(restartGameCallback func()) *HUD {
	return &HUD{
		query: donburi.NewQuery(
			filter.Contains(
				component.TagHUD,
				component.Sprite,
			),
		),
		restartGameCallback: restartGameCallback,
	}
}

// Update will convert the integer game score into a list of path operations for the Render system to
// draw to the screen. This could be converted to a script that is triggered by a 'scored' or 'died' event.
func (h *HUD) Update(w donburi.World) {
	if h.game == nil {
		h.game = component.MustFindGame(w)
	}

	e, ok := h.query.First(w)
	if !ok {
		panic(errNoHUDFound)
	}

	if h.gameOver == nil {
		children, ok := transform.GetChildren(e)
		if !ok || len(children) == 0 {
			panic(errNoHUDFound)
		}

		h.gameOver = component.Sprite.Get(children[0])
	}

	if h.game.GameOver {
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			h.restartGameCallback()
			return
		}

		h.gameOver.Hidden = false
		return
	}

	var scoreTextPath []path.Op
	var scoreSizeX float64
	// if I ever do this vector.Path drawing again, I should move this loop into a dedicated 'text' or 'font' package
	// since I've had to duplicate this exact loop a few times.
	for i, r := range strconv.Itoa(h.game.Score) {
		nextRunePoints, ok := archetype.FontSpritePoints[r]
		if !ok {
			continue
		}

		for _, op := range nextRunePoints {
			scoreTextPath = append(scoreTextPath, op.Copy(op.Target().Add(dmath.NewVec2(scoreSizeX+float64(8*i), 0))))
		}

		scoreSizeX += calculateFontRuneSize(nextRunePoints).X
	}

	var playerLivesSizeX float64
	for i := 0; i < h.game.PlayerLives; i++ {
		for _, op := range archetype.PlayerSpritePoints {
			scoreTextPath = append(scoreTextPath, op.Copy(op.Target().Add(dmath.NewVec2(playerLivesSizeX+float64(4*i), 48))))
		}

		playerLivesSizeX += 32
	}

	component.Sprite.Get(e).PathOps = scoreTextPath
}

// calculateFontRuneSize determines the X, Y size for a given path operations slice. Characters like ' don't have the
// normal width so using a constant width and height would end up putting extra space in the middle of a word where
// it shouldn't be.
func calculateFontRuneSize(fontRunePath []path.Op) dmath.Vec2 {
	target := fontRunePath[0].Target()
	minX, maxX, minY, maxY := target.X, target.X, target.Y, target.Y
	for _, op := range fontRunePath {
		target = op.Target()
		minX = math.Min(minX, target.X)
		maxX = math.Max(maxX, target.X)
		minY = math.Min(minY, target.Y)
		maxY = math.Max(maxY, target.Y)
	}

	return dmath.NewVec2(maxX-minX, maxY-minY)
}

package game

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/milk9111/asteroids/scene"
)

type Scene interface {
	Update()
	Draw(screen *ebiten.Image)
}

type Game struct {
	scene        Scene
	worldWidth   int
	worldHeight  int
	screenWidth  int
	screenHeight int

	scaleFactor float64
}

type Config struct {
	Quick        bool
	ScreenWidth  int
	ScreenHeight int
}

func NewGame(config Config) *Game {
	// assets.MustLoadAssets()

	return &Game{
		// scene:  scene.NewEntrance(config.ScreenWidth, config.ScreenHeight),
		worldWidth:   config.ScreenWidth,
		worldHeight:  config.ScreenHeight,
		screenWidth:  config.ScreenWidth,
		screenHeight: config.ScreenHeight,
	}
}

func (g *Game) Update() error {
	if g.scene == nil {
		g.scene = scene.NewEntrance(g.worldWidth, g.worldHeight, g.screenWidth, g.screenHeight, g.scaleFactor)
	}

	g.scene.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.scene.Draw(screen)
}

func (g *Game) LayoutF(logicWinWidth, logicWinHeight float64) (float64, float64) {
	scale := ebiten.Monitor().DeviceScaleFactor()
	screenWidth := math.Ceil(logicWinWidth * scale)
	screenHeight := math.Ceil(logicWinHeight * scale)
	g.scaleFactor = (screenWidth * screenHeight) / (float64(g.worldWidth) * float64(g.worldHeight))
	g.screenWidth = int(screenWidth)
	g.screenHeight = int(screenHeight)
	// g.width = int(canvasWidth)
	// g.height = int(canvasHeight)
	return float64(g.worldWidth), float64(g.worldHeight)
}

func (g *Game) Layout(_, _ int) (int, int) {
	panic("game: shouldn't have called Layout")
}

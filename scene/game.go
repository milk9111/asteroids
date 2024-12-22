package scene

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/milk9111/asteroids/archetype"
	"github.com/milk9111/asteroids/component"
	"github.com/milk9111/asteroids/engine"
	"github.com/milk9111/asteroids/system"
	"github.com/yohamta/donburi"
	dmath "github.com/yohamta/donburi/features/math"
)

type System interface {
	Update(w donburi.World)
}

type Drawable interface {
	Draw(w donburi.World, screen *ebiten.Image)
}

type Debuggable interface {
	DebugDraw(w donburi.World, screen *ebiten.Image)
}

type Entrance struct {
	worldWidth   int
	worldHeight  int
	screenWidth  int
	screenHeight int
	scaleFactor  float64

	world       donburi.World
	game        *component.GameData
	systems     []System
	drawables   []Drawable
	debuggables []Debuggable

	count int

	isStepping bool
	step       int

	isRestarting bool
}

func NewEntrance(worldWidth, worldHeight, screenWidth, screenHeight int, scaleFactor float64) *Entrance {
	e := &Entrance{
		worldWidth:   worldWidth,
		worldHeight:  worldHeight,
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
		scaleFactor:  scaleFactor,
	}

	e.loadLevel()

	return e
}

func (e *Entrance) restartGame() {
	e.isRestarting = true
}

func (e *Entrance) pauseGame() {
	e.game.Paused = !e.game.Paused
}

func (e *Entrance) nextStep() {
	if !e.game.Paused {
		return
	}

	e.isStepping = true
	e.step = e.count + 1
}

func (e *Entrance) loadLevel() {
	render := system.NewRender(e.worldWidth, e.worldHeight)
	debug := system.NewDebug(e.worldWidth, e.worldHeight, e.nextStep, e.pauseGame)
	collisionEffect := system.NewCollisionEffect()

	e.systems = []System{
		system.NewVelocity(),
		system.NewInput(),
		system.NewTravel(),
		system.NewWave(),
		system.NewSpawn(),
		system.NewCameraBounds(),
		system.NewTTL(),
		system.NewCollisionCheck(),
		collisionEffect,
		system.NewScore(),
		system.NewPlayerDeath(),
		system.NewHUD(e.restartGame),
		system.NewAudio(),
		render,
		debug,
		system.NewDestroy(),
	}

	e.drawables = []Drawable{
		render,
	}

	e.debuggables = []Debuggable{
		debug,
		collisionEffect,
	}

	e.world = e.createWorld()
}

func (e *Entrance) createWorld() donburi.World {
	world := donburi.NewWorld()

	archetype.NewCamera(world, dmath.Vec2{
		X: float64(e.worldWidth),
		Y: float64(e.worldHeight),
	})

	game := world.Entry(world.Create(component.Game))
	component.Game.SetValue(game, component.GameData{
		Debugging:   false,
		Score:       0,
		PlayerLives: 3,
		Settings: component.Settings{
			WorldWidth:        e.worldWidth,
			WorldHeight:       e.worldHeight,
			ScreenWidth:       e.screenWidth,
			ScreenHeight:      e.screenHeight,
			TileSize:          32,
			ScreenScaleFactor: e.scaleFactor,
		},
	})

	e.game = component.Game.Get(game)

	debug := world.Entry(world.Create(component.Debug))
	component.Debug.SetValue(debug, component.DebugData{
		Enabled: true,
	})

	archetype.NewPlayer(world, dmath.Vec2{
		X: float64(e.worldWidth) / 2,
		Y: float64(e.worldHeight) / 2,
	})

	wave := world.Entry(world.Create(component.Wave))
	component.Wave.SetValue(wave, component.WaveData{
		Number:            1,
		WaveCooldownTimer: engine.NewTimer(2 * time.Second),
	})

	archetype.NewHUD(world, e.game)

	world.Create(component.ScoreQueue)

	return world
}

func (e *Entrance) Update() {
	if e.isRestarting {
		e.isRestarting = false
		e.loadLevel()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySlash) {
		e.game.Debugging = !e.game.Debugging
	}

	for _, s := range e.systems {
		if e.game.Paused && !e.isStepping {
			if _, ok := s.(*system.Debug); !ok {
				continue
			}
		}

		if _, ok := s.(*system.Debug); !ok && e.isStepping {
			fmt.Printf("")
		}

		s.Update(e.world)
	}

	if e.isStepping && e.step == e.count {
		e.isStepping = false
	}

	e.count++
}

func (e *Entrance) Draw(screen *ebiten.Image) {
	screen.Clear()
	for _, d := range e.drawables {
		d.Draw(e.world, screen)
	}

	if e.game.Debugging {
		for _, d := range e.debuggables {
			d.DebugDraw(e.world, screen)
		}
	}
}

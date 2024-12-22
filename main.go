package main

import (
	"bytes"
	"flag"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/milk9111/asteroids/assets"
	"github.com/milk9111/asteroids/game"
)

func init() {
	audio.NewContext(44100)
}

func main() {
	forceMonitorFlag := flag.Bool("q", false, "force game window to third monitor (for testing)")
	flag.Parse()

	// couldn't get the screen size to dynamically resize so decided to hardcode it here
	config := game.Config{
		ScreenWidth:  1024,
		ScreenHeight: 768,
	}

	// setup has 3 monitors and this forces the game window to open up on the 3rd one
	monitors := ebiten.AppendMonitors(nil)
	if forceMonitorFlag != nil && *forceMonitorFlag && len(monitors) == 3 {
		ebiten.SetMonitor(monitors[2])
	}

	icon16x16, _, err := image.Decode(bytes.NewReader(assets.Icon16x16_png))
	if err != nil {
		panic(err)
	}

	icon32x32, _, err := image.Decode(bytes.NewReader(assets.Icon32x32_png))
	if err != nil {
		panic(err)
	}

	icon48x48, _, err := image.Decode(bytes.NewReader(assets.Icon48x48_png))
	if err != nil {
		panic(err)
	}

	ebiten.SetWindowTitle("Asteroids")
	ebiten.SetWindowIcon([]image.Image{
		icon16x16,
		icon32x32,
		icon48x48,
	})
	ebiten.SetWindowSize(config.ScreenWidth, config.ScreenHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.MaximizeWindow()

	err = ebiten.RunGame(game.NewGame(config))
	if err != nil {
		log.Fatal(err)
	}
}

package assets

import (
	_ "embed"
)

var (
	//go:embed shoot.wav
	Shoot_wav []byte

	//go:embed asteroid_destroyed.wav
	AsteroidDestroyed_wav []byte

	//go:embed game_over.wav
	GameOver_wav []byte

	//go:embed icon_16x16.png
	Icon16x16_png []byte

	//go:embed icon_32x32.png
	Icon32x32_png []byte

	//go:embed icon_48x48.png
	Icon48x48_png []byte
)

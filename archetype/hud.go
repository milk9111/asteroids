package archetype

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/milk9111/asteroids/component"
	"github.com/milk9111/asteroids/engine/path"
	"github.com/yohamta/donburi"
	dmath "github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"golang.org/x/image/colornames"
)

var (
	FontSpritePoints = map[rune][]path.Op{
		'0': {
			path.NewMoveToOp(dmath.NewVec2(0, 0)),
			path.NewLineToOp(dmath.NewVec2(24, 0)),
			path.NewLineToOp(dmath.NewVec2(24, 31)),
			path.NewLineToOp(dmath.NewVec2(0, 31)),
			path.NewCloseOp(),
		},
		'1': {
			path.NewMoveToOp(dmath.NewVec2(0, 3)),
			path.NewLineToOp(dmath.NewVec2(4, 0)),
			// vector.Path vertices calculation is making the bottom of '1' slightly shorter than the other numbers so making it a bit longer to account for that.
			path.NewLineToOp(dmath.NewVec2(4, 31.5)),
		},
		'2': {
			path.NewMoveToOp(dmath.NewVec2(0, 0)),
			path.NewLineToOp(dmath.NewVec2(24, 0)),
			path.NewLineToOp(dmath.NewVec2(24, 15.5)),
			path.NewLineToOp(dmath.NewVec2(0, 15.5)),
			path.NewLineToOp(dmath.NewVec2(0, 31)),
			path.NewLineToOp(dmath.NewVec2(24, 31)),
		},
		'3': {
			path.NewMoveToOp(dmath.NewVec2(0, 0)),
			path.NewLineToOp(dmath.NewVec2(24, 0)),
			path.NewLineToOp(dmath.NewVec2(24, 31)),
			path.NewLineToOp(dmath.NewVec2(0, 31)),
			path.NewMoveToOp(dmath.NewVec2(24, 16)),
			path.NewLineToOp(dmath.NewVec2(0, 16)),
		},
		'4': {
			path.NewMoveToOp(dmath.NewVec2(0, 0)),
			path.NewLineToOp(dmath.NewVec2(0, 16)),
			path.NewLineToOp(dmath.NewVec2(24, 16)),
			path.NewMoveToOp(dmath.NewVec2(24, 0)),
			path.NewLineToOp(dmath.NewVec2(24, 31)),
		},
		'5': {
			path.NewMoveToOp(dmath.NewVec2(24, 0)),
			path.NewLineToOp(dmath.NewVec2(0, 0)),
			path.NewLineToOp(dmath.NewVec2(0, 16)),
			path.NewLineToOp(dmath.NewVec2(24, 16)),
			path.NewLineToOp(dmath.NewVec2(24, 31)),
			path.NewLineToOp(dmath.NewVec2(0, 31)),
		},
		'6': {
			path.NewMoveToOp(dmath.NewVec2(0, 0)),
			path.NewLineToOp(dmath.NewVec2(0, 31)),
			path.NewLineToOp(dmath.NewVec2(24, 31)),
			path.NewLineToOp(dmath.NewVec2(24, 16)),
			path.NewLineToOp(dmath.NewVec2(0, 16)),
		},
		'7': {
			path.NewMoveToOp(dmath.NewVec2(0, 0)),
			path.NewLineToOp(dmath.NewVec2(24, 0)),
			path.NewLineToOp(dmath.NewVec2(24, 31)),
		},
		'8': {
			path.NewMoveToOp(dmath.NewVec2(0, 0)),
			path.NewLineToOp(dmath.NewVec2(24, 0)),
			path.NewLineToOp(dmath.NewVec2(24, 31)),
			path.NewLineToOp(dmath.NewVec2(0, 31)),
			path.NewLineToOp(dmath.NewVec2(0, 0)),
			path.NewMoveToOp(dmath.NewVec2(0, 16)),
			path.NewLineToOp(dmath.NewVec2(24, 16)),
		},
		'9': {
			path.NewMoveToOp(dmath.NewVec2(24, 16)),
			path.NewLineToOp(dmath.NewVec2(0, 16)),
			path.NewLineToOp(dmath.NewVec2(0, 0)),
			path.NewLineToOp(dmath.NewVec2(24, 0)),
			path.NewLineToOp(dmath.NewVec2(24, 31)),
		},
		'A': {
			path.NewMoveToOp(dmath.NewVec2(0, 31)),
			path.NewLineToOp(dmath.NewVec2(0, 8)),
			path.NewLineToOp(dmath.NewVec2(12, 0)),
			path.NewLineToOp(dmath.NewVec2(24, 10)),
			path.NewLineToOp(dmath.NewVec2(24, 31)),
			path.NewMoveToOp(dmath.NewVec2(0, 16)),
			path.NewLineToOp(dmath.NewVec2(24, 16)),
		},
		'D': {
			path.NewMoveToOp(dmath.NewVec2(0, 0)),
			path.NewLineToOp(dmath.NewVec2(12, 0)),
			path.NewLineToOp(dmath.NewVec2(24, 12)),
			path.NewLineToOp(dmath.NewVec2(24, 20)),
			path.NewLineToOp(dmath.NewVec2(12, 31)),
			path.NewLineToOp(dmath.NewVec2(0, 31)),
			path.NewCloseOp(),
		},
		'E': {
			path.NewMoveToOp(dmath.NewVec2(24, 0)),
			path.NewLineToOp(dmath.NewVec2(0, 0)),
			path.NewLineToOp(dmath.NewVec2(0, 31)),
			path.NewLineToOp(dmath.NewVec2(24, 31)),
			path.NewMoveToOp(dmath.NewVec2(0, 16)),
			path.NewLineToOp(dmath.NewVec2(20, 16)),
		},
		'G': {
			path.NewMoveToOp(dmath.NewVec2(12, 16)),
			path.NewLineToOp(dmath.NewVec2(24, 16)),
			path.NewLineToOp(dmath.NewVec2(24, 31)),
			path.NewLineToOp(dmath.NewVec2(0, 31)),
			path.NewLineToOp(dmath.NewVec2(0, 0)),
			path.NewLineToOp(dmath.NewVec2(24, 0)),
			path.NewLineToOp(dmath.NewVec2(24, 10)),
		},
		'I': {
			path.NewMoveToOp(dmath.NewVec2(0, 0)),
			path.NewLineToOp(dmath.NewVec2(24, 0)),
			path.NewMoveToOp(dmath.NewVec2(12, 0)),
			path.NewLineToOp(dmath.NewVec2(12, 31)),
			path.NewMoveToOp(dmath.NewVec2(0, 31)),
			path.NewLineToOp(dmath.NewVec2(24, 31)),
		},
		'M': {
			path.NewMoveToOp(dmath.NewVec2(0, 31)),
			path.NewLineToOp(dmath.NewVec2(0, 0)),
			path.NewLineToOp(dmath.NewVec2(12, 16)),
			path.NewLineToOp(dmath.NewVec2(24, 0)),
			path.NewLineToOp(dmath.NewVec2(24, 31)),
		},
		'O': {
			path.NewMoveToOp(dmath.NewVec2(0, 0)),
			path.NewLineToOp(dmath.NewVec2(24, 0)),
			path.NewLineToOp(dmath.NewVec2(24, 31)),
			path.NewLineToOp(dmath.NewVec2(0, 31)),
			path.NewCloseOp(),
		},
		'V': {
			path.NewMoveToOp(dmath.NewVec2(0, 0)),
			path.NewLineToOp(dmath.NewVec2(0, 16)),
			path.NewLineToOp(dmath.NewVec2(12, 31)),
			path.NewLineToOp(dmath.NewVec2(24, 16)),
			path.NewLineToOp(dmath.NewVec2(24, 0)),
		},
		'R': {
			path.NewMoveToOp(dmath.NewVec2(0, 31)),
			path.NewLineToOp(dmath.NewVec2(0, 0)),
			path.NewLineToOp(dmath.NewVec2(24, 0)),
			path.NewLineToOp(dmath.NewVec2(24, 16)),
			path.NewLineToOp(dmath.NewVec2(0, 16)),
			path.NewLineToOp(dmath.NewVec2(24, 31)),
		},
		'P': {
			path.NewMoveToOp(dmath.NewVec2(0, 31)),
			path.NewLineToOp(dmath.NewVec2(0, 0)),
			path.NewLineToOp(dmath.NewVec2(24, 0)),
			path.NewLineToOp(dmath.NewVec2(24, 16)),
			path.NewLineToOp(dmath.NewVec2(0, 20)),
		},
		'S': {
			path.NewMoveToOp(dmath.NewVec2(24, 0)),
			path.NewLineToOp(dmath.NewVec2(0, 0)),
			path.NewLineToOp(dmath.NewVec2(0, 16)),
			path.NewLineToOp(dmath.NewVec2(24, 16)),
			path.NewLineToOp(dmath.NewVec2(24, 31)),
			path.NewLineToOp(dmath.NewVec2(0, 31)),
		},
		'T': {
			path.NewMoveToOp(dmath.NewVec2(0, 0)),
			path.NewLineToOp(dmath.NewVec2(24, 0)),
			path.NewMoveToOp(dmath.NewVec2(12, 0)),
			path.NewLineToOp(dmath.NewVec2(12, 31)),
		},
		'\'': {
			path.NewMoveToOp(dmath.NewVec2(0, 0)),
			path.NewLineToOp(dmath.NewVec2(0, 8)),
		},
		' ': {
			path.NewMoveToOp(dmath.NewVec2(0, 0)),
			path.NewMoveToOp(dmath.NewVec2(12, 0)),
		},
	}
)

func NewHUD(w donburi.World, game *component.GameData) *donburi.Entry {
	hud := w.Entry(w.Create(
		component.TagHUD,
		component.Sprite,
		transform.Transform,
	))

	hudPosition := dmath.NewVec2(float64(game.Settings.WorldWidth)/8.0, float64(game.Settings.WorldHeight)/16.0)
	transform.Transform.Get(hud).LocalPosition = hudPosition
	component.Sprite.Get(hud).Layer = component.SpriteLayerUI

	gameOverBackground := w.Entry(w.Create(component.Sprite, transform.Transform))
	gameOverTransform := transform.Transform.Get(gameOverBackground)
	gameOverTransform.LocalPosition = dmath.NewVec2(0, 192) // didn't have time to figure out the proper math
	gameOverBackgroundSprite := component.Sprite.Get(gameOverBackground)
	gameOverBackgroundSprite.Image = ebiten.NewImage(game.Settings.WorldWidth, game.Settings.WorldHeight/2)
	clr := component.NewColorOverrideFromColor(colornames.Black)
	gameOverBackgroundSprite.ColorOverride = clr
	gameOverBackgroundSprite.Image.Fill(clr)
	gameOverBackgroundSprite.Size = dmath.NewVec2(float64(game.Settings.WorldWidth), float64(game.Settings.WorldHeight)/2)
	gameOverBackgroundSprite.Hidden = true
	gameOverBackgroundSprite.Layer = component.SpriteLayerUI

	gameOverText := w.Entry(w.Create(component.Sprite, transform.Transform))

	gameOverTextSprite := component.Sprite.Get(gameOverText)
	text := "GAME OVER"
	var gameOverTextSizeX float64
	var gameOverTextSizeY float64
	for i, r := range text {
		p, ok := FontSpritePoints[r]
		if !ok {
			continue
		}

		scaledP := make([]path.Op, len(p))
		for j, op := range p {
			scaledOp := op.Copy(op.Target().MulScalar(2))
			scaledP[j] = scaledOp
			gameOverTextSprite.PathOps = append(gameOverTextSprite.PathOps, op.Copy(scaledOp.Target().Add(dmath.NewVec2(gameOverTextSizeX+float64(8*i), 0))))
		}

		fontRuneSize := calculateFontRuneSize(scaledP)
		gameOverTextSizeX += fontRuneSize.X
		gameOverTextSizeY = fontRuneSize.Y
	}

	gameOverTextSizeY += 20

	subText := "PRESS 'R' TO RESTART"
	var subTextSizeX float64
	for i, r := range subText {
		padding := 8

		var fontRuneSizeX float64
		p, ok := FontSpritePoints[r]
		if ok {
			scaledP := make([]path.Op, len(p))
			for j, op := range p {
				scaledOp := op.Copy(op.Target().MulScalar(0.8))
				scaledP[j] = scaledOp
				gameOverTextSprite.PathOps = append(gameOverTextSprite.PathOps, op.Copy(scaledOp.Target().Add(dmath.NewVec2(subTextSizeX+float64(padding*i), gameOverTextSizeY))))
			}

			fontRuneSizeX = calculateFontRuneSize(scaledP).X
		}

		subTextSizeX += fontRuneSizeX
	}

	gameOverTextPosition := dmath.NewVec2(294, 329.2) // same here
	gameOverTextSprite.Layer = component.SpriteLayerUI

	transform.Transform.Get(gameOverText).LocalPosition = gameOverTextPosition

	transform.AppendChild(gameOverBackground, gameOverText, true)

	transform.AppendChild(hud, gameOverBackground, true)

	return hud
}

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

package component

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/milk9111/asteroids/engine/path"
	"github.com/yohamta/donburi"
	dmath "github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
)

type SpriteLayer int

const (
	SpriteLayerBackground SpriteLayer = iota
	SpriteLayerEntity
	SpriteLayerUI
)

type SpritePivot int

const (
	SpritePivotCenter SpritePivot = iota
	SpritePivotTopLeft
)

type SpriteData struct {
	PathOps []path.Op
	Image   *ebiten.Image
	Layer   SpriteLayer `yaml:"layer"`
	Pivot   SpritePivot `yaml:"pivot"`
	Size    dmath.Vec2

	Hidden bool

	ColorOverride *ColorOverride
}

func (s *SpriteData) PivotPoint(position dmath.Vec2) dmath.Vec2 {
	if s.Pivot == SpritePivotTopLeft {
		return position
	}

	return dmath.NewVec2(position.X-(s.Size.X/2), position.Y-(s.Size.Y/2))
}

type ColorOverride struct {
	R, G, B, A float32
}

func NewColorOverride(r, g, b, a float32) *ColorOverride {
	return &ColorOverride{R: r, G: g, B: b, A: a}
}

func NewColorOverrideFromColor(color color.Color) *ColorOverride {
	r, g, b, a := color.RGBA()
	return &ColorOverride{R: float32(r), G: float32(g), B: float32(b), A: float32(a)}
}

func (c *ColorOverride) RGBA() (r, g, b, a uint32) {
	return uint32(c.R), uint32(c.G), uint32(c.B), uint32(c.A)
}

func (s *SpriteData) Show() {
	s.Hidden = false
}

func (s *SpriteData) Hide() {
	s.Hidden = true
}

var Sprite = donburi.NewComponentType[SpriteData]()

func WorldHidden(entry *donburi.Entry) bool {
	s := Sprite.Get(entry)

	p, ok := transform.GetParent(entry)
	if !ok {
		return s.Hidden
	}

	hidden := WorldHidden(p) || s.Hidden
	return hidden
}

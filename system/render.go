package system

import (
	"image"
	"image/color"
	"slices"
	"sort"

	"golang.org/x/image/colornames"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/milk9111/asteroids/archetype"
	"github.com/milk9111/asteroids/component"
	"github.com/milk9111/asteroids/engine"
	"github.com/milk9111/asteroids/engine/path"
	"github.com/yohamta/donburi"
	dmath "github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
)

var (
	whiteImage = ebiten.NewImage(3, 3)

	// whiteSubImage is an internal sub image of whiteImage.
	// Use whiteSubImage at DrawTriangles instead of whiteImage in order to avoid bleeding edges.
	whiteSubImage = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func init() {
	whiteImage.Fill(color.White)
}

type Render struct {
	query  *donburi.Query
	world  *ebiten.Image
	camera *donburi.Entry
	game   *component.GameData
}

func NewRender(w, h int) *Render {
	return &Render{
		query: donburi.NewQuery(
			filter.And(
				filter.Contains(transform.Transform, component.Sprite),
			),
		),
		world: ebiten.NewImage(w, h),
	}
}

func (r *Render) Update(w donburi.World) {
	if r.camera == nil {
		r.camera = archetype.MustFindCamera(w)
	}

	if r.game == nil {
		r.game = component.MustFindGame(w)
	}
}

func (r *Render) Draw(w donburi.World, screen *ebiten.Image) {
	if r.camera == nil || r.game == nil {
		return
	}

	r.world.Clear()

	r.world.Fill(colornames.Black)

	var entries []*donburi.Entry
	r.query.Each(w, func(entry *donburi.Entry) {
		entries = append(entries, entry)
	})

	byLayer := make(map[int][]*donburi.Entry)
	for _, entry := range entries {
		layer := int(component.Sprite.Get(entry).Layer)
		if _, ok := byLayer[layer]; !ok {
			byLayer[layer] = []*donburi.Entry{}
		}

		byLayer[layer] = append(byLayer[layer], entry)
	}

	layers := make([]int, len(byLayer))
	i := 0
	for layer := range byLayer {
		layers[i] = layer
		i++
	}

	sort.Ints(layers)

	for _, layer := range layers {
		layerEntries := byLayer[layer]
		slices.SortFunc(layerEntries, func(a, b *donburi.Entry) int {
			aSprite := component.Sprite.Get(a)
			bSprite := component.Sprite.Get(b)

			if aSprite.Hidden || bSprite.Hidden {
				return -1
			}

			return int(aSprite.PivotPoint(transform.Transform.Get(a).LocalPosition).Y - bSprite.PivotPoint(transform.Transform.Get(b).LocalPosition).Y)
		})

		path := &vector.Path{}
		var vertices []ebiten.Vertex
		var indices []uint16

		for _, entry := range byLayer[layer] {
			if component.WorldHidden(entry) {
				continue
			}

			sprite := component.Sprite.Get(entry)

			position := transform.WorldPosition(entry)
			scale := transform.WorldScale(entry)
			rotation := transform.WorldRotation(entry)
			pivot := sprite.PivotPoint(position)

			if sprite.Image == nil {
				r.drawPathOps(path, sprite.PathOps, pivot, scale, dmath.NewVec2(sprite.Size.X/2, sprite.Size.Y/2), rotation, &vertices, &indices)
			} else {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(position.X, position.Y)
				r.world.DrawImage(sprite.Image, op)
			}
		}

		for i, vt := range vertices {
			vt.SrcX = 1
			vt.SrcY = 1
			vt.ColorR = float32(colornames.White.R) / float32(0xff)
			vt.ColorG = float32(colornames.White.G) / float32(0xff)
			vt.ColorB = float32(colornames.White.B) / float32(0xff)
			vt.ColorA = 1

			vertices[i] = vt
		}

		op := &ebiten.DrawTrianglesOptions{}
		op.AntiAlias = true
		op.FillRule = ebiten.FillRuleNonZero
		r.world.DrawTriangles(vertices, indices, whiteSubImage, op)
	}

	screen.DrawImage(r.world, &ebiten.DrawImageOptions{})
}

func (r *Render) drawPathOps(path *vector.Path, ops []path.Op, worldPos, scale, rotationPivot dmath.Vec2, rotation float64, vertices *[]ebiten.Vertex, indices *[]uint16) {
	if len(ops) == 0 {
		return
	}

	for _, op := range ops {
		op.Do(path, func(target dmath.Vec2) dmath.Vec2 {
			return engine.RotateAndScaleVec2(target, rotationPivot, scale, rotation).Add(worldPos)
		})
	}

	sop := &vector.StrokeOptions{}
	sop.Width = strokeWidth
	sop.LineJoin = vector.LineJoinRound
	*vertices, *indices = path.AppendVerticesAndIndicesForStroke(*vertices, *indices, sop)
}

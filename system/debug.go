package system

import (
	"fmt"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/milk9111/asteroids/component"
	"github.com/milk9111/asteroids/engine"
	"github.com/yohamta/donburi"
	dmath "github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
	"golang.org/x/image/colornames"
)

const (
	strokeWidth float32 = 2
)

type Debug struct {
	query             *donburi.Query
	debug             *component.DebugData
	game              *component.GameData
	wave              *component.WaveData
	world             *ebiten.Image
	nextStepCallback  func()
	pauseGameCallback func()

	hoveredEntry  *donburi.Entry
	selectedEntry *donburi.Entry

	holdStepTimer *engine.Timer
}

func NewDebug(w, h int, nextStepCallback func(), pauseGameCallback func()) *Debug {
	return &Debug{
		query: donburi.NewQuery(
			filter.Contains(transform.Transform, component.Sprite, component.Collider),
		),
		world:             ebiten.NewImage(w, h),
		holdStepTimer:     engine.NewTimer(250 * time.Millisecond),
		nextStepCallback:  nextStepCallback,
		pauseGameCallback: pauseGameCallback,
	}
}

func (d *Debug) Update(w donburi.World) {
	if d.debug == nil {
		d.debug = component.MustFindDebug(w)
	}

	if d.game == nil {
		d.game = component.MustFindGame(w)
	}

	if d.wave == nil {
		d.wave = component.MustFindWave(w)
	}

	isHoldingStep := ebiten.IsKeyPressed(ebiten.KeyE)
	if isHoldingStep {
		d.holdStepTimer.Update()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyE) || (isHoldingStep && d.holdStepTimer.IsReady()) {
		d.nextStepCallback()
	}

	if inpututil.IsKeyJustReleased(ebiten.KeyE) {
		d.holdStepTimer.Reset()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		d.pauseGameCallback()
	}

	isMouseJustPressed := inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0)
	cursorX, cursorY := ebiten.CursorPosition()
	cursorRect := engine.NewRect(dmath.NewVec2(float64(cursorX), float64(cursorY)), 0, 0, 0)
	d.query.Each(w, func(e *donburi.Entry) {
		r := RectFromEntry(e)
		cursorIntersects := r.Intersects(cursorRect)
		if cursorIntersects {
			d.hoveredEntry = e

			if isMouseJustPressed {
				d.selectedEntry = e
			}
		}

		if !cursorIntersects && e == d.hoveredEntry {
			d.hoveredEntry = nil
		}
	})
}

func (d *Debug) DebugDraw(w donburi.World, screen *ebiten.Image) {
	if d.debug == nil || !d.debug.Enabled {
		return
	}

	allCount := w.Len()
	x, y := ebiten.Monitor().Size()

	debugPrints := []string{
		fmt.Sprintf("Entities: %v", allCount),
		fmt.Sprintf("TPS: %.2f, FPS: %.2f", ebiten.ActualTPS(), ebiten.ActualFPS()),
		fmt.Sprintf("Canvas size: %d x %d", d.game.Settings.WorldWidth, d.game.Settings.WorldHeight),
		fmt.Sprintf("Screen size: %d x %d", x, y),
		fmt.Sprintf("Level %d - Remaining Hazards: %d", d.wave.Number, d.wave.RemainingHazards),
		fmt.Sprintf("Score %d - Lives %d", d.game.Score, d.game.PlayerLives),
	}

	if d.selectedEntry != nil {
		if d.selectedEntry.Valid() {
			selectedEntryDebugPrints, _ := getDebugPrintAtsForEntry(d.selectedEntry)
			debugPrints = append(debugPrints, fmt.Sprintf("\nENTITY %d\n--------", d.selectedEntry.Entity().Id()))
			debugPrints = append(debugPrints, fmt.Sprintf("type: %s", component.ID.Get(d.selectedEntry).Type))
			debugPrints = append(debugPrints, selectedEntryDebugPrints...)
		} else {
			debugPrints = append(debugPrints, "\n(destroyed)")
		}
	}

	ebitenutil.DebugPrint(screen, strings.Join(debugPrints, "\n"))

	d.world.Clear()
	d.query.Each(w, func(entry *donburi.Entry) {
		debugPrintAts, pos := getDebugPrintAtsForEntry(entry)

		ebitenutil.DebugPrintAt(d.world, fmt.Sprintf("%d", entry.Entity().Id()), int(pos.X), int(pos.Y))

		for i, debugPrintAt := range debugPrintAts {
			ebitenutil.DebugPrintAt(d.world, debugPrintAt, int(pos.X), int(pos.Y)+30+15*i)
		}

		if entry == d.hoveredEntry || entry == d.selectedEntry {
			var path vector.Path
			path.MoveTo(float32(pos.X), float32(pos.Y))
			r := 16.0
			if entry == d.selectedEntry {
				r = 8
			}
			path.Arc(float32(pos.X), float32(pos.Y), float32(r), 0, 360, vector.Clockwise)
			path.Close()

			sop := &vector.StrokeOptions{}
			sop.Width = 3
			sop.LineJoin = vector.LineJoinRound
			vertices, indices := path.AppendVerticesAndIndicesForStroke([]ebiten.Vertex{}, []uint16{}, sop)

			color := colornames.Lightgreen
			if entry == d.selectedEntry {
				color = colornames.Lightblue
			}

			for i := range vertices {
				vertices[i].SrcX = 1
				vertices[i].SrcY = 1
				vertices[i].ColorR = float32(color.R) / float32(0xff)
				vertices[i].ColorG = float32(color.G) / float32(0xff)
				vertices[i].ColorB = float32(color.B) / float32(0xff)
				vertices[i].ColorA = 1
			}

			op := &ebiten.DrawTrianglesOptions{}
			op.AntiAlias = true
			op.FillRule = ebiten.FillRuleFillAll
			d.world.DrawTriangles(vertices, indices, whiteSubImage, op)
		}
	})

	op := &ebiten.DrawImageOptions{}

	screen.DrawImage(d.world, op)
}

func getDebugPrintAtsForEntry(entry *donburi.Entry) ([]string, dmath.Vec2) {
	var sprite *component.SpriteData
	if entry.HasComponent(component.Sprite) {
		sprite = component.Sprite.Get(entry)
	}

	position := transform.WorldPosition(entry)

	pivot := position
	rotationPivot := pivot
	if sprite != nil {
		halfW, halfH := sprite.Size.X/2, sprite.Size.Y/2
		rotationPivot.X += halfW
		rotationPivot.Y += halfH
		switch sprite.Pivot {
		case component.SpritePivotCenter:
			pivot.X += halfW
			pivot.Y += halfH
			rotationPivot.X -= halfW
			rotationPivot.Y -= halfH
		}
	}

	debugPrintAts := []string{
		fmt.Sprintf("pos: %.0f, %.0f", position.X, position.Y),
		fmt.Sprintf("rot pos: %.0f, %.0f", rotationPivot.X, rotationPivot.Y),
		fmt.Sprintf("rot: %.2f", transform.WorldRotation(entry)),
	}

	if entry.HasComponent(component.Velocity) {
		velocity := component.Velocity.Get(entry)
		debugPrintAts = append(debugPrintAts, fmt.Sprintf("vel: %.2f, %.2f", velocity.Velocity().X, velocity.Velocity().Y))
	}

	if entry.HasComponent(component.Traveler) {
		traveler := component.Traveler.Get(entry)
		debugPrintAts = append(debugPrintAts, fmt.Sprintf("dir: %.2f, %.2f - speed: %.2f", traveler.Direction.X, traveler.Direction.Y, traveler.Speed))
	}

	return debugPrintAts, rotationPivot
}

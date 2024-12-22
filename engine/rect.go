package engine

import (
	"math"

	dmath "github.com/yohamta/donburi/features/math"
)

type Rect struct {
	Center         dmath.Vec2
	PositionOffset dmath.Vec2
	Width          float64
	Height         float64
	Angle          float64 // Rotation in radians

	offsets []dmath.Vec2
}

func NewRect(center dmath.Vec2, width, height, angle float64) Rect {
	hw := width / 2
	hh := height / 2

	return Rect{
		Center: center,
		Width:  width,
		Height: height,
		Angle:  angle,
		offsets: []dmath.Vec2{
			dmath.NewVec2(-hw, -hh),
			dmath.NewVec2(hw, -hh),
			dmath.NewVec2(hw, hh),
			dmath.NewVec2(-hw, hh),
		},
	}
}

func (r Rect) Intersects(r2 Rect) bool {
	vertices1 := r.GetVertices()
	vertices2 := r2.GetVertices()

	edges := append(getEdges(vertices1), getEdges(vertices2)...)

	for _, edge := range edges {
		axis := dmath.NewVec2(-edge.Y, edge.X)
		if projectionCheck(vertices1, vertices2, axis) {
			return false
		}
	}

	return true
}

func (r Rect) GetVertices() []dmath.Vec2 {
	cosTheta := math.Cos(r.Angle)
	sinTheta := math.Sin(r.Angle)

	vertices := make([]dmath.Vec2, 4)
	for i, offset := range r.offsets {
		vertices[i] = dmath.NewVec2(r.Center.X+offset.X*cosTheta-offset.Y*sinTheta, r.Center.Y+offset.X*sinTheta+offset.Y*cosTheta)
	}

	return vertices
}

func getEdges(vertices []dmath.Vec2) []dmath.Vec2 {
	edges := make([]dmath.Vec2, len(vertices))
	for i := 0; i < len(vertices); i++ {
		next := (i + 1) % len(vertices)
		edges[i] = vertices[next].Sub(vertices[i])
	}

	return edges
}

func projectionCheck(vertices1, vertices2 []dmath.Vec2, axis dmath.Vec2) bool {
	proj1Min, proj1Max := projectVertices(vertices1, axis)
	proj2Min, proj2Max := projectVertices(vertices2, axis)

	return proj1Max < proj2Min || proj2Max < proj1Min
}

func projectVertices(vertices []dmath.Vec2, axis dmath.Vec2) (float64, float64) {
	axisLength := math.Hypot(axis.X, axis.Y)
	axis = axis.DivScalar(axisLength)

	projection := vertices[0].Dot(&axis)

	min, max := projection, projection
	for _, vertex := range vertices[1:] {
		projection = vertex.Dot(&axis)
		min = math.Min(min, projection)
		max = math.Max(max, projection)
	}

	return min, max
}

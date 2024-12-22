package engine

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
)

var (
	worldRight = math.NewVec2(1, 0)
	worldUp    = math.NewVec2(0, -1)

	vec2Zero = math.NewVec2(0, 0)
	vec2One  = math.NewVec2(1, 1)
)

func WorldRight() math.Vec2 {
	return worldRight
}

func WorldUp() math.Vec2 {
	return worldUp
}

func Vec2Zero() math.Vec2 {
	return vec2Zero
}

func Vec2One() math.Vec2 {
	return vec2One
}

func FinalWorldPosition(entry *donburi.Entry) math.Vec2 {
	t := transform.Transform.Get(entry)
	return PositionToFinalWorldPosition(t.LocalPosition, t.LocalScale.X)
}

func PositionToFinalWorldPosition(pos math.Vec2, scaleFactor float64) math.Vec2 {
	return pos.MulScalar(scaleFactor)
}

func MidPoint(p1, p2 math.Vec2) math.Vec2 {
	return math.NewVec2(
		(p1.X+p2.X)/2,
		(p1.Y+p2.Y)/2,
	)
}

package engine

import (
	"math"
	"math/rand/v2"

	dmath "github.com/yohamta/donburi/features/math"
)

func RotateVec2(vec2 dmath.Vec2, rotation float64) dmath.Vec2 {
	return dmath.NewVec2(vec2.X*math.Cos(rotation)-vec2.Y*math.Sin(rotation), vec2.Y*math.Cos(rotation)+vec2.X*math.Sin(rotation))
}

func RotateAndScaleVec2(point, rotationPivot, scale dmath.Vec2, rotation float64) dmath.Vec2 {
	p := rotationPivot.MulScalar(-1)
	p = point.Add(p)
	p = p.Mul(scale)
	p = RotateVec2(p, rotation)
	p = p.Add(rotationPivot)
	return p
}

func RandomRotation() float64 {
	return float64(rand.Int() % 360)
}

func DegToRad(x float64) float64 {
	return x * (math.Pi / 180)
}

func RadToDeg(x float64) float64 {
	return x * (180 / math.Pi)
}

func DegToDirection(deg float64) dmath.Vec2 {
	rad := DegToRad(deg)
	return dmath.NewVec2(math.Cos(rad), math.Sin(rad))
}

func RadToDirection(rad float64) dmath.Vec2 {
	return dmath.NewVec2(math.Cos(rad), math.Sin(rad))
}

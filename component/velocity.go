package component

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
)

type VelocityData struct {
	velocity         math.Vec2
	rotationVelocity float64
}

func NewVelocityData(velocity math.Vec2, rotationVelocity float64) VelocityData {
	return VelocityData{
		velocity:         velocity,
		rotationVelocity: rotationVelocity,
	}
}

func (v *VelocityData) Velocity() math.Vec2 {
	return v.velocity
}

func (v *VelocityData) RotationVelocity() float64 {
	return v.rotationVelocity
}

func (v *VelocityData) Reset() {
	v.velocity = math.NewVec2(0, 0)
	v.rotationVelocity = 0
}

func (v *VelocityData) AddVelocity(velocity math.Vec2) {
	v.velocity = v.velocity.Add(velocity)
}

func (v *VelocityData) AddRotationalVelocity(rotationalVelocity float64) {
	v.rotationVelocity += rotationalVelocity
}

var Velocity = donburi.NewComponentType[VelocityData]()

package physics

import "spacewar/geometry"

type Object struct {
	Center      geometry.PointFloatXY
	SpeedX      float32
	SpeedY      float32
	Radius      int
	RotateAngle float32
}

func (o *Object) UpdatePosition() {
	o.Center.X += o.SpeedX
	o.Center.Y += o.SpeedY
}

func (o *Object) ChangeInertia(deltaSpeedX, deltaSpeedY float32) {
	o.SpeedX += deltaSpeedX
	o.SpeedY += deltaSpeedY
}

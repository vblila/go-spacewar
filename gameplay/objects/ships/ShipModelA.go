package ships

import (
	"spacewar/geometry"
	"spacewar/physics"
)

type ShipModelA struct {
}

func (s *ShipModelA) BodyPolygon(po physics.Object) []geometry.PointFloatXY {
	halfLen := float32(po.Radius / 2)
	thirdLen := float32(po.Radius / 3)

	return []geometry.PointFloatXY{
		{po.Center.X, po.Center.Y - halfLen},
		{po.Center.X - halfLen, po.Center.Y + thirdLen},
		{po.Center.X - thirdLen, po.Center.Y + halfLen},
		{po.Center.X + thirdLen, po.Center.Y + halfLen},
		{po.Center.X + halfLen, po.Center.Y + thirdLen},
		{po.Center.X, po.Center.Y - halfLen},
	}
}

func (s *ShipModelA) rotatePower() float32 {
	return 0.1
}

func (s *ShipModelA) enginePower() float32 {
	return 0.05
}

func (s *ShipModelA) gunPower() float32 {
	return 8
}

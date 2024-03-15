package ships

import (
	"spacewar/geometry"
	"spacewar/physics"
)

type ShipModelB struct {
}

func (s *ShipModelB) BodyPolygon(po physics.Object) []geometry.PointFloatXY {
	halfLen := float32(po.Radius / 2)
	thirdLen := float32(po.Radius / 3)

	return []geometry.PointFloatXY{
		{po.Center.X, po.Center.Y - halfLen},
		{po.Center.X - thirdLen, po.Center.Y + halfLen},
		{po.Center.X + thirdLen, po.Center.Y + halfLen},
		{po.Center.X, po.Center.Y - halfLen},
	}
}

func (s *ShipModelB) rotatePower() float32 {
	return 0.1
}

func (s *ShipModelB) enginePower() float32 {
	return 0.05
}

func (s *ShipModelB) gunPower() float32 {
	return 16
}

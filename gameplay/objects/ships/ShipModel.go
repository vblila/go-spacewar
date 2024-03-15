package ships

import (
	"spacewar/geometry"
	"spacewar/physics"
)

type ShipModel interface {
	BodyPolygon(po physics.Object) []geometry.PointFloatXY

	rotatePower() float32
	enginePower() float32
	gunPower() float32
}

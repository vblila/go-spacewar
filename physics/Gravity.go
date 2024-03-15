package physics

import (
	"math"
	"spacewar/config"
	"spacewar/geometry"
)

type Gravity struct {
	Power float32
}

func (g *Gravity) GetInfluence(starRadius int, starXY geometry.PointFloatXY, objectXY geometry.PointFloatXY) (deltaSpeedX float32, deltaSpeedY float32) {
	distance := geometry.Distance(objectXY, starXY)
	gravity := float32(math.Max(float64(starRadius)-10, 1)) * g.Power / (distance * distance)

	deltaSpeedX = (starXY.X - objectXY.X) * gravity * config.GetGameSpeed()
	deltaSpeedY = (starXY.Y - objectXY.Y) * gravity * config.GetGameSpeed()

	return
}

func (g *Gravity) DoInfluence(starRadius int, starXY geometry.PointFloatXY, o *Object) {
	o.ChangeInertia(g.GetInfluence(starRadius, starXY, o.Center))
}

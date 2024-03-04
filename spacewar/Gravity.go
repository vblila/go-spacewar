package spacewar

import "math"

type Gravity struct {
	Power float32
}

func (g *Gravity) Influence(planetRadius int, planetXY PointXY, objectXY PointXY) (deltaSpeedX float32, deltaSpeedY float32) {
	distance := math.Hypot((float64)(objectXY.X-planetXY.X), (float64)(objectXY.Y-planetXY.Y))
	gravity := (float32)(planetRadius-MinPlanetRadius+1) * g.Power / (float32)(distance*distance)

	deltaSpeedX = (float32)(planetXY.X-objectXY.X) * gravity
	deltaSpeedY = (float32)(planetXY.Y-objectXY.Y) * gravity

	return
}

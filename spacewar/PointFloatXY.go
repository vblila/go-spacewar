package spacewar

import "math"

type PointFloatXY struct {
	X float32
	Y float32
}

func (p *PointFloatXY) GetPointXY() PointXY {
	return PointXY{(int)(p.X), (int)(p.Y)}
}

func (p *PointFloatXY) Rotate(center PointFloatXY, angle float32) {
	sin, cos := (float32)(math.Sin((float64)(angle))), (float32)(math.Cos((float64)(angle)))

	x := cos*(p.X-center.X) - sin*(p.Y-center.Y) + center.X
	y := sin*(p.X-center.X) + cos*(p.Y-center.Y) + center.Y

	p.X, p.Y = x, y
}

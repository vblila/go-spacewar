package spacewar

import "github.com/gonutz/prototype/draw"

type Bomb struct {
	Center PointFloatXY
	SpeedX float32
	SpeedY float32
}

func (b *Bomb) UpdatePosition() {
	b.Center.X += b.SpeedX
	b.Center.Y += b.SpeedY
}

func (b *Bomb) Render(window draw.Window) {
	centerXY := b.Center.GetPointXY()
	window.FillEllipse(centerXY.X-BombRadius, centerXY.Y-BombRadius, BombRadius*2, BombRadius*2, draw.White)
}

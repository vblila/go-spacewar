package spacewar

import "github.com/gonutz/prototype/draw"

type Planet struct {
	Center PointXY
	Radius int
}

func (p *Planet) Render(window draw.Window) {
	window.FillEllipse(p.Center.X-p.Radius, p.Center.Y-p.Radius, p.Radius*2, p.Radius*2, draw.Black)
	window.DrawEllipse(p.Center.X-p.Radius, p.Center.Y-p.Radius, p.Radius*2, p.Radius*2, draw.White)
}

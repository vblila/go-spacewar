package objects

import (
	"spacewar/graphics"
	"spacewar/physics"
)

type Star struct {
	PO physics.Object
}

func (p *Star) Render(gw *graphics.GLWindow) {
	gw.FillCircle(p.PO.Center.X, p.PO.Center.Y, float32(p.PO.Radius), graphics.Color{A: 1})
	gw.DrawCircle(p.PO.Center.X, p.PO.Center.Y, float32(p.PO.Radius), gw.Color)
}

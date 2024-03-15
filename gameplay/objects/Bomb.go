package objects

import (
	"spacewar/config"
	"spacewar/graphics"
	"spacewar/physics"
)

type Bomb struct {
	PO physics.Object
}

func (b *Bomb) Render(gw *graphics.GLWindow) {
	gw.DrawCircle(b.PO.Center.X, b.PO.Center.Y, config.BombRadius, gw.Color)
	gw.FillCircle(b.PO.Center.X, b.PO.Center.Y, config.BombRadius, gw.Color)
}

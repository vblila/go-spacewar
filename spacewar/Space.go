package spacewar

import (
	"github.com/gonutz/prototype/draw"
	"math/rand"
)

type Space struct {
	Points []PointXY
}

func (s *Space) InitSpace(count int, width int, height int) {
	s.Points = make([]PointXY, count)
	for i := 0; i < count; i++ {
		s.Points[i].X = rand.Intn(width)
		s.Points[i].Y = rand.Intn(height)
	}
}

func (s *Space) Render(window draw.Window) {
	for i := 0; i < len(s.Points); i++ {
		window.DrawPoint(s.Points[i].X, s.Points[i].Y, draw.White)
	}
}

package objects

import (
	"math/rand"
	"spacewar/config"
	"spacewar/graphics"
	"spacewar/physics"
)

type Space struct {
	Points        []physics.Object
	width, height int
}

func (s *Space) InitSpace(count int, width int, height int) {
	s.width, s.height = width, height

	s.Points = make([]physics.Object, count)
	for i := 0; i < count; i++ {
		s.Points[i].Center.X = float32(rand.Intn(width))
		s.Points[i].Center.Y = float32(rand.Intn(height))

		s.Points[i].SpeedX = rand.Float32() * 6 * config.GetGameSpeed()
	}
}

func (s *Space) Render(gw *graphics.GLWindow) {
	for i := 0; i < len(s.Points); i++ {
		gw.DrawPoint(s.Points[i].Center.X, s.Points[i].Center.Y, gw.Color)
	}
}

func (s *Space) Fly() {
	for i := 0; i < len(s.Points); i++ {
		s.Points[i].UpdatePosition()
		if int(s.Points[i].Center.X) > s.width {
			s.Points[i].Center.X = 0
		}
	}
}

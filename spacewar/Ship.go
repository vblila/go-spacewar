package spacewar

import (
	"github.com/gonutz/prototype/draw"
	"time"
)

type Ship struct {
	Center      PointFloatXY
	SpeedX      float32
	SpeedY      float32
	RotateAngle float32
	Power       float32

	lastFireTime time.Time
	accelerating bool
}

func (s *Ship) UpdatePosition() {
	s.Center.X += s.SpeedX
	s.Center.Y += s.SpeedY
}

func (s *Ship) RotateLeft() {
	s.RotateAngle -= s.Power
}

func (s *Ship) RotateRight() {
	s.RotateAngle += s.Power
}

func (s *Ship) Accelerate() {
	// Вычисляем точку, в которую устремлен вектор скорости
	// При rotateAngle = 0 вектор скорости направлен вверх
	speedPoint := s.Center
	speedPoint.Y -= s.Power

	// Поворачиваем точку вектора скорости по направлению кабины корабля
	speedPoint.Rotate(s.Center, s.RotateAngle)

	s.SpeedX += speedPoint.X - s.Center.X
	s.SpeedY += speedPoint.Y - s.Center.Y

	s.accelerating = true
}

func (s *Ship) Fire() *Bomb {
	// Пушка позволяет выпускать один снаряд в 300 мс
	if s.lastFireTime.After(time.Now().Add(-time.Millisecond * 300)) {
		return nil
	}

	bomb := &Bomb{Center: s.Center}

	// Положение пушки находится у кабины, снаряд вылетает отсюда
	bomb.Center.Y -= (float32)(ShipRadius / 2)

	// Вычисляем точку, в которую устремлен вектор выстрела
	// При rotateAngle = 0 вектор выстрела направлен от кабины вверх
	firePoint := bomb.Center
	firePoint.Y -= ShipGunPower

	// Поворачиваем точку вектора выстрела по направлению кабины корабля
	firePoint.Rotate(s.Center, s.RotateAngle)
	bomb.Center.Rotate(s.Center, s.RotateAngle)

	bomb.SpeedX = firePoint.X - bomb.Center.X
	bomb.SpeedY = firePoint.Y - bomb.Center.Y

	// Бомба учитывает инерцию корабля
	bomb.SpeedX += s.SpeedX
	bomb.SpeedY += s.SpeedY

	s.lastFireTime = time.Now()

	return bomb
}

func (s *Ship) Render(window draw.Window) {
	halfLen := (float32)(ShipRadius / 2)
	thirdLen := (float32)(ShipRadius / 3)

	shipPoints := []PointFloatXY{
		{s.Center.X, s.Center.Y - halfLen},
		{s.Center.X - thirdLen, s.Center.Y + halfLen},
		{s.Center.X + thirdLen, s.Center.Y + halfLen},
	}

	var flamePoints []PointFloatXY

	if s.accelerating {
		flamePoints = []PointFloatXY{
			{s.Center.X - 7, s.Center.Y + halfLen},
			{s.Center.X, s.Center.Y + ShipRadius},
			{s.Center.X + 7, s.Center.Y + halfLen},
		}
	}

	// Разворачиваем все точки вокруг центра
	for i := 0; i < len(shipPoints); i++ {
		shipPoints[i].Rotate(s.Center, s.RotateAngle)
	}
	for i := 0; i < len(flamePoints); i++ {
		flamePoints[i].Rotate(s.Center, s.RotateAngle)
	}

	// Рисуем точки корабля
	DrawPolygon(shipPoints, window)
	cabin := shipPoints[0].GetPointXY()
	window.FillEllipse(cabin.X-2, cabin.Y-3, 6, 6, draw.White)

	// Рисуем пламя двигателя
	DrawPolygon(flamePoints, window)

	s.accelerating = false
}

package ships

import (
	"math"
	"spacewar/config"
	"spacewar/gameplay/objects"
	"spacewar/geometry"
	"spacewar/graphics"
	"spacewar/physics"
	"time"
)

type Ship struct {
	PO    physics.Object
	model ShipModel

	gunRecoveryExpTime       time.Time
	hyperEngineExpTime       time.Time
	hyperStopExpTime         time.Time
	accelerationFlameExpTime time.Time

	Damage int
}

func NewShip(model ShipModel) *Ship {
	ship := &Ship{model: model}
	ship.PO.Radius = config.ShipRadius
	return ship
}

func (s *Ship) RotateLeft() {
	s.PO.RotateAngle -= s.model.rotatePower() * config.GetGameSpeed()
}

func (s *Ship) RotateRight() {
	s.PO.RotateAngle += s.model.rotatePower() * config.GetGameSpeed()
}

func (s *Ship) RotateToAngle(targetAngle float32) {
	myRotateAngle := geometry.NormalizeAngle2Pi(s.PO.RotateAngle)

	if targetAngle > myRotateAngle {
		if targetAngle-myRotateAngle > math.Pi {
			s.RotateLeft()
		} else {
			s.RotateRight()
		}
	} else {
		if myRotateAngle-targetAngle > math.Pi {
			s.RotateRight()
		} else {
			s.RotateLeft()
		}
	}
}

func (s *Ship) IsTurnedToTarget(targetAngle float32) bool {
	if math.Abs(float64(geometry.NormalizeAngle2Pi(s.PO.RotateAngle)-targetAngle)) < 0.1 {
		return true
	}

	return false
}

func (s *Ship) Accelerate() {
	enginePower := s.model.enginePower()
	if time.Now().Before(s.hyperEngineExpTime) {
		enginePower *= 4
	}

	// Вычисляем точку, в которую устремлен вектор скорости
	// При rotateAngle = 0 вектор скорости направлен вверх
	speedPoint := s.PO.Center
	speedPoint.Y -= enginePower * config.GetGameSpeed()

	// Поворачиваем точку вектора скорости по направлению кабины корабля
	speedPoint.Rotate(s.PO.Center, s.PO.RotateAngle)

	// Задаем ускорение
	s.PO.ChangeInertia(speedPoint.X-s.PO.Center.X, speedPoint.Y-s.PO.Center.Y)

	s.accelerationFlameExpTime = time.Now().Add(time.Millisecond * time.Duration(30.0/config.GetGameSpeed()))
}

func (s *Ship) ActivateHyperEngine() {
	s.hyperEngineExpTime = time.Now().Add(time.Millisecond * time.Duration(150.0/config.GetGameSpeed()))
}

func (s *Ship) HyperStop() {
	// Мощность корабля не позволяет часто делать гипер остановку
	if time.Now().Before(s.hyperStopExpTime) {
		return
	}

	s.PO.SpeedX, s.PO.SpeedY = 0, 0
	s.hyperStopExpTime = time.Now().Add(time.Millisecond * time.Duration(600.0/config.GetGameSpeed()))
}

func (s *Ship) Fire() *objects.Bomb {
	if time.Now().Before(s.gunRecoveryExpTime) {
		return nil
	}

	bomb := &objects.Bomb{}
	bomb.PO.Center = s.PO.Center

	// Положение пушки находится у кабины, снаряд вылетает отсюда
	bomb.PO.Center.Y -= float32(s.PO.Radius / 2)

	// Вычисляем точку, в которую устремлен вектор выстрела
	// При rotateAngle = 0 вектор выстрела направлен от кабины вверх
	firePoint := bomb.PO.Center
	firePoint.Y -= s.model.gunPower()

	// Поворачиваем точку вектора выстрела по направлению кабины корабля
	firePoint.Rotate(s.PO.Center, s.PO.RotateAngle)
	bomb.PO.Center.Rotate(s.PO.Center, s.PO.RotateAngle)

	// Задаем вектор движения бомбы
	bomb.PO.ChangeInertia((firePoint.X-bomb.PO.Center.X)*config.GetGameSpeed(), (firePoint.Y-bomb.PO.Center.Y)*config.GetGameSpeed())

	// Бомба учитывает инерцию корабля
	bomb.PO.ChangeInertia(s.PO.SpeedX, s.PO.SpeedY)

	// Пушка не может часто стрелять
	s.gunRecoveryExpTime = time.Now().Add(time.Millisecond * time.Duration(300.0/config.GetGameSpeed()))

	return bomb
}

func (s *Ship) Lives() int {
	return config.MaxShipDamage - s.Damage
}

func (s *Ship) AccelerationFlamePolygon() []geometry.PointFloatXY {
	halfLen := float32(s.PO.Radius / 2)

	return []geometry.PointFloatXY{
		{s.PO.Center.X - 7, s.PO.Center.Y + halfLen},
		{s.PO.Center.X, s.PO.Center.Y + float32(s.PO.Radius)},
		{s.PO.Center.X + 7, s.PO.Center.Y + halfLen},
	}
}

func (s *Ship) RotatedBodyPolygon() []geometry.PointFloatXY {
	bodyPoints := s.model.BodyPolygon(s.PO)

	// Разворачиваем все точки вокруг центра
	for i := 0; i < len(bodyPoints); i++ {
		bodyPoints[i].Rotate(s.PO.Center, s.PO.RotateAngle)
	}

	return bodyPoints
}

func (s *Ship) Render(gw *graphics.GLWindow) {
	bodyPoints := s.RotatedBodyPolygon()

	var flamePoints []geometry.PointFloatXY
	if time.Now().Before(s.accelerationFlameExpTime) {
		flamePoints = s.AccelerationFlamePolygon()
	}

	// Рисуем точки корабля
	gw.DrawPolygon(bodyPoints, gw.Color)
	gw.FillCircle(bodyPoints[0].X, bodyPoints[0].Y, 3, gw.Color)

	// Рисуем пламя двигателя
	for i := 0; i < len(flamePoints); i++ {
		flamePoints[i].Rotate(s.PO.Center, s.PO.RotateAngle)
	}
	gw.DrawPolygon(flamePoints, gw.Color)

	// Если сработал гипер стоп, рисуем его
	if time.Now().Before(s.hyperStopExpTime) {
		hyperStopPoint := s.PO.Center
		hyperStopPoint.Y += 7
		hyperStopPoint.Rotate(s.PO.Center, s.PO.RotateAngle)

		gw.DrawCircle(hyperStopPoint.X, hyperStopPoint.Y, float32(s.PO.Radius), gw.Color)
	}
}

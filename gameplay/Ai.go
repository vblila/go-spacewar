package gameplay

import (
	"math"
	"math/rand"
	"spacewar/config"
	"spacewar/geometry"
)

type Ai struct {
	game *Game
}

func (ai *Ai) DoIterationLogic() {
	g := ai.game

	enemyToUserAngle := geometry.VectorAngle(g.Enemy.PO.Center, g.User.PO.Center)
	enemyToUserDistance := geometry.Distance(g.Enemy.PO.Center, g.User.PO.Center)

	// AI избегает сильного влияния нейтронных звёзд, чтобы не врезаться
	wasAntigravityManeuver := ai.antigravityManeuver()

	// AI пытается стрелять в корабль пользователя, если кабина направлена в его сторону
	fireProbability := int(20000 / enemyToUserDistance)
	if g.Enemy.IsTurnedToTarget(enemyToUserAngle) && ai.toBeOrNotToBe(fireProbability) {
		g.AddBomb(g.Enemy.Fire())
	}

	// Если AI не пытается облететь нейтронную звёзду, он поворачивает в сторону корабля пользователя и летит к нему
	if !wasAntigravityManeuver {
		g.Enemy.RotateToAngle(enemyToUserAngle)

		var accelerationProbability int
		if enemyToUserDistance > 300 {
			accelerationProbability = 1000
		} else {
			accelerationProbability = int(enemyToUserDistance * enemyToUserDistance / 1000)
		}
		if ai.toBeOrNotToBe(accelerationProbability) {
			g.Enemy.Accelerate()
		}
	}

	// Когда AI далеко улетает, он активирует гипер остановку
	if g.IsOutOfScene(g.Enemy.PO.Center.PointXY(), g.Enemy.PO.Radius*2) {
		g.Enemy.HyperStop()
	}
}

func (ai *Ai) antigravityManeuver() bool {
	g := ai.game

	if len(g.Stars) == 0 {
		return false
	}

	doing := false

	isTooCloseToStar := false

	// Вычислим общий вектор гравитации по всем нейтронным звёздам
	var totalDeltaSpeedX, totalDeltaSpeedY float32
	var gravityAngle float32

	for i := 0; i < len(g.Stars); i++ {
		deltaSpeedX, deltaSpeedY := g.Gravity.GetInfluence(g.Stars[i].PO.Radius, g.Stars[i].PO.Center, g.Enemy.PO.Center)

		// Корабль слишком близко к нейтронной звезде, гравитационный манёвр нужно сделать вокруг неё
		if geometry.Distance(g.Enemy.PO.Center, g.Stars[i].PO.Center)-float32(g.Stars[i].PO.Radius) < 100 {
			isTooCloseToStar = true

			totalDeltaSpeedX = g.Stars[i].PO.Center.X - g.Enemy.PO.Center.X
			totalDeltaSpeedY = g.Stars[i].PO.Center.Y - g.Enemy.PO.Center.Y

			break
		}

		totalDeltaSpeedX += deltaSpeedX
		totalDeltaSpeedY += deltaSpeedY
	}

	// При большой близости к нейтронной звёзде AI делает гипер остановку и запускает гипер двигатель
	if isTooCloseToStar {
		g.Enemy.HyperStop()
		g.Enemy.ActivateHyperEngine()
	}

	gravityAngle = geometry.VectorAngle(g.Enemy.PO.Center, geometry.PointFloatXY{X: g.Enemy.PO.Center.X + totalDeltaSpeedX, Y: g.Enemy.PO.Center.Y + totalDeltaSpeedY})

	// При сильном влиянии AI пытается облететь нейтронную звёзду
	gravityPowerThreshold := float64(0.02 * config.GetGameSpeed())
	if isTooCloseToStar || math.Abs(float64(totalDeltaSpeedX)) > gravityPowerThreshold || math.Abs(float64(totalDeltaSpeedY)) > gravityPowerThreshold {
		doing = true

		// Разворачиваться и лететь в противоположную стороны от вектора гравитации - плохая идея
		// Лучше повернуть на угол > 90 градусов по отношению к вектору гравитации, чтобы совершить гравитационный манёвр
		// AI вычисляет ближайшее направление из возможных: antigravityManeuverAngle1, antigravityManeuverAngle2
		antigravityManeuverAngle1 := geometry.NormalizeAngle2Pi(gravityAngle - math.Pi/2)
		antigravityManeuverAngle2 := geometry.NormalizeAngle2Pi(gravityAngle + math.Pi/2)

		var antigravityManeuverAngle float32
		nearestTarget := geometry.NearestRotateAngle(g.Enemy.PO.RotateAngle, antigravityManeuverAngle1, antigravityManeuverAngle2)
		if nearestTarget == 1 {
			antigravityManeuverAngle = antigravityManeuverAngle1
		} else {
			antigravityManeuverAngle = antigravityManeuverAngle2
		}

		g.Enemy.RotateToAngle(antigravityManeuverAngle)

		// Как только AI развернет корабль в нужном направлении, он включит двигатель, чтобы выполнить гравитационный маневр
		if g.Enemy.IsTurnedToTarget(antigravityManeuverAngle) && (ai.toBeOrNotToBe(500) || isTooCloseToStar) {
			g.Enemy.Accelerate()
		}
	}

	return doing
}

func (ai *Ai) toBeOrNotToBe(threshold int) bool {
	if rand.Intn(1000) <= threshold {
		return true
	}

	return false
}

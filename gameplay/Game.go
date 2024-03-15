package gameplay

import (
	"math/rand"
	"spacewar/config"
	"spacewar/embedded"
	"spacewar/gameplay/objects"
	"spacewar/gameplay/objects/ships"
	"spacewar/geometry"
	"spacewar/graphics"
	"spacewar/physics"
)

const (
	GameIsNotInitialized      = iota
	GameStateMainMenu         = iota
	GameStateControlSettings  = iota
	GameStatePlay             = iota
	GameStateGameOverCrashed  = iota
	GameStateGameOverKamikaze = iota
	GameStateGameOverShipLost = iota
	GameStateGameOverVictory  = iota
)

type Game struct {
	gw    *graphics.GLWindow
	State uint8
	Ai    *Ai

	freestyleMode bool

	// Объекты игрового процесса
	Space   objects.Space
	Gravity physics.Gravity
	Stars   []objects.Star
	User    *ships.Ship
	Enemy   *ships.Ship
	Bombs   [config.BombsLimit]*objects.Bomb

	UserScore  uint16
	EnemyScore uint16
}

func (g *Game) start(freestyleMode bool) {
	g.Gravity = physics.Gravity{Power: 0.1}
	g.Stars = make([]objects.Star, rand.Intn(config.MaxStars+1))
	g.Bombs = [config.BombsLimit]*objects.Bomb{}

	// Создаем корабль пользователя
	g.User = ships.NewShip(&ships.ShipModelA{})
	g.User.PO.Center.X, g.User.PO.Center.Y = float32(g.gw.Width/6), float32(g.gw.Height/6)

	g.freestyleMode = freestyleMode

	if g.freestyleMode {
		g.Ai = nil
		g.Enemy = nil
	} else {
		// Создаем корабль компьютера
		g.Enemy = ships.NewShip(&ships.ShipModelB{})
		g.Enemy.PO.Center.X, g.Enemy.PO.Center.Y = float32(g.gw.Width-g.gw.Width/6), float32(g.gw.Height-g.gw.Height/6)

		// Включаем AI
		g.Ai = &Ai{game: g}
	}

	// Создаем нейтронные звёзды случайным образом
	for i := 0; i < len(g.Stars); i++ {
		g.Stars[i].PO.Radius = config.MinStarRadius + rand.Intn(config.MaxStarRadius-config.MinStarRadius)/len(g.Stars)

		centerOffsetX := -rand.Intn(g.gw.Width/4) + rand.Intn(g.gw.Width/4)
		centerOffsetY := -rand.Intn(g.gw.Width/4) + rand.Intn(g.gw.Width/4)

		g.Stars[i].PO.Center.X, g.Stars[i].PO.Center.Y = float32(g.gw.Width/2+centerOffsetX), float32(g.gw.Height/2+centerOffsetY)
	}

	g.State = GameStatePlay
}

func (g *Game) updatePositions() {
	g.User.PO.UpdatePosition()

	if g.Enemy != nil {
		g.Enemy.PO.UpdatePosition()
	}

	for i := 0; i < len(g.Bombs); i++ {
		if g.Bombs[i] != nil {
			g.Bombs[i].PO.UpdatePosition()
		}
	}
}

func (g *Game) doGravityInfluence() {
	for i := 0; i < len(g.Stars); i++ {
		// Нейтронные звёзды притягивают корабль пользователя
		g.Gravity.DoInfluence(g.Stars[i].PO.Radius, g.Stars[i].PO.Center, &g.User.PO)

		// Нейтронные звёзды притягивают корабль AI
		if g.Enemy != nil {
			g.Gravity.DoInfluence(g.Stars[i].PO.Radius, g.Stars[i].PO.Center, &g.Enemy.PO)
		}

		// Нейтронные звёзды притягивают бомбы
		for j := 0; j < len(g.Bombs); j++ {
			if g.Bombs[j] == nil {
				continue
			}
			g.Gravity.DoInfluence(g.Stars[i].PO.Radius, g.Stars[i].PO.Center, &g.Bombs[j].PO)
		}
	}
}

func (g *Game) checkBounds() {
	// Корабли врезается в нейтронные звёзды
	for i := 0; i < len(g.Stars); i++ {
		if geometry.Distance(g.User.PO.Center, g.Stars[i].PO.Center) < float32(g.Stars[i].PO.Radius+g.User.PO.Radius/2) {
			g.State = GameStateGameOverCrashed
			g.updateScore()
			return
		}

		if g.Enemy != nil {
			if geometry.Distance(g.Enemy.PO.Center, g.Stars[i].PO.Center) < float32(g.Stars[i].PO.Radius+g.Enemy.PO.Radius/2) {
				g.State = GameStateGameOverVictory
				g.updateScore()
				return
			}
		}
	}

	// Бомбы врезаются в корабли, планеты и покидают границы экрана
	for i := 0; i < len(g.Bombs); i++ {
		if g.Bombs[i] == nil {
			continue
		}

		// Бомбы вышли за пределы экрана
		if g.IsOutOfScene(g.Bombs[i].PO.Center.PointXY(), 0) {
			g.Bombs[i] = nil
			continue
		}

		// Бомбы врезаются в нейтронную звёзду
		for j := 0; j < len(g.Stars); j++ {
			if geometry.Distance(g.Bombs[i].PO.Center, g.Stars[j].PO.Center) < float32(g.Stars[j].PO.Radius) {
				g.Bombs[i] = nil
				break
			}
		}
		if g.Bombs[i] == nil {
			continue
		}

		// Бомбы врезаются в корабль пользователя
		if geometry.IsInPolygon(g.Bombs[i].PO.Center, g.User.RotatedBodyPolygon()) {
			g.User.Damage++
			if g.User.Lives() <= 0 {
				g.State = GameStateGameOverCrashed
				g.updateScore()
				return
			}

			g.Bombs[i] = nil
			continue
		}

		// Бомбы врезаются в корабль AI
		if g.Enemy != nil {
			if geometry.IsInPolygon(g.Bombs[i].PO.Center, g.Enemy.RotatedBodyPolygon()) {
				g.Enemy.Damage++
				if g.Enemy.Lives() <= 0 {
					g.State = GameStateGameOverVictory
					g.updateScore()
					return
				}

				g.Bombs[i] = nil
				continue
			}
		}
	}

	// Столкновение кораблей
	if g.Enemy != nil {
		if geometry.Distance(g.User.PO.Center, g.Enemy.PO.Center) < float32(g.Enemy.PO.Radius) {
			g.State = GameStateGameOverKamikaze
			g.updateScore()
			return
		}
	}

	// Корабль вышел за пределы экрана
	if g.IsOutOfScene(g.User.PO.Center.PointXY(), 200) {
		g.State = GameStateGameOverShipLost
		g.updateScore()
		return
	}
}

func (g *Game) AddBomb(bomb *objects.Bomb) {
	if bomb == nil {
		return
	}
	// Кольцевой буфер с выпущенными бомбами
	for i := 0; i < len(g.Bombs); i++ {
		if g.Bombs[i] == nil {
			g.Bombs[i] = bomb
			return
		}
	}
}

func (g *Game) IsOutOfScene(point geometry.PointXY, threshold int) bool {
	return point.X < -threshold || point.Y < -threshold || point.X > g.gw.Width+threshold || point.Y > g.gw.Height+threshold
}

func (g *Game) updateScore() {
	if g.freestyleMode || g.State == GameStateGameOverKamikaze {
		return
	}

	if g.State == GameStateGameOverVictory {
		g.UserScore++
	}

	if g.State == GameStateGameOverCrashed || g.State == GameStateGameOverShipLost {
		g.EnemyScore++
	}
}

func (g *Game) Update60Fps(gw *graphics.GLWindow) {
	switch g.State {
	case GameIsNotInitialized:
		g.gw = gw

		// Создаем космос
		g.Space.InitSpace(100, gw.Width, gw.Height)

		// Загрузим в память шрифты, с которыми будем работать
		gw.FontsBuffer.Init()
		gw.FontsBuffer.Load(embedded.FontOpenGostTypeA, 120, embedded.GetFontReader(embedded.FontOpenGostTypeA))
		gw.FontsBuffer.Load(embedded.FontPrimeMono, 54, embedded.GetFontReader(embedded.FontPrimeMono))
		gw.FontsBuffer.Load(embedded.FontPrimeMono, 24, embedded.GetFontReader(embedded.FontPrimeMono))
		gw.FontsBuffer.Load(embedded.FontPrimeMono, 14, embedded.GetFontReader(embedded.FontPrimeMono))

		g.State = GameStateMainMenu
		break
	case GameStateMainMenu:
		screenMainMenu(g)
		break
	case GameStateControlSettings:
		screenControlSettings(g)
		break
	case GameStatePlay:
		screenPlay(g)
		break
	case GameStateGameOverCrashed, GameStateGameOverShipLost, GameStateGameOverVictory, GameStateGameOverKamikaze:
		screenGameOver(g)
		break
	}
}

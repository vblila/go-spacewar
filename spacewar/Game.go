package spacewar

import (
	"github.com/gonutz/prototype/draw"
	"math"
	"math/rand"
)

const (
	GameStateIsNotInitialized = iota
	GameStateInProgress       = iota
	GameStateGameOverCrashed  = iota
	GameStateGameOverShipLost = iota

	ShipGunPower     = 4
	BombRadius       = 2
	ScreenBombsLimit = 10

	ShipRadius      = 36
	MinPlanetRadius = 15
	MaxPlanetRadius = 100
)

type Game struct {
	Width, Height int

	window draw.Window

	State   uint8
	Space   Space
	Gravity Gravity
	Me      Ship
	Planets []Planet

	// Кольцевой буфер с выпущенными бомбами
	Bombs [ScreenBombsLimit]*Bomb
}

func NewGame(width int, height int, window draw.Window) Game {
	return Game{Width: width, Height: height, window: window}
}

func (g *Game) Restart() {
	g.Space = Space{}
	g.Gravity = Gravity{Power: 0.1}
	g.Planets = make([]Planet, 1+rand.Intn(2))
	g.Bombs = [ScreenBombsLimit]*Bomb{}
	g.Me = Ship{}

	// Создаем космос
	g.Space.InitSpace(100, g.Width, g.Height)

	// Создаем планеты случайным образом
	for i := 0; i < len(g.Planets); i++ {
		g.Planets[i].Radius = MinPlanetRadius + rand.Intn(MaxPlanetRadius-ShipRadius/2+1)

		centerOffsetX := -rand.Intn(g.Width/4) + rand.Intn(g.Width/4)
		centerOffsetY := -rand.Intn(g.Width/4) + rand.Intn(g.Width/4)

		g.Planets[i].Center.X, g.Planets[i].Center.Y = g.Width/2+centerOffsetX, g.Height/2+centerOffsetY
	}

	// Создаем корабль
	g.Me.Center.X, g.Me.Center.Y = (float32)(g.Width/6), (float32)(g.Height/6)
	g.Me.Power = 0.05

	g.State = GameStateInProgress
}

func (g *Game) UpdatePositions() {
	g.Me.UpdatePosition()
	for i := 0; i < len(g.Bombs); i++ {
		if g.Bombs[i] != nil {
			g.Bombs[i].UpdatePosition()
		}
	}
}

func (g *Game) GravityInfluence() {
	meXY := g.Me.Center.GetPointXY()

	for i := 0; i < len(g.Planets); i++ {
		var deltaSpeedX, deltaSpeedY float32

		// Планеты притягивают корабль
		deltaSpeedX, deltaSpeedY = g.Gravity.Influence(g.Planets[i].Radius, g.Planets[i].Center, meXY)
		g.Me.SpeedX += deltaSpeedX
		g.Me.SpeedY += deltaSpeedY

		// Планеты притягивают бомбы
		for j := 0; j < len(g.Bombs); j++ {
			bomb := g.Bombs[j]
			if bomb == nil {
				continue
			}

			deltaSpeedX, deltaSpeedY = g.Gravity.Influence(g.Planets[i].Radius, g.Planets[i].Center, bomb.Center.GetPointXY())
			bomb.SpeedX += deltaSpeedX
			bomb.SpeedY += deltaSpeedY
		}
	}
}

func (g *Game) CheckBounds() {
	meCenterXY := g.Me.Center.GetPointXY()

	// Корабль врезается в планету
	for i := 0; i < len(g.Planets); i++ {
		if math.Hypot((float64)(meCenterXY.X-g.Planets[i].Center.X), (float64)(meCenterXY.Y-g.Planets[i].Center.Y)) < (float64)(g.Planets[i].Radius+ShipRadius/2) {
			g.State = GameStateGameOverCrashed
		}
	}

	// Бомбы врезается в планету
	for i := 0; i < len(g.Planets); i++ {
		for j := 0; j < len(g.Bombs); j++ {
			if g.Bombs[j] == nil {
				continue
			}
			bombXY := g.Bombs[j].Center.GetPointXY()
			if math.Hypot((float64)(bombXY.X-g.Planets[i].Center.X), (float64)(bombXY.Y-g.Planets[i].Center.Y)) < (float64)(g.Planets[i].Radius) {
				g.Bombs[j] = nil
			}
		}
	}

	// Корабль вышел за пределы экрана
	if meCenterXY.X < -200 || meCenterXY.Y < -200 || meCenterXY.X > g.Width+200 || meCenterXY.Y > g.Height+200 {
		g.State = GameStateGameOverShipLost
	}

	// Бомбы вышли за пределы экрана
	for i := 0; i < len(g.Bombs); i++ {
		if g.Bombs[i] == nil {
			continue
		}
		bombXY := g.Bombs[i].Center.GetPointXY()
		if bombXY.X < 0 || bombXY.Y < 0 || bombXY.X > g.Width || bombXY.Y > g.Height {
			g.Bombs[i] = nil
		}
	}
}

func (g *Game) AddBomb(bomb *Bomb) {
	if bomb == nil {
		return
	}

	for i := 0; i < len(g.Bombs); i++ {
		if g.Bombs[i] == nil {
			g.Bombs[i] = bomb
			return
		}
	}
}

func (g *Game) Render() {
	g.Space.Render(g.window)
	g.Me.Render(g.window)

	for i := 0; i < len(g.Bombs); i++ {
		if g.Bombs[i] != nil {
			g.Bombs[i].Render(g.window)
		}
	}

	for i := 0; i < len(g.Planets); i++ {
		g.Planets[i].Render(g.window)
	}
}

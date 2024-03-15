package config

const (
	MinGameSpeedLevel = 1
	MaxGameSpeedLevel = 5
	BombRadius        = 2
	BombsLimit        = 40
	ShipRadius        = 36
	MaxStars          = 2
	MinStarRadius     = 35
	MaxStarRadius     = 100
	MaxShipDamage     = 5
)

var gameSpeed float32
var gameSpeedLevel int

func init() {
	gameSpeedLevel = 1
	SetNextGameSpeedLevel()
}

func GetGameSpeedLevel() int {
	return gameSpeedLevel
}

func SetNextGameSpeedLevel() {
	gameSpeedLevel++
	if gameSpeedLevel > MaxGameSpeedLevel {
		gameSpeedLevel = MinGameSpeedLevel
	}

	gameSpeed = float32(gameSpeedLevel)*0.1 + 0.1
}

func GetGameSpeed() float32 {
	return gameSpeed
}

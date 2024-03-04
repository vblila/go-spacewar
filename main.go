package main

import (
	"github.com/gonutz/prototype/draw"
	"log"
	"math/rand"
	"os"
	"spacewar/spacewar"
	"time"
)

var game spacewar.Game

func main() {
	rand.Seed(time.Now().UnixNano())

	e := draw.RunWindow("SPACE WAR", 640, 480, update60Fps)
	if e != nil {
		log.Panic(e.Error())
	}
}

func initGame(window draw.Window) {
	window.ShowCursor(false)
	window.SetFullscreen(true)

	videoMode := spacewar.GetVideoMode(window)
	game = spacewar.NewGame(videoMode.Width, videoMode.Height, window)
	game.Restart()
}

func update60Fps(window draw.Window) {
	if game.State == spacewar.GameStateIsNotInitialized {
		initGame(window)
	}

	if window.IsKeyDown(draw.KeyLeft) {
		game.Me.RotateLeft()
	}
	if window.IsKeyDown(draw.KeyRight) {
		game.Me.RotateRight()
	}
	if window.IsKeyDown(draw.KeyUp) {
		game.Me.Accelerate()
	}
	if window.IsKeyDown(draw.KeySpace) {
		game.AddBomb(game.Me.Fire())
	}

	if window.WasKeyPressed(draw.KeyEscape) {
		os.Exit(0)
	}

	if game.State == spacewar.GameStateInProgress {
		game.GravityInfluence()
		game.UpdatePositions()
		game.Render()
		game.CheckBounds()
	} else if game.State == spacewar.GameStateGameOverCrashed {
		window.DrawScaledText("GAME OVER", game.Width/2-170, game.Height/2-30, 4, draw.White)
		window.DrawScaledText("Press ENTER for restart", game.Width/2-320, game.Height/2+70, 3, draw.White)

		if window.WasKeyPressed(draw.KeyEnter) {
			game.Restart()
		}
	} else if game.State == spacewar.GameStateGameOverShipLost {
		window.DrawScaledText("LOST IN SPACE", game.Width/2-230, game.Height/2-30, 4, draw.White)
		window.DrawScaledText("Press ENTER for restart", game.Width/2-320, game.Height/2+70, 3, draw.White)

		if window.WasKeyPressed(draw.KeyEnter) {
			game.Restart()
		}
	}
}

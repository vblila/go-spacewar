package main

import (
	"math/rand"
	"runtime"
	"spacewar/gameplay"
	"spacewar/graphics"
	"time"
)

func init() {
	// Проект использует GLFW и GL библиотеки, которые используют CGo,
	// поэтому необходимо работать с ними в одном и том же потоке
	runtime.LockOSThread()
}

func main() {
	rand.Seed(time.Now().UnixNano())

	game := gameplay.Game{}

	gw := graphics.GLWindow{Color: graphics.Color{R: 0.4118, G: 0.7255, B: 0.8039, A: 1}}
	gw.Init(640, 480, "SPACEWAR 2024", game.Update60Fps)
}

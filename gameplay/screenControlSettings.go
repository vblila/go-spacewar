package gameplay

import (
	"spacewar/embedded"
	"spacewar/graphics"
)

func screenControlSettings(g *Game) {
	g.Space.Render(g.gw)
	g.Space.Fly()

	startY := g.gw.Height/2 - 320

	g.gw.DrawTextWithAlignCenter(startY, "SPACEWAR 2024", embedded.FontOpenGostTypeA, 120, g.gw.Color)
	g.gw.DrawTextWithAlignCenter(startY+100, "INSPIRED BY THE VIDEO GAME DEVELOPED BY STEVE RUSSELL IN 1962", embedded.FontPrimeMono, 14, g.gw.Color)

	controllerLeft := g.gw.DrawTextWithAlignCenter(startY+190, "UP    - ACCELERATE", embedded.FontPrimeMono, 24, g.gw.Color)
	g.gw.DrawText(controllerLeft, startY+220, "LEFT  - ROTATE LEFT", embedded.FontPrimeMono, 24, g.gw.Color)
	g.gw.DrawText(controllerLeft, startY+250, "RIGHT - ROTATE RIGHT", embedded.FontPrimeMono, 24, g.gw.Color)
	g.gw.DrawText(controllerLeft, startY+280, "SPACE - FIRE", embedded.FontPrimeMono, 24, g.gw.Color)

	g.gw.DrawTextWithAlignCenter(startY+375, "PRESS ESCAPE FOR MENU", embedded.FontPrimeMono, 24, g.gw.Color)

	if g.gw.WasKeyPressed(graphics.KeyEscape) {
		g.State = GameStateMainMenu
	}
}

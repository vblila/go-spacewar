package gameplay

import (
	"os"
	"spacewar/config"
	"spacewar/embedded"
	"spacewar/graphics"
	"strconv"
)

func screenMainMenu(g *Game) {
	g.Space.Render(g.gw)
	g.Space.Fly()

	startY := g.gw.Height/2 - 320

	g.gw.DrawTextWithAlignCenter(startY, "SPACEWAR 2024", embedded.FontOpenGostTypeA, 120, g.gw.Color)
	g.gw.DrawTextWithAlignCenter(startY+100, "INSPIRED BY THE VIDEO GAME DEVELOPED BY STEVE RUSSELL IN 1962", embedded.FontPrimeMono, 14, g.gw.Color)

	menuLeft := g.gw.DrawTextWithAlignCenter(startY+190, "PRESS 1: FREE FLIGHT     ", embedded.FontPrimeMono, 24, g.gw.Color)
	g.gw.DrawText(menuLeft, startY+220, "PRESS 2: BATTLE WITH AI", embedded.FontPrimeMono, 24, g.gw.Color)
	g.gw.DrawText(menuLeft, startY+250, "PRESS 3: CONTROL SETTINGS", embedded.FontPrimeMono, 24, g.gw.Color)
	g.gw.DrawText(menuLeft, startY+280, "PRESS 4: GAMEPLAY SPEED "+strconv.Itoa(config.GetGameSpeedLevel()), embedded.FontPrimeMono, 24, g.gw.Color)

	g.gw.DrawTextWithAlignCenter(startY+375, "PRESS ESCAPE FOR EXIT", embedded.FontPrimeMono, 24, g.gw.Color)

	g.gw.DrawTextWithAlignCenter(startY+500, "DEVELOPED BY VLADIMIR LILA", embedded.FontPrimeMono, 14, g.gw.Color)
	g.gw.DrawTextWithAlignCenter(startY+525, "https://github.com/vblila/go-spacewar", embedded.FontPrimeMono, 14, g.gw.Color)

	if g.gw.IsKeyDown(graphics.Key1) {
		g.start(true)
	} else if g.gw.IsKeyDown(graphics.Key2) {
		g.start(false)
	} else if g.gw.IsKeyDown(graphics.Key3) {
		g.State = GameStateControlSettings
	} else if g.gw.WasKeyPressed(graphics.Key4) {
		config.SetNextGameSpeedLevel()
	} else if g.gw.WasKeyPressed(graphics.KeyEscape) {
		os.Exit(0)
	}
}

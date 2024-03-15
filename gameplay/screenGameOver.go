package gameplay

import (
	"spacewar/embedded"
	"spacewar/graphics"
)

func screenGameOver(g *Game) {
	g.Space.Render(g.gw)
	g.Space.Fly()

	if g.State == GameStateGameOverCrashed || g.State == GameStateGameOverKamikaze {
		g.gw.DrawTextWithAlignCenter(g.gw.Height/2-100, "WAS DESTROYED", embedded.FontPrimeMono, 54, g.gw.Color)
	}

	if g.State == GameStateGameOverShipLost {
		g.gw.DrawTextWithAlignCenter(g.gw.Height/2-100, "LOST IN SPACE", embedded.FontPrimeMono, 54, g.gw.Color)
	}

	if g.State == GameStateGameOverVictory {
		g.gw.DrawTextWithAlignCenter(g.gw.Height/2-100, "VICTORY", embedded.FontPrimeMono, 54, g.gw.Color)
	}

	g.gw.DrawTextWithAlignCenter(g.gw.Height/2-20, "PRESS ESCAPE FOR MENU", embedded.FontPrimeMono, 24, g.gw.Color)

	if g.gw.WasKeyPressed(graphics.KeyEscape) {
		g.State = GameStateMainMenu
	}
}

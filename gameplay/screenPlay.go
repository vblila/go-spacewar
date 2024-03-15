package gameplay

import (
	"fmt"
	"spacewar/embedded"
	"spacewar/graphics"
	"strings"
)

func screenPlay(g *Game) {
	if g.gw.WasKeyPressed(graphics.KeyEscape) {
		g.State = GameStateMainMenu
		return
	}

	if g.gw.IsKeyDown(graphics.KeyLeft) {
		g.User.RotateLeft()
	}
	if g.gw.IsKeyDown(graphics.KeyRight) {
		g.User.RotateRight()
	}
	if g.gw.IsKeyDown(graphics.KeyUp) {
		g.User.Accelerate()
	}
	if g.gw.IsKeyDown(graphics.KeySpace) {
		g.AddBomb(g.User.Fire())
	}

	// 1. Просчитываем и применяем влияние гравитации
	g.doGravityInfluence()

	// 2. AI продумавает действия на каждой итерации
	if g.Ai != nil {
		g.Ai.DoIterationLogic()
	}

	// 3. Обновляем положение всех объектов в игре
	g.updatePositions()

	// 4. Контроль границы всех объектов в игре
	g.checkBounds()

	// 5. Рисуем все объекты игрового процесса с учётом новых позиций
	g.Space.Render(g.gw)
	g.User.Render(g.gw)
	if g.Enemy != nil {
		g.Enemy.Render(g.gw)
	}

	for i := 0; i < len(g.Bombs); i++ {
		if g.Bombs[i] != nil {
			g.Bombs[i].Render(g.gw)
		}
	}

	for i := 0; i < len(g.Stars); i++ {
		g.Stars[i].Render(g.gw)
	}

	if !g.freestyleMode {
		g.gw.DrawText(10, 5, strings.Repeat("O", g.User.Lives()), embedded.FontPrimeMono, 24, g.gw.Color)

		enemyLivesText := strings.Repeat("O", g.Enemy.Lives())
		enemyLivesWidth, _ := g.gw.TextMetrics(enemyLivesText, embedded.FontPrimeMono, 24)
		g.gw.DrawText(g.gw.Width-enemyLivesWidth-10, 5, enemyLivesText, embedded.FontPrimeMono, 24, g.gw.Color)

		g.gw.DrawTextWithAlignCenter(5, fmt.Sprintf("%d:%d", g.UserScore, g.EnemyScore), embedded.FontPrimeMono, 24, g.gw.Color)
	}
}

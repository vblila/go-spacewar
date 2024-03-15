package graphics

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"math"
	"spacewar/geometry"
	"time"
)

type Update60FpsFunction func(gw *GLWindow)

const (
	Key1      = int(glfw.Key1)
	Key2      = int(glfw.Key2)
	Key3      = int(glfw.Key3)
	Key4      = int(glfw.Key4)
	KeyEscape = int(glfw.KeyEscape)

	KeyUp    = int(glfw.KeyUp)
	KeyLeft  = int(glfw.KeyLeft)
	KeyRight = int(glfw.KeyRight)
	KeySpace = int(glfw.KeySpace)
)

type GLWindow struct {
	w *glfw.Window

	FontsBuffer FontsBuffer
	pressedKey  glfw.Key

	Color Color

	Width  int
	Height int
}

func (gw *GLWindow) Init(width int, height int, title string, update Update60FpsFunction) {
	var err error
	err = glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	gw.Width = width
	gw.Height = height

	gw.w, err = glfw.CreateWindow(gw.Width, gw.Height, title, nil, nil)
	if err != nil {
		panic(err)
	}

	// Для Windows очень важно сперва сделать контекст окна текущим, и только потом вызывать gl.Init()
	gw.w.MakeContextCurrent()

	err = gl.Init()
	if err != nil {
		panic(err)
	}

	gw.w.SetInputMode(glfw.CursorMode, glfw.CursorHidden)

	// Устанавливать Viewport только через колбек на изменение размеров окна
	gw.w.SetSizeCallback(func(_ *glfw.Window, width, height int) {
		gw.Width, gw.Height = width, height
		gl.MatrixMode(gl.PROJECTION)
		gl.LoadIdentity()
		gl.Ortho(0, float64(width), float64(height), 0, -1, 1)
		gl.Viewport(0, 0, int32(width), int32(height))
		gl.MatrixMode(gl.MODELVIEW)

		// Включаем сглаживание
		gl.LineWidth(1.5)
		gl.PointSize(2.5)
		gl.Enable(gl.POINT_SMOOTH)
		gl.Hint(gl.POINT_SMOOTH_HINT, gl.NICEST)
		gl.Enable(gl.BLEND)
		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

		// Отключаем VSync
		glfw.SwapInterval(0)

	})

	gw.w.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if action == glfw.Release {
			gw.pressedKey = key
		}
	})

	gw.setFullscreen()

	lastUpdateTime := float64(0)
	for !gw.w.ShouldClose() {
		now := glfw.GetTime()

		// Поддерживаем 60 FPS, 1 фрейм в секуду 1/60 = 0.016666667
		if now-lastUpdateTime >= 0.016666667 {
			lastUpdateTime = now

			gl.ClearColor(0, 0, 0, 1)
			gl.Clear(gl.COLOR_BUFFER_BIT)
			update(gw)
			glfw.PollEvents()
			gw.w.SwapBuffers()
		} else {
			time.Sleep(time.Millisecond)
		}
	}
}

func (gw *GLWindow) setFullscreen() {
	monitor := gw.getMonitor(gw.w.GetPos())
	mode := monitor.GetVideoMode()
	gw.w.SetMonitor(monitor, 0, 0, mode.Width, mode.Height, 60)
	gw.Width, gw.Height = mode.Width, mode.Height
}

func (gw *GLWindow) getMonitor(winX, winY int) *glfw.Monitor {
	for _, m := range glfw.GetMonitors() {
		x, y, w, h := m.GetWorkarea()
		if x <= winX && winX < x+w && y <= winY && winY < y+h {
			return m
		}
	}
	return glfw.GetPrimaryMonitor()
}

func (gw *GLWindow) IsKeyDown(key int) bool {
	return gw.w.GetKey(glfw.Key(key)) == glfw.Press
}

func (gw *GLWindow) WasKeyPressed(key int) bool {
	pressedKey := gw.pressedKey
	if pressedKey == glfw.Key(key) {
		gw.pressedKey = 1
		return true
	}

	return false
}

// DrawPolygon Рисуем по точкам
func (gw *GLWindow) DrawPolygon(points []geometry.PointFloatXY, color Color) {
	for i := 0; i < len(points)-1; i++ {
		gw.DrawLine(points[i].X, points[i].Y, points[i+1].X, points[i+1].Y, color)
	}
}

func (gw *GLWindow) DrawCircle(x, y, radius float32, color Color) {
	theta := 0.017453293 // math.Pi * 2 / 360
	tanTheta := float32(math.Tan(theta))
	cosTheta := float32(math.Cos(theta))

	var i, j float32
	i = radius
	j = 0

	gl.Begin(gl.POINTS)
	gl.Color4f(color.R, color.G, color.B, color.A)

	for index := 0; index < 360; index++ {
		gl.Vertex2f(x+i, y+j)

		ti, tj := -j, i

		i += ti * tanTheta
		i *= cosTheta
		j += tj * tanTheta
		j *= cosTheta
	}

	gl.End()
}

func (gw *GLWindow) FillCircle(x, y, r float32, color Color) {
	gl.Begin(gl.TRIANGLE_FAN)
	gl.Color4f(color.R, color.G, color.B, color.A)

	for i := 0.0; i <= 360; i++ {
		gl.Vertex2f(
			r*float32(math.Cos(math.Pi*i/180.0))+x,
			r*float32(math.Sin(math.Pi*i/180.0))+y,
		)
	}
	gl.End()
}

func (gw *GLWindow) DrawPoint(x, y float32, color Color) {
	gl.Begin(gl.POINTS)
	gl.Color4f(color.R, color.G, color.B, color.A)
	gl.Vertex2f(x+0.5, y+0.5)
	gl.End()
}

func (gw *GLWindow) DrawLine(x, y, x2, y2 float32, color Color) {
	if x == x2 && y == y2 {
		gw.DrawPoint(x, y, color)
		return
	}

	gl.Begin(gl.LINES)
	gl.Color4f(color.R, color.G, color.B, color.A)
	gl.Vertex2f(x, y)
	gl.Vertex2f(x2, y2)
	gl.End()
}

func (gw *GLWindow) DrawText(x, y int, text string, fontName string, fontScale int, color Color) {
	font := gw.FontsBuffer.Get(fontName, fontScale)
	gl.Color4f(color.R, color.G, color.B, color.A)
	font.Printf(float32(x), float32(y), text)
}

func (gw *GLWindow) TextMetrics(text string, fontName string, fontScale int) (width, height int) {
	font := gw.FontsBuffer.Get(fontName, fontScale)
	return font.Metrics(text)
}

func (gw *GLWindow) DrawTextWithAlignCenter(y int, text string, fontName string, fontScale int, color Color) (startX int) {
	font := gw.FontsBuffer.Get(fontName, fontScale)
	width, _ := font.Metrics(text)

	startX = gw.Width/2 - width/2

	gl.Color4f(color.R, color.G, color.B, color.A)
	font.Printf(float32(startX), float32(y), text)

	return
}

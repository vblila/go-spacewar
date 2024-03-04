package spacewar

import (
	"github.com/gonutz/glfw/v3.3/glfw"
	"github.com/gonutz/prototype/draw"
)

// DrawPolygon Рисуем по точкам. Последнюю соединяем с первой
func DrawPolygon(points []PointFloatXY, window draw.Window) {
	for i := 0; i < len(points); i++ {
		var p, nextP PointXY

		p = points[i].GetPointXY()
		if i == len(points)-1 {
			nextP = points[0].GetPointXY()
		} else {
			nextP = points[i+1].GetPointXY()
		}

		window.DrawLine(p.X, p.Y, nextP.X, nextP.Y, draw.White)
	}
}

func GetVideoMode(window draw.Window) *glfw.VidMode {
	winX, winY := window.Size()
	for _, m := range glfw.GetMonitors() {
		x, y, w, h := m.GetWorkarea()
		if x <= winX && winX < x+w && y <= winY && winY < y+h {
			return m.GetVideoMode()
		}
	}
	return glfw.GetPrimaryMonitor().GetVideoMode()
}

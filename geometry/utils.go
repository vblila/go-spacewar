package geometry

import "math"

// NormalizeAngle2Pi normalize to [0; 2*pi]
func NormalizeAngle2Pi(angle float32) float32 {
	normalized := angle - (float32)((math.Pi*2)*math.Floor((float64)(angle/(math.Pi*2))))
	if normalized < 0 {
		return 2*math.Pi - normalized
	}

	return normalized
}

// VectorAngle [0; 2*pi]
func VectorAngle(object PointFloatXY, target PointFloatXY) float32 {
	var deltaX, deltaY float32
	deltaX, deltaY = target.X-object.X, target.Y-object.Y

	if deltaY == 0 {
		if deltaX < 0 {
			return 3 * math.Pi / 2
		} else {
			return math.Pi / 2
		}
	}

	if deltaX == 0 {
		if deltaY < 0 {
			return 0
		} else {
			return math.Pi
		}
	}

	var vectorAngle float64
	vectorAngleTg := math.Abs((float64)(deltaY)) / (float64)(deltaX)

	if deltaX > 0 && deltaY > 0 {
		vectorAngle = math.Pi/2 + math.Atan(vectorAngleTg)
		return (float32)(vectorAngle)
	}

	if deltaX < 0 && deltaY > 0 {
		vectorAngle = 3*math.Pi/2 + math.Atan(vectorAngleTg)
		return (float32)(vectorAngle)
	}

	if deltaX > 0 && deltaY < 0 {
		vectorAngle = math.Pi/2 - math.Atan(vectorAngleTg)
		return (float32)(vectorAngle)
	}

	if deltaX < 0 && deltaY < 0 {
		vectorAngle = 3*math.Pi/2 - math.Atan(vectorAngleTg)
		return (float32)(vectorAngle)
	}

	return (float32)(vectorAngle)
}

func Distance(a PointFloatXY, b PointFloatXY) float32 {
	return (float32)(math.Hypot((float64)(a.X-b.X), (float64)(a.Y-b.Y)))
}

// NearestRotateAngle return: 1 - target1Angle, 2 - target2Angle
func NearestRotateAngle(objectAngle float32, target1Angle float32, target2Angle float32) uint8 {
	target1AbsDiff := math.Abs(float64(objectAngle - target1Angle))
	target1AbsDiffMinAngle := math.Min(2*math.Pi-target1AbsDiff, target1AbsDiff)

	target2AbsDiff := math.Abs(float64(objectAngle - target2Angle))
	target2AbsDiffMinAngle := math.Min(2*math.Pi-target2AbsDiff, target2AbsDiff)

	if target1AbsDiffMinAngle < target2AbsDiffMinAngle {
		return 1
	}

	return 2
}

func IsInPolygon(point PointFloatXY, polygon []PointFloatXY) bool {
	result := false

	max := len(polygon) - 1
	for i, j := 0, max; i <= max; i++ {
		if (polygon[i].Y > point.Y) != (polygon[j].Y > point.Y) &&
			(point.X < (polygon[j].X-polygon[i].X)*(point.Y-polygon[i].Y)/(polygon[j].Y-polygon[i].Y)+polygon[i].X) {
			result = !result
		}
		j = i
	}

	return result
}

package game

import (
	"math"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
)

type line struct {
	x1, y1, x2, y2 float64
}

func (l line) Angle() float64 {
	return math.Atan2(l.y2-l.y1, l.x2-l.x1)
}

func (l line) Dist() float64 {
	return math.Sqrt(math.Pow(l.x2-l.x1, 2) + (math.Pow(l.y2-l.y1, 2)))
}

type object struct {
	walls []line
}

func (o object) points() [][2]float64 {
	// Get one of the endpoints for all segments,
	// + the startpoint of the first one, for non-closed paths
	var points [][2]float64
	for _, wall := range o.walls {
		points = append(points, [2]float64{wall.x2, wall.y2})
	}
	p := [2]float64{o.walls[0].x1, o.walls[0].y1}
	if p[0] != points[len(points)-1][0] && p[1] != points[len(points)-1][1] {
		points = append(points, [2]float64{o.walls[0].x1, o.walls[0].y1})
	}
	return points
}

func newRay(x, y, length, angle float64) line {
	return line{
		x1: x,
		y1: y,
		x2: x + length*math.Cos(angle),
		y2: y + length*math.Sin(angle),
	}
}

// intersection calculates the intersection of given two lines.
func intersection(l1, l2 line) (float64, float64, bool) {
	// https://en.wikipedia.org/wiki/Line%E2%80%93line_intersection#Given_two_points_on_each_line
	denom := (l1.x1-l1.x2)*(l2.y1-l2.y2) - (l1.y1-l1.y2)*(l2.x1-l2.x2)
	tNum := (l1.x1-l2.x1)*(l2.y1-l2.y2) - (l1.y1-l2.y1)*(l2.x1-l2.x2)
	uNum := -((l1.x1-l1.x2)*(l1.y1-l2.y1) - (l1.y1-l1.y2)*(l1.x1-l2.x1))

	if denom == 0 {
		return 0, 0, false
	}

	t := tNum / denom
	if t > 1 || t < 0 {
		return 0, 0, false
	}

	u := uNum / denom
	if u > 1 || u < 0 {
		return 0, 0, false
	}

	x := l1.x1 + t*(l1.x2-l1.x1)
	y := l1.y1 + t*(l1.y2-l1.y1)
	return x, y, true
}

// rayCasting returns a slice of line originating from point cx, cy and intersecting with objects
func rayCasting(cx, cy float64, objects []object) []line {
	const rayLength = 1000 // something large enough to reach all objects

	var rays []line
	for _, obj := range objects {
		// Cast two rays per point
		for _, p := range obj.points() {
			l := line{cx, cy, p[0], p[1]}
			angle := l.Angle()

			for _, offset := range []float64{-0.005, 0.005} {
				points := [][2]float64{}
				ray := newRay(cx, cy, rayLength, angle+offset)

				// Unpack all objects
				for _, o := range objects {
					for _, wall := range o.walls {
						if px, py, ok := intersection(ray, wall); ok {
							points = append(points, [2]float64{px, py})
						}
					}
				}

				// Find the point closest to start of ray
				min := math.Inf(1)
				minI := -1
				for i, p := range points {
					d2 := (cx-p[0])*(cx-p[0]) + (cy-p[1])*(cy-p[1])
					if d2 < min {
						min = d2
						minI = i
					}
				}
				rays = append(rays, line{cx, cy, points[minI][0], points[minI][1]})
			}
		}
	}

	// Sort rays based on angle, otherwise light triangles will not come out right
	sort.Slice(rays, func(i int, j int) bool {
		return rays[i].Angle() < rays[j].Angle()
	})
	return rays
}

func rayVertices(x1, y1, x2, y2, x3, y3 float64) []ebiten.Vertex {
	return []ebiten.Vertex{
		{DstX: float32(x1), DstY: float32(y1), SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
		{DstX: float32(x2), DstY: float32(y2), SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
		{DstX: float32(x3), DstY: float32(y3), SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
	}
}

func rect(x, y, w, h float64) []line {
	return []line{
		{x, y, x, y + h},
		{x, y + h, x + w, y + h},
		{x + w, y + h, x + w, y},
		{x + w, y, x, y},
	}
}

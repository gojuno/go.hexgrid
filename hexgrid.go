package hexgrid

import (
	"math"

	morton "github.com/gojuno/go.morton"
)

type Point struct {
	x float64
	y float64
}

type Hex struct {
	q int64
	r int64
}

type FractionalHex struct {
	q float64
	r float64
}

type Orientation struct {
	f          [4]float64
	b          [4]float64
	startAngle float64
	sinuses    [6]float64
	cosinuses  [6]float64
}

type Grid struct {
	orientation Orientation
	origin      Point
	size        Point
	mort        *morton.Morton64
}

type Region struct {
	grid   *Grid
	hexes  []Hex
	lookup map[int64]int
}

var PointyOrientation = Orientation{
	f:          [4]float64{math.Sqrt(3.0), math.Sqrt(3.0) / 2.0, 0.0, 3.0 / 2.0},
	b:          [4]float64{math.Sqrt(3.0) / 3.0, -1.0 / 3.0, 0.0, 2.0 / 3.0},
	startAngle: 0.5}

var FlatOrientation = Orientation{
	f:          [4]float64{3.0 / 2.0, 0.0, math.Sqrt(3.0) / 2.0, math.Sqrt(3.0)},
	b:          [4]float64{2.0 / 3.0, 0.0, -1.0 / 3.0, math.Sqrt(3.0) / 3.0},
	startAngle: 0.0}

func init() {
	prehashAngles(&PointyOrientation)
	prehashAngles(&FlatOrientation)
}

func prehashAngles(orientation *Orientation) {
	for i := 0; i < 6; i++ {
		angle := 2.0 * math.Pi * (float64(i) + orientation.startAngle) / 6.0
		orientation.sinuses[i] = math.Sin(angle)
		orientation.cosinuses[i] = math.Cos(angle)
	}
}

func round(val float64) int64 {
	if val < 0 {
		return int64(val - 0.5)
	}
	return int64(val + 0.5)
}

func max(a int64, b int64) int64 {
	if a >= b {
		return a
	}
	return b
}

func min(a int64, b int64) int64 {
	if a <= b {
		return a
	}
	return b
}

func MakePoint(x float64, y float64) Point {
	return Point{x: x, y: y}
}

func (point Point) X() float64 {
	return point.x
}

func (point Point) Y() float64 {
	return point.y
}

func MakeHex(q int64, r int64) Hex {
	return Hex{q: q, r: r}
}

func (hex Hex) Q() int64 {
	return hex.q
}

func (hex Hex) R() int64 {
	return hex.r
}

func (hex Hex) S() int64 {
	return -(hex.q + hex.r)
}

func MakeFractionalHex(q float64, r float64) FractionalHex {
	return FractionalHex{q: q, r: r}
}

func (fhex FractionalHex) Q() float64 {
	return fhex.q
}

func (fhex FractionalHex) R() float64 {
	return fhex.r
}

func (fhex FractionalHex) S() float64 {
	return -(fhex.q + fhex.r)
}

func (fhex FractionalHex) ToHex() Hex {
	q := round(fhex.Q())
	r := round(fhex.R())
	s := round(fhex.S())
	qDiff := math.Abs(float64(q) - fhex.Q())
	rDiff := math.Abs(float64(r) - fhex.R())
	sDiff := math.Abs(float64(s) - fhex.S())

	if qDiff > rDiff && qDiff > sDiff {
		q = -(r + s)
	} else if rDiff > sDiff {
		r = -(q + s)
	}

	return Hex{q: q, r: r}
}

func MakeGrid(orientation Orientation, origin Point, size Point, mort *morton.Morton64) *Grid {
	return &Grid{orientation: orientation, origin: origin, size: size, mort: mort}
}

func (grid *Grid) HexToCode(hex Hex) int64 {
	return grid.mort.SPack(hex.Q(), hex.R())
}

func (grid *Grid) HexFromCode(code int64) Hex {
	qr := grid.mort.SUnpack(code)
	return MakeHex(qr[0], qr[1])
}

func (grid *Grid) HexAt(point Point) Hex {
	x := (point.X() - grid.origin.X()) / grid.size.X()
	y := (point.Y() - grid.origin.Y()) / grid.size.Y()
	q := grid.orientation.b[0]*x + grid.orientation.b[1]*y
	r := grid.orientation.b[2]*x + grid.orientation.b[3]*y
	return MakeFractionalHex(q, r).ToHex()
}

func (grid *Grid) HexCenter(hex Hex) Point {
	x := (grid.orientation.f[0]*float64(hex.Q())+grid.orientation.f[1]*float64(hex.R()))*grid.size.X() + grid.origin.X()
	y := (grid.orientation.f[2]*float64(hex.Q())+grid.orientation.f[3]*float64(hex.R()))*grid.size.Y() + grid.origin.Y()
	return MakePoint(x, y)
}

func (grid *Grid) HexCorners(hex Hex) [6]Point {
	var corners [6]Point
	center := grid.HexCenter(hex)
	for i := 0; i < 6; i++ {
		x := grid.size.X()*grid.orientation.cosinuses[i] + center.X()
		y := grid.size.Y()*grid.orientation.sinuses[i] + center.Y()
		corners[i] = MakePoint(x, y)
	}
	return corners
}

func (grid *Grid) HexNeighbors(hex Hex, layers int64) []Hex {
	total := (layers + 1) * layers * 3
	neighbors := make([]Hex, total)
	i := 0
	for q := -layers; q <= layers; q++ {
		r1 := max(-layers, -q-layers)
		r2 := min(layers, -q+layers)
		for r := r1; r <= r2; r++ {
			if q == 0 && r == 0 {
				continue
			}
			neighbors[i] = MakeHex(q+hex.Q(), r+hex.R())
			i += 1
		}
	}
	return neighbors
}

func pointInGeometry(geometry []Point, point Point) bool {
	contains := intersectsWithRaycast(point, geometry[len(geometry)-1], geometry[0])
	for i := 1; i < len(geometry); i++ {
		if intersectsWithRaycast(point, geometry[i-1], geometry[i]) {
			contains = !contains
		}
	}
	return contains
}

/* from https://github.com/kellydunn/golang-geo */
func intersectsWithRaycast(point Point, start Point, end Point) bool {
	if start.Y() > end.Y() {
		start, end = end, start
	}

	for point.Y() == start.Y() || point.Y() == end.Y() {
		newY := math.Nextafter(point.Y(), math.Inf(1))
		point = MakePoint(point.X(), newY)
	}

	if point.Y() < start.Y() || point.Y() > end.Y() {
		return false
	}

	if start.X() > end.X() {
		if point.X() > start.X() {
			return false
		}
		if point.X() < end.X() {
			return true
		}
	} else {
		if point.X() > end.X() {
			return false
		}
		if point.X() < start.X() {
			return true
		}
	}

	raySlope := (point.Y() - start.Y()) / (point.X() - start.X())
	diagSlope := (end.Y() - start.Y()) / (end.X() - start.X())

	return raySlope >= diagSlope
}

func (grid *Grid) MakeRegion(geometry []Point) *Region {
	if geometry[0] == geometry[len(geometry)-1] {
		geometry = geometry[:len(geometry)-1]
	}

	hex := grid.HexAt(geometry[0])
	q1 := hex.Q()
	q2 := hex.Q()
	r1 := hex.R()
	r2 := hex.R()

	for i := 1; i < len(geometry); i++ {
		hex = grid.HexAt(geometry[i])
		q1 = min(q1, hex.Q())
		q2 = max(q2, hex.Q())
		r1 = min(r1, hex.R())
		r2 = max(r2, hex.R())
	}

	q1 -= 1
	q2 += 1
	r1 -= 1
	r2 += 1

	hexes := make([]Hex, 0)

	for q := q1; q <= q2; q++ {
		for r := r1; r <= r2; r++ {
			hex := MakeHex(q, r)
			corners := grid.HexCorners(hex)
			add := false
			for c := 0; c < 6; c++ {
				if pointInGeometry(geometry, corners[c]) {
					add = true
					break
				}
			}

			if !add {
				for i := 0; i < len(geometry); i++ {
					if pointInGeometry(corners[:], geometry[i]) {
						add = true
						break
					}
				}
			}

			if add {
				hexes = append(hexes, hex)
			}
		}
	}

	lookup := make(map[int64]int)
	for i := 0; i < len(hexes); i++ {
		lookup[grid.HexToCode(hexes[i])] = i
	}

	return &Region{grid: grid, hexes: hexes, lookup: lookup}
}

func (region *Region) Hexes() []Hex {
	return region.hexes
}

func (region *Region) Contains(hex Hex) bool {
	_, contains := region.lookup[region.grid.HexToCode(hex)]
	return contains
}

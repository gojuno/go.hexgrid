package hexgrid

import (
	"math"
	"testing"

	morton "github.com/gojuno/go.morton"
)

func validateHex(t *testing.T, e Hex, r Hex) {
	if e.Q() != r.Q() || e.R() != r.R() {
		t.Errorf("expected hex{q: %d, r: %d} but got hex{q: %d, r: %d}", e.Q(), e.R(), r.Q(), r.R())
	}
}

func TestFlat(t *testing.T) {
	grid := MakeGrid(OrientationFlat, MakePoint(10, 20), MakePoint(20, 10), morton.Make64(2, 32))
	validateHex(t, MakeHex(0, 37), grid.HexAt(MakePoint(13, 666)))
	validateHex(t, MakeHex(22, -11), grid.HexAt(MakePoint(666, 13)))
	validateHex(t, MakeHex(-1, -39), grid.HexAt(MakePoint(-13, -666)))
	validateHex(t, MakeHex(-22, 9), grid.HexAt(MakePoint(-666, -13)))
}

func TestPointy(t *testing.T) {
	grid := MakeGrid(OrientationPointy, MakePoint(10, 20), MakePoint(20, 10), morton.Make64(2, 32))
	validateHex(t, MakeHex(-21, 43), grid.HexAt(MakePoint(13, 666)))
	validateHex(t, MakeHex(19, 0), grid.HexAt(MakePoint(666, 13)))
	validateHex(t, MakeHex(22, -46), grid.HexAt(MakePoint(-13, -666)))
	validateHex(t, MakeHex(-19, -2), grid.HexAt(MakePoint(-666, -13)))
}

func validatePoint(t *testing.T, e Point, r Point, precision float64) {
	if math.Abs(e.X()-r.X()) > precision || math.Abs(e.Y()-r.Y()) > precision {
		t.Errorf("expected point{x: %f, y: %f} but got point{x: %f, y: %f}", e.X(), e.Y(), r.X(), r.Y())
	}
}

func TestCoordinatesFlat(t *testing.T) {
	grid := MakeGrid(OrientationFlat, MakePoint(10, 20), MakePoint(20, 10), morton.Make64(2, 32))
	hex := grid.HexAt(MakePoint(666, 666))
	validatePoint(t, MakePoint(670.00000, 660.85880), grid.HexCenter(hex), 0.00001)
	expectedCorners := [6]Point{
		MakePoint(690.00000, 660.85880),
		MakePoint(680.00000, 669.51905),
		MakePoint(660.00000, 669.51905),
		MakePoint(650.00000, 660.85880),
		MakePoint(660.00000, 652.19854),
		MakePoint(680.00000, 652.19854)}
	corners := grid.HexCorners(hex)
	for i := 0; i < 6; i++ {
		validatePoint(t, expectedCorners[i], corners[i], 0.00001)
	}
}

func TestCoordinatesPointy(t *testing.T) {
	grid := MakeGrid(OrientationPointy, MakePoint(10, 20), MakePoint(20, 10), morton.Make64(2, 32))
	hex := grid.HexAt(MakePoint(666, 666))
	validatePoint(t, MakePoint(650.85880, 665.00000), grid.HexCenter(hex), 0.00001)
	expectedCorners := [6]Point{
		MakePoint(668.17930, 670.00000),
		MakePoint(650.85880, 675.00000),
		MakePoint(633.53829, 670.00000),
		MakePoint(633.53829, 660.00000),
		MakePoint(650.85880, 655.00000),
		MakePoint(668.17930, 660.00000)}
	corners := grid.HexCorners(hex)
	for i := 0; i < 6; i++ {
		validatePoint(t, expectedCorners[i], corners[i], 0.00001)
	}
}

func TestNeighbors(t *testing.T) {
	grid := MakeGrid(OrientationFlat, MakePoint(10, 20), MakePoint(20, 10), morton.Make64(2, 32))
	hex := grid.HexAt(MakePoint(666, 666))
	expectedNeighbors := [18]int64{
		920, 922, 944, 915, 921, 923, 945, 916, 918,
		926, 948, 917, 919, 925, 927, 960, 962, 968}
	neighbors := grid.HexNeighbors(hex, 2)
	for i := 0; i < len(neighbors); i++ {
		code := grid.HexToCode(neighbors[i])
		if expectedNeighbors[i] != code {
			t.Errorf("expected code %d but got %d", expectedNeighbors[i], code)
		}
	}
}

func TestRegion(t *testing.T) {
	grid := MakeGrid(OrientationFlat, MakePoint(10, 20), MakePoint(20, 10), morton.Make64(2, 32))
	geometry := [6]Point{
		MakePoint(20, 19.99999), MakePoint(20, 40), MakePoint(40, 60),
		MakePoint(60, 40), MakePoint(50, 30), MakePoint(40, 40)}
	region := grid.MakeRegion(geometry[:])
	hexes := region.Hexes()
	expectedHexCodes := []int64{0, 2, 1, 3, 9, 4}
	if len(hexes) != len(expectedHexCodes) {
		t.Errorf("expected %d hexes but got %d", len(expectedHexCodes), len(hexes))
		return
	}
	for i := 0; i < len(hexes); i++ {
		code := grid.HexToCode(hexes[i])
		if expectedHexCodes[i] != code {
			t.Errorf("expected hex with code %d but got %d", expectedHexCodes[i], code)
			return
		}
	}
}

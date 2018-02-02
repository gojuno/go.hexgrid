# HexGrid [![GoDoc](https://godoc.org/github.com/gojuno/go.hexgrid?status.svg)](http://godoc.org/github.com/gojuno/go.hexgrid) [![Build Status](https://travis-ci.org/gojuno/go.hexgrid.svg?branch=master)](https://travis-ci.org/gojuno/go.hexgrid)

## Basics

Configurable hex grid on abstract surface.

## Examples

```
import "github.com/gojuno/go.morton"
import "github.com/gojuno/go.hexgrid"

center := hexgrid.MakePoint(0, 0)
size := hexgrid.MakePoint(20, 10)
grid := hexgrid.MakeGrid(hexgrid.OrientationFlat, center, size, morton.Make64(2, 32))
hex := grid.HexAt(hexgrid.MakePoint(50, 50))
code := grid.HexToCode(hex)
restoredHex := grid.HexFromCode(code)
neighbors := grid.HexNeighbors(hex, 2)
points := []hexgrid.Point{hexgrid.MakePoint(0, 0), hexgrid.MakePoint(0, 10), hexgrid.MakePoint(10, 10), hexgrid.MakePoint(10, 0)}
region := grid.MakeRegion(points)
hexesInRegion := region.Hexes()
```
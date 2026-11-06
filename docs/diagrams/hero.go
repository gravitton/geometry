package main

import (
	"math"

	geom "github.com/gravitton/geometry"
)

const tileSize = 1.5

func hero(t theme) []*panel {
	p := newPanel(t, "")
	scene := geom.Bx(geom.Pt(0.0, 0.0), geom.Pt(45.0, 15.0))
	cells, walls := tilesOf([]string{
		"..............................",
		"..............................",
		"..............................",
		".............######...........",
		"..............................",
		".....####.....................",
		"..............................",
		"....................#.........",
		"#...................#........#",
		"###########....###############",
	})

	spikes := geom.Pol([]point{
		geom.Pt(16.5, 15.0), geom.Pt(16.5, 14.4), geom.Pt(17.25, 13.65), geom.Pt(18.0, 14.4), geom.Pt(18.75, 13.65),
		geom.Pt(19.5, 14.4), geom.Pt(20.25, 13.65), geom.Pt(21.0, 14.4), geom.Pt(21.75, 13.65), geom.Pt(22.5, 14.4),
		geom.Pt(22.5, 15.0),
	})
	solids := []geom.Shape[float64]{spikes}
	for _, wall := range walls {
		solids = append(solids, wall)
	}

	view := geom.Sect(geom.Pt(39.75, 2.25), 48.0, math.Pi/2+0.25, math.Pi/2-0.15)
	facing := view.Angle + view.Sweep/2
	drone := geom.RegPol(view.Center, geom.SzU(0.7), 5, 0, facing)
	light := geom.Pol(sight(view, scene, solids))
	look := geom.Ry(view.Center, geom.VectorFromAngle(facing, 1.0))

	hero := geom.Cap(geom.Pt(29.25, 12.15), geom.Pt(29.25, 13.05), 0.45)
	coins := []geom.Circle[float64]{
		geom.Circ(geom.Pt(35.25, 12.75), 0.375),
		geom.Circ(geom.Pt(36.75, 12.75), 0.375),
		geom.Circ(geom.Pt(38.25, 12.75), 0.375),
	}

	for _, cell := range cells {
		p.area(cell.Polygon().Points, style{fill: t.cell})
	}
	p.area(light.Points, style{fill: t.accentSoft})
	for _, cell := range cells {
		p.area(cell.Polygon().Points, style{stroke: t.edge, width: 1})
	}
	vision := geom.Sect(view.Center, 6.0, view.Angle, view.Sweep)
	p.area(vision.Polygon(tolerance).Points, style{stroke: t.accent, width: 1, dashed: true})
	for _, wall := range walls {
		p.area(wall.Polygon().Points, style{fill: t.wall, stroke: t.cell, width: 1.5})
	}
	p.area(spikes.Points, style{fill: t.coolSoft, stroke: t.cool, width: 1.5})

	if hit, ok := raycast(look, solids); ok {
		p.line(look.Origin, hit.Point, style{stroke: t.accent, width: 2})
		bounce := geom.Ry(hit.Point.Add(hit.Normal.Multiply(0.01)), look.Vector.Reflect(hit.Normal))
		if next, ok := raycast(bounce, solids); ok {
			p.line(hit.Point, next.Point, style{stroke: t.accent, width: 2, dashed: true})
		} else if exit, ok := bounce.ClipBox(scene); ok {
			p.line(hit.Point, exit.End, style{stroke: t.accent, width: 2, dashed: true})
		}
		p.arrow(hit.Point, hit.Point.Add(hit.Normal.Multiply(1.5)), t.leaf, false)
		p.dot(hit.Point, t.accent)
	}
	for _, coin := range coins {
		if light.Contains(coin.Center) {
			p.fill(coin, t.accent, t.accent)
		} else {
			p.fill(coin, "none", t.accent)
		}
	}
	p.fill(hero, t.coolSoft, t.cool)
	p.fill(drone, t.accent, t.accent)

	return []*panel{p}
}

func tilesOf(rows []string) (cells, walls []geom.Box[float64]) {
	size := geom.SzU(tileSize)
	for y, row := range rows {
		for x, tile := range row {
			box := geom.BoxFromMin(geom.Pt(float64(x)*tileSize, float64(y)*tileSize), size)
			if tile == '#' {
				walls = append(walls, box)
			} else {
				cells = append(cells, box)
			}
		}
	}

	return cells, walls
}

func sight(view geom.Sector[float64], scene geom.Box[float64], solids []geom.Shape[float64]) []point {
	outline := []point{view.Center}
	for angle := view.Angle; angle <= view.Angle+view.Sweep; angle += 0.002 {
		ray := geom.Ry(view.Center, geom.VectorFromAngle(angle, 1.0))
		reach := float64(view.Radius)
		if exit, ok := ray.ClipBox(scene); ok {
			reach = min(reach, exit.Length())
		}
		for _, solid := range solids {
			if hit, ok := geom.Raycast(ray, solid); ok {
				reach = min(reach, hit.Fraction)
			}
		}
		outline = append(outline, ray.PointAt(reach))
	}

	return outline
}

func raycast(ray geom.Ray[float64], solids []geom.Shape[float64]) (geom.Hit[float64], bool) {
	nearest, found := geom.Hit[float64]{}, false
	for _, solid := range solids {
		if hit, ok := geom.Raycast(ray, solid); ok && (!found || hit.Fraction < nearest.Fraction) {
			nearest, found = hit, true
		}
	}

	return nearest, found
}

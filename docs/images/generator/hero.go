package main

import (
	"fmt"
	"math"

	geom "github.com/gravitton/geometry"
)

const tileSize = 0.9

var heroMap = []string{
	"...............",
	"..#.....#......",
	"..#.....#....G.",
	"..#..~~.#......",
	"....S~~........",
	"..#..~~...##...",
	"..#.......#....",
	"..........#....",
	"...............",
}

type world struct {
	bounds geom.Box[float64]
	ground []geom.Box[float64]
	rough  []geom.Box[float64]
	walls  []geom.Box[float64]
	scout  point
	goal   point
}

func hero(t theme) []*panel {
	w := worldOf(heroMap)

	return []*panel{w.view(t), w.sweep(t), w.clearance(t)}
}

func (w world) view(t theme) *panel {
	p := w.backdrop(t, "Raycast(ray, wall)")
	light := geom.Pol(sight(geom.Sect(w.scout, 6*tileSize, 0, 2*math.Pi), w.bounds, w.solids()))
	p.area(light.Points, style{fill: t.accentSoft})
	w.edges(p)
	w.draw(p)
	p.fill(geom.RegPol(w.scout, geom.SzU(0.35), 5, 0, -math.Pi/2), t.accent, t.accent)

	return p
}

func (w world) sweep(t theme) *panel {
	p := w.backdrop(t, "Sweep(scout, motion, wall)")
	w.draw(p)
	scout := geom.Circ(w.scout, 0.3)
	motion := w.goal.Subtract(w.scout)
	p.arrow(w.scout, w.goal, t.muted, true)
	if hit, ok := sweep(scout, motion, w.walls); ok {
		p.fill(scout.Translate(motion.Multiply(hit.Fraction)), t.accentSoft, t.accent)
		p.arrow(hit.Point, hit.Point.Add(hit.Normal.Multiply(1.2)), t.leaf, false)
		p.dot(hit.Point, t.accent)
	}
	p.fill(scout, t.accent, t.accent)
	p.fill(geom.Circ(w.goal, 0.3), t.cool, t.cool)

	return p
}

func (w world) clearance(t theme) *panel {
	p := w.backdrop(t, "wall.Grow(radius)")
	for _, wall := range w.walls {
		p.area(wall.Grow(0.3).Polygon().Points, style{fill: t.coolSoft, stroke: t.cool, width: 1, dashed: true})
	}
	w.draw(p)
	p.fill(geom.Circ(w.scout, 0.3), t.accent, t.accent)
	p.fill(geom.Circ(w.goal, 0.3), t.cool, t.cool)

	return p
}

func (w world) backdrop(t theme, caption string) *panel {
	p := newPanel(t, caption)
	for _, cell := range w.ground {
		p.area(cell.Polygon().Points, style{fill: t.cell, stroke: t.edge, width: 1})
	}
	for _, cell := range w.rough {
		p.area(cell.Polygon().Points, style{fill: t.leafSoft, stroke: t.edge, width: 1})
	}

	return p
}

func (w world) edges(p *panel) {
	for _, cell := range append(append([]geom.Box[float64]{}, w.ground...), w.rough...) {
		p.area(cell.Polygon().Points, style{stroke: p.theme.edge, width: 1})
	}
}

func (w world) draw(p *panel) {
	for _, wall := range w.walls {
		p.area(wall.Polygon().Points, style{fill: p.theme.wall, stroke: p.theme.cell, width: 1})
	}
}

func (w world) solids() []geom.Shape[float64] {
	solids := make([]geom.Shape[float64], len(w.walls))
	for i, wall := range w.walls {
		solids[i] = wall
	}

	return solids
}

func worldOf(rows []string) world {
	size := geom.SzU(tileSize)
	w := world{bounds: geom.Bx(geom.Pt(0.0, 0.0), geom.Pt(float64(len(rows[0]))*tileSize, float64(len(rows))*tileSize))}
	for y, row := range rows {
		for x, tile := range row {
			box := geom.BoxFromMin(geom.Pt(float64(x)*tileSize, float64(y)*tileSize), size)
			switch tile {
			case '#':
				w.walls = append(w.walls, box)
			case '~':
				w.rough = append(w.rough, box)
			case 'S':
				w.ground, w.scout = append(w.ground, box), box.Center()
			case 'G':
				w.ground, w.goal = append(w.ground, box), box.Center()
			case '.':
				w.ground = append(w.ground, box)
			default:
				panic(fmt.Sprintf("diagrams: unknown tile %q on the hero map", tile))
			}
		}
	}

	return w
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

func sweep(body geom.Circle[float64], motion geom.Vector[float64], walls []geom.Box[float64]) (geom.Hit[float64], bool) {
	nearest, found := geom.Hit[float64]{}, false
	for _, wall := range walls {
		if hit, ok := geom.Sweep[float64](body, motion, wall); ok && (!found || hit.Fraction < nearest.Fraction) {
			nearest, found = hit, true
		}
	}

	return nearest, found
}

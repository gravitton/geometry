package main

import (
	"fmt"
	"math"

	geom "github.com/gravitton/geometry"
)

const (
	tileSize = 0.9
	bodySize = 0.3
)

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

	return []*panel{w.view(t), w.sweep(t), w.curve(t)}
}

func (w world) view(t theme) *panel {
	p := w.backdrop(t, "Raycast(ray, wall), Reflect(normal)")
	light := geom.Pol(sight(geom.Sect(w.scout, 6*tileSize, 0, 2*math.Pi), w.bounds, w.solids()))
	p.area(light.Points, style{fill: t.accentSoft})
	w.edges(p)
	w.draw(p)
	angle, hits := w.aim()
	from := w.scout
	for _, hit := range hits {
		p.line(from, hit.Point, style{stroke: t.accent, width: 2})
		p.dot(hit.Point, t.accent)
		from = hit.Point
	}
	last := hits[len(hits)-1]
	p.arrow(last.Point, last.Point.Add(last.Normal.Multiply(0.9)), t.leaf, false)
	p.fill(geom.Circ(w.goal, bodySize), t.cool, t.cool)
	p.fill(geom.RegPol(w.scout, geom.SzU(0.35), 5, 0, angle), t.accent, t.accent)

	return p
}

func (w world) sweep(t theme) *panel {
	p := w.backdrop(t, "Sweep(scout, motion, wall)")
	w.draw(p)
	scout := geom.Circ(w.scout, bodySize)
	motion := w.goal.Subtract(w.scout)
	p.arrow(w.scout, w.goal, t.muted, true)
	if hit, ok := sweep(scout, motion, w.walls); ok {
		p.fill(scout.Translate(motion.Multiply(hit.Fraction)), t.accentSoft, t.accent)
		p.arrow(hit.Point, hit.Point.Add(hit.Normal.Multiply(1.2)), t.leaf, false)
		p.dot(hit.Point, t.accent)
	}
	p.fill(scout, t.accent, t.accent)
	p.fill(geom.Circ(w.goal, bodySize), t.cool, t.cool)

	return p
}

func (w world) curve(t theme) *panel {
	p := w.backdrop(t, "curve.IntersectionSegment(edge)")
	w.draw(p)
	blocked := geom.CubBez(w.scout, geom.Pt(6.0, 1.2), geom.Pt(9.6, 1.6), w.goal)
	p.path(blocked.Polyline(tolerance).Points, style{stroke: t.muted, width: 1.5, dashed: true})
	for _, wall := range w.walls {
		for edge := range wall.Edges() {
			for _, at := range blocked.IntersectionSegment(edge) {
				p.dot(at, t.accent)
			}
		}
	}
	clear := geom.CubBez(w.scout, geom.Pt(8.4, 5.2), geom.Pt(9.6, 2.8), w.goal)
	p.path([]point{clear.Start, clear.Control1, clear.Control2, clear.End}, style{stroke: t.muted, width: 1, dashed: true})
	p.path(clear.Polyline(tolerance).Points, style{stroke: t.cool, width: 2.5})
	p.ring(clear.Control1, t.muted)
	p.ring(clear.Control2, t.muted)
	p.fill(geom.Circ(w.scout, bodySize), t.accent, t.accent)
	p.fill(geom.Circ(w.goal, bodySize), t.cool, t.cool)

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

func (w world) fence() []geom.Shape[float64] {
	b := w.bounds
	return []geom.Shape[float64]{
		geom.Bx(geom.Pt(b.Min.X-1, b.Min.Y-1), geom.Pt(b.Max.X+1, b.Min.Y)),
		geom.Bx(geom.Pt(b.Min.X-1, b.Max.Y), geom.Pt(b.Max.X+1, b.Max.Y+1)),
		geom.Bx(geom.Pt(b.Min.X-1, b.Min.Y), geom.Pt(b.Min.X, b.Max.Y)),
		geom.Bx(geom.Pt(b.Max.X, b.Min.Y), geom.Pt(b.Max.X+1, b.Max.Y)),
	}
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

func raycast(ray geom.Ray[float64], solids []geom.Shape[float64]) (geom.Hit[float64], bool) {
	nearest, found := geom.Hit[float64]{}, false
	for _, solid := range solids {
		if hit, ok := geom.Raycast(ray, solid); ok && (!found || hit.Fraction < nearest.Fraction) {
			nearest, found = hit, true
		}
	}

	return nearest, found
}

func (w world) aim() (float64, []geom.Hit[float64]) {
	for bounces := 2; bounces <= 4; bounces++ {
		for step := range 3600 {
			angle := float64(step) * 2 * math.Pi / 3600
			hits := w.reflect(geom.Ry(w.scout, geom.VectorFromAngle(angle, 1.0)), bounces)
			if len(hits) == bounces+1 && geom.Circ(w.goal, bodySize).Contains(hits[bounces].Point) {
				return angle, hits
			}
		}
	}

	panic(fmt.Sprintf("diagrams: no ray from %s reaches the goal at %s in four bounces", w.scout, w.goal))
}

func (w world) reflect(ray geom.Ray[float64], bounces int) []geom.Hit[float64] {
	goal := geom.Circ(w.goal, bodySize)
	targets := append(append(w.solids(), w.fence()...), goal)
	hits := make([]geom.Hit[float64], 0, bounces+1)
	for range bounces + 1 {
		hit, ok := raycast(ray, targets)
		if !ok {
			break
		}
		hits = append(hits, hit)
		if goal.Contains(hit.Point) {
			break
		}
		ray = geom.Ry(hit.Point.Add(hit.Normal.Multiply(0.01)), ray.Vector.Reflect(hit.Normal))
	}

	return hits
}

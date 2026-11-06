package main

import (
	"fmt"
	"math"

	geom "github.com/gravitton/geometry"
)

const (
	unit      = 20.0
	tolerance = 0.02
)

type panel struct {
	canvas
}

func newPanel(t theme, caption string) *panel {
	return &panel{canvas: newCanvas(t, caption)}
}

func paper(t theme, caption string, area geom.Box[float64]) *panel {
	if !area.Min.IsZero() || math.Mod(area.Max.X, 5) != 0 || math.Mod(area.Max.Y, 5) != 0 {
		panic(fmt.Sprintf("diagrams: paper %s is not a multiple of 5 from the origin", area))
	}

	p := newPanel(t, caption)
	for x := area.Min.X; x <= area.Max.X; x++ {
		p.line(geom.Pt(x, area.Min.Y), geom.Pt(x, area.Max.Y), style{stroke: t.edge, width: ruleWidth(x)})
	}
	for y := area.Min.Y; y <= area.Max.Y; y++ {
		p.line(geom.Pt(area.Min.X, y), geom.Pt(area.Max.X, y), style{stroke: t.edge, width: ruleWidth(y)})
	}

	return p
}

func ruleWidth(at float64) float64 {
	if math.Mod(at, 5) == 0 {
		return 1
	}

	return 0.4
}

func (p *panel) area(points []point, s style) {
	p.polygon(scaledOf(points), s)
}

func (p *panel) path(points []point, s style) {
	p.polyline(scaledOf(points), s)
}

func (p *panel) line(a, b point, s style) {
	p.canvas.line(scaled(a), scaled(b), s)
}

func (p *panel) dot(at point, color string) {
	p.circle(scaled(at), 4, style{fill: color})
}

func (p *panel) ring(at point, color string) {
	p.circle(scaled(at), 4, style{fill: p.theme.cell, stroke: color, width: 1.5})
}

func (p *panel) arrow(from, to point, color string, dashed bool) {
	p.canvas.arrow(scaled(from), scaled(to), style{stroke: color, width: 2, dashed: dashed})
}

func (p *panel) text(at point, anchor, text, color string) {
	p.canvas.text(scaled(at), font{size: 12, anchor: anchor, halo: true}, span{text, color})
}

func (p *panel) named(at point, text string, dx, dy float64) {
	p.text(at.AddXY(dx, dy), anchorOf(dx), text, p.theme.ink)
}

func (p *panel) result(at point, text string, dx, dy float64) {
	p.dot(at, p.theme.accent)
	if text != "" {
		p.text(at.AddXY(dx, dy), anchorOf(dx), text, p.theme.accent)
	}
}

func (p *panel) fill(shape geom.Shape[float64], fill, stroke string) {
	if segment, ok := shape.(geom.Segment[float64]); ok {
		p.line(segment.Start, segment.End, style{stroke: stroke, width: 2.5})
		return
	}

	p.area(outline(shape), style{fill: fill, stroke: stroke, width: 2})
}

func (p *panel) trace(shape geom.Shape[float64], color string) {
	p.area(outline(shape), style{stroke: color, width: 1.5, dashed: true})
}

func (p *panel) beam(part geom.Segment[float64], ok bool, color string, dashed bool) {
	if ok {
		p.line(part.Start, part.End, style{stroke: color, width: 2, dashed: dashed})
	}
}

func scaled(at point) point {
	return at.Multiply(unit)
}

func scaledOf(points []point) []point {
	result := make([]point, len(points))
	for i, at := range points {
		result[i] = scaled(at)
	}

	return result
}

func anchorOf(dx float64) string {
	switch {
	case dx > 0:
		return "start"
	case dx < 0:
		return "end"
	default:
		return "middle"
	}
}

func outline(shape geom.Shape[float64]) []point {
	switch s := shape.(type) {
	case geom.Circle[float64]:
		return s.Polygon(tolerance).Points
	case geom.Capsule[float64]:
		return s.Polygon(tolerance).Points
	case geom.Sector[float64]:
		return s.Polygon(tolerance).Points
	case geom.Segment[float64]:
		return []point{s.Start, s.End}
	case geom.Polygon[float64]:
		return s.Points
	case geom.Rectangle[float64]:
		return s.Polygon().Points
	case geom.RegularPolygon[float64]:
		return s.Polygon().Points
	case geom.Box[float64]:
		return s.Polygon().Points
	case geom.Ellipse[float64]:
		return s.Polygon(tolerance).Points
	default:
		panic(fmt.Sprintf("diagrams: no outline for %T", shape))
	}
}

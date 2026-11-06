package main

import (
	"fmt"
	"math"

	geom "github.com/gravitton/geometry"
)

var figures = []figure{
	{"hero", hero},
	{"shapes", shapes},
	{"paths", paths},
	{"vectors", vectors},
	{"sizes", sizes},
	{"rectangles", rectangles},
	{"boxes", boxes},
	{"circles", circles},
	{"ellipses", ellipses},
	{"polygons", polygons},
	{"regular-polygons", regularPolygons},
	{"intersections", intersections},
	{"casts", casts},
	{"lines", lines},
	{"collisions", collisions},
	{"curves", curves},
	{"coordinates", coordinates},
	{"matrices", matrices},
	{"numbers", numbers},
	{"rounding", rounding},
	{"pixels", pixels},
}

var notch = geom.Pol([]point{{X: 1, Y: 1}, {X: 9, Y: 1}, {X: 9, Y: 9}, {X: 5, Y: 4}, {X: 1, Y: 9}})

func shapes(t theme) []*panel {
	p := paper(t, "", geom.Bx(geom.Pt(0.0, 0.0), geom.Pt(45.0, 15.0)))
	kinds := []struct {
		name  string
		shape geom.Shape[float64]
		at    point
	}{
		{"Circle", geom.Circ(geom.Pt(4.5, 4.0), 3.0), geom.Pt(4.5, 7.0)},
		{"Capsule", geom.Cap(geom.Pt(11.5, 5.5), geom.Pt(15.5, 2.5), 1.5), geom.Pt(13.5, 7.0)},
		{"Sector", geom.Sect(geom.Pt(22.5, 6.5), 5.0, -3*math.Pi/4, math.Pi/2), geom.Pt(22.5, 7.0)},
		{"Segment", geom.Seg(geom.Pt(28.5, 6.0), geom.Pt(34.5, 1.0)), geom.Pt(31.5, 7.0)},
		{"Polygon", geom.Pol([]point{{X: 37.5, Y: 1}, {X: 43.5, Y: 1}, {X: 43.5, Y: 6}, {X: 40.5, Y: 3.5}, {X: 37.5, Y: 6}}), geom.Pt(40.5, 7.0)},
		{"Rectangle", geom.Rect(geom.Pt(9.0, 11.0), geom.Sz(6.0, 3.0), 0.4), geom.Pt(9.0, 14.0)},
		{"RegularPolygon", geom.RegularPolygonWithOrientation(geom.Pt(18.0, 11.0), geom.SzU(2.5), 6, geom.OrientationFlatTop), geom.Pt(18.0, 14.0)},
		{"Box", geom.Bx(geom.Pt(24.0, 9.0), geom.Pt(30.0, 13.0)), geom.Pt(27.0, 14.0)},
		{"Ellipse", geom.Ell(geom.Pt(36.0, 11.0), geom.Sz(3.5, 1.8), 0.5), geom.Pt(36.0, 14.0)},
	}
	for _, kind := range kinds {
		p.fill(kind.shape, t.coolSoft, t.cool)
		if segment, ok := kind.shape.(geom.Segment[float64]); ok {
			p.dot(segment.Start, t.cool)
			p.dot(segment.End, t.cool)
		}
		p.text(kind.at.AddXY(0, 0.4), "middle", kind.name, t.muted)
	}

	return []*panel{p}
}

func paths(t theme) []*panel {
	p := paper(t, "", geom.Bx(geom.Pt(0.0, 0.0), geom.Pt(45.0, 10.0)))

	ray := geom.Ry(geom.Pt(1.0, 6.0), geom.Vec(2.0, -1.0))
	if part, ok := ray.ClipBox(geom.Bx(geom.Pt(0.0, 0.0), geom.Pt(7.0, 8.0))); ok {
		p.arrow(part.Start, part.End, t.cool, false)
	}
	p.dot(ray.Origin, t.cool)

	line := geom.Ln(geom.Pt(11.0, 4.0), geom.Vec(2.0, 1.0))
	part, ok := line.ClipBox(geom.Bx(geom.Pt(8.0, 1.0), geom.Pt(14.0, 7.0)))
	p.beam(part, ok, t.cool, false)
	p.dot(line.Origin, t.cool)

	polyline := geom.PolLn([]point{{X: 16, Y: 7}, {X: 18, Y: 1}, {X: 20, Y: 5}, {X: 22, Y: 1}})
	p.path(polyline.Points, style{stroke: t.cool, width: 2.5})
	for vertex := range polyline.Vertices() {
		p.dot(vertex, t.cool)
	}

	arc := geom.Ac(geom.Pt(27.0, 4.0), 3.0, -math.Pi, 1.25*math.Pi)
	p.trace(arc.Circle(), t.edge)
	p.path(arc.Polyline(tolerance).Points, style{stroke: t.cool, width: 2.5})
	p.dot(arc.Start(), t.cool)
	p.dot(arc.End(), t.cool)

	quadratic := geom.QuadBez(geom.Pt(31.0, 7.0), geom.Pt(34.0, -1.0), geom.Pt(37.0, 7.0))
	p.path([]point{quadratic.Start, quadratic.Control, quadratic.End}, style{stroke: t.muted, width: 1, dashed: true})
	p.path(quadratic.Polyline(tolerance).Points, style{stroke: t.cool, width: 2.5})
	p.ring(quadratic.Control, t.muted)

	cubic := geom.CubBez(geom.Pt(38.0, 6.0), geom.Pt(39.0, 0.0), geom.Pt(43.0, 8.0), geom.Pt(44.0, 2.0))
	p.path([]point{cubic.Start, cubic.Control1, cubic.Control2, cubic.End}, style{stroke: t.muted, width: 1, dashed: true})
	p.path(cubic.Polyline(tolerance).Points, style{stroke: t.cool, width: 2.5})
	p.ring(cubic.Control1, t.muted)
	p.ring(cubic.Control2, t.muted)

	for _, name := range []struct {
		x    float64
		text string
	}{{3.5, "Ray"}, {11.0, "Line"}, {19.0, "Polyline"}, {27.0, "Arc"}, {34.0, "QuadraticBezier"}, {41.0, "CubicBezier"}} {
		p.text(geom.Pt(name.x, 8.8), "middle", name.text, t.muted)
	}

	return []*panel{p}
}

func vectors(t theme) []*panel {
	p := paper(t, "", geom.Bx(geom.Pt(0.0, 0.0), geom.Pt(15.0, 15.0)))
	a, b := geom.Pt(2.0, 8.0), geom.Pt(8.0, 5.0)
	v := b.Subtract(a)
	bounce := v.Reflect(geom.Vec(0.0, 1.0))

	p.line(geom.Pt(5.0, 5.0), geom.Pt(15.0, 5.0), style{stroke: t.muted, width: 2.5})
	p.arrow(a, b, t.cool, false)
	p.arrow(a, a.Add(v.Normal()), t.leaf, false)
	p.arrow(b, b.Add(bounce), t.cool, true)
	p.result(a.Lerp(b, 0.25), "a.Lerp(b, 0.25)", 0.5, 0.6)
	p.result(a.RotateAround(b, -math.Pi/2), "a.RotateAround(b, -π/2)", 0.5, 0.6)
	p.dot(a, t.ink)
	p.dot(b, t.ink)
	p.named(a, "a", -0.5, 0)
	p.named(b, "b", 0, -0.8)
	p.named(a.Lerp(b, 0.6), "v", 0, -0.7)
	p.named(a.Add(v.Normal()), "v.Normal()", 0.5, 0)
	p.named(b.Add(bounce), "v.Reflect(⟨0,1⟩)", 0, 0.8)

	return []*panel{p}
}

func sizes(t theme) []*panel {
	p := paper(t, "", geom.Bx(geom.Pt(0.0, 0.0), geom.Pt(20.0, 10.0)))
	origin := geom.Pt(1.0, 1.0)
	screen, slot := geom.Sz(16.0, 9.0), geom.SzU(6.0)

	p.trace(geom.BoxFromMin(origin, screen), t.muted)
	p.fill(geom.BoxFromMin(origin, screen.Fill(slot)), t.leafSoft, t.leaf)
	p.fill(geom.BoxFromMin(origin, slot), t.coolSoft, t.cool)
	p.fill(geom.BoxFromMin(origin, screen.Fit(slot)), t.accentSoft, t.accent)
	p.named(origin.AddXY(16, 9), "screen", 0, 0.7)
	p.named(origin.AddXY(6, 6), "slot", -0.3, -0.5)
	p.text(origin.AddXY(3, 1.7), "middle", "Fit", t.accent)
	p.text(origin.AddXY(8.5, 5.4), "middle", "Fill", t.leaf)

	return []*panel{p}
}

func rectangles(t theme) []*panel {
	p := paper(t, "", geom.Bx(geom.Pt(0.0, 0.0), geom.Pt(30.0, 10.0)))
	r := geom.Rect(geom.Pt(6.0, 5.5), geom.Sz(7.0, 3.5), 0.4)

	p.trace(r.Bounds(), t.muted)
	p.trace(r.BoundingCircle(), t.edge)
	p.fill(r, t.coolSoft, t.cool)
	for _, corner := range []struct {
		name string
		at   point
	}{{"TopLeft", r.TopLeft()}, {"TopRight", r.TopRight()}, {"BottomRight", r.BottomRight()}, {"BottomLeft", r.BottomLeft()}} {
		outward := corner.at.Subtract(r.Center).Normalize().Multiply(0.6)
		p.dot(corner.at, t.cool)
		p.text(corner.at.Add(outward), anchorOf(outward.X), corner.name, t.ink)
	}
	p.result(r.Anchor(geom.Top), "Anchor(Top)", 0, -1.3)
	p.named(r.Center, "r", 0, 0)
	p.text(geom.Pt(r.Bounds().Min.X, r.Bounds().Max.Y+0.7), "start", "Bounds", t.muted)

	screen := geom.RectangleFromMinMax(geom.Pt(15.0, 1.0), geom.Pt(27.0, 10.0))
	w := geom.RectangleFromCenter(geom.Pt(26.0, 1.5), geom.Sz(5.0, 3.0))
	clamped := w.Clamp(screen)
	p.fill(screen, t.cell, t.muted)
	p.trace(w, t.muted)
	p.fill(clamped, t.accentSoft, t.accent)
	p.arrow(w.Center, clamped.Center, t.accent, false)
	p.named(screen.Anchor(geom.BottomLeft), "screen", 0.4, -0.6)
	p.named(w.Center.AddXY(2.5, 0), "w", 0.3, 0)
	p.text(clamped.Center.AddXY(0, 2.2), "middle", "w.Clamp(screen)", t.accent)

	return []*panel{p}
}

func boxes(t theme) []*panel {
	p := paper(t, "", geom.Bx(geom.Pt(0.0, 0.0), geom.Pt(25.0, 15.0)))
	view := geom.BoxFromSize(geom.Sz(20.0, 12.0))
	hud := geom.Bx(geom.Pt(6.0, 2.0), geom.Pt(0.0, 0.0))
	tip := geom.BoxFromMin(geom.Pt(17.0, 8.0), geom.Sz(7.0, 3.0))

	p.fill(view, t.cell, t.muted)
	p.trace(view.Inset(geom.PadU(1.0)), t.edge)
	p.fill(hud, t.coolSoft, t.cool)
	p.fill(hud.AlignTo(geom.TopRight, view.Anchor(geom.TopRight)), t.accentSoft, t.accent)
	p.trace(tip, t.muted)
	p.fill(tip.Clamp(view), t.accentSoft, t.accent)
	if part, ok := view.ClipBox(tip); ok {
		p.area(part.Polygon().Points, style{stroke: t.leaf, width: 2.5})
	}
	p.named(geom.Pt(10.0, 12.7), "view", 0, 0)
	p.named(geom.Pt(3.0, 1.0), "hud", 0, 0)
	p.text(geom.Pt(17.0, 2.7), "middle", "hud.AlignTo(TopRight, …)", t.accent)
	p.named(geom.Pt(24.4, 9.5), "tip", 0.3, 0)
	p.text(geom.Pt(16.5, 7.3), "middle", "tip.Clamp(view)", t.accent)
	p.text(geom.Pt(18.5, 12.7), "middle", "view.ClipBox(tip)", t.leaf)

	return []*panel{p}
}

func circles(t theme) []*panel {
	c, wheel := geom.Circ(geom.Pt(3.0, 5.0), 2.0), geom.Circ(geom.Pt(8.0, 6.0), 1.0)
	area := geom.Bx(geom.Pt(0.0, 0.0), geom.Pt(10.0, 10.0))

	tangents := paper(t, "c.Tangents(from)", area)
	from := geom.Pt(8.5, 1.5)
	tangents.fill(c, t.coolSoft, t.cool)
	for at := range c.Tangents(from).All() {
		tangents.line(from, at, style{stroke: t.accent, width: 2})
		tangents.dot(at, t.accent)
	}
	tangents.dot(from, t.ink)
	tangents.named(from, "from", 0.5, 0)
	tangents.named(c.Center, "c", 0, 0)

	outer := paper(t, "c.OuterTangents(wheel)", area)
	outer.fill(c, t.coolSoft, t.cool)
	outer.fill(wheel, t.leafSoft, t.leaf)
	for belt := range c.OuterTangents(wheel).All() {
		outer.line(belt.Start, belt.End, style{stroke: t.accent, width: 2})
	}

	inner := paper(t, "c.InnerTangents(wheel)", area)
	inner.fill(c, t.coolSoft, t.cool)
	inner.fill(wheel, t.leafSoft, t.leaf)
	for belt := range c.InnerTangents(wheel).All() {
		inner.line(belt.Start, belt.End, style{stroke: t.accent, width: 2})
	}

	union := paper(t, "c.Union(wheel)", area)
	union.trace(c.Union(wheel), t.accent)
	union.fill(c, t.coolSoft, t.cool)
	union.fill(wheel, t.leafSoft, t.leaf)

	return []*panel{tangents, outer, inner, union}
}

func ellipses(t theme) []*panel {
	c := geom.Circ(geom.Pt(4.0, 5.0), 2.5)
	e := c.Ellipse().Transform(geom.ScaleMatrix(2.0, 1.0))

	ellipse := paper(t, "o.Ellipse().Transform(ScaleMatrix(2, 1))", geom.Bx(geom.Pt(0.0, 0.0), geom.Pt(15.0, 10.0)))
	ellipse.trace(c, t.muted)
	ellipse.fill(e, t.coolSoft, t.cool)
	inside, outside := geom.Pt(6.0, 4.0), geom.Pt(10.0, 8.5)
	focusA, focusB := e.Foci()
	ellipse.dot(focusA, t.cool)
	ellipse.dot(focusB, t.cool)
	ellipse.result(inside, "Contains", 0.4, -0.5)
	ellipse.ring(outside, t.ink)
	ellipse.line(outside, e.Nearest(outside), style{stroke: t.accent, width: 2, dashed: true})
	ellipse.dot(e.Nearest(outside), t.accent)
	ellipse.text(outside.AddXY(0.6, 0.2), "start", fmt.Sprintf("DistanceTo %.2f", e.DistanceTo(outside)), t.accent)
	ellipse.named(focusB, "Foci", 0, 0.8)

	body := geom.Cap(geom.Pt(4.0, 2.0), geom.Pt(4.0, 8.0), 1.5)
	capsule := paper(t, "body.DistanceTo(side)", geom.Bx(geom.Pt(0.0, 0.0), geom.Pt(10.0, 10.0)))
	capsule.fill(body, t.coolSoft, t.cool)
	capsule.line(body.Segment.Start, body.Segment.End, style{stroke: t.cool, width: 1, dashed: true})
	side := geom.Pt(8.5, 5.0)
	capsule.ring(side, t.ink)
	capsule.line(side, body.Nearest(side), style{stroke: t.accent, width: 2, dashed: true})
	capsule.dot(body.Nearest(side), t.accent)
	capsule.text(side.Lerp(body.Nearest(side), 0.5).AddXY(0, -0.6), "middle", fmt.Sprintf("%g", body.DistanceTo(side)), t.accent)
	capsule.named(side, "side", 0, 0.8)

	return []*panel{ellipse, capsule}
}

func polygons(t theme) []*panel {
	area := geom.Bx(geom.Pt(0.0, 0.0), geom.Pt(10.0, 10.0))

	hull := paper(t, "p.ConvexHull()", area)
	hull.fill(notch, t.coolSoft, t.cool)
	hull.trace(notch.ConvexHull(), t.accent)

	triangles := paper(t, "p.Triangulate()", area)
	indices := notch.Triangulate()
	for i := 0; i+2 < len(indices); i += 3 {
		triangle := []point{notch.Points[indices[i]], notch.Points[indices[i+1]], notch.Points[indices[i+2]]}
		triangles.area(triangle, style{fill: t.coolSoft, stroke: t.cool, width: 1.5})
	}

	offset := paper(t, "p.Grow(0.5), p.Shrink(0.5)", area)
	offset.area(notch.Grow(0.5).Points, style{stroke: t.leaf, width: 2})
	offset.fill(notch, t.coolSoft, t.cool)
	offset.area(notch.Shrink(0.5).Points, style{stroke: t.accent, width: 2})

	split := paper(t, "p.Split(cut)", area)
	cut := geom.Ln(geom.Pt(3.0, 0.0), geom.Vec(1.0, 2.0))
	front, back := notch.Split(cut)
	split.fill(front, t.coolSoft, t.cool)
	split.fill(back, t.leafSoft, t.leaf)
	part, ok := cut.ClipBox(area)
	split.beam(part, ok, t.accent, true)

	return []*panel{hull, triangles, offset, split}
}

func regularPolygons(t theme) []*panel {
	area := geom.Bx(geom.Pt(0.0, 0.0), geom.Pt(10.0, 10.0))
	numbered := func(caption string, rp geom.RegularPolygon[float64]) *panel {
		p := paper(t, caption, area)
		p.fill(rp, t.coolSoft, t.cool)
		i := 0
		for vertex := range rp.Vertices() {
			p.text(rp.Center.Add(vertex.Subtract(rp.Center).Multiply(1.22)), "middle", fmt.Sprint(i), t.muted)
			i++
		}
		p.result(rp.Anchor(geom.Top), "", 0, 0)

		return p
	}

	flat := numbered("FlatTop", geom.RegularPolygonWithOrientation(geom.Pt(5.0, 5.0), geom.SzU(4.0), 6, geom.OrientationFlatTop))
	pointy := numbered("PointyTop", geom.RegularPolygonWithOrientation(geom.Pt(5.0, 5.0), geom.SzU(4.0), 6, geom.OrientationPointyTop))

	stretched := paper(t, "Sz(4.5, 2.5)", area)
	wide := geom.RegularPolygonWithOrientation(geom.Pt(5.0, 5.0), geom.Sz(4.5, 2.5), 6, geom.OrientationFlatTop)
	stretched.trace(wide.Ellipse(), t.muted)
	stretched.fill(wide, t.coolSoft, t.cool)

	base := geom.RegPol(geom.Pt(5.0, 5.0), geom.Sz(4.5, 2.5), 5, 0, 0)
	phased := paper(t, "phase or angle 0.6", area)
	phased.area(base.Polygon().Points, style{stroke: t.edge, width: 1.5})
	phased.area(geom.RegPol(geom.Pt(5.0, 5.0), geom.Sz(4.5, 2.5), 5, 0, 0.6).Polygon().Points, style{stroke: t.accent, width: 2})
	phased.fill(base.Rotate(0.6), t.coolSoft, t.cool)

	return []*panel{flat, pointy, stretched, phased}
}

func intersections(t theme) []*panel {
	area := geom.Bx(geom.Pt(0.0, 0.0), geom.Pt(10.0, 10.0))

	circles := paper(t, "a.IntersectionCircle(b)", area)
	a, b := geom.Circ(geom.Pt(4.0, 5.0), 3.0), geom.Circ(geom.Pt(7.0, 6.0), 2.5)
	circles.fill(a, t.coolSoft, t.cool)
	circles.fill(b, t.leafSoft, t.leaf)
	for at := range a.IntersectionCircle(b).All() {
		circles.dot(at, t.accent)
	}

	clipPolygon := paper(t, "s.ClipPolygon(p)", area)
	s := geom.Seg(geom.Pt(0.0, 6.0), geom.Pt(10.0, 6.0))
	clipPolygon.fill(notch, t.coolSoft, t.cool)
	clipPolygon.line(s.Start, s.End, style{stroke: t.muted, width: 1.5, dashed: true})
	for _, part := range s.ClipPolygon(notch) {
		clipPolygon.line(part.Start, part.End, style{stroke: t.accent, width: 3.5})
	}

	clipSector := paper(t, "s.ClipSector(cone)", area)
	cone := geom.Sect(geom.Pt(5.0, 5.0), 4.0, -math.Pi/4, 1.5*math.Pi)
	beam := geom.Seg(geom.Pt(0.0, 3.5), geom.Pt(10.0, 3.5))
	clipSector.fill(cone, t.coolSoft, t.cool)
	clipSector.line(beam.Start, beam.End, style{stroke: t.muted, width: 1.5, dashed: true})
	for _, part := range beam.ClipSector(cone) {
		clipSector.line(part.Start, part.End, style{stroke: t.accent, width: 3.5})
	}

	clipRectangle := paper(t, "r.ClipRectangle(q)", area)
	r := geom.Rect(geom.Pt(4.0, 4.5), geom.Sz(5.0, 3.5), 0.3)
	q := geom.Rect(geom.Pt(6.5, 6.0), geom.Sz(4.0, 3.0), 0.3)
	clipRectangle.fill(r, t.coolSoft, t.cool)
	clipRectangle.fill(q, t.leafSoft, t.leaf)
	if overlap, ok := r.ClipRectangle(q); ok {
		clipRectangle.fill(overlap, t.accentSoft, t.accent)
	}

	return []*panel{circles, clipPolygon, clipSector, clipRectangle}
}

func casts(t theme) []*panel {
	area := geom.Bx(geom.Pt(0.0, 0.0), geom.Pt(15.0, 10.0))

	raycast := paper(t, "ray.HitRectangle(wall)", area)
	wall := geom.Rect(geom.Pt(10.0, 5.0), geom.Sz(2.0, 6.0), 0.35)
	ray := geom.Ry(geom.Pt(1.0, 7.0), geom.Vec(3.0, -1.0))
	raycast.fill(wall, t.coolSoft, t.cool)
	if hit, ok := ray.HitRectangle(wall); ok {
		raycast.arrow(ray.Origin, hit.Point, t.accent, false)
		raycast.arrow(hit.Point, hit.Point.Add(hit.Normal.Multiply(2)), t.leaf, false)
		raycast.arrow(hit.Point, hit.Point.Add(ray.Vector.Reflect(hit.Normal).Resize(4)), t.accent, true)
		raycast.dot(hit.Point, t.accent)
		raycast.text(hit.Point.Add(hit.Normal.Multiply(2)).AddXY(0, -0.7), "middle", "Normal", t.leaf)
	}
	raycast.dot(ray.Origin, t.ink)
	raycast.named(ray.Origin, "ray", 0, 0.8)

	sweep := paper(t, "Sweep(ball, motion, block)", area)
	ball := geom.Circ(geom.Pt(2.0, 2.0), 1.0)
	motion := geom.Vec(11.0, 5.0)
	block := geom.Bx(geom.Pt(8.0, 4.0), geom.Pt(11.0, 9.0))
	sweep.fill(block, t.coolSoft, t.cool)
	sweep.trace(ball, t.muted)
	sweep.trace(ball.Translate(motion), t.edge)
	sweep.arrow(ball.Center, ball.Center.Add(motion), t.muted, true)
	if hit, ok := geom.Sweep[float64](ball, motion, block); ok {
		at := ball.Translate(motion.Multiply(hit.Fraction))
		sweep.fill(at, t.accentSoft, t.accent)
		sweep.arrow(hit.Point, hit.Point.Add(hit.Normal.Multiply(2)), t.leaf, false)
		sweep.dot(hit.Point, t.accent)
		sweep.text(at.Center.AddXY(0, -1.6), "middle", fmt.Sprintf("Fraction %.2f", hit.Fraction), t.accent)
	}

	return []*panel{raycast, sweep}
}

func lines(t theme) []*panel {
	area := geom.Bx(geom.Pt(0.0, 0.0), geom.Pt(15.0, 10.0))
	p := paper(t, "", area)
	wall := geom.LineThrough(geom.Pt(1.0, 8.0), geom.Pt(13.0, 2.0))
	a := geom.Pt(10.0, 8.0)

	part, ok := wall.ClipBox(area)
	p.beam(part, ok, t.cool, false)
	p.dot(wall.Origin, t.cool)
	p.dot(geom.Pt(13.0, 2.0), t.cool)
	p.line(a, wall.Reflect(a), style{stroke: t.muted, width: 1, dashed: true})
	p.dot(a, t.ink)
	p.result(wall.Nearest(a), "Nearest(a)", 0.5, 0.4)
	p.result(wall.Reflect(a), "Reflect(a)", -0.5, 0)
	p.named(a, fmt.Sprintf("a, Side %d", wall.Side(a)), 0.5, 0)
	p.named(geom.Pt(3.0, 7.0), "wall", 0, -0.8)

	return []*panel{p}
}

func collisions(t theme) []*panel {
	area := geom.Bx(geom.Pt(0.0, 0.0), geom.Pt(10.0, 10.0))
	floor := geom.Bx(geom.Pt(1.0, 6.0), geom.Pt(9.0, 8.0))

	collide := func(caption string, shape geom.Shape[float64]) *panel {
		p := paper(t, caption, area)
		p.fill(floor, t.coolSoft, t.cool)
		p.fill(shape, t.accentSoft, t.accent)
		if contact, ok := geom.Collide(shape, geom.Shape[float64](floor)); ok {
			p.trace(moved(shape, contact.Separation()), t.accent)
			for at := range contact.Points.All() {
				p.arrow(at, at.Add(contact.Normal.Multiply(2)), t.leaf, false)
				p.dot(at, t.ink)
			}
		}

		return p
	}

	round := collide("Collide(ball, floor)", geom.Circ(geom.Pt(4.0, 5.0), 1.5))
	flat := collide("Collide(crate, floor)", geom.RectangleFromCenter(geom.Pt(5.0, 5.0), geom.Sz(3.0, 2.5)))

	nearest := paper(t, "Nearest(turret, target)", area)
	turret := geom.RegularPolygonWithOrientation(geom.Pt(2.5, 5.0), geom.SzU(1.8), 5, geom.OrientationPointyTop)
	target := geom.Cap(geom.Pt(7.0, 2.0), geom.Pt(8.0, 7.0), 1.0)
	nearest.fill(turret, t.coolSoft, t.cool)
	nearest.fill(target, t.leafSoft, t.leaf)
	on, at := geom.Nearest[float64](turret, target)
	nearest.line(on, at, style{stroke: t.accent, width: 2, dashed: true})
	nearest.dot(on, t.accent)
	nearest.dot(at, t.accent)
	nearest.text(on.Lerp(at, 0.5).AddXY(0, -0.7), "middle", fmt.Sprintf("%.2f", geom.Distance[float64](turret, target)), t.accent)

	return []*panel{round, flat, nearest}
}

func moved(shape geom.Shape[float64], by geom.Vector[float64]) geom.Shape[float64] {
	switch s := shape.(type) {
	case geom.Circle[float64]:
		return s.Translate(by)
	case geom.Rectangle[float64]:
		return s.Translate(by)
	default:
		panic(fmt.Sprintf("diagrams: cannot move %T", shape))
	}
}

func curves(t theme) []*panel {
	area := geom.Bx(geom.Pt(0.0, 0.0), geom.Pt(10.0, 10.0))

	ring := paper(t, "Ac(c, 4, -π/2, 1.3π)", area)
	arc := geom.Ac(geom.Pt(5.0, 5.0), 4.0, -math.Pi/2, 1.3*math.Pi)
	ring.trace(arc.Circle(), t.edge)
	ring.path(arc.Polyline(tolerance).Points, style{stroke: t.accent, width: 5})
	ring.result(arc.PointAt(0.5), "PointAt(0.5)", 0, -0.9)

	vision := paper(t, "cone.Contains(player)", area)
	cone := geom.Sect(geom.Pt(1.0, 5.0), 8.5, -math.Pi/6, math.Pi/3)
	vision.fill(cone, t.coolSoft, t.cool)
	for _, player := range []point{geom.Pt(6.0, 4.0), geom.Pt(8.0, 6.5), geom.Pt(5.0, 8.5), geom.Pt(8.0, 1.0)} {
		if cone.Contains(player) {
			vision.dot(player, t.accent)
		} else {
			vision.ring(player, t.muted)
		}
	}

	every := paper(t, "route.Every(1)", area)
	route := geom.PolLn([]point{{X: 1, Y: 8}, {X: 4, Y: 2}, {X: 6, Y: 7}, {X: 9, Y: 2}})
	every.path(route.Points, style{stroke: t.cool, width: 2})
	for at := range route.Every(1) {
		every.dot(at, t.accent)
	}

	flatten := paper(t, "curve.Polyline(0.2)", area)
	curve := geom.CubBez(geom.Pt(1.0, 9.0), geom.Pt(1.0, 1.0), geom.Pt(9.0, 1.0), geom.Pt(9.0, 9.0))
	coarse := curve.Polyline(0.2)
	flatten.path([]point{curve.Start, curve.Control1, curve.Control2, curve.End}, style{stroke: t.muted, width: 1, dashed: true})
	flatten.path(curve.Polyline(tolerance).Points, style{stroke: t.edge, width: 6})
	flatten.path(coarse.Points, style{stroke: t.cool, width: 2})
	for vertex := range coarse.Vertices() {
		flatten.dot(vertex, t.cool)
	}
	flatten.result(curve.PointAt(0.5), "PointAt(0.5)", 0, 0.9)

	return []*panel{ring, vision, every, flatten}
}

func coordinates(t theme) []*panel {
	axes := paper(t, "+Y down, +angle", geom.Bx(geom.Pt(0.0, 0.0), geom.Pt(10.0, 10.0)))
	origin := geom.Pt(1.0, 1.0)
	axes.arrow(origin, geom.Pt(9.0, 1.0), t.ink, false)
	axes.arrow(origin, geom.Pt(1.0, 9.0), t.ink, false)
	axes.text(geom.Pt(9.0, 1.8), "middle", "+X", t.ink)
	axes.text(geom.Pt(1.8, 9.0), "start", "+Y", t.ink)
	turn := geom.Ac(origin, 4.0, 0, math.Pi/3)
	axes.path(turn.Polyline(tolerance).Points, style{stroke: t.accent, width: 2})
	axes.arrow(turn.PointAt(0.85), turn.End(), t.accent, false)
	axes.text(turn.PointAt(0.5).Add(turn.PointAt(0.5).Subtract(origin).Resize(1)), "start", "+angle", t.accent)

	directions := paper(t, "Directions()", geom.Bx(geom.Pt(0.0, 0.0), geom.Pt(15.0, 10.0)))
	center := geom.Pt(7.5, 5.0)
	for _, direction := range geom.Directions() {
		directions.arrow(center, center.Add(direction.Vector(2.0)), t.cool, false)
		directions.text(center.Add(direction.Vector(3.1)), "middle", direction.String(), t.ink)
	}

	winding := paper(t, "Vertices wind clockwise", geom.Bx(geom.Pt(0.0, 0.0), geom.Pt(10.0, 10.0)))
	r := geom.RectangleFromCenter(geom.Pt(5.0, 5.0), geom.Sz(6.0, 4.0))
	winding.fill(r, t.coolSoft, t.cool)
	vertices := r.Polygon().Points
	for i, vertex := range vertices {
		next := vertices[(i+1)%len(vertices)]
		winding.arrow(vertex.Lerp(next, 0.3), vertex.Lerp(next, 0.7), t.accent, false)
		winding.text(r.Center.Add(vertex.Subtract(r.Center).Multiply(1.2)), "middle", fmt.Sprint(i), t.muted)
	}

	return []*panel{axes, directions, winding}
}

func matrices(t theme) []*panel {
	area := geom.Bx(geom.Pt(0.0, 0.0), geom.Pt(10.0, 10.0))
	arrow := geom.Pol([]point{{X: 3, Y: 4}, {X: 6, Y: 4}, {X: 6, Y: 3}, {X: 8, Y: 5}, {X: 6, Y: 7}, {X: 6, Y: 6}, {X: 3, Y: 6}})

	transformed := func(caption string, m geom.Matrix[float64]) *panel {
		p := paper(t, caption, area)
		p.trace(arrow, t.muted)
		p.fill(arrow.Transform(m), t.coolSoft, t.cool)
		p.dot(geom.Pt(5.0, 5.0), t.ink)

		return p
	}

	return []*panel{
		transformed("RotationMatrixAround(5, 5, π/4)", geom.RotationMatrixAround(5.0, 5.0, math.Pi/4)),
		transformed("ScaleMatrixAround(5, 5, 1.5, 0.5)", geom.ScaleMatrixAround(5.0, 5.0, 1.5, 0.5)),
		transformed("ShearMatrix(0.5, 0)", geom.ShearMatrix(0.5, 0.0)),
	}
}

func numbers(t theme) []*panel {
	p := paper(t, "", geom.Bx(geom.Pt(0.0, 0.0), geom.Pt(35.0, 10.0)))
	axis := func(y float64, ticks []float64) {
		p.line(geom.Pt(2.5, y), geom.Pt(32.5, y), style{stroke: t.ink, width: 1.5})
		for _, tick := range ticks {
			x := geom.Remap(tick, ticks[0], ticks[len(ticks)-1], 2.5, 32.5)
			p.line(geom.Pt(x, y-0.25), geom.Pt(x, y+0.25), style{stroke: t.ink, width: 1.5})
			p.text(geom.Pt(x, y+0.9), "middle", fmt.Sprint(tick), t.muted)
		}
	}
	axis(2.5, []float64{-1, -0.5, 0, 0.5, 1})
	axis(7.0, []float64{-30, -15, 0, 15, 30})
	stick := geom.Remap(0.5, -1.0, 1.0, 2.5, 32.5)
	angle := geom.Remap(geom.Remap(0.5, -1.0, 1.0, -30.0, 30.0), -30.0, 30.0, 2.5, 32.5)
	p.arrow(geom.Pt(stick, 2.8), geom.Pt(angle, 6.7), t.accent, false)
	p.dot(geom.Pt(stick, 2.5), t.accent)
	p.dot(geom.Pt(angle, 7.0), t.accent)
	p.text(geom.Pt(stick+0.6, 4.8), "start", "Remap(0.5, -1, 1, -30, 30)", t.accent)

	return []*panel{p}
}

func rounding(t theme) []*panel {
	p := paper(t, "", geom.Bx(geom.Pt(0.0, 0.0), geom.Pt(10.0, 10.0)))
	for _, from := range []point{geom.Pt(1.4, 2.6), geom.Pt(3.5, 1.5), geom.Pt(6.6, 3.2), geom.Pt(2.5, 7.5), geom.Pt(7.2, 6.7)} {
		to := from.Int().Float()
		p.ring(from, t.ink)
		p.result(to, to.Int().String(), 0.3, 0.5)
	}

	return []*panel{p}
}

func pixels(t theme) []*panel {
	p := paper(t, "", geom.Bx(geom.Pt(0.0, 0.0), geom.Pt(10.0, 10.0)))
	f := geom.RectangleFromMin(geom.Pt(1.6, 2.3), geom.Sz(5.8, 5.6))
	r := f.Int()
	img := r.Rectangle()
	for x := img.Min.X; x < img.Max.X; x++ {
		for y := img.Min.Y; y < img.Max.Y; y++ {
			p.area(geom.BoxFromMin(geom.Pt(float64(x)+0.08, float64(y)+0.08), geom.SzU(0.84)).Polygon().Points, style{fill: t.accentSoft})
		}
	}
	p.trace(f, t.ink)
	p.area(r.Float().Polygon().Points, style{stroke: t.cool, width: 2})
	bounds := r.Bounds()
	for x := bounds.Min.X; x <= bounds.Max.X; x++ {
		for y := bounds.Min.Y; y <= bounds.Max.Y; y++ {
			p.dot(geom.Pt(float64(x), float64(y)), t.cool)
		}
	}
	p.named(f.Anchor(geom.TopLeft), "f", -0.3, -0.4)
	p.named(r.Float().Anchor(geom.BottomRight), "f.Int()", 0.3, 0.5)

	return []*panel{p}
}

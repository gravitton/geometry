package geom

import (
	"fmt"
	"testing"
)

func TestAssertPoint(t *testing.T) {
	testAssertPoint(t, Pt(1, 2), 1, 2, true)
	testAssertPoint(t, Pt(1, 2), 9, 2, false)
	testAssertPoint(t, Pt(1, 2), 1, 9, false)
	testAssertPoint(t, Pt(1.0, 2.0), 1.0000001, 2.0, true) // within Delta
}

func TestAssertVector(t *testing.T) {
	testAssertVector(t, Vec(1, 2), 1, 2, true)
	testAssertVector(t, Vec(1, 2), 9, 2, false)
	testAssertVector(t, Vec(1, 2), 1, 9, false)
}

func TestAssertSize(t *testing.T) {
	testAssertSize(t, Sz(1, 2), 1, 2, true)
	testAssertSize(t, Sz(1, 2), 9, 2, false)
	testAssertSize(t, Sz(1, 2), 1, 9, false)
}

func TestAssertCircle(t *testing.T) {
	testAssertCircle(t, Circ(Pt(1, 2), 3), 1, 2, 3, true)
	testAssertCircle(t, Circ(Pt(1, 2), 3), 9, 2, 3, false) // center differs
	testAssertCircle(t, Circ(Pt(1, 2), 3), 1, 2, 9, false) // radius differs
}

func TestAssertLine(t *testing.T) {
	testAssertLine(t, Ln(Pt(1, 2), Pt(3, 4)), 1, 2, 3, 4, true)
	testAssertLine(t, Ln(Pt(1, 2), Pt(3, 4)), 9, 2, 3, 4, false) // start differs
	testAssertLine(t, Ln(Pt(1, 2), Pt(3, 4)), 1, 2, 9, 4, false) // end differs
}

func TestAssertRect(t *testing.T) {
	testAssertRect(t, Rect(Pt(1, 2), Sz(3, 4)), 1, 2, 3, 4, true)
	testAssertRect(t, Rect(Pt(1, 2), Sz(3, 4)), 9, 2, 3, 4, false) // center differs
	testAssertRect(t, Rect(Pt(1, 2), Sz(3, 4)), 1, 2, 9, 4, false) // size differs
}

func TestAssertVertices(t *testing.T) {
	points := []Point[int]{Pt(1, 2), Pt(3, 4)}

	testAssertVertices(t, points, points, true)
	testAssertVertices(t, points, []Point[int]{Pt(1, 2)}, false)           // length differs
	testAssertVertices(t, points, []Point[int]{Pt(1, 2), Pt(9, 4)}, false) // element differs
	testAssertVertices(t, []Point[int]{}, []Point[int]{}, true)
}

func TestAssertPolygon(t *testing.T) {
	points := []Point[int]{Pt(1, 2), Pt(3, 4)}

	testAssertPolygon(t, Pol(points), points, true)
	testAssertPolygon(t, Pol(points), []Point[int]{Pt(1, 2), Pt(9, 4)}, false)
}

func TestAssertRegularPolygon(t *testing.T) {
	rp := RegPol(Pt(1, 2), Sz(3, 4), 5, 0.5)

	testAssertRegularPolygon(t, rp, 1, 2, 3, 4, 5, 0.5, true)
	testAssertRegularPolygon(t, rp, 9, 2, 3, 4, 5, 0.5, false) // center differs
	testAssertRegularPolygon(t, rp, 1, 2, 9, 4, 5, 0.5, false) // size differs
	testAssertRegularPolygon(t, rp, 1, 2, 3, 4, 9, 0.5, false) // n differs
	testAssertRegularPolygon(t, rp, 1, 2, 3, 4, 5, 9.5, false) // angle differs
}

func TestAssertPadding(t *testing.T) {
	p := Pad(1, 2, 3, 4)

	testAssertPadding(t, p, 1, 2, 3, 4, true)
	testAssertPadding(t, p, 9, 2, 3, 4, false) // top differs
	testAssertPadding(t, p, 1, 9, 3, 4, false) // right differs
	testAssertPadding(t, p, 1, 2, 9, 4, false) // bottom differs
	testAssertPadding(t, p, 1, 2, 3, 9, false) // left differs
}

func TestAssertMatrix(t *testing.T) {
	m := Mat(1, 2, 3, 4, 5, 6)

	testAssertMatrix(t, m, m, true)
	testAssertMatrix(t, m, Mat(9, 2, 3, 4, 5, 6), false) // A differs
	testAssertMatrix(t, m, Mat(1, 9, 3, 4, 5, 6), false) // B differs
	testAssertMatrix(t, m, Mat(1, 2, 9, 4, 5, 6), false) // C differs
	testAssertMatrix(t, m, Mat(1, 2, 3, 9, 5, 6), false) // D differs
	testAssertMatrix(t, m, Mat(1, 2, 3, 4, 9, 6), false) // E differs
	testAssertMatrix(t, m, Mat(1, 2, 3, 4, 5, 9), false) // F differs
}

type logger struct {
	LastError string
}

func (m *logger) Helper() {
}

func (m *logger) Errorf(format string, args ...any) {
	m.LastError = fmt.Sprintf(format, args...)
}

func (m *logger) Clear() {
	m.LastError = ""
}

var tt = &logger{}

func testAssertPoint[T Number](t *testing.T, p Point[T], x, y T, result bool) {
	t.Helper()

	tt.Clear()
	if AssertPoint(tt, p, x, y) != result {
		t.Errorf("AssertPoint(%#v,%#v,%#v) should return %#v: %s", p, x, y, result, tt.LastError)
	}
}

func testAssertVector[T Number](t *testing.T, v Vector[T], x, y T, result bool) {
	t.Helper()

	tt.Clear()
	if AssertVector(tt, v, x, y) != result {
		t.Errorf("AssertVector(%#v,%#v,%#v) should return %#v: %s", v, x, y, result, tt.LastError)
	}
}

func testAssertSize[T Number](t *testing.T, s Size[T], w, h T, result bool) {
	t.Helper()

	tt.Clear()
	if AssertSize(tt, s, w, h) != result {
		t.Errorf("AssertSize(%#v,%#v,%#v) should return %#v: %s", s, w, h, result, tt.LastError)
	}
}

func testAssertCircle[T Number](t *testing.T, c Circle[T], x, y, radius T, result bool) {
	t.Helper()

	tt.Clear()
	if AssertCircle(tt, c, x, y, radius) != result {
		t.Errorf("AssertCircle(%#v,%#v,%#v,%#v) should return %#v: %s", c, x, y, radius, result, tt.LastError)
	}
}

func testAssertLine[T Number](t *testing.T, l Line[T], sx, sy, ex, ey T, result bool) {
	t.Helper()

	tt.Clear()
	if AssertLine(tt, l, sx, sy, ex, ey) != result {
		t.Errorf("AssertLine(%#v,%#v,%#v,%#v,%#v) should return %#v: %s", l, sx, sy, ex, ey, result, tt.LastError)
	}
}

func testAssertRect[T Number](t *testing.T, r Rectangle[T], cx, cy, w, h T, result bool) {
	t.Helper()

	tt.Clear()
	if AssertRect(tt, r, cx, cy, w, h) != result {
		t.Errorf("AssertRect(%#v,%#v,%#v,%#v,%#v) should return %#v: %s", r, cx, cy, w, h, result, tt.LastError)
	}
}

func testAssertVertices[T Number](t *testing.T, vertices, points []Point[T], result bool) {
	t.Helper()

	tt.Clear()
	if AssertVertices(tt, vertices, points) != result {
		t.Errorf("AssertVertices(%#v,%#v) should return %#v: %s", vertices, points, result, tt.LastError)
	}
}

func testAssertPolygon[T Number](t *testing.T, p Polygon[T], vertices []Point[T], result bool) {
	t.Helper()

	tt.Clear()
	if AssertPolygon(tt, p, vertices) != result {
		t.Errorf("AssertPolygon(%#v,%#v) should return %#v: %s", p, vertices, result, tt.LastError)
	}
}

func testAssertRegularPolygon[T Number](t *testing.T, p RegularPolygon[T], x, y, w, h T, n int, angle float64, result bool) {
	t.Helper()

	tt.Clear()
	if AssertRegularPolygon(tt, p, x, y, w, h, n, angle) != result {
		t.Errorf("AssertRegularPolygon(%#v,%#v,%#v,%#v,%#v,%#v,%#v) should return %#v: %s", p, x, y, w, h, n, angle, result, tt.LastError)
	}
}

func testAssertPadding[T Number](t *testing.T, p Padding[T], top, right, bottom, left T, result bool) {
	t.Helper()

	tt.Clear()
	if AssertPadding(tt, p, top, right, bottom, left) != result {
		t.Errorf("AssertPadding(%#v,%#v,%#v,%#v,%#v) should return %#v: %s", p, top, right, bottom, left, result, tt.LastError)
	}
}

func testAssertMatrix[T Number](t *testing.T, m, expected Matrix[T], result bool) {
	t.Helper()

	tt.Clear()
	if AssertMatrix(tt, m, expected) != result {
		t.Errorf("AssertMatrix(%#v,%#v) should return %#v: %s", m, expected, result, tt.LastError)
	}
}

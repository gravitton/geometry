package geom

import (
	"fmt"
	"testing"
)

func TestAssertPoint(t *testing.T) {
	testAssert(t, AssertPoint, Pt(1, 2), Pt(1, 2), true)
	testAssert(t, AssertPoint, Pt(1, 2), Pt(9, 2), false)
	testAssert(t, AssertPoint, Pt(1, 2), Pt(1, 9), false)
	testAssert(t, AssertPoint, Pt(1.0, 2.0), Pt(1.0000001, 2.0), true) // within Delta
}

func TestAssertVector(t *testing.T) {
	testAssert(t, AssertVector, Vec(1, 2), Vec(1, 2), true)
	testAssert(t, AssertVector, Vec(1, 2), Vec(9, 2), false)
	testAssert(t, AssertVector, Vec(1, 2), Vec(1, 9), false)
}

func TestAssertSize(t *testing.T) {
	testAssert(t, AssertSize, Sz(1, 2), Sz(1, 2), true)
	testAssert(t, AssertSize, Sz(1, 2), Sz(9, 2), false)
	testAssert(t, AssertSize, Sz(1, 2), Sz(1, 9), false)
}

func TestAssertCircle(t *testing.T) {
	c := Circ(Pt(1, 2), 3)

	testAssert(t, AssertCircle, c, Circ(Pt(1, 2), 3), true)
	testAssert(t, AssertCircle, c, Circ(Pt(9, 2), 3), false) // center differs
	testAssert(t, AssertCircle, c, Circ(Pt(1, 2), 9), false) // radius differs
}

func TestAssertLine(t *testing.T) {
	l := Ln(Pt(1, 2), Pt(3, 4))

	testAssert(t, AssertLine, l, Ln(Pt(1, 2), Pt(3, 4)), true)
	testAssert(t, AssertLine, l, Ln(Pt(9, 2), Pt(3, 4)), false) // start differs
	testAssert(t, AssertLine, l, Ln(Pt(1, 2), Pt(9, 4)), false) // end differs
}

func TestAssertRect(t *testing.T) {
	r := Rect(Pt(1, 2), Sz(3, 4))

	testAssert(t, AssertRect, r, Rect(Pt(1, 2), Sz(3, 4)), true)
	testAssert(t, AssertRect, r, Rect(Pt(9, 2), Sz(3, 4)), false) // center differs
	testAssert(t, AssertRect, r, Rect(Pt(1, 2), Sz(9, 4)), false) // size differs
}

func TestAssertVertices(t *testing.T) {
	points := []Point[int]{Pt(1, 2), Pt(3, 4)}

	testAssert(t, AssertVertices, points, points, true)
	testAssert(t, AssertVertices, points, []Point[int]{Pt(1, 2)}, false)           // length differs
	testAssert(t, AssertVertices, points, []Point[int]{Pt(1, 2), Pt(9, 4)}, false) // element differs
	testAssert(t, AssertVertices, []Point[int]{}, []Point[int]{}, true)
}

func TestAssertPolygon(t *testing.T) {
	p := Pol([]Point[int]{Pt(1, 2), Pt(3, 4)})

	testAssert(t, AssertPolygon, p, Pol([]Point[int]{Pt(1, 2), Pt(3, 4)}), true)
	testAssert(t, AssertPolygon, p, Pol([]Point[int]{Pt(1, 2), Pt(9, 4)}), false)
}

func TestAssertRegularPolygon(t *testing.T) {
	p := RegPol(Pt(1, 2), Sz(3, 4), 5, 0.5)

	testAssert(t, AssertRegularPolygon, p, RegPol(Pt(1, 2), Sz(3, 4), 5, 0.5), true)
	testAssert(t, AssertRegularPolygon, p, RegPol(Pt(9, 2), Sz(3, 4), 5, 0.5), false) // center differs
	testAssert(t, AssertRegularPolygon, p, RegPol(Pt(1, 2), Sz(9, 4), 5, 0.5), false) // size differs
	testAssert(t, AssertRegularPolygon, p, RegPol(Pt(1, 2), Sz(3, 4), 9, 0.5), false) // n differs
	testAssert(t, AssertRegularPolygon, p, RegPol(Pt(1, 2), Sz(3, 4), 5, 9.5), false) // angle differs
}

func TestAssertPadding(t *testing.T) {
	p := Pad(1, 2, 3, 4)

	testAssert(t, AssertPadding, p, Pad(1, 2, 3, 4), true)
	testAssert(t, AssertPadding, p, Pad(9, 2, 3, 4), false) // top differs
	testAssert(t, AssertPadding, p, Pad(1, 9, 3, 4), false) // right differs
	testAssert(t, AssertPadding, p, Pad(1, 2, 9, 4), false) // bottom differs
	testAssert(t, AssertPadding, p, Pad(1, 2, 3, 9), false) // left differs
}

func TestAssertMatrix(t *testing.T) {
	m := Mat(1, 2, 3, 4, 5, 6)

	testAssert(t, AssertMatrix, m, Mat(1, 2, 3, 4, 5, 6), true)
	testAssert(t, AssertMatrix, m, Mat(9, 2, 3, 4, 5, 6), false) // A differs
	testAssert(t, AssertMatrix, m, Mat(1, 9, 3, 4, 5, 6), false) // B differs
	testAssert(t, AssertMatrix, m, Mat(1, 2, 9, 4, 5, 6), false) // C differs
	testAssert(t, AssertMatrix, m, Mat(1, 2, 3, 9, 5, 6), false) // D differs
	testAssert(t, AssertMatrix, m, Mat(1, 2, 3, 4, 9, 6), false) // E differs
	testAssert(t, AssertMatrix, m, Mat(1, 2, 3, 4, 5, 9), false) // F differs
}

type logger struct {
	LastError string
}

func (m *logger) Helper() {
}

func (m *logger) Errorf(format string, args ...any) {
	m.LastError = fmt.Sprintf(format, args...)
}

// testAssert runs one of the Assert helpers against a fresh recorder and checks it
// reported the expected outcome. Every helper takes (actual, expected) of the same
// type, so one signature covers them all.
func testAssert[V any](t *testing.T, assertion func(Testing, V, V, ...string) bool, actual, expected V, result bool) {
	t.Helper()

	recorder := &logger{}
	if assertion(recorder, actual, expected) != result {
		t.Errorf("assert(%#v, %#v) should return %#v: %s", actual, expected, result, recorder.LastError)
	}
}

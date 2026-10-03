package geomtest

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gravitton/assert"
	geom "github.com/gravitton/geometry"
)

func TestAssertNumber(t *testing.T) {
	t.Run("int is exact", func(t *testing.T) {
		assertHelper(t, AssertNumber, 1, 1, true)
		assertHelper(t, AssertNumber, 1, 9, false)
	})
	t.Run("float is relative to the magnitude", func(t *testing.T) {
		assertHelper(t, AssertNumber, 1.0, 1.0000001, true)
		assertHelper(t, AssertNumber, 1.0, 1.1, false)

		// the tolerance scales, so a large value tolerates a large absolute difference
		assertHelper(t, AssertNumber, 1e6, 1e6+0.5, true)
		assertHelper(t, AssertNumber, 1e6, 1e6+2.0, false)
	})
}

func TestAssertAngle(t *testing.T) {
	t.Run("equal", func(t *testing.T) {
		assertHelper(t, AssertAngle, geom.Pi/2, geom.Pi/2, true)
	})
	t.Run("a full turn does not matter", func(t *testing.T) {
		assertHelper(t, AssertAngle, geom.Pi/2, geom.Pi/2+2*geom.Pi, true)
		assertHelper(t, AssertAngle, 0.0, -2*geom.Pi, true)
	})
	t.Run("differs", func(t *testing.T) {
		assertHelper(t, AssertAngle, geom.Pi/2, geom.Pi, false)
	})
}

func TestAssertPoint(t *testing.T) {
	t.Run("equal", func(t *testing.T) {
		assertHelper(t, AssertPoint, geom.Pt(1, 2), geom.Pt(1, 2), true)
	})
	t.Run("one component differs", func(t *testing.T) {
		assertHelper(t, AssertPoint, geom.Pt(1, 2), geom.Pt(9, 2), false)
		assertHelper(t, AssertPoint, geom.Pt(1, 2), geom.Pt(1, 9), false)
	})
	t.Run("float within the tolerance", func(t *testing.T) {
		assertHelper(t, AssertPoint, geom.Pt(1.0, 2.0), geom.Pt(1.0000001, 2.0), true)
	})
}

func TestAssertVector(t *testing.T) {
	t.Run("equal", func(t *testing.T) {
		assertHelper(t, AssertVector, geom.Vec(1, 2), geom.Vec(1, 2), true)
	})
	t.Run("one component differs", func(t *testing.T) {
		assertHelper(t, AssertVector, geom.Vec(1, 2), geom.Vec(9, 2), false)
		assertHelper(t, AssertVector, geom.Vec(1, 2), geom.Vec(1, 9), false)
	})
}

func TestAssertSize(t *testing.T) {
	t.Run("equal", func(t *testing.T) {
		assertHelper(t, AssertSize, geom.Sz(1, 2), geom.Sz(1, 2), true)
	})
	t.Run("one component differs", func(t *testing.T) {
		assertHelper(t, AssertSize, geom.Sz(1, 2), geom.Sz(9, 2), false)
		assertHelper(t, AssertSize, geom.Sz(1, 2), geom.Sz(1, 9), false)
	})
}

func TestAssertCircle(t *testing.T) {
	c := geom.Circ(geom.Pt(1, 2), 3)

	t.Run("equal", func(t *testing.T) {
		assertHelper(t, AssertCircle, c, geom.Circ(geom.Pt(1, 2), 3), true)
	})
	t.Run("center differs", func(t *testing.T) {
		assertHelper(t, AssertCircle, c, geom.Circ(geom.Pt(9, 2), 3), false)
	})
	t.Run("radius differs", func(t *testing.T) {
		assertHelper(t, AssertCircle, c, geom.Circ(geom.Pt(1, 2), 9), false)
	})
}

func TestAssertEllipse(t *testing.T) {
	e := geom.Ell(geom.Pt(1, 2), geom.Sz(3, 4), 0)

	t.Run("equal", func(t *testing.T) {
		assertHelper(t, AssertEllipse, e, geom.Ell(geom.Pt(1, 2), geom.Sz(3, 4), 0), true)
	})
	t.Run("center differs", func(t *testing.T) {
		assertHelper(t, AssertEllipse, e, geom.Ell(geom.Pt(9, 2), geom.Sz(3, 4), 0), false)
	})
	t.Run("size differs", func(t *testing.T) {
		assertHelper(t, AssertEllipse, e, geom.Ell(geom.Pt(1, 2), geom.Sz(9, 4), 0), false)
	})
	t.Run("angle is compared normalized", func(t *testing.T) {
		assertHelper(t, AssertEllipse, e.Rotate(geom.Pi/2), e.Rotate(-3*geom.Pi/2), true)
		assertHelper(t, AssertEllipse, e, e.Rotate(geom.Pi/2), false)
	})
}

func TestAssertSegment(t *testing.T) {
	s := geom.Seg(geom.Pt(1, 2), geom.Pt(3, 4))

	t.Run("equal", func(t *testing.T) {
		assertHelper(t, AssertSegment, s, geom.Seg(geom.Pt(1, 2), geom.Pt(3, 4)), true)
	})
	t.Run("start differs", func(t *testing.T) {
		assertHelper(t, AssertSegment, s, geom.Seg(geom.Pt(9, 2), geom.Pt(3, 4)), false)
	})
	t.Run("end differs", func(t *testing.T) {
		assertHelper(t, AssertSegment, s, geom.Seg(geom.Pt(1, 2), geom.Pt(9, 4)), false)
	})
}

func TestAssertSegments(t *testing.T) {
	segments := []geom.Segment[int]{geom.Seg(geom.Pt(1, 2), geom.Pt(3, 4)), geom.Seg(geom.Pt(5, 6), geom.Pt(7, 8))}

	t.Run("equal", func(t *testing.T) {
		assertHelper(t, AssertSegments, segments, segments, true)
	})
	t.Run("length differs", func(t *testing.T) {
		assertHelper(t, AssertSegments, segments, segments[:1], false)
	})
	t.Run("element differs", func(t *testing.T) {
		assertHelper(t, AssertSegments, segments, []geom.Segment[int]{segments[0], geom.Seg(geom.Pt(9, 6), geom.Pt(7, 8))}, false)
	})
	t.Run("empty slices are equal", func(t *testing.T) {
		assertHelper(t, AssertSegments, []geom.Segment[int]{}, []geom.Segment[int]{}, true)
	})
}

func TestAssertRay(t *testing.T) {
	r := geom.Ry(geom.Pt(1, 2), geom.Vec(3, 4))

	t.Run("equal", func(t *testing.T) {
		assertHelper(t, AssertRay, r, geom.Ry(geom.Pt(1, 2), geom.Vec(3, 4)), true)
	})
	t.Run("origin differs", func(t *testing.T) {
		assertHelper(t, AssertRay, r, geom.Ry(geom.Pt(9, 2), geom.Vec(3, 4)), false)
	})
	t.Run("direction differs", func(t *testing.T) {
		assertHelper(t, AssertRay, r, geom.Ry(geom.Pt(1, 2), geom.Vec(3, 9)), false)
	})
}

func TestAssertRectangle(t *testing.T) {
	r := geom.Rect(geom.Pt(1, 2), geom.Sz(3, 4))

	t.Run("equal", func(t *testing.T) {
		assertHelper(t, AssertRectangle, r, geom.Rect(geom.Pt(1, 2), geom.Sz(3, 4)), true)
	})
	t.Run("center differs", func(t *testing.T) {
		assertHelper(t, AssertRectangle, r, geom.Rect(geom.Pt(9, 2), geom.Sz(3, 4)), false)
	})
	t.Run("size differs", func(t *testing.T) {
		assertHelper(t, AssertRectangle, r, geom.Rect(geom.Pt(1, 2), geom.Sz(9, 4)), false)
	})
	t.Run("angle is compared normalized", func(t *testing.T) {
		assertHelper(t, AssertRectangle, r.Rotate(geom.Pi/2), r.Rotate(-3*geom.Pi/2), true)
		assertHelper(t, AssertRectangle, r, r.Rotate(geom.Pi/2), false)
	})
}

func TestAssertBox(t *testing.T) {
	b := geom.BoxFromMinMax(geom.Pt(1, 2), geom.Pt(3, 4))

	t.Run("equal", func(t *testing.T) {
		assertHelper(t, AssertBox, b, geom.BoxFromMinMax(geom.Pt(1, 2), geom.Pt(3, 4)), true)
	})
	t.Run("min differs", func(t *testing.T) {
		assertHelper(t, AssertBox, b, geom.BoxFromMinMax(geom.Pt(0, 2), geom.Pt(3, 4)), false)
	})
	t.Run("max differs", func(t *testing.T) {
		assertHelper(t, AssertBox, b, geom.BoxFromMinMax(geom.Pt(1, 2), geom.Pt(3, 9)), false)
	})
}

func TestAssertPolygon(t *testing.T) {
	p := geom.Pol([]geom.Point[int]{geom.Pt(1, 2), geom.Pt(3, 4)})

	t.Run("equal", func(t *testing.T) {
		assertHelper(t, AssertPolygon, p, geom.Pol([]geom.Point[int]{geom.Pt(1, 2), geom.Pt(3, 4)}), true)
	})
	t.Run("vertex differs", func(t *testing.T) {
		assertHelper(t, AssertPolygon, p, geom.Pol([]geom.Point[int]{geom.Pt(1, 2), geom.Pt(9, 4)}), false)
	})
}

func TestAssertVertices(t *testing.T) {
	points := []geom.Point[int]{geom.Pt(1, 2), geom.Pt(3, 4)}

	t.Run("equal", func(t *testing.T) {
		assertHelper(t, AssertVertices, points, points, true)
	})
	t.Run("length differs", func(t *testing.T) {
		assertHelper(t, AssertVertices, points, []geom.Point[int]{geom.Pt(1, 2)}, false)
	})
	t.Run("element differs", func(t *testing.T) {
		assertHelper(t, AssertVertices, points, []geom.Point[int]{geom.Pt(1, 2), geom.Pt(9, 4)}, false)
	})
	t.Run("empty slices are equal", func(t *testing.T) {
		assertHelper(t, AssertVertices, []geom.Point[int]{}, []geom.Point[int]{}, true)
	})
}

func TestAssertRegularPolygon(t *testing.T) {
	p := geom.RegPol(geom.Pt(1, 2), geom.Sz(3, 4), 5, 0.5, 0)

	t.Run("equal", func(t *testing.T) {
		assertHelper(t, AssertRegularPolygon, p, geom.RegPol(geom.Pt(1, 2), geom.Sz(3, 4), 5, 0.5, 0), true)
	})
	t.Run("angle is compared normalized", func(t *testing.T) {
		assertHelper(t, AssertRegularPolygon, p, geom.RegPol(geom.Pt(1, 2), geom.Sz(3, 4), 5, 0.5+2*geom.Pi, 0), true)
		assertHelper(t, AssertRegularPolygon, geom.RegPol(geom.Pt(1, 2), geom.Sz(3, 4), 5, 0, 0), geom.RegPol(geom.Pt(1, 2), geom.Sz(3, 4), 5, -1e-9, 0), true)
	})
	t.Run("one field differs", func(t *testing.T) {
		assertHelper(t, AssertRegularPolygon, p, geom.RegPol(geom.Pt(9, 2), geom.Sz(3, 4), 5, 0.5, 0), false, "center")
		assertHelper(t, AssertRegularPolygon, p, geom.RegPol(geom.Pt(1, 2), geom.Sz(9, 4), 5, 0.5, 0), false, "size")
		assertHelper(t, AssertRegularPolygon, p, geom.RegPol(geom.Pt(1, 2), geom.Sz(3, 4), 9, 0.5, 0), false, "n")
		assertHelper(t, AssertRegularPolygon, p, geom.RegPol(geom.Pt(1, 2), geom.Sz(3, 4), 5, 9.5, 0), false, "angle")
		assertHelper(t, AssertRegularPolygon, p, geom.RegPol(geom.Pt(1, 2), geom.Sz(3, 4), 5, 0.5, 1), false, "phase")
	})
}

func TestAssertPadding(t *testing.T) {
	p := geom.Pad(1, 2, 3, 4)

	t.Run("equal", func(t *testing.T) {
		assertHelper(t, AssertPadding, p, geom.Pad(1, 2, 3, 4), true)
	})
	t.Run("one side differs", func(t *testing.T) {
		assertHelper(t, AssertPadding, p, geom.Pad(9, 2, 3, 4), false, "top")
		assertHelper(t, AssertPadding, p, geom.Pad(1, 9, 3, 4), false, "right")
		assertHelper(t, AssertPadding, p, geom.Pad(1, 2, 9, 4), false, "bottom")
		assertHelper(t, AssertPadding, p, geom.Pad(1, 2, 3, 9), false, "left")
	})
}

func TestAssertMatrix(t *testing.T) {
	m := geom.Mat(1, 2, 3, 4, 5, 6)

	t.Run("equal", func(t *testing.T) {
		assertHelper(t, AssertMatrix, m, geom.Mat(1, 2, 3, 4, 5, 6), true)
	})
	t.Run("one component differs", func(t *testing.T) {
		assertHelper(t, AssertMatrix, m, geom.Mat(9, 2, 3, 4, 5, 6), false, "A")
		assertHelper(t, AssertMatrix, m, geom.Mat(1, 9, 3, 4, 5, 6), false, "B")
		assertHelper(t, AssertMatrix, m, geom.Mat(1, 2, 9, 4, 5, 6), false, "C")
		assertHelper(t, AssertMatrix, m, geom.Mat(1, 2, 3, 9, 5, 6), false, "D")
		assertHelper(t, AssertMatrix, m, geom.Mat(1, 2, 3, 4, 9, 6), false, "E")
		assertHelper(t, AssertMatrix, m, geom.Mat(1, 2, 3, 4, 5, 9), false, "F")
	})
}

// assertHelper runs one of the Assert helpers against a fresh recorder and checks it
// reported the expected outcome. Every helper takes (actual, expected) of the same
// type, so one signature covers them all.
func assertHelper[V any](t *testing.T, assertion func(assert.Testing, V, V, ...string) bool, actual, expected V, result bool, labels ...string) {
	t.Helper()

	recorder := &logger{}
	if assertion(recorder, actual, expected) != result {
		t.Errorf("%sassert(%#v, %#v) should return %#v: %s", strings.Join(labels, ""), actual, expected, result, recorder.LastError)
	}
}

type logger struct {
	LastError string
}

func (m *logger) Helper() {
}

func (m *logger) Errorf(format string, args ...any) {
	m.LastError = fmt.Sprintf(format, args...)
}

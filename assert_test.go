package geom

import (
	"fmt"
	"strings"
	"testing"
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

func TestAssertPoint(t *testing.T) {
	t.Run("equal", func(t *testing.T) {
		assertHelper(t, AssertPoint, Pt(1, 2), Pt(1, 2), true)
	})
	t.Run("one component differs", func(t *testing.T) {
		assertHelper(t, AssertPoint, Pt(1, 2), Pt(9, 2), false)
		assertHelper(t, AssertPoint, Pt(1, 2), Pt(1, 9), false)
	})
	t.Run("float within the tolerance", func(t *testing.T) {
		assertHelper(t, AssertPoint, Pt(1.0, 2.0), Pt(1.0000001, 2.0), true)
	})
}

func TestAssertVector(t *testing.T) {
	t.Run("equal", func(t *testing.T) {
		assertHelper(t, AssertVector, Vec(1, 2), Vec(1, 2), true)
	})
	t.Run("one component differs", func(t *testing.T) {
		assertHelper(t, AssertVector, Vec(1, 2), Vec(9, 2), false)
		assertHelper(t, AssertVector, Vec(1, 2), Vec(1, 9), false)
	})
}

func TestAssertSize(t *testing.T) {
	t.Run("equal", func(t *testing.T) {
		assertHelper(t, AssertSize, Sz(1, 2), Sz(1, 2), true)
	})
	t.Run("one component differs", func(t *testing.T) {
		assertHelper(t, AssertSize, Sz(1, 2), Sz(9, 2), false)
		assertHelper(t, AssertSize, Sz(1, 2), Sz(1, 9), false)
	})
}

func TestAssertCircle(t *testing.T) {
	c := Circ(Pt(1, 2), 3)

	t.Run("equal", func(t *testing.T) {
		assertHelper(t, AssertCircle, c, Circ(Pt(1, 2), 3), true)
	})
	t.Run("center differs", func(t *testing.T) {
		assertHelper(t, AssertCircle, c, Circ(Pt(9, 2), 3), false)
	})
	t.Run("radius differs", func(t *testing.T) {
		assertHelper(t, AssertCircle, c, Circ(Pt(1, 2), 9), false)
	})
}

func TestAssertLine(t *testing.T) {
	l := Ln(Pt(1, 2), Pt(3, 4))

	t.Run("equal", func(t *testing.T) {
		assertHelper(t, AssertLine, l, Ln(Pt(1, 2), Pt(3, 4)), true)
	})
	t.Run("start differs", func(t *testing.T) {
		assertHelper(t, AssertLine, l, Ln(Pt(9, 2), Pt(3, 4)), false)
	})
	t.Run("end differs", func(t *testing.T) {
		assertHelper(t, AssertLine, l, Ln(Pt(1, 2), Pt(9, 4)), false)
	})
}

func TestAssertRect(t *testing.T) {
	r := Rect(Pt(1, 2), Sz(3, 4))

	t.Run("equal", func(t *testing.T) {
		assertHelper(t, AssertRect, r, Rect(Pt(1, 2), Sz(3, 4)), true)
	})
	t.Run("center differs", func(t *testing.T) {
		assertHelper(t, AssertRect, r, Rect(Pt(9, 2), Sz(3, 4)), false)
	})
	t.Run("size differs", func(t *testing.T) {
		assertHelper(t, AssertRect, r, Rect(Pt(1, 2), Sz(9, 4)), false)
	})
}

func TestAssertPolygon(t *testing.T) {
	p := Pol([]Point[int]{Pt(1, 2), Pt(3, 4)})

	t.Run("equal", func(t *testing.T) {
		assertHelper(t, AssertPolygon, p, Pol([]Point[int]{Pt(1, 2), Pt(3, 4)}), true)
	})
	t.Run("vertex differs", func(t *testing.T) {
		assertHelper(t, AssertPolygon, p, Pol([]Point[int]{Pt(1, 2), Pt(9, 4)}), false)
	})
}

func TestAssertVertices(t *testing.T) {
	points := []Point[int]{Pt(1, 2), Pt(3, 4)}

	t.Run("equal", func(t *testing.T) {
		assertHelper(t, AssertVertices, points, points, true)
	})
	t.Run("length differs", func(t *testing.T) {
		assertHelper(t, AssertVertices, points, []Point[int]{Pt(1, 2)}, false)
	})
	t.Run("element differs", func(t *testing.T) {
		assertHelper(t, AssertVertices, points, []Point[int]{Pt(1, 2), Pt(9, 4)}, false)
	})
	t.Run("empty slices are equal", func(t *testing.T) {
		assertHelper(t, AssertVertices, []Point[int]{}, []Point[int]{}, true)
	})
}

func TestAssertRegularPolygon(t *testing.T) {
	p := RegPol(Pt(1, 2), Sz(3, 4), 5, 0.5)

	t.Run("equal", func(t *testing.T) {
		assertHelper(t, AssertRegularPolygon, p, RegPol(Pt(1, 2), Sz(3, 4), 5, 0.5), true)
	})
	t.Run("angle is compared normalized", func(t *testing.T) {
		assertHelper(t, AssertRegularPolygon, p, RegPol(Pt(1, 2), Sz(3, 4), 5, 0.5+2*Pi), true)
	})
	t.Run("one field differs", func(t *testing.T) {
		assertHelper(t, AssertRegularPolygon, p, RegPol(Pt(9, 2), Sz(3, 4), 5, 0.5), false, "center")
		assertHelper(t, AssertRegularPolygon, p, RegPol(Pt(1, 2), Sz(9, 4), 5, 0.5), false, "size")
		assertHelper(t, AssertRegularPolygon, p, RegPol(Pt(1, 2), Sz(3, 4), 9, 0.5), false, "n")
		assertHelper(t, AssertRegularPolygon, p, RegPol(Pt(1, 2), Sz(3, 4), 5, 9.5), false, "angle")
	})
}

func TestAssertPadding(t *testing.T) {
	p := Pad(1, 2, 3, 4)

	t.Run("equal", func(t *testing.T) {
		assertHelper(t, AssertPadding, p, Pad(1, 2, 3, 4), true)
	})
	t.Run("one side differs", func(t *testing.T) {
		assertHelper(t, AssertPadding, p, Pad(9, 2, 3, 4), false, "top")
		assertHelper(t, AssertPadding, p, Pad(1, 9, 3, 4), false, "right")
		assertHelper(t, AssertPadding, p, Pad(1, 2, 9, 4), false, "bottom")
		assertHelper(t, AssertPadding, p, Pad(1, 2, 3, 9), false, "left")
	})
}

func TestAssertMatrix(t *testing.T) {
	m := Mat(1, 2, 3, 4, 5, 6)

	t.Run("equal", func(t *testing.T) {
		assertHelper(t, AssertMatrix, m, Mat(1, 2, 3, 4, 5, 6), true)
	})
	t.Run("one component differs", func(t *testing.T) {
		assertHelper(t, AssertMatrix, m, Mat(9, 2, 3, 4, 5, 6), false, "A")
		assertHelper(t, AssertMatrix, m, Mat(1, 9, 3, 4, 5, 6), false, "B")
		assertHelper(t, AssertMatrix, m, Mat(1, 2, 9, 4, 5, 6), false, "C")
		assertHelper(t, AssertMatrix, m, Mat(1, 2, 3, 9, 5, 6), false, "D")
		assertHelper(t, AssertMatrix, m, Mat(1, 2, 3, 4, 9, 6), false, "E")
		assertHelper(t, AssertMatrix, m, Mat(1, 2, 3, 4, 5, 9), false, "F")
	})
}

// assertHelper runs one of the Assert helpers against a fresh recorder and checks it
// reported the expected outcome. Every helper takes (actual, expected) of the same
// type, so one signature covers them all.
func assertHelper[V any](t *testing.T, assertion func(Testing, V, V, ...string) bool, actual, expected V, result bool, labels ...string) {
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

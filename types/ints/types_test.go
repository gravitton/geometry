package ints

import (
	"testing"

	"github.com/gravitton/assert"
	geom "github.com/gravitton/geometry"
)

func TestConstructors(t *testing.T) {
	// float inputs are rounded to the nearest integer
	assert.Equal(t, Pt(1.6, -2.4), geom.Pt(2, -2))
	assert.Equal(t, Vec(1.6, -2.4), geom.Vec(2, -2))
	assert.Equal(t, Sz(1.6, 2.4), geom.Sz(2, 2))

	assert.Equal(t, Circ(geom.Pt(1.6, 2.4), 3.5), geom.Circ(geom.Pt(2, 2), 4))
	assert.Equal(t, Ell(geom.Pt(1.6, 2.4), geom.Sz(3.4, 4.6), 0.5), geom.Ell(geom.Pt(2, 2), geom.Sz(3, 5), 0.5))
	assert.Equal(t, Ac(geom.Pt(1.6, 2.4), 2.6, 0.5, 1.5), geom.Ac(geom.Pt(2, 2), 3, 0.5, 1.5))
	assert.Equal(t, Sect(geom.Pt(1.6, 2.4), 2.6, 0.5, 1.5), geom.Sect(geom.Pt(2, 2), 3, 0.5, 1.5))
	assert.Equal(t, Cap(geom.Pt(1.6, 2.4), geom.Pt(3.4, 4.6), 0.6), geom.Cap(geom.Pt(2, 2), geom.Pt(3, 5), 1))
	assert.Equal(t, Seg(geom.Pt(1.6, 2.4), geom.Pt(3.4, 4.6)), geom.Seg(geom.Pt(2, 2), geom.Pt(3, 5)))
	assert.Equal(t, QuadBez(geom.Pt(1.6, 2.4), geom.Pt(3.4, 4.6), geom.Pt(5.5, 6.4)), geom.QuadBez(geom.Pt(2, 2), geom.Pt(3, 5), geom.Pt(6, 6)))
	assert.Equal(t, CubBez(geom.Pt(1.6, 2.4), geom.Pt(3.4, 4.6), geom.Pt(5.5, 6.4), geom.Pt(7.4, 8.6)), geom.CubBez(geom.Pt(2, 2), geom.Pt(3, 5), geom.Pt(6, 6), geom.Pt(7, 9)))
	assert.Equal(t, Ry(geom.Pt(1.6, 2.4), geom.Vec(3.4, 4.6)), geom.Ry(geom.Pt(2, 2), geom.Vec(3, 5)))
	assert.Equal(t, Ln(geom.Pt(1.6, 2.4), geom.Vec(3.4, 4.6)), geom.Ln(geom.Pt(2, 2), geom.Vec(3, 5)))
	assert.Equal(t, Rect(geom.Pt(1.6, 2.4), geom.Sz(3.4, 4.6)), geom.Rect(geom.Pt(2, 2), geom.Sz(3, 5)))
	assert.Equal(t, Bx(geom.Pt(3.4, 4.6), geom.Pt(1.6, 2.4)), geom.Bx(geom.Pt(2, 2), geom.Pt(3, 5)))

	assert.True(t, Pol([]geom.Point[float64]{geom.Pt(1.6, 2.4), geom.Pt(3.4, 4.6)}).Equal(geom.Pol([]geom.Point[int]{geom.Pt(2, 2), geom.Pt(3, 5)})))

	assert.True(t, PolLn([]geom.Point[float64]{geom.Pt(1.6, 2.4), geom.Pt(3.4, 4.6)}).Equal(geom.PolLn([]geom.Point[int]{geom.Pt(2, 2), geom.Pt(3, 5)})))

	assert.Equal(t, RegPol(geom.Pt(1.6, 2.4), geom.Sz(3.4, 4.6), 6, 0.5, 0), geom.RegPol(geom.Pt(2, 2), geom.Sz(3, 5), 6, 0.5, 0))

	assert.Equal(t, Pad(1.6, 2.4, 3.5, 4.4), geom.Pad(2, 2, 4, 4))
	assert.Equal(t, Mat(1.6, 2.4, 3.5, 4.4, 5.5, 6.4), geom.Mat(2, 2, 4, 4, 6, 6))
	assert.Equal(t, IdentityMatrix(), geom.IdentityMatrix[int]())
}

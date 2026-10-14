package floats

import (
	"testing"

	"github.com/gravitton/assert"
	geom "github.com/gravitton/geometry"
)

func TestConstructors(t *testing.T) {
	// integer inputs are widened to float64
	assert.Equal(t, Pt(1, -2), geom.Pt(1.0, -2.0))
	assert.Equal(t, Vec(1, -2), geom.Vec(1.0, -2.0))
	assert.Equal(t, Sz(1, 2), geom.Sz(1.0, 2.0))
	assert.Equal(t, Pad(4, 2, 1, 3), geom.Pad(4.0, 2.0, 1.0, 3.0))
	assert.Equal(t, Mat(1, 2, 3, 4, 5, 6), geom.Mat(1.0, 2.0, 3.0, 4.0, 5.0, 6.0))
	assert.Equal(t, IdentityMatrix(), geom.IdentityMatrix[float64]())

	assert.Equal(t, Circ(geom.Pt(1, 2), 3), geom.Circ(geom.Pt(1.0, 2.0), 3.0))
	assert.Equal(t, Cap(geom.Pt(1, 2), geom.Pt(3, 4), 5), geom.Cap(geom.Pt(1.0, 2.0), geom.Pt(3.0, 4.0), 5.0))
	assert.Equal(t, Sect(geom.Pt(1, 2), 3, 0.5, 1.5), geom.Sect(geom.Pt(1.0, 2.0), 3.0, 0.5, 1.5))
	assert.Equal(t, Ac(geom.Pt(1, 2), 3, 0.5, 1.5), geom.Ac(geom.Pt(1.0, 2.0), 3.0, 0.5, 1.5))
	assert.Equal(t, Seg(geom.Pt(1, 2), geom.Pt(3, 4)), geom.Seg(geom.Pt(1.0, 2.0), geom.Pt(3.0, 4.0)))
	assert.Equal(t, Ry(geom.Pt(1, 2), geom.Vec(3, 4)), geom.Ry(geom.Pt(1.0, 2.0), geom.Vec(3.0, 4.0)))
	assert.Equal(t, Ln(geom.Pt(1, 2), geom.Vec(3, 4)), geom.Ln(geom.Pt(1.0, 2.0), geom.Vec(3.0, 4.0)))
	assert.True(t, PolLn([]geom.Point[int]{geom.Pt(1, 2), geom.Pt(3, 4)}).Equal(geom.PolLn([]geom.Point[float64]{geom.Pt(1.0, 2.0), geom.Pt(3.0, 4.0)})))
	assert.Equal(t, QuadBez(geom.Pt(1, 2), geom.Pt(3, 4), geom.Pt(5, 6)), geom.QuadBez(geom.Pt(1.0, 2.0), geom.Pt(3.0, 4.0), geom.Pt(5.0, 6.0)))
	assert.Equal(t, CubBez(geom.Pt(1, 2), geom.Pt(3, 4), geom.Pt(5, 6), geom.Pt(7, 8)), geom.CubBez(geom.Pt(1.0, 2.0), geom.Pt(3.0, 4.0), geom.Pt(5.0, 6.0), geom.Pt(7.0, 8.0)))
	assert.True(t, Pol([]geom.Point[int]{geom.Pt(1, 2), geom.Pt(3, 4)}).Equal(geom.Pol([]geom.Point[float64]{geom.Pt(1.0, 2.0), geom.Pt(3.0, 4.0)})))
	assert.Equal(t, Rect(geom.Pt(1, 2), geom.Sz(3, 4), 0), geom.Rect(geom.Pt(1.0, 2.0), geom.Sz(3.0, 4.0), 0))
	assert.Equal(t, RegPol(geom.Pt(1, 2), geom.Sz(3, 4), 6, 0.5, 0), geom.RegPol(geom.Pt(1.0, 2.0), geom.Sz(3.0, 4.0), 6, 0.5, 0))
	assert.Equal(t, Bx(geom.Pt(3, 4), geom.Pt(1, 2)), geom.Bx(geom.Pt(1.0, 2.0), geom.Pt(3.0, 4.0)))
	assert.Equal(t, Ell(geom.Pt(1, 2), geom.Sz(3, 4), 0.5), geom.Ell(geom.Pt(1.0, 2.0), geom.Sz(3.0, 4.0), 0.5))

	// fractional values survive the round trip
	assert.Equal(t, Pt(1.5, -2.5), geom.Pt(1.5, -2.5))
}

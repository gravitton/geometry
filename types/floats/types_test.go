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

	assert.Equal(t, Circ(geom.Pt(1, 2), 3), geom.Circ(geom.Pt(1.0, 2.0), 3.0))
	assert.Equal(t, Ln(geom.Pt(1, 2), geom.Pt(3, 4)), geom.Ln(geom.Pt(1.0, 2.0), geom.Pt(3.0, 4.0)))
	assert.Equal(t, Rect(geom.Pt(1, 2), geom.Sz(3, 4)), geom.Rect(geom.Pt(1.0, 2.0), geom.Sz(3.0, 4.0)))

	assert.True(t, Pol([]geom.Point[int]{geom.Pt(1, 2), geom.Pt(3, 4)}).Equal(geom.Pol([]geom.Point[float64]{geom.Pt(1.0, 2.0), geom.Pt(3.0, 4.0)})))

	assert.Equal(t, RegPol(geom.Pt(1, 2), geom.Sz(3, 4), 6, 0.5), geom.RegPol(geom.Pt(1.0, 2.0), geom.Sz(3.0, 4.0), 6, 0.5))

	assert.Equal(t, Pad(1, 2, 3, 4), geom.Pad(1.0, 2.0, 3.0, 4.0))
	assert.Equal(t, Mat(1, 2, 3, 4, 5, 6), geom.Mat(1.0, 2.0, 3.0, 4.0, 5.0, 6.0))
	assert.Equal(t, IdentityMatrix(), geom.IdentityMatrix[float64]())

	// fractional values survive the round trip
	assert.Equal(t, Pt(1.5, -2.5), geom.Pt(1.5, -2.5))
}

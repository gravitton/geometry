package geom_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/gravitton/assert"
	. "github.com/gravitton/geometry"
	"github.com/gravitton/geometry/geomtest"
)

func TestHit_Equal(t *testing.T) {
	hit := Hit[float64]{Point: Pt(1.0, 2.0), Normal: Vec(0.0, -1.0), Fraction: 0.5}

	t.Run("the same hit, within the tolerance", func(t *testing.T) {
		assert.True(t, hit.Equal(hit))
		assert.True(t, hit.Equal(Hit[float64]{Point: Pt(1.0+Delta/2, 2.0), Normal: Vec(0.0, -1.0), Fraction: 0.5 + Delta/2}))
	})
	t.Run("another point, normal or fraction", func(t *testing.T) {
		assert.False(t, hit.Equal(Hit[float64]{Point: Pt(1.5, 2.0), Normal: Vec(0.0, -1.0), Fraction: 0.5}))
		assert.False(t, hit.Equal(Hit[float64]{Point: Pt(1.0, 2.0), Normal: Vec(0.0, 1.0), Fraction: 0.5}))
		assert.False(t, hit.Equal(Hit[float64]{Point: Pt(1.0, 2.0), Normal: Vec(0.0, -1.0), Fraction: 0.75}))
	})
	t.Run("a cast answers with a hit equal to itself", func(t *testing.T) {
		hit, ok := Seg(Pt(-5, 0), Pt(5, 0)).HitCircle(Circ(Pt(0, 0), 2))

		assert.True(t, ok)
		assert.True(t, hit.Equal(Hit[int]{Point: Pt(-2, 0), Normal: Vec(-1.0, 0.0), Fraction: 0.3}))
	})
}

func TestHit_IsZero(t *testing.T) {
	t.Run("a cast that reaches nothing", func(t *testing.T) {
		hit, ok := Seg(Pt(-5, 9), Pt(5, 9)).HitCircle(Circ(Pt(0, 0), 2))

		assert.False(t, ok)
		assert.True(t, hit.IsZero())
	})
	t.Run("a hit", func(t *testing.T) {
		assert.False(t, Hit[int]{Point: Pt(-2, 0), Normal: Vec(-1.0, 0.0), Fraction: 0.3}.IsZero())
		assert.False(t, Hit[int]{Fraction: 0.3}.IsZero())
	})
}

func TestHit_Cast(t *testing.T) {
	hit := Hit[float64]{Point: Pt(1.5, -2.4), Normal: Vec(0.6, -0.8), Fraction: 0.25}

	t.Run("rounds the point and keeps the normal and the fraction", func(t *testing.T) {
		geomtest.AssertHit(t, Hit[float64].Cast[int32](hit), Hit[int32]{Point: Pt[int32](2, -2), Normal: Vec(0.6, -0.8), Fraction: 0.25})
		geomtest.AssertHit(t, Hit[float64].Cast[float32](hit), Hit[float32]{Point: Pt[float32](1.5, -2.4), Normal: Vec(0.6, -0.8), Fraction: 0.25})
	})
}

func TestHit_Int(t *testing.T) {
	geomtest.AssertHit(t, Hit[float64]{Point: Pt(1.5, -2.4), Normal: Vec(0.6, -0.8), Fraction: 0.25}.Int(), Hit[int]{Point: Pt(2, -2), Normal: Vec(0.6, -0.8), Fraction: 0.25})
}

func TestHit_Float(t *testing.T) {
	geomtest.AssertHit(t, Hit[int]{Point: Pt(2, -2), Normal: Vec(0.6, -0.8), Fraction: 0.25}.Float(), Hit[float64]{Point: Pt(2.0, -2.0), Normal: Vec(0.6, -0.8), Fraction: 0.25})
}

func TestHit_String(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		assert.Equal(t, Hit[int]{Point: Pt(-2, 0), Normal: Vec(-1.0, 0.0), Fraction: 0.3}.String(), "Hit((-2,0);⟨-1.00,0.00⟩;0.30)")
	})
	t.Run("float", func(t *testing.T) {
		assert.Equal(t, Hit[float64]{Point: Pt(1.5, 2), Normal: Vec(0.6, -0.8), Fraction: 2}.String(), "Hit((1.50,2.00);⟨0.60,-0.80⟩;2.00)")
	})
	t.Run("allocates the string alone", func(t *testing.T) {
		hit := Hit[int]{Point: Pt(-2, 0), Normal: Vec(-1.0, 0.0), Fraction: 0.3}

		assertStringAllocation(t, hit.String)
	})
}

func TestHit_JSON(t *testing.T) {
	hit := Hit[int]{Point: Pt(-2, 0), Normal: Vec(-1.0, 0.0), Fraction: 3}

	data, err := json.Marshal(hit)
	assert.NoError(t, err)
	assert.Equal(t, string(data), `{"p":{"x":-2,"y":0},"n":{"x":-1,"y":0},"f":3}`)

	var decoded Hit[int]
	assert.NoError(t, json.Unmarshal(data, &decoded))
	assert.Equal(t, decoded, hit)
}

func ExampleHit() {
	hit, _ := Seg(Pt(-5, 0), Pt(5, 0)).HitCircle(Circ(Pt(0, 0), 2))

	fmt.Println(hit)
	// Output: Hit((-2,0);⟨-1.00,0.00⟩;0.30)
}

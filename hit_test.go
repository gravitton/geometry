package geom_test

import (
	"encoding/json"
	"testing"

	"github.com/gravitton/assert"
	. "github.com/gravitton/geometry"
)

func TestHit_JSON(t *testing.T) {
	hit := Hit[int]{Point: Pt(-2, 0), Normal: Vec(-1.0, 0.0), Fraction: 3}

	data, err := json.Marshal(hit)
	assert.NoError(t, err)
	assert.Equal(t, string(data), `{"p":{"x":-2,"y":0},"n":{"x":-1,"y":0},"f":3}`)

	var decoded Hit[int]
	assert.NoError(t, json.Unmarshal(data, &decoded))
	assert.Equal(t, decoded, hit)
}

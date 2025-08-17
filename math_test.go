package geom

import (
	"github.com/gravitton/assert"
	"testing"
)

func TestMidpoint(t *testing.T) {
	assert.Equal(t, Midpoint(1, 4), 2)
	assert.Equal(t, Midpoint(1, 5), 3)
	assert.Equal(t, Midpoint(1, 6), 3)

	assert.Equal(t, Midpoint(1.0, 6.0), 3.5)
}

func TestLerp(t *testing.T) {
	assert.Equal(t, Lerp(1, 4, 0.25), 1)
	assert.Equal(t, Lerp(1, 5, 0.25), 2)
	assert.Equal(t, Lerp(1, 6, 0.25), 2)

	assert.Equal(t, Lerp(1.0, 6.0, 0.25), 2.25)
	assert.Equal(t, Lerp(1.0, 6.0, 0.75), 4.75)

}

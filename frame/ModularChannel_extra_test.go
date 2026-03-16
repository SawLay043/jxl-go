package frame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAbsInt32(t *testing.T) {
	assert.Equal(t, int32(5), absInt32(5))
	assert.Equal(t, int32(5), absInt32(-5))
	assert.Equal(t, int32(0), absInt32(0))
}

func TestNewModularChannelFromChannel(t *testing.T) {
	orig := NewModularChannelWithAllParams(2, 2, 0, 0, false)
	orig.decoded = true
	orig.allocate()
	orig.buffer[0][0] = 1
	orig.buffer[1][1] = 2

	copy := NewModularChannelFromChannel(*orig)
	assert.Equal(t, orig.size, copy.size)
	assert.Equal(t, orig.decoded, copy.decoded)
	assert.NotNil(t, copy.buffer)
	assert.Equal(t, int32(1), copy.buffer[0][0])
	assert.Equal(t, int32(2), copy.buffer[1][1])
}

func TestModularChannel_EdgePredictions(t *testing.T) {
	mc := NewModularChannelWithAllParams(4, 4, 0, 0, false)
	mc.allocate()
	// Fill buffer
	for y := uint32(0); y < 4; y++ {
		for x := uint32(0); x < 4; x++ {
			mc.buffer[y][x] = int32(y*4 + x)
		}
	}

	assert.Equal(t, int32(3), mc.northEastEast(1, 1)) // y=1, x=1. NEE is y=0, x=3. buffer[0][3]=3.
	assert.Equal(t, int32(8), mc.westWest(2, 2))     // y=2, x=2. WW is y=2, x=0. buffer[2][0]=8.
}

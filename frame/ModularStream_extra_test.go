package frame

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInverseHorizontalSqueeze_ExtraWidth(t *testing.T) {
	// orig: 1x2, res: 1x1 -> channel: 1x3
	orig := NewModularChannelWithAllParams(1, 2, 0, 0, false)
	orig.allocate()
	orig.buffer[0][0] = 10
	orig.buffer[0][1] = 50

	res := NewModularChannelWithAllParams(1, 1, 0, 0, false)
	res.allocate()
	res.buffer[0][0] = 2

	outputInfo := NewModularChannelWithAllParams(1, 3, 0, 0, false)
	output, err := inverseHorizontalSqueeze(outputInfo, orig, res)
	require.NoError(t, err)
	require.NotNil(t, output)
	assert.Equal(t, uint32(3), output.size.Width)
	
	// With the fix: left = avg = 10. nextAvg = 50.
	// a=10, b=10, c=50.
	// a <= b <= c -> true.
	// x = (4*10 - 3*50 - 10 - 6) / 12 = -126 / 12 = -10.
	// d = 2*(10-10)=0. e=2*(10-50)=-80.
	// x+(x&1) = -10 + 0 = -10. -10 < 0 -> true. x = d-1 = -1.
	// x-(x&1) = -10 - 0 = -10. -10 < -80 -> false.
	// returns -1.
	// diff = 2 + (-1) = 1.
	// first = 10 + 1/2 = 10.
	// buffer[0][0] = 10, buffer[0][1] = 10 - 1 = 9
	assert.Equal(t, int32(10), output.buffer[0][0])
	assert.Equal(t, int32(9), output.buffer[0][1])
	assert.Equal(t, int32(50), output.buffer[0][2])
}

func TestTendancy(t *testing.T) {
	testCases := []struct {
		left, avg, nextAvg int32
		expected           int32
	}{
		{15, 10, 5, 3},
		{5, 10, 15, -3},
		{10, 5, 10, 0},
	}
	for _, tc := range testCases {
		assert.Equal(t, tc.expected, tendancy(tc.left, tc.avg, tc.nextAvg), "tendancy(%d, %d, %d)", tc.left, tc.avg, tc.nextAvg)
	}
}

func TestInverseHorizontalSqueeze_Corrupted(t *testing.T) {
	orig := NewModularChannelWithAllParams(2, 1, 0, 0, false)
	res := NewModularChannelWithAllParams(2, 1, 0, 0, false)
	// Output width should be 2, but we pass 3
	outputInfo := NewModularChannelWithAllParams(2, 3, 0, 0, false)
	_, err := inverseHorizontalSqueeze(outputInfo, orig, res)
	assert.Error(t, err)
}

func TestInverseVerticalSqueeze_Corrupted(t *testing.T) {
	orig := NewModularChannelWithAllParams(1, 2, 0, 0, false)
	res := NewModularChannelWithAllParams(1, 2, 0, 0, false)
	// Output height should be 2, but we pass 3
	outputInfo := NewModularChannelWithAllParams(3, 2, 0, 0, false)
	_, err := inverseVerticalSqueeze(outputInfo, orig, res)
	assert.Error(t, err)
}

func TestInverseHorizontalSqueeze(t *testing.T) {
	// orig: 2x1, res: 2x1 -> channel: 2x2
	orig := NewModularChannelWithAllParams(2, 1, 0, 0, false)
	orig.allocate()
	orig.buffer[0][0] = 10
	orig.buffer[1][0] = 20

	res := NewModularChannelWithAllParams(2, 1, 0, 0, false)
	res.allocate()
	res.buffer[0][0] = 2
	res.buffer[1][0] = 4

	outputInfo := NewModularChannelWithAllParams(2, 2, 0, 0, false)
	output, err := inverseHorizontalSqueeze(outputInfo, orig, res)
	require.NoError(t, err)
	require.NotNil(t, output)
	assert.Equal(t, uint32(2), output.size.Width)
	assert.Equal(t, uint32(2), output.size.Height)
	
	assert.Equal(t, int32(11), output.buffer[0][0])
	assert.Equal(t, int32(9), output.buffer[0][1])
}

func TestInverseVerticalSqueeze(t *testing.T) {
	// orig: 1x2, res: 1x2 -> channel: 2x2
	orig := NewModularChannelWithAllParams(1, 2, 0, 0, false)
	orig.allocate()
	orig.buffer[0][0] = 10
	orig.buffer[0][1] = 20

	res := NewModularChannelWithAllParams(1, 2, 0, 0, false)
	res.allocate()
	res.buffer[0][0] = 2
	res.buffer[0][1] = 4

	outputInfo := NewModularChannelWithAllParams(2, 2, 0, 0, false)
	output, err := inverseVerticalSqueeze(outputInfo, orig, res)
	require.NoError(t, err)
	require.NotNil(t, output)
	assert.Equal(t, uint32(2), output.size.Width)
	assert.Equal(t, uint32(2), output.size.Height)

	assert.Equal(t, int32(11), output.buffer[0][0])
	assert.Equal(t, int32(9), output.buffer[1][0])
}

func TestApplyTransforms_Squeeze(t *testing.T) {
	ms := &ModularStream{
		transforms: []TransformInfo{
			{tr: SQUEEZE},
		},
		squeezeMap: make(map[int][]SqueezeParam),
		channels: []*ModularChannel{
			NewModularChannelWithAllParams(2, 1, 0, 0, false), // orig
			NewModularChannelWithAllParams(2, 1, 0, 0, false), // res
		},
	}
	ms.channels[0].allocate()
	ms.channels[0].buffer[0][0] = 10
	ms.channels[0].buffer[1][0] = 20
	ms.channels[1].allocate()
	ms.channels[1].buffer[0][0] = 2
	ms.channels[1].buffer[1][0] = 4

	ms.squeezeMap[0] = []SqueezeParam{
		{horizontal: true, numC: 1, beginC: 0, inPlace: true},
	}

	err := ms.applyTransforms()
	require.NoError(t, err)
	assert.Len(t, ms.channels, 1)
	assert.Equal(t, uint32(2), ms.channels[0].size.Width)
	assert.Equal(t, int32(11), ms.channels[0].buffer[0][0])
}

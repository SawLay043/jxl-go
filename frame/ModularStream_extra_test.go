package frame

import (
	"testing"

	"github.com/kpfaulkner/jxl-go/bundle"
	"github.com/stretchr/testify/assert"
)

func TestApplyTransforms_RCT_AllTypes(t *testing.T) {
	for rctType := 0; rctType <= 6; rctType++ {
		t.Run(string(rune('0'+rctType)), func(t *testing.T) {
			ms := &ModularStream{
				transforms: []TransformInfo{
					{
						tr:      RCT,
						beginC:  0,
						rctType: rctType,
					},
				},
				channels: []*ModularChannel{
					NewModularChannelWithAllParams(1, 1, 0, 0, false),
					NewModularChannelWithAllParams(1, 1, 0, 0, false),
					NewModularChannelWithAllParams(1, 1, 0, 0, false),
				},
			}
			for i := 0; i < 3; i++ {
				ms.channels[i].allocate()
				ms.channels[i].buffer[0][0] = int32(i + 10) // 10, 11, 12
			}

			err := ms.applyTransforms()
			assert.Nil(t, err)
			assert.NotNil(t, ms.channels[0].buffer[0][0])
		})
	}
}


func TestApplyTransforms_Palette(t *testing.T) {
	// Setup Palette transform
	// PALETTE transform in applyTransforms expects:
	// ms.channels[0] to be the palette (meta channel)
	// ms.channels[1...] to be the indexed channels
	
	palette := NewModularChannelWithAllParams(3, 2, -1, -1, false) // 3 components, 2 colors
	palette.allocate()
	palette.buffer[0][0] = 10 // Component 0, Color 0
	palette.buffer[0][1] = 20 // Component 0, Color 1
	palette.buffer[1][0] = 30 // Component 1, Color 0
	palette.buffer[1][1] = 40 // Component 1, Color 1
	palette.buffer[2][0] = 50 // Component 2, Color 0
	palette.buffer[2][1] = 60 // Component 2, Color 1

	indexed := NewModularChannelWithAllParams(1, 1, 0, 0, false)
	indexed.allocate()
	indexed.buffer[0][0] = 1 // Color 1

	ms := &ModularStream{
		frame: NewFakeFramer(MODULAR),
		transforms: []TransformInfo{
			{
				tr:        PALETTE,
				beginC:    0,
				numC:      3,
				nbColours: 2,
				nbDeltas:  0,
			},
		},
		channels: []*ModularChannel{palette, indexed},
	}
	// Inject bitDepth
	ms.frame.(*FakeFramer).imageHeader.BitDepth = &bundle.BitDepthHeader{BitsPerSample: 8}

	err := ms.applyTransforms()
	assert.Nil(t, err)
	
	// Palette transform replaces indexed channel with numC channels
	// and removes the palette channel from the beginning.
	assert.Equal(t, 3, len(ms.channels))
	assert.Equal(t, int32(20), ms.channels[0].buffer[0][0])
	assert.Equal(t, int32(40), ms.channels[1].buffer[0][0])
	assert.Equal(t, int32(60), ms.channels[2].buffer[0][0])
}

func TestApplyTransforms_SqueezeHorizontal(t *testing.T) {
	// Setup Squeeze Horizontal
	orig := NewModularChannelWithAllParams(1, 1, 0, 1, false)
	orig.allocate()
	orig.buffer[0][0] = 10
	
	residu := NewModularChannelWithAllParams(1, 1, 0, 1, false)
	residu.allocate()
	residu.buffer[0][0] = 2

	ms := &ModularStream{
		transforms: []TransformInfo{
			{
				tr: SQUEEZE,
				sp: []SqueezeParam{{horizontal: true, inPlace: true, beginC: 0, numC: 1}},
			},
		},
		channels: []*ModularChannel{orig, residu},
		squeezeMap: map[int][]SqueezeParam{
			0: {{horizontal: true, inPlace: true, beginC: 0, numC: 1}},
		},
	}

	err := ms.applyTransforms()
	assert.Nil(t, err)
	
	// Inverse horizontal squeeze: 
	// diff = residu + tendancy(...) = 2 + 0 = 2
	// first = avg + diff/2 = 10 + 1 = 11
	// second = first - diff = 11 - 2 = 9
	assert.Equal(t, 1, len(ms.channels))
	assert.Equal(t, uint32(2), ms.channels[0].size.Width)
	assert.Equal(t, int32(11), ms.channels[0].buffer[0][0])
	assert.Equal(t, int32(9), ms.channels[0].buffer[0][1])
}

func TestApplyTransforms_SqueezeVertical(t *testing.T) {
	// Setup Squeeze Vertical
	orig := NewModularChannelWithAllParams(1, 1, 1, 0, false)
	orig.allocate()
	orig.buffer[0][0] = 20
	
	residu := NewModularChannelWithAllParams(1, 1, 1, 0, false)
	residu.allocate()
	residu.buffer[0][0] = 4

	ms := &ModularStream{
		transforms: []TransformInfo{
			{
				tr: SQUEEZE,
				sp: []SqueezeParam{{horizontal: false, inPlace: true, beginC: 0, numC: 1}},
			},
		},
		channels: []*ModularChannel{orig, residu},
		squeezeMap: map[int][]SqueezeParam{
			0: {{horizontal: false, inPlace: true, beginC: 0, numC: 1}},
		},
	}

	err := ms.applyTransforms()
	assert.Nil(t, err)
	
	// Inverse vertical squeeze: 
	// diff = residu + tendancy(...) = 4 + 0 = 4
	// first = avg + diff/2 = 20 + 2 = 22
	// second = first - diff = 22 - 4 = 18
	assert.Equal(t, 1, len(ms.channels))
	assert.Equal(t, uint32(2), ms.channels[0].size.Height)
	assert.Equal(t, int32(22), ms.channels[0].buffer[0][0])
	assert.Equal(t, int32(18), ms.channels[0].buffer[1][0])
}

func TestInverseHorizontalSqueeze_Tendancy(t *testing.T) {
	orig := NewModularChannelWithAllParams(1, 2, 0, 1, false)
	orig.allocate()
	orig.buffer[0][0] = 10
	orig.buffer[0][1] = 20
	
	residu := NewModularChannelWithAllParams(1, 2, 0, 1, false)
	residu.allocate()
	residu.buffer[0][0] = 2
	residu.buffer[0][1] = 4

	outputInfo := NewModularChannelWithAllParams(1, 4, 0, 0, false)
	
	out, err := inverseHorizontalSqueeze(outputInfo, orig, residu)
	assert.Nil(t, err)
	assert.NotNil(t, out)
	assert.Equal(t, uint32(4), out.size.Width)
}

func TestInverseVerticalSqueeze_Tendancy(t *testing.T) {
	orig := NewModularChannelWithAllParams(2, 1, 1, 0, false)
	orig.allocate()
	orig.buffer[0][0] = 10
	orig.buffer[1][0] = 20
	
	residu := NewModularChannelWithAllParams(2, 1, 1, 0, false)
	residu.allocate()
	residu.buffer[0][0] = 2
	residu.buffer[1][0] = 4

	outputInfo := NewModularChannelWithAllParams(4, 1, 0, 0, false)
	
	out, err := inverseVerticalSqueeze(outputInfo, orig, residu)
	assert.Nil(t, err)
	assert.NotNil(t, out)
	assert.Equal(t, uint32(4), out.size.Height)
}

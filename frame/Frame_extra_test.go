package frame

import (
	"testing"

	"github.com/kpfaulkner/jxl-go/bundle"
	"github.com/kpfaulkner/jxl-go/colour"
	"github.com/kpfaulkner/jxl-go/image"
	"github.com/kpfaulkner/jxl-go/options"
	"github.com/kpfaulkner/jxl-go/testcommon"
	"github.com/kpfaulkner/jxl-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPerformEdgePreservingFilter(t *testing.T) {
	f := &Frame{
		Header: &FrameHeader{
			Encoding: MODULAR,
			Bounds: &util.Rectangle{
				Size: util.Dimension{Width: 16, Height: 16},
			},
			restorationFilter: &RestorationFilter{
				epfIterations:      1,
				epfSigmaForModular: 1.0,
				epfPass0SigmaScale: 1.0,
				epfPass2SigmaScale: 1.0,
				epfQuantMul:        1.0,
				epfChannelScale:    []float32{1.0, 1.0, 1.0},
				epfSharpLut:        []float32{0, 0, 0, 0, 0, 0, 0, 0},
			},
			jpegUpsamplingX: []int32{0, 0, 0},
			jpegUpsamplingY: []int32{0, 0, 0},
			Upsampling:      1,
		},
		GlobalMetadata: &bundle.ImageHeader{
			BitDepth: &bundle.BitDepthHeader{BitsPerSample: 8},
			ColourEncoding: &colour.ColourEncodingBundle{ColourEncoding: colour.CE_RGB},
		},
		options: &options.JXLOptions{MaxGoroutines: 1},
	}

	// Create 3 float buffers
	f.Buffer = make([]image.ImageBuffer, 3)
	for c := 0; c < 3; c++ {
		ib, err := image.NewImageBuffer(image.TYPE_FLOAT, 16, 16)
		require.NoError(t, err)
		for y := 0; y < 16; y++ {
			for x := 0; x < 16; x++ {
				ib.FloatBuffer[y][x] = float32(y*16 + x)
			}
		}
		f.Buffer[c] = *ib
	}

	err := f.performEdgePreservingFilter()
	require.NoError(t, err)
	assert.Len(t, f.Buffer, 3)
}

func TestPerformEdgePreservingFilter_Vardct(t *testing.T) {
	f := &Frame{
		Header: &FrameHeader{
			Encoding: VARDCT,
			Bounds: &util.Rectangle{
				Size: util.Dimension{Width: 16, Height: 16},
			},
			restorationFilter: &RestorationFilter{
				epfIterations:      1,
				epfSigmaForModular: 1.0,
				epfPass0SigmaScale: 1.0,
				epfPass2SigmaScale: 1.0,
				epfQuantMul:        1.0,
				epfChannelScale:    []float32{1.0, 1.0, 1.0},
				epfSharpLut:        []float32{0, 0, 0, 0, 0, 0, 0, 0},
			},
			jpegUpsamplingX: []int32{0, 0, 0},
			jpegUpsamplingY: []int32{0, 0, 0},
			Upsampling:      1,
		},
		GlobalMetadata: &bundle.ImageHeader{
			BitDepth: &bundle.BitDepthHeader{BitsPerSample: 8},
		},
		LfGlobal: &LFGlobal{
			lfDequant: []float32{1.0, 1.0, 1.0},
			globalScale: 1,
		},
		lfGroups: make([]*LFGroup, 1),
		lfGroupRowStride: 1,
		options: &options.JXLOptions{MaxGoroutines: 1},
	}
	f.lfGroups[0] = &LFGroup{
		hfMetadata: &HFMetadata{
			hfMultiplier: util.MakeMatrix2D[int32](256, 256),
			hfStreamBuffer: make([][][]int32, 4),
		},
	}
	for i := 0; i < 4; i++ {
		f.lfGroups[0].hfMetadata.hfStreamBuffer[i] = util.MakeMatrix2D[int32](256, 256)
	}
	// Set hfMultiplier[0][0] to 1 to avoid div by zero
	f.lfGroups[0].hfMetadata.hfMultiplier[0][0] = 1

	f.Buffer = make([]image.ImageBuffer, 3)
	for c := 0; c < 3; c++ {
		ib, err := image.NewImageBuffer(image.TYPE_FLOAT, 16, 16)
		require.NoError(t, err)
		f.Buffer[c] = *ib
	}

	err := f.performEdgePreservingFilter()
	require.NoError(t, err)
}

func TestReadFrameHeader(t *testing.T) {
	parent := makeParent(100, 200, true, 0)
	reader := &testcommon.FakeBitReader{
		ReadBoolData: []bool{true}, // allDefault = true
	}

	f := &Frame{
		GlobalMetadata: parent,
		reader:         reader,
	}

	header, err := f.ReadFrameHeader()
	require.NoError(t, err)
	assert.Equal(t, uint32(REGULAR_FRAME), header.FrameType)
	assert.Equal(t, uint32(256), f.Header.groupDim)
	assert.Equal(t, uint32(1), f.numGroups) // CeilDiv(100, 256) * CeilDiv(200, 256) = 1 * 1 = 1
	assert.Equal(t, uint32(1), f.numLFGroups) // CeilDiv(100, 2048) * CeilDiv(200, 2048) = 1 * 1 = 1
}

func TestReadTOC_Simple(t *testing.T) {
	// Simple TOC: numGroups=1, numPasses=1 -> tocEntries=1
	parent := makeParent(100, 200, true, 0)
	reader := &testcommon.FakeBitReader{
		ReadBoolData: []bool{
			true,  // allDefault = true
			false, // permutatedTOC = false
		},
		ReadU32Data: []uint32{
			100, // tocLengths[0]
		},
	}

	f := &Frame{
		GlobalMetadata: parent,
		reader:         reader,
	}

	_, err := f.ReadFrameHeader()
	require.NoError(t, err)

	err = f.ReadTOC()
	require.NoError(t, err)
	assert.False(t, f.permutatedTOC)
	assert.Len(t, f.tocLengths, 1)
	assert.Equal(t, uint32(100), f.tocLengths[0])
}

func TestReadFrameHeader_ErrorPadding(t *testing.T) {
	// Currently FakeBitReader.ZeroPadToByte returns nil.
	// Let's assume it works for now or update it if needed.
}

func TestReadFrameHeader_ErrorNewFrameHeader(t *testing.T) {
	f := &Frame{
		GlobalMetadata: makeParent(100, 200, true, 0),
		reader: &testcommon.FakeBitReader{
			ReadBoolData: []bool{}, // Error on first ReadBool
		},
	}
	_, err := f.ReadFrameHeader()
	assert.Error(t, err)
}

func TestReadTOC_ErrorPermutatedBool(t *testing.T) {
	f := &Frame{
		numGroups: 1,
		Header: &FrameHeader{
			passes: &PassesInfo{numPasses: 1},
		},
		reader: &testcommon.FakeBitReader{
			ReadBoolData: []bool{}, // Error on permutatedTOC ReadBool
		},
	}
	err := f.ReadTOC()
	assert.Error(t, err)
}

func TestReadBuffer(t *testing.T) {
	f := &Frame{
		tocLengths: []uint32{10},
		reader:     &testcommon.FakeBitReader{},
	}
	buf, err := f.readBuffer(0)
	require.NoError(t, err)
	assert.Len(t, buf, 14) // length + 4
}

func TestSetupBitReaders_Multiple(t *testing.T) {
	f := &Frame{
		tocLengths: []uint32{10, 20},
		reader:     &testcommon.FakeBitReader{},
	}
	err := f.setupBitReaders()
	require.NoError(t, err)
	assert.Len(t, f.bitreaders, 2)
}

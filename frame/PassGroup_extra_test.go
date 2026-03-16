package frame

import (
	"testing"

	"github.com/kpfaulkner/jxl-go/testcommon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPassGroupWithReader_Modular(t *testing.T) {
	parent := makeParent(100, 200, true, 0)
	reader := &testcommon.FakeBitReader{
		ReadBoolData: []bool{true}, // allDefault = true
	}

	f := &Frame{
		GlobalMetadata: parent,
		reader:         reader,
	}

	_, err := f.ReadFrameHeader()
	require.NoError(t, err)
	// Change encoding to MODULAR
	f.Header.Encoding = MODULAR
	f.lfGroups = make([]*LFGroup, f.numLFGroups)
	for i := uint32(0); i < f.numLFGroups; i++ {
		f.lfGroups[i] = &LFGroup{}
	}

	// NewPassGroupWithReader will call:
	// 1. NewModularStreamWithStreamIndex -> returns early if replacedChannels is empty
	// 2. stream.decodeChannels -> returns early if channels is empty
	
	pg, err := NewPassGroupWithReader(reader, f, 0, 0, nil)
	require.NoError(t, err)
	require.NotNil(t, pg)
	assert.Equal(t, uint32(0), pg.groupID)
	assert.Equal(t, uint32(0), pg.passID)
	assert.Nil(t, pg.hfCoefficients)
}

func TestNewPassGroupWithReader_VardctError(t *testing.T) {
	parent := makeParent(100, 200, true, 0)
	reader := &testcommon.FakeBitReader{
		ReadBoolData: []bool{true}, // allDefault = true
	}

	f := &Frame{
		GlobalMetadata: parent,
		reader:         reader,
	}

	_, err := f.ReadFrameHeader()
	require.NoError(t, err)
	// Encoding is VARDCT by default
	f.hfGlobal = &HFGlobal{numHFPresets: 1}
	f.LfGlobal = &LFGlobal{hfBlockCtx: &HFBlockContext{}}
	f.lfGroups = make([]*LFGroup, f.numLFGroups)
	for i := uint32(0); i < f.numLFGroups; i++ {
		f.lfGroups[i] = &LFGroup{}
	}

	// NewPassGroupWithReader will call NewHFCoefficientsWithReader
	// which will fail because reader is empty
	pg, err := NewPassGroupWithReader(&testcommon.FakeBitReader{}, f, 0, 0, nil)
	assert.Error(t, err)
	assert.Nil(t, pg)
}

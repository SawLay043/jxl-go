package frame

import (
	"testing"

	"github.com/kpfaulkner/jxl-go/jxlio"
	"github.com/kpfaulkner/jxl-go/testcommon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHFBlockContextWithReader_ErrorPaths(t *testing.T) {
	t.Run("Error on useDefault", func(t *testing.T) {
		reader := &testcommon.FakeBitReader{
			ReadBoolData: []bool{}, // Error
		}
		hf, err := NewHFBlockContextWithReader(reader, fakeReadClusterMap)
		assert.Error(t, err)
		assert.Nil(t, hf)
	})

	t.Run("Error on nbLFThresh", func(t *testing.T) {
		reader := &testcommon.FakeBitReader{
			ReadBoolData: []bool{false},
			ReadBitsData: []uint64{}, // Error on first nbLFThresh
		}
		hf, err := NewHFBlockContextWithReader(reader, fakeReadClusterMap)
		assert.Error(t, err)
		assert.Nil(t, hf)
	})

	t.Run("Error on lfThresholds value", func(t *testing.T) {
		reader := &testcommon.FakeBitReader{
			ReadBoolData: []bool{false},
			ReadBitsData: []uint64{1, 0, 0}, // nbLFThresh = [1, 0, 0]
			ReadU32Data:  []uint32{},        // Error on ReadU32
		}
		hf, err := NewHFBlockContextWithReader(reader, fakeReadClusterMap)
		assert.Error(t, err)
		assert.Nil(t, hf)
	})

	t.Run("Error on nbQfThread", func(t *testing.T) {
		reader := &testcommon.FakeBitReader{
			ReadBoolData: []bool{false},
			ReadBitsData: []uint64{0, 0, 0}, // nbLFThresh = [0, 0, 0]
			// nbQfThread read next
		}
		hf, err := NewHFBlockContextWithReader(reader, fakeReadClusterMap)
		assert.Error(t, err)
		assert.Nil(t, hf)
	})

	t.Run("Error on qfThresholds value", func(t *testing.T) {
		reader := &testcommon.FakeBitReader{
			ReadBoolData: []bool{false},
			ReadBitsData: []uint64{0, 0, 0, 1}, // nbLFThresh=[0,0,0], nbQfThread=1
			ReadU32Data:  []uint32{},           // Error on qfThresholds ReadU32
		}
		hf, err := NewHFBlockContextWithReader(reader, fakeReadClusterMap)
		assert.Error(t, err)
		assert.Nil(t, hf)
	})

	t.Run("HF block size too large", func(t *testing.T) {
		// bSize = 39 * (nbQfThread + 1) * (nbLFThresh[0] + 1) * (nbLFThresh[1] + 1) * (nbLFThresh[2] + 1)
		// Max allowed 39 * 64
		// If nbQfThread = 15, nbLFThresh = [1, 1, 1]
		// bSize = 39 * 16 * 2 * 2 * 2 = 39 * 128 > 39 * 64
		reader := &testcommon.FakeBitReader{
			ReadBoolData: []bool{false},
			ReadBitsData: []uint64{1, 1, 1, 15}, // nbLFThresh=[1,1,1], nbQfThread=15
			ReadU32Data:  []uint32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, 
		}
		hf, err := NewHFBlockContextWithReader(reader, fakeReadClusterMap)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "HF block size too large")
		assert.Nil(t, hf)
	})

	t.Run("Error on readClusterMap", func(t *testing.T) {
		reader := &testcommon.FakeBitReader{
			ReadBoolData: []bool{false},
			ReadBitsData: []uint64{0, 0, 0, 0}, // nbLFThresh=[0,0,0], nbQfThread=0
		}
		errReadClusterMap := func(reader jxlio.BitReader, clusterMap []int, maxClusters int) (int, error) {
			return 0, assert.AnError
		}
		hf, err := NewHFBlockContextWithReader(reader, errReadClusterMap)
		assert.Error(t, err)
		assert.Nil(t, hf)
	})
}

func TestNewHFBlockContextWithReader_UseDefault(t *testing.T) {
	// THIS TEST SHOULD FAIL IF THE BUG IS PRESENT (unless the test provides data it shouldn't read)
	reader := &testcommon.FakeBitReader{
		ReadBoolData: []bool{true},
	}
	hf, err := NewHFBlockContextWithReader(reader, fakeReadClusterMap)
	require.NoError(t, err)
	require.NotNil(t, hf)

	assert.Equal(t, int32(15), hf.numClusters)
	assert.Equal(t, int32(1), hf.numLFContexts)
	assert.Len(t, hf.clusterMap, 39)
}

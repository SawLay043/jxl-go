package frame

import (
	"testing"

	"github.com/kpfaulkner/jxl-go/testcommon"
	"github.com/stretchr/testify/assert"
)

func TestNewFrameHeaderWithReader_ErrorPaths(t *testing.T) {
	parent := makeParent(100, 200, true, 0)

	t.Run("Error on ReadBits(2) frameType", func(t *testing.T) {
		reader := &testcommon.FakeBitReader{
			ReadBoolData: []bool{false}, // allDefault = false
			ReadBitsData: []uint64{},    // Error on frameType
		}
		fh, err := NewFrameHeaderWithReader(reader, parent)
		assert.Error(t, err)
		assert.Nil(t, fh)
	})

	t.Run("Error on ReadBits(1) encoding", func(t *testing.T) {
		reader := &testcommon.FakeBitReader{
			ReadBoolData: []bool{false},
			ReadBitsData: []uint64{0}, // frameType = 0
			// Error on encoding
		}
		fh, err := NewFrameHeaderWithReader(reader, parent)
		assert.Error(t, err)
		assert.Nil(t, fh)
	})

	t.Run("Error on ReadU64 flags", func(t *testing.T) {
		reader := &testcommon.FakeBitReader{
			ReadBoolData: []bool{false},
			ReadBitsData: []uint64{0, 0}, // frameType=0, encoding=0
			// Error on ReadU64 flags
		}
		fh, err := NewFrameHeaderWithReader(reader, parent)
		assert.Error(t, err)
		assert.Nil(t, fh)
	})

	t.Run("Error on DoYCbCr ReadBool", func(t *testing.T) {
		parentNoXyb := makeParent(100, 200, false, 0)
		reader := &testcommon.FakeBitReader{
			ReadBoolData: []bool{false}, // allDefault = false
			ReadBitsData: []uint64{0, 0}, // frameType=0, encoding=0
			ReadU64Data:  []uint64{0},    // flags = 0
			// Error on DoYCbCr ReadBool
		}
		fh, err := NewFrameHeaderWithReader(reader, parentNoXyb)
		assert.Error(t, err)
		assert.Nil(t, fh)
	})

	t.Run("Error on jpeg mode ReadBits", func(t *testing.T) {
		parentNoXyb := makeParent(100, 200, false, 0)
		reader := &testcommon.FakeBitReader{
			ReadBoolData: []bool{false, true}, // allDefault=false, DoYCbCr=true
			ReadBitsData: []uint64{0, 0, 0},    // frameType=0, encoding=0, flags=0
			ReadU64Data:  []uint64{0},          // flags = 0
			// Error on first jpeg mode ReadBits
		}
		fh, err := NewFrameHeaderWithReader(reader, parentNoXyb)
		assert.Error(t, err)
		assert.Nil(t, fh)
	})

	t.Run("Error on upsampling ReadBits", func(t *testing.T) {
		reader := &testcommon.FakeBitReader{
			ReadBoolData: []bool{false},
			ReadBitsData: []uint64{0, 0}, // frameType=0, encoding=0
			ReadU64Data:  []uint64{0},    // flags=0
			// Error on upsampling ReadBits
		}
		fh, err := NewFrameHeaderWithReader(reader, parent)
		assert.Error(t, err)
		assert.Nil(t, fh)
	})

	t.Run("Error on ecUpsampling ReadBits", func(t *testing.T) {
		parentExtra := makeParent(100, 200, true, 1)
		reader := &testcommon.FakeBitReader{
			ReadBoolData: []bool{false},
			ReadBitsData: []uint64{0, 0, 0}, // frameType=0, encoding=0, upsampling=0
			ReadU64Data:  []uint64{0},       // flags=0
			// Error on ecUpsampling[0] ReadBits
		}
		fh, err := NewFrameHeaderWithReader(reader, parentExtra)
		assert.Error(t, err)
		assert.Nil(t, fh)
	})

	t.Run("Error on groupSizeShift ReadBits", func(t *testing.T) {
		reader := &testcommon.FakeBitReader{
			ReadBoolData: []bool{false},
			ReadBitsData: []uint64{0, 1, 0}, // frameType=0, encoding=MODULAR, upsampling=0
			ReadU64Data:  []uint64{0},       // flags=0
			// Error on groupSizeShift ReadBits
		}
		fh, err := NewFrameHeaderWithReader(reader, parent)
		assert.Error(t, err)
		assert.Nil(t, fh)
	})

	t.Run("Error on xqmScale ReadBits", func(t *testing.T) {
		reader := &testcommon.FakeBitReader{
			ReadBoolData: []bool{false},
			ReadBitsData: []uint64{0, 0, 0}, // frameType=0, encoding=VARDCT, upsampling=0
			ReadU64Data:  []uint64{0},       // flags=0
			// Error on xqmScale ReadBits
		}
		fh, err := NewFrameHeaderWithReader(reader, parent)
		assert.Error(t, err)
		assert.Nil(t, fh)
	})

	t.Run("Error on bqmScale ReadBits", func(t *testing.T) {
		reader := &testcommon.FakeBitReader{
			ReadBoolData: []bool{false},
			ReadBitsData: []uint64{0, 0, 0, 3}, // frameType=0, encoding=VARDCT, upsampling=0, xqmScale=3
			ReadU64Data:  []uint64{0},          // flags=0
			// Error on bqmScale ReadBits
		}
		fh, err := NewFrameHeaderWithReader(reader, parent)
		assert.Error(t, err)
		assert.Nil(t, fh)
	})

	t.Run("Error on LfLevel ReadBits", func(t *testing.T) {
		reader := &testcommon.FakeBitReader{
			ReadBoolData: []bool{false},
			ReadBitsData: []uint64{1, 0, 0, 3, 2}, // frameType=LF_FRAME, encoding=VARDCT, upsampling=0, xqmScale=3, bqmScale=2
			ReadU64Data:  []uint64{0},             // flags=0
			ReadU32Data:  []uint32{1},             // numPasses=1
			// Error on LfLevel ReadBits
		}
		fh, err := NewFrameHeaderWithReader(reader, parent)
		assert.Error(t, err)
		assert.Nil(t, fh)
	})

	t.Run("Error on haveCrop ReadBool", func(t *testing.T) {
		reader := &testcommon.FakeBitReader{
			ReadBoolData: []bool{false}, // allDefault = false
			ReadBitsData: []uint64{0, 0, 0, 3, 2},
			ReadU64Data:  []uint64{0},
			ReadU32Data:  []uint32{1},
			// Error on haveCrop ReadBool
		}
		fh, err := NewFrameHeaderWithReader(reader, parent)
		assert.Error(t, err)
		assert.Nil(t, fh)
	})

	t.Run("Error on x0 ReadU32", func(t *testing.T) {
		reader := &testcommon.FakeBitReader{
			ReadBoolData: []bool{false, true, true}, // allDefault=false, haveCrop=true, IsLast=true
			ReadBitsData: []uint64{0, 0, 0, 3, 2},
			ReadU64Data:  []uint64{0},
			ReadU32Data:  []uint32{1, 0}, // numPasses=1, BlendingInfo mode=0
			// Error on x0 ReadU32
		}
		fh, err := NewFrameHeaderWithReader(reader, parent)
		assert.Error(t, err)
		assert.Nil(t, fh)
	})

	t.Run("Error on width ReadU32", func(t *testing.T) {
		reader := &testcommon.FakeBitReader{
			ReadBoolData: []bool{false, true, true}, // allDefault=false, haveCrop=true, IsLast=true
			ReadBitsData: []uint64{0, 0, 0, 3, 2},
			ReadU64Data:  []uint64{0},
			ReadU32Data:  []uint32{1, 0, 0, 0}, // numPasses=1, BlendingInfo mode=0, x0=0, y0=0
			// Error on width ReadU32
		}
		fh, err := NewFrameHeaderWithReader(reader, parent)
		assert.Error(t, err)
		assert.Nil(t, fh)
	})

	t.Run("Error on IsLast ReadBool", func(t *testing.T) {
		reader := &testcommon.FakeBitReader{
			ReadBoolData: []bool{false, false}, // allDefault=false, haveCrop=false
			ReadBitsData: []uint64{0, 0, 0, 3, 2},
			ReadU64Data:  []uint64{0},
			ReadU32Data:  []uint32{1, 0}, // numPasses=1, BlendingInfo mode=0
			// Error on IsLast ReadBool
		}
		fh, err := NewFrameHeaderWithReader(reader, parent)
		assert.Error(t, err)
		assert.Nil(t, fh)
	})

	t.Run("Error on SaveAsReference ReadBits", func(t *testing.T) {
		reader := &testcommon.FakeBitReader{
			ReadBoolData: []bool{false, false, false}, // allDefault=false, haveCrop=false, IsLast=false
			ReadBitsData: []uint64{0, 0, 0, 3, 2},
			ReadU64Data:  []uint64{0},
			ReadU32Data:  []uint32{1, 0}, // numPasses=1, BlendingInfo mode=0
			// Error on SaveAsReference ReadBits
		}
		fh, err := NewFrameHeaderWithReader(reader, parent)
		assert.Error(t, err)
		assert.Nil(t, fh)
	})

	t.Run("Error on SaveBeforeCT ReadBool", func(t *testing.T) {
		reader := &testcommon.FakeBitReader{
			ReadBoolData: []bool{false, false, false}, // allDefault=false, haveCrop=false, IsLast=false
			ReadBitsData: []uint64{2, 0, 0, 1},        // FrameType=REFERENCE_ONLY, encoding=0, upsampling=0, saveAsReference=1
			ReadU64Data:  []uint64{0},
			// Error on SaveBeforeCT ReadBool
		}
		fh, err := NewFrameHeaderWithReader(reader, parent)
		assert.Error(t, err)
		assert.Nil(t, fh)
	})

	t.Run("Error on nameLen ReadU32", func(t *testing.T) {
		reader := &testcommon.FakeBitReader{
			ReadBoolData: []bool{false, false, true}, // allDefault=false, haveCrop=false, IsLast=true
			ReadBitsData: []uint64{0, 0, 0, 3, 2},
			ReadU64Data:  []uint64{0},
			ReadU32Data:  []uint32{1, 0}, // numPasses=1, BlendingInfo mode=0
			// Error on nameLen ReadU32
		}
		fh, err := NewFrameHeaderWithReader(reader, parent)
		assert.Error(t, err)
		assert.Nil(t, fh)
	})

	t.Run("Error on name byte ReadByte", func(t *testing.T) {
		reader := &testcommon.FakeBitReader{
			ReadBoolData: []bool{false, false, true}, // allDefault=false, haveCrop=false, IsLast=true
			ReadBitsData: []uint64{0, 0, 0, 3, 2},
			ReadU64Data:  []uint64{0},
			ReadU32Data:  []uint32{1, 0, 1}, // numPasses=1, BlendingInfo mode=0, nameLen=1
			// Error on ReadByte
		}
		fh, err := NewFrameHeaderWithReader(reader, parent)
		assert.Error(t, err)
		assert.Nil(t, fh)
	})
}

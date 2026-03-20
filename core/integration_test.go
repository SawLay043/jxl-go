package core

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"math"
	"testing"

	"github.com/kpfaulkner/jxl-go/options"
	"github.com/stretchr/testify/assert"
)

func hashImage(img *JXLImage) string {
	h := md5.New()
	for _, ib := range img.Buffer {
		if ib.IsInt() {
			for y := 0; y < int(ib.Height); y++ {
				for x := 0; x < int(ib.Width); x++ {
					binary.Write(h, binary.LittleEndian, ib.IntBuffer[y][x])
				}
			}
		} else {
			for y := 0; y < int(ib.Height); y++ {
				for x := 0; x < int(ib.Width); x++ {
					bits := math.Float32bits(ib.FloatBuffer[y][x])
					binary.Write(h, binary.LittleEndian, bits)
				}
			}
		}
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func TestIntegrationDecodes(t *testing.T) {
	testCases := []struct {
		filename     string
		expectedHash string
	}{
		{
			filename:     "../testdata/unittest.jxl",
			expectedHash: "c154a3a4419b4883d69ab49d39d00278",
		},
		{
			filename:     "../testdata/tiny2.jxl",
			expectedHash: "dcb08821dea984caac28d994bfcba326",
		},
		{
			filename:     "../testdata/lossless.jxl",
			expectedHash: "c154a3a4419b4883d69ab49d39d00278",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.filename, func(t *testing.T) {
			br := GenerateTestBitReader(t, tc.filename)
			opts := options.NewJXLOptions(nil)
			decoder := NewJXLCodestreamDecoder(br, opts)

			img, err := decoder.decode()
			assert.Nil(t, err)
			assert.NotNil(t, img)

			actualHash := hashImage(img)
			if tc.expectedHash == "TODO" {
				t.Logf("Hash for %s: %s", tc.filename, actualHash)
			} else {
				assert.Equal(t, tc.expectedHash, actualHash, "Hash mismatch for %s", tc.filename)
			}
		})
	}
}

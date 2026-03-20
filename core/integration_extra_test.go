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

func hashImageExtra(img *JXLImage) string {
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

func TestIntegrationDecodesExtra(t *testing.T) {
	testCases := []struct {
		filename     string
		expectedHash string
	}{
		{
			filename:     "../testdata/grayscale.jxl",
			expectedHash: "e99f5d4caf40647622d8285337d8493b",
		},
		/*
		{
			filename:     "../testdata/patches.jxl",
			expectedHash: "EXPECTED_HASH",
		},
		*/
		{
			filename:     "../testdata/alpha-triangles.jxl",
			expectedHash: "ae06b22fb7541f7ebc360383c4b0c7c5",
		},
		/*
		{
			filename:     "../testdata/white.jxl",
			expectedHash: "EXPECTED_HASH",
		},
		*/
		{
			filename:     "../testdata/sunset_logo.jxl",
			expectedHash: "d691865dea655dcadf8499b6f45fa0bd",
		},
		{
			filename:     "../testdata/art.jxl",
			expectedHash: "cfa0d5bf4d80242bc68b4ad3824667d9",
		},
		{
			filename:     "../testdata/quilt.jxl",
			expectedHash: "5830480191d8d4bffae662086161eaa3",
		},
		{
			filename:     "../testdata/lenna.jxl",
			expectedHash: "e2cd199b6f3e838dfe1faee7634ecf01",
		},
		{
			filename:     "../testdata/church.jxl",
			expectedHash: "51e3fffb70989be99282fefeb289cffe",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.filename, func(t *testing.T) {
			br := GenerateTestBitReader(t, tc.filename)
			opts := options.NewJXLOptions(nil)
			decoder := NewJXLCodestreamDecoder(br, opts)

			img, err := decoder.decode()
			if err != nil {
				t.Fatalf("Failed to decode %s: %v", tc.filename, err)
			}
			assert.NotNil(t, img)

			actualHash := hashImageExtra(img)
			assert.Equal(t, tc.expectedHash, actualHash, "Hash for %s", tc.filename)
		})
	}
}

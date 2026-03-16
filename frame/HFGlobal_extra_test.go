package frame

import (
	"testing"
)

func TestHFGlobal_DebugFunctions(t *testing.T) {
	hfg := &HFGlobal{
		weights: make([][][][]float32, 17),
	}
	for i := 0; i < 17; i++ {
		hfg.weights[i] = make([][][]float32, 1)
		hfg.weights[i][0] = make([][]float32, 1)
		hfg.weights[i][0][0] = make([]float32, 1)
		hfg.weights[i][0][0][0] = 1.0
	}

	// Just ensure they run without panic
	hfg.totalWeights()
	hfg.displayWeights()
	hfg.displaySpecificWeights(0, 0, 0)
}

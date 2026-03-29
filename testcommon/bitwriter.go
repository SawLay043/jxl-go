package testcommon

// BitWriter is a helper to write bits into a byte slice.
// Used for crafting mock bitstreams in tests.
type BitWriter struct {
	data []byte
	byte byte
	bits int
}

func NewBitWriter() *BitWriter {
	return &BitWriter{}
}

func (bw *BitWriter) WriteBit(bit uint8) {
	if bit != 0 {
		bw.byte |= (1 << bw.bits)
	}
	bw.bits++
	if bw.bits == 8 {
		bw.data = append(bw.data, bw.byte)
		bw.byte = 0
		bw.bits = 0
	}
}

func (bw *BitWriter) WriteBits(val uint64, numBits int) {
	for i := 0; i < numBits; i++ {
		bw.WriteBit(uint8((val >> i) & 1))
	}
}

func (bw *BitWriter) WriteU8(val int) {
	if val == 0 {
		bw.WriteBit(0)
		return
	}
	bw.WriteBit(1)
	n := 0
	for (1 << (n + 1)) <= val {
		n++
	}
	bw.WriteBits(uint64(n), 3)
	bw.WriteBits(uint64(val-(1<<n)), n)
}

func (bw *BitWriter) Bytes() []byte {
	if bw.bits > 0 {
		return append(bw.data, bw.byte)
	}
	return bw.data
}

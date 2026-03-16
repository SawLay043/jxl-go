
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

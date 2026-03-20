
func TestNewRestorationFilterWithReader_Error(t *testing.T) {
	reader := &testcommon.FakeBitReader{
		ReadBoolData: []bool{}, // EOF on first ReadBool
	}
	_, err := NewRestorationFilterWithReader(reader, VARDCT)
	assert.Error(t, err)
}
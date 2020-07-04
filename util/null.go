package util

// NullReader is an io.Reader that always returns nulls.
type NullReader struct{}

// Read implements io.Reader.
func (i NullReader) Read(p []byte) (n int, err error) {
	for i := range p {
		p[i] = 0
	}
	return len(p), nil
}

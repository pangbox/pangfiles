package util

import "io"

// NullInputPadder pads all short reads to return zeros. This is useful for
// doing streaming encryption using XTEA, for example.
type NullInputPadder struct {
	Reader io.Reader
}

// Read implements io.Reader.
func (i NullInputPadder) Read(p []byte) (n int, err error) {
	n, err = i.Reader.Read(p)
	if n > 0 && n < len(p) {
		for ; n < len(p); n++ {
			p[n] = 0
		}
	}
	return n, err
}

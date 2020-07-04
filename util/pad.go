package util

import (
	"errors"
	"io"
)

var (
	errInvalidReadLength = errors.New("invalid read length")
)

// NullInputPadder pads the final short read with null bytes.
type NullInputPadder struct {
	Reader io.Reader
}

// Read implements io.Reader.
func (i NullInputPadder) Read(p []byte) (n int, err error) {
	// Perform initial read attempt.
	n, err = i.Reader.Read(p)
	if n < 0 || n > len(p) {
		panic(errInvalidReadLength)
	}
	if err != nil && err != io.EOF {
		return n, err
	}

	// For short reads, repeatedly read until EOF or full read.
	for n > 0 && n < len(p) {
		m, err := i.Reader.Read(p[n:])
		if m < 0 || m > len(p[n:]) {
			panic(errInvalidReadLength)
		}
		n += m
		if err == io.EOF || n == len(p) {
			// Break on either EOF, or full read.
			break
		} else if err != nil {
			return n, err
		}
	}

	// Only pad if we got any data - otherwise just return empty.
	if n > 0 {
		// Fill buffer with zeros, if still short after EOF.
		for ; n < len(p); n++ {
			p[n] = 0
		}
	}

	return n, err
}

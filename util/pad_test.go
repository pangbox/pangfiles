package util

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ReplayReader [][]byte

func (r *ReplayReader) Read(p []byte) (n int, err error) {
	// If no more buffers remain: return EOF.
	if len(*r) == 0 {
		return 0, io.EOF
	}

	// Read bytes from current buffer.
	n = copy(p, (*r)[0])

	// Advance current buffer.
	(*r)[0] = (*r)[0][n:]

	// If done: advance to next buffer.
	if len((*r)[0]) == 0 {
		(*r) = (*r)[1:]
	}

	return n, err
}

func readall(t *testing.T, r io.Reader, rlen int) []byte {
	data := []byte{}
	rbuf := make([]byte, rlen)
	n := 0
	for {
		for i := 0; i < rlen; i++ {
			rbuf[i] = 0
		}
		m, err := r.Read(rbuf[:rlen])
		if m < 0 || m > rlen {
			t.Fatalf("Invalid read for size %d: %d", rlen, m)
		}
		data = append(data, rbuf[:m]...)
		n += m
		if err == io.EOF {
			return data
		} else if err != nil {
			t.Fatal(err)
		}
	}
}

func TestPadNullInitialEOF(t *testing.T) {
	eof := &ReplayReader{}
	pad := &NullInputPadder{Reader: eof}
	assert.Equal(t, []byte{}, readall(t, pad, 8))
}

func TestPadNullSingleBytePartials(t *testing.T) {
	eof := &ReplayReader{{0}, {1}, {2}, {3}, {4}, {5}, {6}}
	pad := &NullInputPadder{Reader: eof}
	assert.Equal(t, []byte{0, 1, 2, 3, 4, 5, 6, 0}, readall(t, pad, 8))
}

func TestPadNullWithZeroPartials(t *testing.T) {
	eof := &ReplayReader{{0, 1, 2}, {}, {}, {3, 4, 5}, {}, {}, {6}}
	pad := &NullInputPadder{Reader: eof}
	assert.Equal(t, []byte{0, 1, 2, 3, 4, 5, 6, 0}, readall(t, pad, 8))
}

func TestPadNullFullReadsOnly(t *testing.T) {
	eof := &ReplayReader{{0, 1, 2, 3, 4, 5, 6, 7}, {0, 1, 2, 3, 4, 5, 6, 7}}
	pad := &NullInputPadder{Reader: eof}
	assert.Equal(t, []byte{0, 1, 2, 3, 4, 5, 6, 7, 0, 1, 2, 3, 4, 5, 6, 7}, readall(t, pad, 8))
}

func TestPadNullFullThenPartial(t *testing.T) {
	eof := &ReplayReader{{0, 1, 2, 3, 4, 5, 6, 7}, {0, 1, 2, 3}}
	pad := &NullInputPadder{Reader: eof}
	assert.Equal(t, []byte{0, 1, 2, 3, 4, 5, 6, 7, 0, 1, 2, 3, 0, 0, 0, 0}, readall(t, pad, 8))
}

func TestPadNullSingleByteChunks(t *testing.T) {
	eof := &ReplayReader{{0}, {1}, {}, {2}, {}, {3}, {}, {4}, {}, {5}, {6}, {7}}
	pad := &NullInputPadder{Reader: eof}
	assert.Equal(t, []byte{0, 1, 2, 3, 4, 5, 6, 7}, readall(t, pad, 1))
}

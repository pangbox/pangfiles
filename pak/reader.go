package pak

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/go-restruct/restruct"
	"github.com/pangbox/pangfiles/crypto/pyxtea"
)

// Errors returned by the reader.
var (
	// ErrInvalidSignature is returned when the pak file contains an invalid
	// signature.
	ErrInvalidSignature = errors.New("invalid signature")
)

// ReadAtSeeker is the combination of io.ReadSeeker and io.ReaderAt.
type ReadAtSeeker interface {
	io.ReadSeeker
	io.ReaderAt
}

// Reader reads data from a pak file.
type Reader struct {
	k pyxtea.Key
	r ReadAtSeeker
	t TrailerData
}

// NewReader returns a new reader.
func NewReader(k pyxtea.Key, r ReadAtSeeker) (*Reader, error) {
	n := Reader{k: k, r: r}

	buf := [TrailerLen]byte{}

	// Read trailer.
	if _, err := r.Seek(-TrailerLen, io.SeekEnd); err != nil {
		return nil, fmt.Errorf("seeking to trailer: %w", err)
	}

	// Read trailer
	if _, err := r.Read(buf[:]); err != nil {
		return nil, fmt.Errorf("reading trailer: %w", err)
	}
	restruct.Unpack(buf[:], binary.LittleEndian, &n.t)
	if n.t.Signature != 0x12 {
		return nil, ErrInvalidSignature
	}

	return &n, nil
}

// ReadFileTable reads the file table entirely. The iteration is stopped if
// callback returns false.
func (r *Reader) ReadFileTable(callback func(path string, entry FileEntryData) bool) error {
	buf := [256]byte{}

	if _, err := r.r.Seek(int64(r.t.FileListOffset), 0); err != nil {
		return fmt.Errorf("seeking to file table: %w", err)
	}

	for i := uint32(0); i < r.t.FileCount; i++ {
		// Read file entry.
		entry := FileEntryData{}
		if _, err := r.r.Read(buf[:14]); err != nil {
			return fmt.Errorf("reading file entry %d: %w", i, err)
		}

		// Handle xtea encryption for the metadata.
		useXTEA := buf[1] >= 4
		if useXTEA {
			tmp := append(buf[2:6], buf[10:14]...)
			if err := pyxtea.Decipher(r.k, tmp); err != nil {
				return fmt.Errorf("decrypting xtea metadata for file entry %d: %w", i, err)
			}
			copy(buf[2:6], tmp[0:4])
			copy(buf[10:14], tmp[4:8])
		}

		// Deserialize metadata.
		if err := restruct.Unpack(buf[:14], binary.LittleEndian, &entry); err != nil {
			return fmt.Errorf("unpacking file entry %d: %w", i, err)
		}

		// Read and, if needed, decrypt, path.
		path := ""
		if useXTEA {
			entry.Compression ^= 0x20
			if _, err := r.r.Read(buf[:int(entry.PathLength)]); err != nil {
				return fmt.Errorf("reading xtea path for file entry %d: %w", i, err)
			}
			if err := pyxtea.Decipher(r.k, buf[:int(entry.PathLength)]); err != nil {
				return fmt.Errorf("decrypting xtea path for file entry %d: %w", i, err)
			}
			path = string(bytes.Trim(buf[:int(entry.PathLength)], "\xCD\x00"))
		} else {
			if _, err := r.r.Read(buf[:int(entry.PathLength)+1]); err != nil {
				return fmt.Errorf("reading legacy path for file entry %d: %w", i, err)
			}
			for j := byte(0); j < entry.PathLength; j++ {
				buf[j] ^= 0x71
			}
			path = string(buf[:int(entry.PathLength)])
		}
		if !callback(path, entry) {
			return nil
		}
	}

	return nil
}

// ReadFile reads an entire file.
func (r *Reader) ReadFile(entry FileEntryData) ([]byte, error) {
	uncompressed, err := decompress(entry, r.r)
	if err != nil {
		return nil, err
	}
	return uncompressed, nil
}

// CalcFileSize calculates the actual filesize of a compressed file.
func (r *Reader) CalcFileSize(entry FileEntryData) (int64, error) {
	return fastfilesize(entry, r.r)
}

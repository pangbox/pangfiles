package pak

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/go-restruct/restruct"
	"github.com/pangbox/pangfiles/crypto/pyxtea"
	"golang.org/x/text/encoding/korean"
)

// Errors returned by the reader.
var (
	// ErrInvalidSignature is returned when the pak file contains an invalid
	// signature.
	ErrInvalidSignature = errors.New("invalid signature")
	// ErrStopIteration is returned if iteration ends early because the
	// callback returned false.
	ErrStopIteration = errors.New("iteration stopped")
)

// ReaderAtLen is an io.ReaderAt that supports returning length.
type ReaderAtLen interface {
	io.ReaderAt
	Len() int
}

// Reader reads data from a pak file.
type Reader struct {
	k pyxtea.Key
	r ReaderAtLen
	t TrailerData
}

// NewReader returns a new reader.
func NewReader(k pyxtea.Key, r ReaderAtLen) (*Reader, error) {
	n := Reader{k: k, r: r}
	buf := [TrailerLen]byte{}
	// Read trailer
	if _, err := r.ReadAt(buf[:], int64(r.Len()-TrailerLen)); err != nil {
		return nil, fmt.Errorf("reading trailer: %w", err)
	}
	if err := restruct.Unpack(buf[:], binary.LittleEndian, &n.t); err != nil {
		return nil, err
	}
	if n.t.Signature != 0x12 {
		return nil, ErrInvalidSignature
	}

	return &n, nil
}

// ReadFileTable reads the file table entirely. The iteration is stopped if
// callback returns false.
func (r *Reader) ReadFileTable(callback func(path string, entry FileEntryData) bool) error {
	var err error
	n := 0
	buf := [256]byte{}
	tmp := [8]byte{}

	decoder := korean.EUCKR.NewDecoder()

	foffset := int64(r.t.FileListOffset)
	for i := uint32(0); i < r.t.FileCount; i++ {
		// Read file entry.
		entry := FileEntryData{}
		if n, err = r.r.ReadAt(buf[:14], foffset); err != nil {
			return fmt.Errorf("reading file entry %d: %w", i, err)
		}
		foffset += int64(n)

		// Handle xtea encryption for the metadata.
		useXTEA := buf[1]&0xF0 == 0x20
		if useXTEA {
			copy(tmp[0:4], buf[2:6])
			copy(tmp[4:8], buf[10:14])
			if err := pyxtea.Decipher(r.k, tmp[0:8]); err != nil {
				return fmt.Errorf("decrypting xtea metadata for file entry %d: %w", i, err)
			}
			copy(buf[2:6], tmp[0:4])
			copy(buf[10:14], tmp[4:8])
		}

		// Deserialize metadata.
		if err := restruct.Unpack(buf[:14], binary.LittleEndian, &entry); err != nil {
			return fmt.Errorf("unpacking file entry %d: %w", i, err)
		}

		// Default to XOR type for legacy entries.
		if entry.Type&EntryTypeMask == 0 {
			entry.Type |= EntryTypeXOR
		}

		path := []byte{}

		switch entry.Type & EntryTypeMask {
		case EntryTypeXOR:
			if n, err = r.r.ReadAt(buf[:int(entry.PathLength)+1], foffset); err != nil {
				return fmt.Errorf("reading legacy path for file entry %d: %w", i, err)
			}
			foffset += int64(n)
			entry.RealFileSize ^= 0x71
			for j := byte(0); j < entry.PathLength; j++ {
				buf[j] ^= 0x71
			}
			path = append(path, buf[:int(entry.PathLength)]...)

		case EntryTypeXTEA:
			if n, err = r.r.ReadAt(buf[:int(entry.PathLength)], foffset); err != nil {
				return fmt.Errorf("reading xtea path for file entry %d: %w", i, err)
			}
			foffset += int64(n)
			if err := pyxtea.Decipher(r.k, buf[:int(entry.PathLength)]); err != nil {
				return fmt.Errorf("decrypting xtea path for file entry %d: %w", i, err)
			}
			path = append(path, bytes.Trim(buf[:int(entry.PathLength)], "\xCD\x00")...)

		case EntryTypeBasic:
			if n, err = r.r.ReadAt(buf[:int(entry.PathLength)+1], foffset); err != nil {
				return fmt.Errorf("reading legacy path for file entry %d: %w", i, err)
			}
			foffset += int64(n)
			path = append(path, buf[:int(entry.PathLength)]...)
		}

		path, err = decoder.Bytes(path)
		if err != nil {
			panic(err)
		}

		if !callback(string(path), entry) {
			return ErrStopIteration
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
	return int64(entry.RealFileSize), nil
}

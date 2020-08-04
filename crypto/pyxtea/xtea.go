package pyxtea

import (
	"encoding/binary"
	"io"

	"github.com/pangbox/pangfiles/util"
)

// numRounds is the number of rounds used for PangYa's variant of XTEA.
const numRounds = 16

// BlockSize is the number of bytes in an XTEA block (for convenience.)
const BlockSize = 8

// Key is a type for XTEA keys.
type Key [4]uint32

// EncryptBlock encrypts a single block of data. Note that XTEA blocks are 8 bytes long.
func EncryptBlock(key Key, buf []byte) {
	data0 := binary.LittleEndian.Uint32(buf[0:4])
	data1 := binary.LittleEndian.Uint32(buf[4:8])
	delta := uint32(0x61C88647)
	sum := uint32(0)
	for i := 0; i < numRounds; i++ {
		data0 += (((data1 << 4) ^ (data1 >> 5)) + data1) ^ (sum + key[sum&3])
		sum -= delta
		data1 += (((data0 << 4) ^ (data0 >> 5)) + data0) ^ (sum + key[(sum>>11)&3])
	}
	binary.LittleEndian.PutUint32(buf[0:4], data0)
	binary.LittleEndian.PutUint32(buf[4:8], data1)
}

// DecryptBlock decrypts a single block of data. Note that XTEA blocks are 8 bytes long.
func DecryptBlock(key Key, buf []byte) {
	data0 := binary.LittleEndian.Uint32(buf[0:4])
	data1 := binary.LittleEndian.Uint32(buf[4:8])
	delta := uint32(0x61C88647)
	sum := uint32(0xE3779B90)
	for i := 0; i < numRounds; i++ {
		data1 -= (((data0 << 4) ^ (data0 >> 5)) + data0) ^ (sum + key[(sum>>11)&3])
		sum += delta
		data0 -= (((data1 << 4) ^ (data1 >> 5)) + data1) ^ (sum + key[sum&3])
	}
	binary.LittleEndian.PutUint32(buf[0:4], data0)
	binary.LittleEndian.PutUint32(buf[4:8], data1)
}

// EncipherStream encrypts a stream of data with XTEA.
func EncipherStream(key Key, r io.Reader, w io.Writer) error {
	buf := [8]byte{}

	for {
		n, err := r.Read(buf[:])
		if err == io.EOF {
			if n != 0 {
				return io.ErrUnexpectedEOF
			}
			return nil
		} else if err != nil {
			return err
		} else if n != 8 {
			return io.ErrUnexpectedEOF
		}

		EncryptBlock(key, buf[:])
		n, err = w.Write(buf[:])
		if err != nil {
			return err
		} else if n != 8 {
			return io.ErrShortWrite
		}
	}
}

// EncipherStreamPadNull encrypts a stream of data with XTEA, inserting null
// bytes to meet XTEA's block alignment requirements.
func EncipherStreamPadNull(key Key, r io.Reader, w io.Writer) error {
	return EncipherStream(key, &util.NullInputPadder{Reader: r}, w)
}

// Encipher encrypts a buffer of data with XTEA.
func Encipher(key Key, buf []byte) error {
	for {
		EncryptBlock(key, buf[0:8])
		if len(buf) == 0 {
			return nil
		} else if len(buf) < 8 {
			return io.ErrUnexpectedEOF
		}
		buf = buf[8:]
	}
}

// DecipherStream decrypts a stream of data with XTEA.
func DecipherStream(key Key, r io.Reader, w io.Writer) error {
	buf := [8]byte{}

	for {
		n, err := r.Read(buf[:])
		if err == io.EOF {
			if n != 0 {
				return io.ErrUnexpectedEOF
			}
			return nil
		} else if err != nil {
			return err
		} else if n != 8 {
			return io.ErrUnexpectedEOF
		}

		DecryptBlock(key, buf[:])
		n, err = w.Write(buf[:])
		if err != nil {
			return err
		} else if n != 8 {
			return io.ErrShortWrite
		}
	}
}

// Decipher decrypts a buffer of data with XTEA.
func Decipher(key Key, buf []byte) error {
	for {
		DecryptBlock(key, buf[0:8])
		buf = buf[8:]
		if len(buf) == 0 {
			return nil
		} else if len(buf) < 8 {
			return io.ErrUnexpectedEOF
		}
	}
}

// DecipherStreamTrimNull decrypts a stream of data with XTEA and trims nulls
// at the end.
func DecipherStreamTrimNull(key Key, r io.Reader, w io.Writer) error {
	buf := [8]byte{}

	nullrun := int64(0)

	for {
		n, err := r.Read(buf[:])
		if err == io.EOF {
			if n != 0 {
				return io.ErrUnexpectedEOF
			}
			return nil
		} else if err != nil {
			return err
		} else if n != 8 {
			return io.ErrUnexpectedEOF
		}

		DecryptBlock(key, buf[:])

		if nullrun == 0 {
			// We're not currently on a null run. Check for length of null suffix.
			for i := 7; i >= 0; i-- {
				if buf[i] != 0 {
					break
				}
				nullrun++
			}
		} else {
			// We are on a null run; check to see if it continues.
			if buf[0] == 0 && buf[1] == 0 && buf[2] == 0 && buf[3] == 0 && buf[4] == 0 && buf[5] == 0 && buf[6] == 0 && buf[7] == 0 {
				// Null run continues, keep buffering.
				nullrun += 8
				continue
			}

			// End of null run; we've hit non-null bytes. Dump null run onto stream.
			io.CopyN(w, util.NullReader{}, nullrun)
			if err != nil {
				return err
			} else if n != 8 {
				return io.ErrShortWrite
			}
			nullrun = 0
		}

		// Write all non-null-suffixed bytes.
		n, err = w.Write(buf[0 : 8-nullrun])
		if err != nil {
			return err
		} else if n != 8 {
			return io.ErrShortWrite
		}
	}
}

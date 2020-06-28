package pyxtea

import (
	"encoding/binary"
	"io"
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
		if len(buf) == 0 {
			return nil
		} else if len(buf) < 8 {
			return io.ErrUnexpectedEOF
		}
		buf = buf[8:]
	}
}

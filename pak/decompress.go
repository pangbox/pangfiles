package pak

import (
	"encoding/binary"
	"io"
)

var valuePad = []uint16{
	0xFF21, 0x834F, 0x675F, 0x0034, 0xF237, 0x815F, 0x4765, 0x0233,
}

func decompress(entry FileEntryData, f io.ReaderAt) ([]byte, error) {
	var out []byte

	if entry.Compression == 0 {
		out = make([]byte, entry.FileSize)
		_, err := f.ReadAt(out, int64(entry.Offset))
		if err != nil {
			return nil, err
		}
		return out, nil
	}

	buf := [2]byte{}

	var counter, seq, realseq byte

	off := int64(entry.Offset)

	readlen := int64(entry.FileSize)
	for j := int64(0); j < readlen; {
		if counter == 0 {
			_, err := f.ReadAt(buf[0:1], off+j)
			if err != nil {
				return []byte{}, err
			}
			seq = buf[0]
			realseq = seq
			j++

			if entry.Compression == 3 {
				seq ^= 0xC8
			}
		} else {
			seq >>= 1
		}

		if seq&1 == 1 {
			_, err := f.ReadAt(buf[0:2], off+j)
			if err != nil {
				return []byte{}, err
			}
			value := binary.LittleEndian.Uint16(buf[0:2])
			j += 2

			if entry.Compression == 3 {
				value ^= valuePad[(realseq>>3)&7]
			}

			off := int(value & 0xFFF)
			size := int((value >> 12) + 2)
			out = append(out, make([]byte, size)...)
			copy(out[len(out)-size:], out[len(out)-off-size:len(out)-off])
		} else {
			_, err := f.ReadAt(buf[0:1], off+j)
			if err != nil {
				return []byte{}, err
			}
			out = append(out, buf[0])
			j++
		}
		counter = (counter + 1) & 7
	}
	return out, nil
}

package pak

import (
	"encoding/binary"
	"io"
)

var valuePad = []uint16{65313, 33615, 26463, 52, 62007, 33119, 18277, 563}

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

			off := int(value & 0x0FFF)
			size := int((value >> 0x0C) + 2)
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

func fastfilesize(entry FileEntryData, f io.ReaderAt) (int64, error) {
	if entry.Compression == 0 {
		return int64(entry.FileSize), nil
	}

	filesize := int64(0)
	buf := [2]byte{}

	var counter, seq, realseq byte

	off := int64(entry.Offset)

	readlen := int64(entry.FileSize)
	for j := int64(0); j < readlen; {
		if counter == 0 {
			_, err := f.ReadAt(buf[0:1], off+j)
			if err != nil {
				return 0, err
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
				return 0, err
			}
			value := binary.LittleEndian.Uint16(buf[0:2])
			j += 2

			if entry.Compression == 3 {
				value ^= valuePad[(realseq>>3)&7]
			}

			filesize += int64((value >> 0x0C) + 2)
		} else {
			filesize++
			j++
		}
		counter = (counter + 1) & 7
	}
	return filesize, nil
}

package pycrc32

import "testing"

func TestCrc(t *testing.T) {
	tests := []struct {
		data     []byte
		checksum uint32
	}{
		{
			data:     []byte{},
			checksum: 0,
		},
		{
			data:     []byte{0x20, 0x20, 0x20, 0x20},
			checksum: 0xfe0fd94e,
		},
		{
			data:     []byte{'0', '0', '0', '1'},
			checksum: 0xfc62a689,
		},
	}

	for _, test := range tests {
		checksum := FileChecksum(test.data)
		if checksum != test.checksum {
			t.Errorf("hashing [% 02x]: expected 0x%08x, got 0x%08x", test.data, test.checksum, checksum)
		}
	}
}

package pyxtea

import (
	"bytes"
	"testing"
)

var allKeys = []Key{KeyUS, KeyJP, KeyTH, KeyEU, KeyID, KeyKR}

func TestEncryptBlock(t *testing.T) {
	tests := []struct {
		key    Key
		input  [8]byte
		output [8]byte
	}{
		{
			KeyUS,
			[8]byte{0, 0, 0, 0, 0, 0, 0, 0},
			[8]byte{0x55, 0x23, 0x8e, 0xcd, 0x5e, 0x56, 0xe5, 0xc7},
		},
	}

	for _, test := range tests {
		buf := [8]byte{}
		copy(buf[:], test.input[:])
		EncryptBlock(test.key, buf[:])
		if !bytes.Equal(buf[:], test.output[:]) {
			t.Errorf("encrypting [% 02x] with %08x: expected [% 02x], got [% 02x]", test.input, test.key, test.output, buf)
		}
	}
}

func TestDecryptBlock(t *testing.T) {
	tests := []struct {
		key    Key
		input  [8]byte
		output [8]byte
	}{
		{
			KeyUS,
			[8]byte{0x55, 0x23, 0x8e, 0xcd, 0x5e, 0x56, 0xe5, 0xc7},
			[8]byte{0, 0, 0, 0, 0, 0, 0, 0},
		},
	}

	for _, test := range tests {
		buf := [8]byte{}
		copy(buf[:], test.input[:])
		DecryptBlock(test.key, buf[:])
		if !bytes.Equal(buf[:], test.output[:]) {
			t.Errorf("decrypting [% 02x] with %08x: expected [% 02x], got [% 02x]", test.input, test.key, test.output, buf)
		}
	}
}

func TestEncryptDecryptBlock(t *testing.T) {
	buffers := [][8]byte{
		{0, 0, 0, 0, 0, 0, 0, 0},
		{1, 1, 1, 1, 1, 1, 1, 1},
		{1, 2, 3, 4, 5, 6, 7, 8},
		{254, 253, 252, 251, 250, 249, 248, 247},
		{254, 254, 254, 254, 254, 254, 254, 254},
		{255, 255, 255, 255, 255, 255, 255, 255},
		{255, 255, 255, 255, 0, 0, 0, 0},
		{10, 20, 30, 40, 50, 60, 70, 80},
		{250, 240, 230, 220, 210, 200, 190, 180},
	}

	for _, test := range buffers {
		buf := [8]byte{}
		for _, key := range allKeys {
			copy(buf[:], test[:])
			EncryptBlock(key, buf[:])
			DecryptBlock(key, buf[:])
			if !bytes.Equal(buf[:], test[:]) {
				t.Errorf("encrypting and decrypting %02x with %08x: corrupted to %02x", test, key, buf)
			}
		}
	}
}

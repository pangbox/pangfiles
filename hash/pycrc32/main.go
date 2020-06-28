// Package pycrc32 includes helpers for PangYa's CRC32 implementation.
package pycrc32

import "hash/crc32"

const (
	// File is the CRC polynomial for calculating file checksums (fcrcs).
	File = 0x04c11db7
)

var (
	// FileTable is the CRC32 table for calculating file checksums (fcrcs).
	FileTable = crc32.MakeTable(File)
)

// FileChecksum returns the file checksum (fcrc) value for a given buffer.
func FileChecksum(data []byte) uint32 {
	return crc32.Checksum(data, FileTable)
}

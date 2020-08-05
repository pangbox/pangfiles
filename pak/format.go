package pak

// TrailerLen is the number of bytes the trailer takes up on-disk.
const TrailerLen = 9

// TrailerData is the data structure at the end of a Pak file.
type TrailerData struct {
	FileListOffset uint32
	FileCount      uint32
	Signature      byte
}

// Enumeration of values that can be combined to form a File Entry type.
const (
	// FileTypeBasic denotes an uncompressed file.
	FileTypeBasic = 0x00
	// FileTypeLz denotes a file compressed with a custom LZ77 scheme.
	FileTypeLz = 0x01
	// FileTypeDir denotes a directory entry.
	FileTypeDir = 0x02
	// FileTypeLz2 denotes an LZ77 compressed file with obfuscation.
	FileTypeLz2 = 0x03
	// FileTypeMask is a mask for grabbing the type of file for an entry.
	FileTypeMask = 0x0F

	// EntryTypeXOR is present on legacy XOR-padded entries. Filenames are
	// XOR'd with 0x71. If the EntryType is zero on-disk, it will be set to
	// XOR padded.
	EntryTypeXOR = 0x10
	// EntryTypeXTEA is present on modern XTEA-encrypted entries. Filenames
	// and some metadata are encrypted with an XTEA cipher using 16 rounds.
	EntryTypeXTEA = 0x20
	// EntryTypeBasic is an unencrypted entry type likely used for debugging.
	EntryTypeBasic = 0x80
	// EntryTypeMask is a mask for grabbing the type of entry obfuscation.
	EntryTypeMask = 0xF0
)

// FileEntryData is the data structure of each file entry in a Pak file.
type FileEntryData struct {
	PathLength     byte
	Type           byte
	Offset         uint32
	PackedFileSize uint32
	RealFileSize   uint32
}

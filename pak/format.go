package pak

// TrailerLen is the number of bytes the trailer takes up on-disk.
const TrailerLen = 9

// TrailerData is the data structure at the end of a Pak file.
type TrailerData struct {
	FileListOffset uint32
	FileCount      uint32
	Signature      byte
}

// FileEntryData is the data structure of each file entry in a Pak file.
type FileEntryData struct {
	PathLength     byte
	Compression    byte
	Offset         uint32
	PackedFileSize uint32
	RealFileSize   uint32
}

package updatelist

import (
	"hash/crc32"
	"io"
	"os"
	"path/filepath"

	"github.com/pangbox/pangfiles/encoding/litexml"
	"github.com/pangbox/pangfiles/hash/pycrc32"
)

const crcBufLen = 1 << 12

// FileInfo represents the per-file information in an UpdateList.
type FileInfo struct {
	Filename   string `attr:"fname"`
	Directory  string `attr:"fdir"`
	Size       int64  `attr:"fsize"`
	Crc        int32  `attr:"fcrc"`
	Date       string `attr:"fdate"`
	Time       string `attr:"ftime"`
	PackedName string `attr:"pname"`
	PackedSize int64  `attr:"psize"`
}

// UpdateFiles contains the files in the updatelist.
type UpdateFiles struct {
	Count int        `attr:"count"`
	Files []FileInfo `tag:"fileinfo"`
}

// Document is the document containing the update list.
type Document struct {
	Info          litexml.DocumentInfo
	PatchVer      string      `tag:"patchVer" attr:"value"`
	PatchNum      int         `tag:"patchNum" attr:"value"`
	UpdateListVer string      `tag:"updatelistVer" attr:"value"`
	UpdateFiles   UpdateFiles `tag:"updatefiles"`
}

// MakeFileInfo generates a FileInfo structure from an OS file.
func MakeFileInfo(basedir, dir string, finfo os.FileInfo, psize int64) (FileInfo, error) {
	name := finfo.Name()
	f, err := os.Open(filepath.Join(basedir, dir, name))
	if err != nil {
		return FileInfo{}, err
	}

	crc := uint32(0)
	flen := 0
	for {
		buf := [crcBufLen]byte{}
		c, err := f.Read(buf[:])
		if c == 0 && err == io.EOF {
			break
		} else if err != nil {
			return FileInfo{}, err
		}
		crc = crc32.Update(crc, pycrc32.FileTable, buf[:c])
		flen += c
	}

	return FileInfo{
		Filename:   name,
		Directory:  dir,
		Size:       finfo.Size(),
		Crc:        int32(crc),
		Date:       finfo.ModTime().Format("2006-01-02"),
		Time:       finfo.ModTime().Format("15:04:05"),
		PackedName: name + ".zip",
		PackedSize: psize,
	}, nil
}

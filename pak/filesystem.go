package pak

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pangbox/pangfiles/crypto/pyxtea"
	"golang.org/x/exp/mmap"
)

var (
	// ErrFuseUnsupported is returned by Mount when Fuse is not supported.
	ErrFuseUnsupported = errors.New("fuse mounting not supported in build")
)

type fsfile struct {
	path   string
	entry  FileEntryData
	reader *Reader
	inode  uint64
	fsize  int64
}

func (f *fsfile) size() (int64, error) {
	if f.fsize == -1 {
		s, err := f.reader.CalcFileSize(f.entry)
		if err != nil {
			return 0, err
		}
		f.fsize = s
	}
	return f.fsize, nil
}

type fsdir struct {
	path  string
	inode uint64
}

// FS is an in-memory filesystem for pak files.
type FS struct {
	filemap map[string]*fsfile
	filetbl []*fsfile
	readers []*Reader
	key     pyxtea.Key

	inodes  uint64
	dirtbl  []*fsdir
	rootdir *fsdir
}

// NewFS returns a new, empty pak filesystem.
func NewFS(key pyxtea.Key) *FS {
	fs := &FS{
		filemap: map[string]*fsfile{},
		key:     key,
	}

	fs.rootdir = fs.adddir("")

	return fs
}

// LoadPaks loads pak files from a series of patterns or paths.
func LoadPaks(key pyxtea.Key, patterns []string) (*FS, error) {
	fs := NewFS(key)
	for _, pattern := range patterns {
		err := fs.LoadPaksFromGlob(pattern)
		if err != nil {
			return nil, err
		}
	}
	return fs, nil
}

func (fs *FS) newinode() uint64 {
	fs.inodes++
	return fs.inodes
}

func basename(path string) string {
	if n := strings.LastIndex(path, "/"); n != -1 {
		return path[n+1:]
	}
	return path
}

func searchfiles(a []*fsfile, fn string) int {
	return sort.Search(len(a), func(i int) bool { return a[i].path >= fn })
}

func searchdirs(a []*fsdir, fn string) int {
	return sort.Search(len(a), func(i int) bool { return a[i].path >= fn })
}

func (fs *FS) adddir(path string) *fsdir {
	i := 0
	if len(fs.dirtbl) > 0 {
		i = searchdirs(fs.dirtbl, path)
		if i < len(fs.dirtbl) && fs.dirtbl[i].path == path {
			return fs.dirtbl[i]
		}
		fs.dirtbl = append(fs.dirtbl, nil)
		copy(fs.dirtbl[i+1:], fs.dirtbl[i:])
	} else {
		fs.dirtbl = append(fs.dirtbl, nil)
	}
	fs.dirtbl[i] = &fsdir{path, fs.newinode()}
	return fs.dirtbl[i]
}

func (fs *FS) addfile(path string, entry FileEntryData, reader *Reader) {
	// Add dirs.
	for i, c := range path {
		if c == '/' {
			fs.adddir(path[:i])
		}
	}

	name := basename(path)
	if n, ok := fs.filemap[name]; ok {
		// Just overwrite.
		n.entry = entry
		n.reader = reader
		n.fsize = -1
	} else {
		// Add file.
		n := &fsfile{path, entry, reader, fs.newinode(), -1}
		if len(fs.filetbl) > 0 {
			i := searchfiles(fs.filetbl, path)
			fs.filetbl = append(fs.filetbl, nil)
			copy(fs.filetbl[i+1:], fs.filetbl[i:])
			fs.filetbl[i] = n
		} else {
			fs.filetbl = append(fs.filetbl, n)
		}
		fs.filemap[name] = n
	}
}

// AddPak adds a new pak on top of the filesystem.
func (fs *FS) AddPak(reader *Reader) error {
	err := reader.ReadFileTable(func(path string, entry FileEntryData) bool {
		// Skip directory entries; we manually construct dirents.
		if entry.Type&FileTypeMask == FileTypeDir {
			return true
		}
		fs.addfile(path, entry, reader)
		return true
	})
	if err != nil {
		return err
	}
	fs.readers = append(fs.readers, reader)
	return nil
}

// AddPakFromFile adds a new pak on the filesystem from a path.
func (fs *FS) AddPakFromFile(path string) error {
	file, err := mmap.Open(path)
	if err != nil {
		return err
	}
	reader, err := NewReader(fs.key, file)
	if err != nil {
		return err
	}
	return fs.AddPak(reader)
}

// LoadPaksFromFiles loads pak files from a list of paths.
func (fs *FS) LoadPaksFromFiles(paths []string) error {
	for _, path := range paths {
		if err := fs.AddPakFromFile(path); err != nil {
			return err
		}
	}
	return nil
}

// LoadPaksFromGlob loads pak files using a glob pattern.
func (fs *FS) LoadPaksFromGlob(pattern string) error {
	paths, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	sort.Strings(paths)
	return fs.LoadPaksFromFiles(paths)
}

// NumFiles returns the number of files in the filesystem.
func (fs *FS) NumFiles() int {
	return len(fs.filetbl)
}

// ReadFileByIndex returns the path and data for a given file index.
func (fs *FS) ReadFileByIndex(index int) (string, []byte, error) {
	if index < 0 || index >= len(fs.filetbl) {
		return "", nil, errors.New("invalid index")
	}

	data, err := fs.filetbl[index].reader.ReadFile(fs.filetbl[index].entry)
	if err != nil {
		return "", nil, err
	}

	return fs.filetbl[index].path, data, nil
}

// NumDirectories returns the number of directories in the filesystem.
func (fs *FS) NumDirectories() int {
	return len(fs.dirtbl)
}

// Directory returns the directory at a given index.
func (fs *FS) Directory(index int) string {
	if index < 0 || index >= len(fs.dirtbl) {
		return ""
	}
	return fs.dirtbl[index].path
}

// Extract extracts the filesystem onto the host disk.
func (fs *FS) Extract(dest string) error {
	for _, dir := range fs.dirtbl {
		if err := os.MkdirAll(filepath.Join(dest, dir.path), 0755); err != nil {
			return err
		}
	}
	for _, file := range fs.filetbl {
		data, err := file.reader.ReadFile(file.entry)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(filepath.Join(dest, file.path), data, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

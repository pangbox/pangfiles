// +build !nofuse
// +build windows cgo

package pak

import (
	"errors"
	"log"
	"strings"

	"github.com/billziss-gh/cgofuse/fuse"
)

// Implementation of pak fuse used on most platforms, CGo required for non-windows.

const (
	// FuseImplementation describes the fuse implementation in use in this build.
	FuseImplementation = "cgofuse"
)

// Mount mounts a pak filesystem via FUSE.
func (fs *FS) Mount(mountpoint string) error {
	fusefs := &cfsfuse{fs: fs}
	host := fuse.NewFileSystemHost(fusefs)
	if !host.Mount(mountpoint, nil) {
		return errors.New("failed to mount filesystem")
	}
	return nil
}

type cfsfuse struct {
	fuse.FileSystemBase
	fs *FS
	fd []cfusefd
}

type cfusefile struct {
	file *fsfile
	dir  *fsdir
	fs   *cfsfuse
}

type cfusefd struct {
	data []byte
}

func (f *cfsfuse) lookup(path string) (*cfusefile, int) {
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}
	if path == "" {
		return &cfusefile{dir: f.fs.rootdir, fs: f}, 0
	}
	if i := searchdirs(f.fs.dirtbl, path); i < len(f.fs.dirtbl) && f.fs.dirtbl[i].path == path {
		return &cfusefile{dir: f.fs.dirtbl[i], fs: f}, 0
	}
	if i := searchfiles(f.fs.filetbl, path); i < len(f.fs.filetbl) && f.fs.filetbl[i].path == path {
		return &cfusefile{file: f.fs.filetbl[i], fs: f}, 0
	}
	return nil, -fuse.ENOENT
}

func (f *cfsfuse) getfattr(d *fsfile, stat *fuse.Stat_t) {
	stat.Ino = d.inode
	stat.Mode = fuse.S_IFREG | 0o444
	size, err := d.size()
	if err != nil {
		log.Printf("Error getting filesize for %q: %s", d.path, err)
	}
	stat.Size = int64(size)
}

func (f *cfsfuse) getdattr(d *fsdir, stat *fuse.Stat_t) {
	stat.Ino = d.inode
	stat.Mode = fuse.S_IFDIR | 0o555
}

func (f *cfsfuse) Open(path string, flags int) (errc int, fh uint64) {
	d, errc := f.lookup(path)
	if errc != 0 || d.file == nil {
		return errc, ^uint64(0)
	}
	data, err := d.file.reader.ReadFile(d.file.entry)
	if err != nil {
		log.Printf("Error reading file for %q: %s", path, err)
	}

	f.fd = append(f.fd, cfusefd{
		data: data,
	})
	return 0, uint64(len(f.fd) - 1)
}

func (f *cfsfuse) Getattr(path string, stat *fuse.Stat_t, fh uint64) (errc int) {
	d, errc := f.lookup(path)
	if errc != 0 {
		log.Println("Getattr", path, "=> not found")
		return errc
	}

	switch {
	case d.dir != nil:
		f.getdattr(d.dir, stat)
	case d.file != nil:
		f.getfattr(d.file, stat)
	default:
		return -fuse.ENOENT
	}
	return 0
}

func (f *cfsfuse) Read(path string, buff []byte, offset int64, fh uint64) (n int) {
	if int(fh) > len(f.fd) || fh < 0 {
		return 0
	}
	fd := f.fd[fh]
	size := len(buff)
	if size > len(fd.data)-int(offset) {
		size = len(fd.data) - int(offset)
	}
	if offset < 0 {
		return 0
	}
	n = copy(buff, fd.data[offset:int(offset)+size])
	return
}

func (f *cfsfuse) Readdir(path string, fill func(name string, stat *fuse.Stat_t, ofst int64) bool, offset int64, fh uint64) (errc int) {
	fill(".", nil, 0)
	fill("..", nil, 0)
	prefix := path
	if len(prefix) > 0 && prefix[0] == '/' {
		prefix = prefix[1:]
	}
	if prefix != "" {
		prefix += "/"
	}
	for i := searchdirs(f.fs.dirtbl, prefix); i < len(f.fs.dirtbl); i++ {
		subdir := f.fs.dirtbl[i]
		if !strings.HasPrefix(subdir.path, prefix) {
			break
		}
		if subdir.path == "" {
			continue
		}
		dirname := subdir.path[len(prefix):]
		if strings.ContainsRune(dirname, '/') {
			continue
		}
		stat := &fuse.Stat_t{}
		f.getdattr(subdir, stat)
		fill(dirname, stat, 0)
	}
	for i := searchfiles(f.fs.filetbl, prefix); i < len(f.fs.filetbl); i++ {
		file := f.fs.filetbl[i]
		if !strings.HasPrefix(file.path, prefix) {
			break
		}
		filename := file.path[len(prefix):]
		if strings.ContainsRune(filename, '/') {
			continue
		}
		stat := &fuse.Stat_t{}
		f.getfattr(file, stat)
		fill(filename, nil, 0)
	}
	return 0
}

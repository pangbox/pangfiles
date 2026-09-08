//go:build !nofuse && (freebsd || linux)
// +build !nofuse
// +build freebsd linux

package pak

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	gofs "github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
)

// Implementation of pak FUSE used on Linux and FreeBSD, no CGO required.

const (
	// FuseImplementation describes the fuse implementation in use in this build.
	FuseImplementation = "gofuse"
)

// Mount mounts a pak filesystem via FUSE.
func (fs *FS) Mount(mountpoint string) error {
	root := newFuseRoot(fs)
	server, err := gofs.Mount(mountpoint, root, fuseOptions(fs))
	if err != nil {
		return err
	}

	interrupt := make(chan os.Signal, 1)
	done := make(chan struct{})
	signal.Notify(interrupt, os.Interrupt)
	go func() {
		select {
		case <-interrupt:
			fmt.Println("Received interrupt, exiting.")
			if err := server.Unmount(); err != nil {
				fmt.Printf("fuse.Unmount: %v\n", err)
			}
		case <-done:
		}
	}()

	server.Wait()
	signal.Stop(interrupt)
	close(done)
	return nil
}

func fuseOptions(fs *FS) *gofs.Options {
	return &gofs.Options{
		MountOptions: fuse.MountOptions{
			FsName:  "pakfs",
			Name:    "pakfs",
			Options: []string{"ro"},
		},
		RootStableAttr: &gofs.StableAttr{Ino: fs.rootdir.inode},
	}
}

func newFuseRoot(fs *FS) *fusedir {
	return &fusedir{dir: fs.rootdir, fs: fs}
}

// fusedir implements a pseudo directory for the purpose of supporting FUSE.
type fusedir struct {
	gofs.Inode
	dir *fsdir
	fs  *FS
}

var _ gofs.NodeGetattrer = (*fusedir)(nil)
var _ gofs.NodeLookuper = (*fusedir)(nil)
var _ gofs.NodeReaddirer = (*fusedir)(nil)

func (d *fusedir) Getattr(_ context.Context, _ gofs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	d.setAttr(&out.Attr)
	return 0
}

func (d *fusedir) Lookup(ctx context.Context, name string, out *fuse.EntryOut) (*gofs.Inode, syscall.Errno) {
	path := d.dir.path
	if path != "" {
		path += "/"
	}
	path += name

	if i := searchdirs(d.fs.dirtbl, path); i < len(d.fs.dirtbl) && d.fs.dirtbl[i].path == path {
		dir := &fusedir{dir: d.fs.dirtbl[i], fs: d.fs}
		dir.setAttr(&out.Attr)
		return d.NewInode(ctx, dir, gofs.StableAttr{
			Mode: fuse.S_IFDIR,
			Ino:  dir.dir.inode,
		}), 0
	}

	if i := searchfiles(d.fs.filetbl, path); i < len(d.fs.filetbl) && d.fs.filetbl[i].path == path {
		file := &fusefile{file: d.fs.filetbl[i]}
		if errno := file.setAttr(&out.Attr); errno != 0 {
			return nil, errno
		}
		return d.NewInode(ctx, file, gofs.StableAttr{
			Mode: fuse.S_IFREG,
			Ino:  file.file.inode,
		}), 0
	}

	return nil, syscall.ENOENT
}

func (d *fusedir) Readdir(_ context.Context) (gofs.DirStream, syscall.Errno) {
	dirents := []fuse.DirEntry{}
	prefix := d.dir.path
	if prefix != "" {
		prefix += "/"
	}

	for i := searchdirs(d.fs.dirtbl, prefix); i < len(d.fs.dirtbl); i++ {
		dir := d.fs.dirtbl[i]
		if !strings.HasPrefix(dir.path, prefix) {
			break
		}
		if dir.path == "" {
			continue
		}
		name := dir.path[len(prefix):]
		if strings.ContainsRune(name, '/') {
			continue
		}
		dirents = append(dirents, fuse.DirEntry{
			Mode: fuse.S_IFDIR,
			Name: name,
			Ino:  dir.inode,
		})
	}

	for i := searchfiles(d.fs.filetbl, prefix); i < len(d.fs.filetbl); i++ {
		file := d.fs.filetbl[i]
		if !strings.HasPrefix(file.path, prefix) {
			break
		}
		name := file.path[len(prefix):]
		if strings.ContainsRune(name, '/') {
			continue
		}
		dirents = append(dirents, fuse.DirEntry{
			Mode: fuse.S_IFREG,
			Name: name,
			Ino:  file.inode,
		})
	}

	return gofs.NewListDirStream(dirents), 0
}

func (d *fusedir) setAttr(out *fuse.Attr) {
	out.Ino = d.dir.inode
	out.Mode = fuse.S_IFDIR | 0o555
}

// fusefile implements a file for FUSE.
type fusefile struct {
	gofs.Inode
	file *fsfile

	sizeOnce sync.Once
	size     int64
	sizeErr  error
	dataOnce sync.Once
	data     []byte
	dataErr  error
}

var _ gofs.NodeGetattrer = (*fusefile)(nil)
var _ gofs.NodeOpener = (*fusefile)(nil)
var _ gofs.NodeReader = (*fusefile)(nil)

func (f *fusefile) Getattr(_ context.Context, _ gofs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	return f.setAttr(&out.Attr)
}

func (f *fusefile) setAttr(out *fuse.Attr) syscall.Errno {
	f.sizeOnce.Do(func() {
		f.size, f.sizeErr = f.file.size()
	})
	if f.sizeErr != nil {
		return gofs.ToErrno(f.sizeErr)
	}

	out.Ino = f.file.inode
	out.Mode = fuse.S_IFREG | 0o444
	out.Size = uint64(f.size)
	return 0
}

func (f *fusefile) Open(_ context.Context, flags uint32) (gofs.FileHandle, uint32, syscall.Errno) {
	if flags&syscall.O_ACCMODE != syscall.O_RDONLY {
		return nil, 0, syscall.EROFS
	}
	if errno := f.load(); errno != 0 {
		return nil, 0, errno
	}
	return nil, fuse.FOPEN_KEEP_CACHE, 0
}

func (f *fusefile) Read(_ context.Context, _ gofs.FileHandle, dest []byte, off int64) (fuse.ReadResult, syscall.Errno) {
	if off < 0 {
		return nil, syscall.EINVAL
	}
	if errno := f.load(); errno != 0 {
		return nil, errno
	}
	if off >= int64(len(f.data)) {
		return fuse.ReadResultData(nil), 0
	}

	end := len(f.data)
	if remaining := int64(len(f.data)) - off; int64(len(dest)) < remaining {
		end = int(off) + len(dest)
	}
	return fuse.ReadResultData(f.data[int(off):end]), 0
}

func (f *fusefile) load() syscall.Errno {
	f.dataOnce.Do(func() {
		f.data, f.dataErr = f.file.reader.ReadFile(f.file.entry)
	})
	return gofs.ToErrno(f.dataErr)
}

// +build freebsd linux

package pak

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"bazil.org/fuse"
	fusefs "bazil.org/fuse/fs"
)

// Mount mounts a pak filesystem via FUSE.
func (fs *FS) Mount(mountpoint string) error {
	c, err := fuse.Mount(
		mountpoint,
		fuse.FSName("pakfs"),
		fuse.Subtype("pakfs"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	i := make(chan os.Signal, 1)
	signal.Notify(i, os.Interrupt)
	go func() {
		<-i
		fmt.Println("Received interrupt, exiting.")
		fuse.Unmount(mountpoint)
		os.Exit(0)
	}()

	return fusefs.Serve(c, fs)
}

// Root implements FUSE
func (fs *FS) Root() (fusefs.Node, error) {
	return &fusedir{fs.rootdir, fs}, nil
}

// fusedir implements a pseudo directory for the purpose of supporting FUSE.
type fusedir struct {
	dir *fsdir
	fs  *FS
}

// Attr implements FUSE
func (d *fusedir) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = d.dir.inode
	a.Mode = os.ModeDir | 0o555
	return nil
}

// Lookup implements FUSE
func (d *fusedir) Lookup(ctx context.Context, name string) (fusefs.Node, error) {
	path := d.dir.path
	if path != "" {
		path += "/"
	}
	path += name

	i := searchdirs(d.fs.dirtbl, path)
	if i < len(d.fs.dirtbl) && d.fs.dirtbl[i].path == path {
		return &fusedir{dir: d.fs.dirtbl[i], fs: d.fs}, nil
	}

	i = searchfiles(d.fs.filetbl, path)
	if i < len(d.fs.filetbl) && d.fs.filetbl[i].path == path {
		return &fusefile{file: d.fs.filetbl[i], fs: d.fs}, nil
	}

	return nil, syscall.ENOENT
}

// ReadDirAll implements FUSE
func (d *fusedir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	dirents := []fuse.Dirent{}

	prefix := d.dir.path
	if prefix != "" {
		prefix += "/"
	}
	i := searchdirs(d.fs.dirtbl, prefix)
	for ; i < len(d.fs.dirtbl); i++ {
		subdir := d.fs.dirtbl[i]
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
		dirents = append(dirents, fuse.Dirent{
			Inode: subdir.inode,
			Name:  dirname,
			Type:  fuse.DT_Dir,
		})
	}

	i = searchfiles(d.fs.filetbl, prefix)
	for ; i < len(d.fs.filetbl); i++ {
		file := d.fs.filetbl[i]
		if !strings.HasPrefix(file.path, prefix) {
			break
		}
		dirname := file.path[len(prefix):]
		if strings.ContainsRune(dirname, '/') {
			continue
		}
		dirents = append(dirents, fuse.Dirent{
			Inode: file.inode,
			Name:  dirname,
			Type:  fuse.DT_File,
		})
	}

	return dirents, nil
}

// fusefile implements a file for FUSE.
type fusefile struct {
	file *fsfile
	fs   *FS
}

// Attr implements FUSE
func (f *fusefile) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = f.file.inode
	a.Mode = 0o444

	size, err := f.file.size()
	if err != nil {
		return err
	}
	a.Size = uint64(size)

	return nil
}

// ReadAll implements FUSE
func (f *fusefile) ReadAll(ctx context.Context) ([]byte, error) {
	data, err := f.file.reader.ReadFile(f.file.entry)
	if err != nil {
		return nil, err
	}
	return data, nil
}

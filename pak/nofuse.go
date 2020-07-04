// +build !freebsd,!linux

package pak

import "errors"

// Mount mounts a pak filesystem via FUSE.
func (fs *FS) Mount(mountpoint string) error {
	return errors.New("fuse mouinting not supported on this platform yet.")
}

// +build nofuse !freebsd,!linux,!windows,!cgo

package pak

// Pak fuse stub for when cgo is not enabled and it is required, or when fuse
// is explicitly disabled with the nofuse build tag.

const (
	// FuseImplementation describes the fuse implementation in use in this build.
	FuseImplementation = "nofuse"
)

// Mount mounts a pak filesystem via FUSE.
func (fs *FS) Mount(mountpoint string) error {
	return ErrFuseUnsupported
}

// +build linux

package block

import (
	"os"
	"syscall"
	"unsafe"
)

// Constants taken from Linux headers to avoid need for cgo
const (
	_BLKGETSIZE64      = 2148012658
	_HDIO_GET_IDENTITY = 0x030d
)

var (
	// Compile-time interface check
	_ devicer = &device{}
)

// A device is a Linux-specific block device.
type device struct {
	dev   *os.File
	fd    uintptr
	ioctl ioctlFunc
}

// openDevice is the operating system-specific entry point for New.
func openDevice(device string, flags int) (*Device, error) {
	dev, err := os.OpenFile(device, flags, 0)
	if err != nil {
		return nil, err
	}

	return newDevice(dev, dev.Fd(), ioctl)
}

// newDevice is the entry point for tests.  It accepts a file, file descriptor,
// and ioctlFunc, which can be swapped out easily for testing.
func newDevice(dev *os.File, fd uintptr, ioctl ioctlFunc) (*Device, error) {
	d := &Device{
		&device{
			dev:   dev,
			fd:    fd,
			ioctl: ioctl,
		},
	}

	// Check the size of the device; normal files and the like will return
	// ENOTTY here
	_, err := d.Size()
	if err == nil {
		return d, nil
	}

	// Error path: close the device
	_ = d.Close()

	// Check for a syscall error
	serr, ok := err.(*os.SyscallError)
	if !ok {
		return nil, err
	}

	// If ioctl() returns ENOTTY, this is not a block device
	if serr.Syscall == "ioctl" && serr.Err == syscall.ENOTTY {
		return nil, ErrNotBlockDevice
	}

	return nil, err
}

// Close closes the file descriptor for a block device.
func (d *device) Close() error {
	d.fd = 0
	return d.dev.Close()
}

// Identify queries a block device for its IDE identification info.
func (d *device) Identify() ([512]byte, error) {
	// TODO(mdlayher): possibly parse and return a struct instead of an array
	b := [512]byte{}
	_, err := d.ioctl(d.fd, _HDIO_GET_IDENTITY, uintptr(unsafe.Pointer(&b[0])))
	return b, err
}

// Size queries a block device for its total size in bytes.
func (d *device) Size() (uint64, error) {
	var size uint64
	_, err := d.ioctl(d.fd, _BLKGETSIZE64, uintptr(unsafe.Pointer(&size)))
	return size, err
}

// ioctlFunc is the signature for a function which can perform the ioctl syscall,
// or a mocked version of it.
type ioctlFunc func(fd uintptr, request int, argp uintptr) (uintptr, error)

// ioctl is a wrapper used to perform the ioctl syscall using the input
// file descriptor, request, and arguments pointer.
//
// ioctl is the default ioctlFunc implementation, and the one used when New
// is called.
func ioctl(fd uintptr, request int, argp uintptr) (uintptr, error) {
	ret, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		fd,
		uintptr(request),
		argp,
	)
	if errno != 0 {
		return 0, os.NewSyscallError("ioctl", errno)
	}

	return ret, nil
}

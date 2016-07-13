// Package block enables programmatic access to block devices.
package block

import (
	"errors"
	"io"
)

var (
	// ErrNotBlockDevice is returned when a device name passed to New is
	// not a block device.
	ErrNotBlockDevice = errors.New("not a block device")

	// ErrNotImplemented is returned when block device functionality is not
	// implemented on the current platform.
	ErrNotImplemented = errors.New("not implemented")
)

// A Device represents a block device.  It can be used to query information
// about the device, seek the device, or read and write data.
type Device struct {
	*device
}

// devicer is an internal interface which maintains parity between operating
// system implementations.
type devicer interface {
	Identify() ([512]byte, error)
	ReadSMART() (*SMARTData, error)
	Size() (uint64, error)

	io.Closer
	io.ReadWriteSeeker
	io.ReaderAt
	io.WriterAt
}

// New attempts to open a block device with the specified flags, and verifies
// that device is a block device.
//
// If device is not a block device, ErrNotBlockDevice is returned.
func New(device string, flags int) (*Device, error) {
	return openDevice(device, flags)
}

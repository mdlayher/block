// +build !linux

package block

var (
	// Compile-time interface check
	_ devicer = &device{}
)

// A device is an unimplemented block device.
type device struct{}

// openDevice is the operating system-specific entry point for New.
func openDevice(device string, flags int) (*Device, error) {
	return nil, ErrNotImplemented
}

// Close is not currently implemented on this platform.
func (d *device) Close() error { return ErrNotImplemented }

// Identify is not currently implemented on this platform.
func (d *device) Identify() ([512]byte, error) { return [512]byte{}, ErrNotImplemented }

// Size is not currently implemented on this platform.
func (d *device) Size() (uint64, error) { return 0, ErrNotImplemented }

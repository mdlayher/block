// +build !linux

package block

var _ devicer = &device{}

// A device is an unimplemented block device.
type device struct{}

// openDevice is the operating system-specific entry point for New.
func openDevice(device string, flags int) (*Device, error) {
	return nil, ErrNotImplemented
}

func (d *device) Close() error                                 { return ErrNotImplemented }
func (d *device) Identify() ([512]byte, error)                 { return [512]byte{}, ErrNotImplemented }
func (d *device) Size() (uint64, error)                        { return 0, ErrNotImplemented }
func (d *device) Read(b []byte) (int, error)                   { return 0, ErrNotImplemented }
func (d *device) ReadAt(b []byte, off int64) (int, error)      { return 0, ErrNotImplemented }
func (d *device) Seek(offset int64, whence int) (int64, error) { return 0, ErrNotImplemented }
func (d *device) Write(b []byte) (int, error)                  { return 0, ErrNotImplemented }
func (d *device) WriteAt(b []byte, off int64) (int, error)     { return 0, ErrNotImplemented }

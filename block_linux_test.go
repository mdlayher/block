// +build linux

package block

import (
	"bytes"
	"os"
	"syscall"
	"testing"
	"unsafe"

	"golang.org/x/sys/unix"
)

// Test for ioctl syscall error ENOTTY, meaning that an input file
// is not a block device, and ErrNotBlockDevice should be returned.
func Test_newDeviceErrNotBlockDevice(t *testing.T) {
	ioctlFn := ioctlSize(t, 0, os.NewSyscallError("ioctl", syscall.ENOTTY))

	_, err := newDevice(nil, 0, ioctlFn)
	if want, got := ErrNotBlockDevice, err; want != got {
		t.Fatalf("unexpected error:\n- want: %v\n-  got: %v", want, got)
	}
}

// Test for syscall error, but not ENOTTY.
func Test_newDeviceSyscallError(t *testing.T) {
	errSys := os.NewSyscallError("ioctl", syscall.ENOSYS)
	ioctlFn := ioctlSize(t, 0, errSys)

	_, err := newDevice(nil, 0, ioctlFn)
	if want, got := errSys, err; want != got {
		t.Fatalf("unexpected error:\n- want: %v\n-  got: %v", want, got)
	}
}

// Test for device OK.
func Test_newDeviceOK(t *testing.T) {
	ioctlFn := ioctlSize(t, 0, nil)

	_, err := newDevice(nil, 0, ioctlFn)
	if err != nil {
		t.Fatal(err)
	}
}

// Test that Device.Size properly returns expected size from an ioctl.
func TestDeviceSizeOK(t *testing.T) {
	const cSize uint64 = 1024
	ioctlFn := ioctlSize(t, cSize, nil)

	size, err := (&Device{&device{ioctl: ioctlFn}}).Size()
	if err != nil {
		t.Fatal(err)
	}

	if want, got := cSize, size; want != got {
		t.Fatalf("unexpected output size:\n- want: %v\n-  got: %v", want, got)
	}
}

// Test that Device.Identify properly returns expected data from an ioctl.
func TestDeviceIdentifyOK(t *testing.T) {
	// Input array and expected bytes
	in := [512]byte{}
	data := []byte{'f', 'o', 'o'}
	copy(in[0:3], data)

	ioctlFn := ioctlIdentify(t, in, nil)

	out, err := (&Device{&device{ioctl: ioctlFn}}).Identify()
	if err != nil {
		t.Fatal(err)
	}

	if want, got := data, out[0:3]; !bytes.Equal(want, got) {
		t.Fatalf("unexpected output bytes:\n- want: %v\n-  got: %v", want, got)
	}
}

// ioctlSize returns an ioctlFunc which expects to be used by Device.Size.  Its
// return values can be customized by size and err.
func ioctlSize(t *testing.T, size uint64, err error) ioctlFunc {
	return func(fd uintptr, request int, argp unsafe.Pointer) (uintptr, error) {
		if want, got := unix.BLKGETSIZE64, request; want != got {
			t.Fatalf("unexpected ioctl request constant:\n- want: %v\n-  got: %v", want, got)
		}

		// This is ugly, but it seems to get the job done
		p := (*uint64)(unsafe.Pointer(argp))
		*p = size
		argp = unsafe.Pointer(p)

		return 0, err
	}
}

// ioctlIdentify returns an ioctlFunc which expects to be used by
// Device.Identify.  Its return values can be customized by data and err.
func ioctlIdentify(t *testing.T, data [512]byte, err error) ioctlFunc {
	return func(fd uintptr, request int, argp unsafe.Pointer) (uintptr, error) {
		if want, got := hdioGetIdentity, request; want != got {
			t.Fatalf("unexpected ioctl request constant:\n- want: %v\n-  got: %v", want, got)
		}

		// This is ugly, but it seems to get the job done
		p := (*[512]byte)(unsafe.Pointer(argp))
		copy((*p)[:], data[:])
		argp = unsafe.Pointer(p)

		return 0, err
	}
}

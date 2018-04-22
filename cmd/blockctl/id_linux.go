//+build linux

package main

import (
	"strings"
	"unsafe"

	"golang.org/x/sys/unix"
)

// canParseID indicates that drive information can be parsed on Linux.
const canParseID = true

// parseID parses an ID structure from raw drive identification bytes.
func parseID(b [512]byte) *ID {
	id := *(*unix.HDDriveID)(unsafe.Pointer(&b[0]))

	return &ID{
		Model:    byteStr(id.Model[:]),
		Serial:   byteStr(id.Serial_no[:]),
		Firmware: byteStr(id.Fw_rev[:]),
	}
}

// byteStr converts bytes to a string and trims any null characters.
func byteStr(b []byte) string {
	return strings.TrimSuffix(strings.TrimSpace(string(b)), "\x00")
}

//+build !linux

package main

// canParseID indicates that drive information cannot be parsed on
// non-Linux platforms.
const canParseID = false

// parseID is not implemented on non-Linux platforms.
func parseID(_ [512]byte) *ID {
	return nil
}

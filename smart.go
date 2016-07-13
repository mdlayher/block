package block

import (
	"encoding/binary"
	"unsafe"
)

// nativeEndian is the native endian byte order for this machine, determined
// on package init.
var nativeEndian binary.ByteOrder

// Courtesy of golang.org/x/net/ipv4
func init() {
	i := uint32(1)
	b := (*[4]byte)(unsafe.Pointer(&i))
	if b[0] == 1 {
		nativeEndian = binary.LittleEndian
	} else {
		nativeEndian = binary.BigEndian
	}
}

func parseSMARTData(raw [516]byte) *SMARTData {
	sd := &SMARTData{}

	// Skip 4 bytes ATA header
	sd.Revision = nativeEndian.Uint16(raw[4:6])

	i := 6
	for c := 0; c < _NR_ATTRIBUTES; c++ {
		sd.Values[c].ID = raw[i]
		i++

		sd.Values[c].Status = nativeEndian.Uint16(raw[i : i+2])
		i += 2

		sd.Values[c].Value = raw[i]
		i++

		copy(sd.Values[c].Vendor[:], raw[i:i+8])
		i += 8
	}

	sd.OfflineStatus = raw[i]
	i++

	sd.Vendor1 = raw[i]
	i++

	sd.OfflineTimeout = nativeEndian.Uint16(raw[i : i+2])
	i += 2

	sd.Vendor2 = raw[i]
	i++

	sd.OfflineCapability = raw[i]
	i++

	sd.SMARTCapability = nativeEndian.Uint16(raw[i : i+2])
	i += 2

	copy(sd.Reserved[:], raw[i:i+16])
	i += 16

	copy(sd.Vendor[:], raw[i:i+125])
	i += 125

	sd.Checksum = raw[i]

	return sd
}

type SMARTData struct {
	Revision          uint16
	Values            [_NR_ATTRIBUTES]Value
	OfflineStatus     uint8
	Vendor1           uint8
	OfflineTimeout    uint16
	Vendor2           uint8
	OfflineCapability uint8
	SMARTCapability   uint16
	Reserved          [16]byte
	Vendor            [125]byte
	Checksum          uint8
}

const _NR_ATTRIBUTES = 30

type Value struct {
	ID     uint8
	Status uint16
	Value  uint8
	Vendor [8]byte
}

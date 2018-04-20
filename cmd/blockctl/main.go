// Command blockctl retrieves information about block devices.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	humanize "github.com/dustin/go-humanize"
	"github.com/mdlayher/block"
)

func main() {
	flag.Parse()
	path := flag.Arg(0)
	if path == "" {
		fmt.Println("usage: blockctl [device]")
		return
	}

	d, err := block.New(path, os.O_RDONLY)
	if err != nil {
		log.Fatalf("failed to open block device %q: %v", path, err)
	}
	defer d.Close()

	size, err := d.Size()
	if err != nil {
		log.Fatalf("failed to get device size: %v", err)
	}

	if !canParseID {
		// Cannot provide more information on non-Linux platforms.
		fmt.Printf("%s: %s\n", path, humanize.Bytes(size))
		return
	}

	b, err := d.Identify()
	if err != nil {
		log.Fatalf("failed to identify device: %v", err)
	}

	id := parseID(b)

	fmt.Printf("%s: %s, model: %q, serial: %q, firmware: %q\n",
		path, humanize.Bytes(size), id.Model, id.Serial, id.Firmware)
}

// ID is information parsed from a drive's raw identification information.
type ID struct {
	Model    string
	Serial   string
	Firmware string
}

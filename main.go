package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
)

func main() {
	chroot := flag.Bool("chroot", false, "Assume in chrooted environment")
	flag.Parse()

	if *chroot {
		if err := phaseTwo(); err != nil {
			log.Fatal(err)
		}
		return
	}

	if err := phaseOne(); err != nil {
		log.Fatal(err)
	}

	clear()
	title("Success! Please reboot your computer!")
}

func blkuuid(dev string) (string, error) {
	out, _, err := getSh("blkid", "-o", "value", dev)
	if err != nil {
		return "", err
	}

	sps := strings.Split(out, "\n")
	return strings.TrimSpace(sps[0]), nil
}

func title(s string) {
	fmt.Printf("##\n# %s\n##\n\n", s)
}

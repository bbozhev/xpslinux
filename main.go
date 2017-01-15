package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
)

func main() {
	chroot := flag.Bool("chroot", false, "Assume in chrooted environment")
	version := flag.Bool("version", false, "Print version")
	flag.Parse()

	if *version {
		fmt.Println("xpsarch v0.1.0")
		return
	}

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

func uid(username string) (int, error) {
	out, _, err := getSh("id", "--user", username)
	if err != nil {
		return -1, err
	}

	return strconv.Atoi(strings.TrimSpace(out))
}

func gid(username string) (int, error) {
	out, _, err := getSh("id", "--group", username)
	if err != nil {
		return -1, err
	}

	return strconv.Atoi(strings.TrimSpace(out))
}

func title(s string) {
	fmt.Printf("##\n# %s\n##\n\n", s)
}

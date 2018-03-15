package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"
)

const (
	devStorage    = "/dev/nvme0n1"
	bootPartition = "/dev/nvme0n1p1"
	rootPartition = "/dev/nvme0n1p2"
	swapPartition = "/dev/nvme0n1p3"
	devCryptroot  = "/dev/mapper/cryptroot"
	lvmRoot				= "/dev/mapper/vg0-root"
	lvmSwap				= "/dev/mapper/vg0-swap"
)

func phaseOne() error {
	if err := sh("setfont", "latarcyrheb-sun32"); err != nil {
		return err
	}
	clear()
	title("Welcome to xpsarch installer")
	time.Sleep(3 * time.Second)

	if !hasInternet() {
		return errors.New("failed to verify internet connection")
	}

	if err := sh("timedatectl", "set-ntp", "true"); err != nil {
		return err
	}

	if err := configDisk(); err != nil {
		return err
	}

	if err := configMirrors(); err != nil {
		return err
	}

	if err := installBootPackages(); err != nil {
		return err
	}

	if err := configFstab(); err != nil {
		return err
	}

	if err := archChroot(); err != nil {
		return err
	}

	return nil
}

func hasInternet() bool {
	conn, err := net.Dial("tcp", "www.archlinux.org:443")
	if err != nil {
		return false
	}
	return conn.Close() == nil
}

func configDisk() error {
	clear()
	title("Configuring disk")

	cmds := [][]string{
		{"parted", "--align", "optimal", devStorage, "mklabel", "gpt"},
		{"parted", "--align", "optimal", devStorage, "mkpart", "ESP", "fat32", "1MiB", "513MiB"},
		{"parted", devStorage, "set", "1", "boot", "on"},
		{"parted", "--align", "optimal", devStorage, "mkpart", "primary", "ext4", "513MiB", "100%"},
	}

	for _, c := range cmds {
		if err := sh(c...); err != nil {
			return err
		}
	}

	cmds = [][]string{
		{"cryptsetup", "--verify-passphrase", "--verbose", "-c aes-xts-plain64", "--key-size", "512", "--hash", "sha512", "--iter-time", "3000", "-y", "--use-random", "luksFormat", rootPartition},
		{"cryptsetup", "open", rootPartition, "cryptroot"},
	}

	fmt.Println("\nEnter disk encryption password")
	for _, c := range cmds {
		if err := sh(c...); err != nil {
			return err
		}
	}

	cmds = [][]string{
		{"pvcreate", devCryptroot},
		{"vgcreate", "vg0", devCryptroot},
		{"lvcreate", "--size", "16G", "vg0", "--name", "swap"},
		{"lvcreate", "-l", "+100%FREE", "vg0", "--name", "root"},
		{"mkfs.ext4", "-F", lvmRoot},
		{"mkswap", lvmSwap},
		{"mkfs.fat", "-F32", bootPartition},

		{"mount", lvmRoot, "/mnt"},
		{"mkdir", "/mnt/boot"},
		{"mount", bootPartition, "/mnt/boot"},
	}

	for _, c := range cmds {
		if err := sh(c...); err != nil {
			return err
		}
	}

	return nil
}

func configMirrors() error {
	clear()
	title("Configuring mirrors")

	country, err := chooseMirrorCountry()
	if err != nil {
		return err
	}

	mirsByCountry, _, err := getSh("grep", "-A1", "--no-group-separator", country, "/etc/pacman.d/mirrorlist")
	if err != nil {
		return err
	}

	err = fwrite("/etc/pacman.d/mirrorlist.national", mirsByCountry)
	if err != nil {
		return err
	}

	fmt.Printf("Searching for fastests mirrors in %s\n", country)
	top, _, err := getSh("rankmirrors", "-n", "5", "/etc/pacman.d/mirrorlist.national")
	if err != nil {
		return err
	}

	return fwrite("/etc/pacman.d/mirrorlist", top)
}

func installBootPackages() error {
	clear()
	title("Installing bootstrap packages")

	pkgs := []string{
		"base",
		"base-devel",

		// Fixes tiny font.
		"terminus-font",
	}

	return sh(append([]string{"pacstrap", "/mnt"}, pkgs...)...)
}

func configFstab() error {
	rootID, err := blkuuid(lvmRoot)
	if err != nil {
		return err
	}

	bootID, err := blkuuid(bootPartition)
	if err != nil {
		return err
	}

	swapID, err := blkuuid(lvmSwap)
	if err != nil {
		return err
	}

	return fwrite("/mnt/etc/fstab", fmt.Sprintf(fstab, rootID, bootID, swapID))
}

func archChroot() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	selfBin := filepath.Base(os.Args[0])
	absp := filepath.Join(dir, selfBin)

	if err := cp(absp, "/mnt"); err != nil {
		return err
	}

	if err := sh("arch-chroot", "/mnt", "/"+selfBin, "-chroot"); err != nil {
		return err
	}

	return os.Remove(filepath.Join("/mnt", selfBin))
}

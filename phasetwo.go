package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func phaseTwo() error {
	if err := sh("setfont", "ter-132n"); err != nil {
		return err
	}
	clear()
	title("Configuring user system")
	time.Sleep(3 * time.Second)

	if err := configTime(); err != nil {
		return err
	}

	if err := configLocale(); err != nil {
		return err
	}

	if err := configNetwork(); err != nil {
		return err
	}

	username, err := configUsers()
	if err != nil {
		return err
	}

	if err := installPackages(); err != nil {
		return err
	}

	if err := installAURs(username); err != nil {
		return err
	}

	if err := configGnome(username); err != nil {
		return err
	}

	if err := configServices(); err != nil {
		return err
	}

	return configBootloader()
}

func configTime() error {
	clear()
	title("Configuring time")

	zi, err := chooseZoneInfo()
	if err != nil {
		return err
	}

	if err := os.Remove("/etc/localtime"); err != nil {
		return err
	}

	err = os.Symlink(filepath.Join("/usr/share/zoneinfo", zi), "/etc/localtime")
	if err != nil {
		return err
	}

	return sh("hwclock", "--systohc", "--utc")
}

func configLocale() error {
	clear()
	title("Configuring locale")

	lc, err := chooseLocale()
	if err != nil {
		return err
	}

	if err := fappend("/etc/locale.gen", lc); err != nil {
		return err
	}

	if err := sh("locale-gen"); err != nil {
		return err
	}

	sps := strings.Split(lc, " ")
	err = fwrite("/etc/locale.conf", fmt.Sprintf(localeconf, sps[0]))
	if err != nil {
		return err
	}

	layout, err := chooseKeyboardLayout()
	if err != nil {
		return err
	}

	return fwrite("/etc/vconsole.conf", fmt.Sprintf(vconsoleconf, layout))
}

func configNetwork() error {
	clear()
	title("Configuring network")

	hname, err := chooseHostname()
	if err != nil {
		return err
	}

	if err := fwrite("/etc/hostname", hname); err != nil {
		return err
	}

	err = fwrite("/etc/hosts", fmt.Sprintf(hosts, hname))
	if err != nil {
		return err
	}

	return sh("systemctl", "enable", "dhcpcd.service")
}

func configUsers() (string, error) {
	clear()
	title("Configuring users")

	fmt.Println("Configuring a new non-root user")
	fmt.Print("Enter a username: ")
	username, err := readLine()
	if err != nil {
		return "", err
	}

	if err := sh("useradd", "--create-home", username); err != nil {
		return "", err
	}

	fmt.Printf("Enter a password for %s\n", username)
	if err := sh("passwd", username); err != nil {
		return "", err
	}

	if err := sh("groupadd", "sudo"); err != nil {
		return "", err
	}

	if err := sh("usermod", "-aG", "sudo", username); err != nil {
		return "", err
	}
	// input group is needed for libinput-gestures.
	if err := sh("usermod", "-aG", "input", username); err != nil {
		return "", err
	}

	err = fwrite("/etc/sudoers", sudoers)
	if err != nil {
		return "", err
	}

	// Disable root login.
	if err := sh("passwd", "--lock", "root"); err != nil {
		return "", err
	}

	return username, nil
}

func installPackages() error {
	clear()
	title("Installing packages")
	time.Sleep(2 * time.Second)

	pkgs := []string{
		"systemd-swap",
		"intel-ucode",
		"xf86-video-intel",
		"util-linux",
		"ufw",
		"dosfstools",
		"lshw",

		// For gestures AUR.
		"xdotool",
		"wmctrl",

		"pulseaudio-alsa",
		"ffmpeg0.10",
		"pulseaudio-bluetooth",

		"bluez",
		"bluez-libs",
		"bluez-utils",
		"bluez-firmware",

		"gdm",
		"gnome",
		"gnome-bluetooth",
		"gnome-user-share",
		"evolution",
		"file-roller",
		"gnome-calendar",
		"gnome-characters",
		"gnome-clocks",
		"gnome-color-manager",
		"gnome-documents",
		"gnome-logs",
		"gnome-music",
		"gnome-photos",
		"gnome-todo",
		"seahorse",
		"gnome-software",

		"ttf-dejavu",
		"adobe-source-code-pro-fonts",
		"ttf-roboto",
		"ttf-inconsolata",
		"noto-fonts-emoji",
		"ttf-liberation",
	}

	cmd := []string{"pacman", "--sync"}
	cmd = append(cmd, pkgs...)
	cmd = append(cmd, "--noconfirm")

	return sh(cmd...)
}

func installAURs(username string) error {
	aurs := []string{
		"https://aur.archlinux.org/cgit/aur.git/snapshot/libinput-gestures.tar.gz",
		"https://aur.archlinux.org/cgit/aur.git/snapshot/systemd-boot-pacman-hook.tar.gz",

		// Need to install deps first before fwupd.
		"https://aur.archlinux.org/cgit/aur.git/snapshot/pesign.tar.gz",
		"https://aur.archlinux.org/cgit/aur.git/snapshot/fwupdate.tar.gz",
		"https://aur.archlinux.org/cgit/aur.git/snapshot/fwupd.tar.gz",
	}

	for _, snap := range aurs {
		if err := installAUR(username, snap); err != nil {
			return err
		}
	}

	return nil
}

func configGnome(username string) error {
	clear()
	title("Configuring GNOME")
	time.Sleep(1 * time.Second)

	asUser := func(name string, cmd []string) []string {
		return append([]string{"sudo", "-u", name}, cmd...)
	}

	cmds := [][]string{
		{"dbus-launch", "gsettings", "set", "org.gnome.desktop.interface", "scaling-factor", "2"},
		{"dbus-launch", "gsettings", "set", "org.gnome.desktop.peripherals.touchpad", "click-method", "fingers"},
		{"dbus-launch", "gsettings", "set", "org.gnome.desktop.peripherals.touchpad", "tap-to-click", "true"},
		{"dbus-launch", "gsettings", "set", "org.gnome.desktop.peripherals.touchpad", "speed", "1"},
		{"dbus-launch", "gsettings", "set", "org.gnome.desktop.peripherals.touchpad", "two-finger-scrolling-enabled", "true"},
		{"dbus-launch", "gsettings", "set", "org.gnome.desktop.peripherals.keyboard", "repeat-interval", "40"},
		{"dbus-launch", "gsettings", "set", "org.gnome.desktop.peripherals.keyboard", "delay", "350"},
	}
	for _, c := range cmds {
		if err := sh(asUser(username, c)...); err != nil {
			return err
		}
		if err := sh(asUser("gdm", c)...); err != nil {
			return err
		}
	}

	return nil
}

func configServices() error {
	clear()
	title("Configuring services")
	time.Sleep(1 * time.Second)

	cmds := [][]string{
		{"systemctl", "enable", "gdm.service"},
		{"systemctl", "enable", "NetworkManager.service"},

		{"systemctl", "enable", "bluetooth.service"},

		// Probably need to configure this...
		{"systemctl", "enable", "systemd-swap"},

		{"systemctl", "enable", "fstrim.timer"},
	}

	for _, c := range cmds {
		if err := sh(c...); err != nil {
			return err
		}
	}

	return nil
}

func configBootloader() error {
	clear()
	fmt.Println("Almost done")
	time.Sleep(2 * time.Second)
	clear()
	title("Configuring bootloader")

	if err := fwrite("/etc/mkinitcpio.conf", mkinitcpioconf); err != nil {
		return err
	}
	if err := sh("mkinitcpio", "-p", "linux"); err != nil {
		return err
	}

	if err := sh("bootctl", "install"); err != nil {
		return err
	}

	rootID, err := blkuuid(rootPartition)
	if err != nil {
		return err
	}

	err = fwrite("/boot/loader/entries/arch.conf", fmt.Sprintf(archconf, rootID, devCryptroot))
	if err != nil {
		return err
	}

	return fwrite("/boot/loader/loader.conf", loaderconf)
}

func installAUR(username, snapURL string) error {
	const (
		aurFile = "xpsaur.tar.gz"
		aurDir  = "xpsaurdir"
	)

	if err := curlo(aurFile, snapURL); err != nil {
		return err
	}

	buildPath := filepath.Join("/home", username, aurDir)
	if err := os.Mkdir(buildPath, 0755); err != nil {
		return err
	}

	err := sh("tar", "-C", buildPath, "-xf", aurFile, "--strip", "1")
	if err != nil {
		return err
	}

	userID, err := uid(username)
	if err != nil {
		return err
	}
	groupID, err := gid(username)
	if err != nil {
		return err
	}

	if err := chownR(buildPath, userID, groupID); err != nil {
		return err
	}

	origPwd, err := os.Getwd()
	if err != nil {
		return err
	}

	if err := os.Chdir(buildPath); err != nil {
		return err
	}

	// makepkg command requires NOPASSWD. Otherwise, it will fail.
	origSudo, err := nopasswdUser(username)
	if err != nil {
		return err
	}
	err = sh("sudo", "-u", username, "makepkg", "-sri", "--noconfirm")
	if err != nil {
		return err
	}
	if err := fwrite("/etc/sudoers", origSudo); err != nil {
		return err
	}

	// Clean up time!
	if err := os.Chdir(origPwd); err != nil {
		return err
	}
	if err := os.RemoveAll(buildPath); err != nil {
		return err
	}
	return os.Remove(aurFile)
}

// nopasswdUser sets NOPASSWD on a given user. It returns the original sudoers
// data so that callers can restore it.
func nopasswdUser(username string) (string, error) {
	orig, err := ioutil.ReadFile("/etc/sudoers")
	if err != nil {
		return "", err
	}

	return string(orig),
		fwrite("/etc/sudoers", fmt.Sprintf(nopasswdSudoers, username))
}

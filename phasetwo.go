package main

import (
	"fmt"
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

	if err := sh("tzselect"); err != nil {
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

	fmt.Println("Enter a new root password")
	if err := sh("passwd"); err != nil {
		return "", err
	}

	fmt.Println("\nConfiguring a new non-root user")
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
		"openssh",
		"lshw",

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

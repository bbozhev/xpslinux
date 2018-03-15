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

	// sudo is needed for sudo.
	if err := sh("usermod", "-aG", "sudo", username); err != nil {
		return "", err
	}
	// input group is needed for libinput-gestures.
	if err := sh("usermod", "-aG", "input", username); err != nil {
		return "", err
	}
	// wheel is needed for Gnome password prompts.
	if err := sh("usermod", "-aG", "wheel", username); err != nil {
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
		"util-linux",
		"ufw",
		"dosfstools",
		"lshw",
		"dmidecode",

		"xf86-video-intel",
		"mesa-libgl",
		"vulkan-intel",
		"libva-intel-driver",

		// For gestures AUR.
		"xdotool",
		"wmctrl",

		"ffmpeg",
		"pulseaudio-alsa",
		"pulseaudio-bluetooth",
		"bluez",
		"bluez-libs",
		"bluez-utils",
		"bluez-firmware",
		"alsa-utils",

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
		"networkmanager",
		"networkmanager-openvpn",
		"networkmanager-openconnect",
		"networkmanager-pptp",

		// Printer
		"cups",
		"cups-pdf",
		"gtk3-print-backends",

		"adobe-source-code-pro-fonts",
		"adobe-source-han-sans-cn-fonts",
		"adobe-source-han-sans-jp-fonts",
		"adobe-source-han-sans-kr-fonts",
		"adobe-source-han-sans-otc-fonts",
		"adobe-source-han-sans-tw-fonts",
		"adobe-source-sans-pro-fonts",
		"noto-fonts-emoji",
		"otf-ipafont",
		"ttf-dejavu",
		"ttf-hanazono",
		"ttf-inconsolata",
		"ttf-liberation",
		"ttf-roboto",
		"ttf-ubuntu-font-family",
	}

	cmd := []string{"pacman", "--sync"}
	cmd = append(cmd, pkgs...)
	cmd = append(cmd, "--noconfirm")

	if err := sh(cmd...); err != nil {
		return err
	}

	err := fwrite("/etc/fonts/conf.d/01-notosans.conf", notosans)
	if err != nil {
		return err
	}

	return sh("fc-cache", "-fv")
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
		{"dbus-launch", "gsettings", "set", "org.gnome.settings-daemon.peripherals.mouse", "double-click", "800"},

		{"dbus-launch", "gsettings", "set", "org.gnome.desktop.peripherals.keyboard", "repeat-interval", "40"},
		{"dbus-launch", "gsettings", "set", "org.gnome.desktop.peripherals.keyboard", "delay", "325"},
	}
	for _, c := range cmds {
		if err := sh(asUser(username, c)...); err != nil {
			return err
		}
		if err := sh(asUser("gdm", c)...); err != nil {
			return err
		}
	}

	cmds = [][]string{
		{"dbus-launch", "gsettings", "set", "org.gnome.desktop.wm.preferences", "button-layout", ":minimize,maximize,close"},
		{"dbus-launch", "gsettings", "set", "org.gnome.desktop.background", "picture-uri", "file:///usr/share/backgrounds/gnome/adwaita-lock.jpg"},
	}
	for _, c := range cmds {
		if err := sh(asUser(username, c)...); err != nil {
			return err
		}
	}

	return nil
}

func configServices() error {
	clear()
	title("Configuring services")
	time.Sleep(1 * time.Second)

	if err := fwrite("/etc/systemd/swap.conf", swapconf); err != nil {
		return err
	}

	cmds := [][]string{
		{"systemctl", "enable", "gdm.service"},
		{"systemctl", "enable", "NetworkManager.service"},
		{"systemctl", "enable", "org.cups.cupsd.service"},
		{"systemctl", "enable", "bluetooth.service"},

		{"systemctl", "enable", "systemd-swap.service"},
		// Periodic TRIM.
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

	rootID, err := blkuuid(lvmRoot)
	if err != nil {
		return err
	}

	err = fwrite("/boot/loader/entries/arch.conf",
		fmt.Sprintf(archconf, rootID, lvmRoot, lvmSwap))
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

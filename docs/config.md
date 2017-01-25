# Configuration

This is a summary of configuration choices that are made on your behalf. If
you're still curious, [phaseone] contains the steps taken before `arch-chroot`,
while [phasetwo] contains the steps take during the chroot.

* 1 LUKS-encrypted root partition
* 1 non-encrypted boot partition
* Root partition uses ext4 filesystem
* Boot partition is UEFI
* Fastest national mirrors
* Terminus as default console font
* Hostname in `/etc/hostname`
* Hostname in `/etc/hosts`
* Enable dhcp service
* Creates non-root user home directory
* `sudo`
* Disables root login
* Installs system utils, graphics and audio drivers, GNOME, and fonts packages
* Color emojis
* Gestures
* Firmware updater
* bootloader updater
* GNOME keyboard and trackpad tweaks
* GNOME minimize and maximize buttons
* Enables gdm, bluetooth, swap, periodic TRIM services
* Custom mkinitcpio for XPS hardware
* systemd-boot UEFI boot manager
* 0 second timeout boot menu
* Install printer packages

[phaseone]: https://github.com/variadico/xpslinux/blob/master/phaseone.go
[phasetwo]: https://github.com/variadico/xpslinux/blob/master/phasetwo.go

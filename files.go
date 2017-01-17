package main

var fstab = `#
# /etc/fstab: static file system information
#
# <file system>	<dir>	<type>	<options>	<dump>	<pass>
# /dev/mapper/cryptroot
UUID=%s	/         	ext4      	relatime,data=ordered	0 1
# /dev/sda1
UUID=%s      	/boot     	vfat      	rw,relatime,fmask=0022,dmask=0022,codepage=437,iocharset=iso8859-1,shortname=mixed,errors=remount-ro	0 2`

var localeconf = `LANG=%[1]s
LC_TIME=%[1]s`

var vconsoleconf = `KEYMAP=%s
FONT=ter-132n`

var hosts = `#
# /etc/hosts: static lookup table for host names
#

#<ip-address>	<hostname.domain.org>	<hostname>
127.0.0.1	localhost.localdomain	localhost
127.0.0.1	%[1]s.localdomain	%[1]s
::1		localhost.localdomain	localhost

# End of file`

var notosans = `<?xml version="1.0" ?>
<!DOCTYPE fontconfig SYSTEM "fonts.dtd">
<fontconfig>
    <match target="scan">
        <test name="family">
            <string>Noto Color Emoji</string>
        </test>
        <edit name="scalable" mode="assign">
            <bool>true</bool>
        </edit>
    </match>

    <match target="pattern">
        <test name="prgname">
            <string>chrome</string>
        </test>
        <edit name="family" mode="prepend_first">
            <string>Noto Color Emoji</string>
        </edit>
    </match>
</fontconfig>`

var sudoers = `## sudoers file.
##
## This file MUST be edited with the 'visudo' command as root.
## Failure to use 'visudo' may result in syntax or file permission errors
## that prevent sudo from running.
##
## See the sudoers man page for the details on how to write a sudoers file.
##

##
## Host alias specification
##
## Groups of machines. These may include host names (optionally with wildcards),
## IP addresses, network numbers or netgroups.
# Host_Alias	WEBSERVERS = www1, www2, www3

##
## User alias specification
##
## Groups of users.  These may consist of user names, uids, Unix groups,
## or netgroups.
# User_Alias	ADMINS = millert, dowdy, mikef

##
## Cmnd alias specification
##
## Groups of commands.  Often used to group related commands together.
# Cmnd_Alias	PROCESSES = /usr/bin/nice, /bin/kill, /usr/bin/renice, \
# 			    /usr/bin/pkill, /usr/bin/top
# Cmnd_Alias	REBOOT = /sbin/halt, /sbin/reboot, /sbin/poweroff

##
## Defaults specification
##
## You may wish to keep some of the following environment variables
## when running commands via sudo.
##
## Locale settings
# Defaults env_keep += "LANG LANGUAGE LINGUAS LC_* _XKB_CHARSET"
##
## Run X applications through sudo; HOME is used to find the
## .Xauthority file.  Note that other programs use HOME to find
## configuration files and this may lead to privilege escalation!
# Defaults env_keep += "HOME"
##
## X11 resource path settings
# Defaults env_keep += "XAPPLRESDIR XFILESEARCHPATH XUSERFILESEARCHPATH"
##
## Desktop path settings
# Defaults env_keep += "QTDIR KDEDIR"
##
## Allow sudo-run commands to inherit the callers' ConsoleKit session
# Defaults env_keep += "XDG_SESSION_COOKIE"
##
## Uncomment to enable special input methods.  Care should be taken as
## this may allow users to subvert the command being run via sudo.
# Defaults env_keep += "XMODIFIERS GTK_IM_MODULE QT_IM_MODULE QT_IM_SWITCHER"
##
## Uncomment to use a hard-coded PATH instead of the user's to find commands
# Defaults secure_path="/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
##
## Uncomment to send mail if the user does not enter the correct password.
# Defaults mail_badpass
##
## Uncomment to enable logging of a command's output, except for
## sudoreplay and reboot.  Use sudoreplay to play back logged sessions.
# Defaults log_output
# Defaults!/usr/bin/sudoreplay !log_output
# Defaults!/usr/local/bin/sudoreplay !log_output
# Defaults!REBOOT !log_output

##
## Runas alias specification
##

##
## User privilege specification
##
root ALL=(ALL) ALL

## Uncomment to allow members of group wheel to execute any command
# %wheel ALL=(ALL) ALL

## Same thing without a password
# %wheel ALL=(ALL) NOPASSWD: ALL

## Uncomment to allow members of group sudo to execute any command
%sudo	ALL=(ALL) ALL

## Uncomment to allow any user to run sudo if they know the password
## of the user they are running the command as (root by default).
# Defaults targetpw  # Ask for the password of the target user
# ALL ALL=(ALL) ALL  # WARNING: only use this together with 'Defaults targetpw'

## Read drop-in files from /etc/sudoers.d
## (the '#' here does not indicate a comment)
#includedir /etc/sudoers.d`

var nopasswdSudoers = `root ALL=(ALL) ALL
%%sudo ALL=(ALL) ALL
%s ALL=(ALL) NOPASSWD:ALL
#includedir /etc/sudoers.d`

var swapconf = `################################################################################
# Defaults are optimized for general usage
################################################################################

################################################################################
# Zswap
#
# Kernel >= 3.11
# Zswap create compress cache between swap and memory for reduce IO
# https://www.kernel.org/doc/Documentation/vm/zswap.txt

zswap_enabled=0
zswap_compressor=lz4
zswap_max_pool_percent=25
zswap_zpool=z3fold

################################################################################
# ZRam
#
# Kernel >= 3.15
# Zram compression streams count for additional information see:
# https://www.kernel.org/doc/Documentation/blockdev/zram.txt

zram_enabled=1
zram_size=$(($ram_size/4))K # This is 1/4 of ram size by default.
zram_streams=$cpu_count
zram_alg=lz4 # lzo lz4 deflate lz4hc 842 - for Linux 4.8.4
zram_prio=32767

################################################################################
# Swap File Universal
# loop + swapfile = support any fs (also btrfs)
swapfu_enabled=0
# File is sparse and dynamically allocated.
swapfu_size=${ram_size}K # Size of swap file.
# File will not be available in fs after script start
# Make sure what script can access to this path during the boot process.
# Full path to swapfile
swapfu_path=/var/swap
swapfu_prio=-1024

################################################################################
# Swap File Chunked
# Allocate swap files dynamically
# Min swap size 256M, Max 16*256M
swapfc_enabled=1
swapfc_chunk_size=256M      # Allocate size of swap chunk
swapfc_max_count=16         # 0 - unlimited
swapfc_free_swap_perc=15    # Add new chunk if free < 15%
swapfc_path=/var/.swapfc/

################################################################################
# Swap devices
# Find and auto swapon all available swap devices
swapd_auto_swapon=1`

var mkinitcpioconf = `# vim:set ft=sh
# MODULES
# The following modules are loaded before any boot hooks are
# run.  Advanced users may wish to specify all system modules
# in this array.  For instance:
#     MODULES="piix ide_disk reiserfs"
MODULES="ext4 algif_skcipher dm_crypt rtsx_pci_sdmmc serio_raw " # NOTICE THE SPACE
MODULES+="atkbd crct10dif_pclmul crc32_pclmul crc32c_intel "     # NOTICE THE SPACE
MODULES+="ghash_clmulni_intel aesni_intel nvme xhci_pci i8042"

# BINARIES
# This setting includes any additional binaries a given user may
# wish into the CPIO image.  This is run last, so it may be used to
# override the actual binaries included by a given hook
# BINARIES are dependency parsed, so you may safely ignore libraries
BINARIES="fsck fsck.ext4 e2fsck cryptsetup"

# FILES
# This setting is similar to BINARIES above, however, files are added
# as-is and are not parsed in any way.  This is useful for config files.
FILES=""

# HOOKS
# This is the most important setting in this file.  The HOOKS control the
# modules and scripts added to the image, and what happens at boot time.
# Order is important, and it is recommended that you do not change the
# order in which HOOKS are added.  Run 'mkinitcpio -H <hook name>' for
# help on a given hook.
# 'base' is _required_ unless you know precisely what you are doing.
# 'udev' is _required_ in order to automatically load modules
# 'filesystems' is _required_ unless you specify your fs modules in MODULES
# Examples:
##   This setup specifies all modules in the MODULES setting above.
##   No raid, lvm2, or encrypted root is needed.
#    HOOKS="base"
#
##   This setup will autodetect all modules for your system and should
##   work as a sane default
#    HOOKS="base udev autodetect block filesystems"
#
##   This setup will generate a 'full' image which supports most systems.
##   No autodetection is done.
#    HOOKS="base udev block filesystems"
#
##   This setup assembles a pata mdadm array with an encrypted root FS.
##   Note: See 'mkinitcpio -H mdadm' for more information on raid devices.
#    HOOKS="base udev block mdadm encrypt filesystems"
#
##   This setup loads an lvm2 volume group on a usb device.
#    HOOKS="base udev block lvm2 filesystems"
#
##   NOTE: If you have /usr on a separate partition, you MUST include the
#    usr, fsck and shutdown hooks.
# HOOKS="base udev block autodetect modconf consolefont encrypt filesystems keyboard fsck"
HOOKS="base udev consolefont encrypt"

# COMPRESSION
# Use this to compress the initramfs image. By default, gzip compression
# is used. Use 'cat' to create an uncompressed image.
#COMPRESSION="gzip"
#COMPRESSION="bzip2"
#COMPRESSION="lzma"
#COMPRESSION="xz"
#COMPRESSION="lzop"
#COMPRESSION="lz4"

# COMPRESSION_OPTIONS
# Additional options for the compressor
#COMPRESSION_OPTIONS=""`

var archconf = `title   arch
linux   /vmlinuz-linux
initrd  /intel-ucode.img
initrd  /initramfs-linux.img
options cryptdevice=UUID=%s:cryptroot root=%s rw quiet vga=current loglevel=3`

var loaderconf = `default  arch
timeout  0
editor   0`

# xpslinux

An Arch Linux installer made for the [Dell XPS 13 9360] for people who want a
computer that Just Works™ ( LOL ) , modified a bit to fit a more secure setup.

[Original Creator](https://github.com/variadico/xpslinux)


## Features

* Minimal XPS configuration
* HiDPI config, font fixes
* Touchpad guestures
* XPS firmware updater
* Disk encryption with LVM
* TMPFS 
* Swap partition 16GBs

More detailed features [here].

## Install

First, create a [bootable USB].

Next, plug it into your Dell XPS. Press the power button, then immediately hit
the F12 button. Select your USB from the boot list.

Then, start running these commands. The text will be tiny, but only until you
run the installer.

```
$ echo "no release has been made which corresponds to the source code so use the following command set:
$ mount -o remount,size=4G /run/archiso/cowspace
$ pacman -Sy go git --noconfirm
$ git clone https://github.com/bbozhev/xpslinux.git
$ cd xpslinux
$ go build
$ ./xpslinux
```

Follow the prompts.

[Dell XPS 13 9360]: https://wiki.archlinux.org/index.php/Dell_XPS_13_(9360)
[bootable USB]: docs/bootable-usb.md
[here]: docs/config.md

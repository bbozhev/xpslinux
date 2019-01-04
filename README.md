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
run the installer, however there is no current build and you will need to build
it yourself due to the fact no pipeline is created for this yet.

```
$ setfont latarcyrheb-sun32                             # to help you see fonts better
$ mount -o remount,size=4G /run/archiso/cowspace        # to give you package installation space
$ wifi-menu                                             # if you're not connected over LAN / cable setup a wifi to have internet access you will need it
$ pacman -Sy go git --noconfirm                         # to not ask you stupid questions and install git and golang
$ git clone https://github.com/bbozhev/xpslinux.git     # to clone the repository
$ cd xpslinux                                           # enter the repo directory
$ go build                                              # build the THING
$ ./xpslinux                                            # execute the THING
```

* Be very careful when entering the user password, as the keyboard is very soft and the current source will cut off the install if you do not provide the confirmation password correctly   

Follow the prompts.

[Dell XPS 13 9360]: https://wiki.archlinux.org/index.php/Dell_XPS_13_(9360)
[bootable USB]: docs/bootable-usb.md
[here]: docs/config.md

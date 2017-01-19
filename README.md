# xpslinux

An Arch Linux installer made for the [Dell XPS 13 9360] for people who want a
computer that Just Works™.

## Features

* Minimal XPS configuration
* HiDPI config, font fixes
* Touchpad guestures
* XPS firmware updater
* Disk encryption

More detailed features [here].

## Install

First, create a [bootable USB].

Next, plug it into your Dell XPS. Press the power button, then immediately hit
the F12 button. Select your USB from the boot list.

Then, start running these commands. The text will be tiny, but only until you
run the installer.

```
$ wifi-menu --obscure  # Select your network, enter password.
$ curl -OL https://github.com/variadico/xpsarch/releases/download/v0.1.0/xpsarch
$ chmod +x xpsarch
$ ./xpsarch
```

Follow the prompts.

[Dell XPS 13 9360]: https://wiki.archlinux.org/index.php/Dell_XPS_13_(9360)
[bootable USB]: docs/bootable-usb.md
[here]: docs/config.md

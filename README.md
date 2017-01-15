# xpsarch

Arch Linux installer for the [Dell XPS 13 9360]

(Alpha)

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

## About

If you're developing software, it's extremely likely that you will have to
interact with Linux at some point in your career. Therefore, learning about
Linux and using Linux tools will be extremely helpful when you have to debug
your app or your app's environment. Server Linux is great.

Unfortunately, desktop [Linux sucks]. You'll experience failures from Wi-Fi not
working to accidentally breaking your whole OS. Developers often use macOS
because it's a decent compromise: you still get a UNIX terminal, but your
computer isn't as buggy.

Unfortunately, macOS is not Linux. Tools you learn on macOS don't always
translate to Linux, so now you have to learn two sets of tools. Learn
`diskutil` on macOS, but `fdisk` on Linux. Sometimes, the commands are the
_same_, but the syntax or behavior is slightly different. `top` on macOS has
different flags than `top` on Linux.

I'm lazy and don't want to learn two sets of things, but I also don't want to
put a bunch of effort in maintaining my computer. Instead of trying to make
Linux work on my computer, I decided to make my computer work on Linux.

### Goals

* Quick and easy installation of Arch Linux.
* Stable configuration specifically for XPS hardware.
* Simple and minimal starting point. Few choices.
* Disk encryption, non-disruptive security improvements.
* Non-disruptive performance optimizations.


[Dell XPS 13 9360]: https://wiki.archlinux.org/index.php/Dell_XPS_13_(9360)
[bootable USB]: docs/bootable-usb.md
[Intel chip]: http://ark.intel.com/products/86068/Intel-Dual-Band-Wireless-AC-8260
[Linux sucks]: https://twitter.com/SwiftOnSecurity/status/817406256583471104

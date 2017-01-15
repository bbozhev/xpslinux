# Create bootable Arch USB

## Download

Download an image from here: https://www.archlinux.org/download

## macOS

First, figure out the name of the drive with `diskutil`. In this example, the
drive is `/dev/disk2`.

```
$ diskutil list
/dev/disk2 (external, physical):
   #:                       TYPE NAME                    SIZE       IDENTIFIER
   0:     FDisk_partition_scheme                        *15.4 GB    disk2
   1:             Windows_FAT_32 NO NAME                 15.4 GB    disk2s1
```

Because macOS auto-mounts USB drives, you have to unmount it before writing blocks with
`dd`.

```
$ diskutil unmountDisk /dev/disk2
diskutil unmountDisk /dev/disk2
Unmount of all volumes on disk2 was successful
```

Now you can use `dd` to copy the ISO onto the USB.

```
$ sudo dd if=archlinux-2017.01.01-dual.iso of=/dev/rdisk2 bs=1m
Password:
867+0 records in
867+0 records out
909115392 bytes transferred in 55.980574 secs (16239837 bytes/sec)
```

Finally, eject the USB.

```
$ diskutil eject /dev/disk2
```


# Dell XPS 13 9360

Note: I replaced the stock Killer wireless chip with an [Intel chip] instead.

```
$ sudo lshw -short
H/W path       Device   Class          Description
==================================================
                        system         XPS 13 9360 (075B)
/0                      bus            0839Y6
/0/0                    memory         64KiB BIOS
/0/38                   memory         16GiB System Memory
/0/38/0                 memory         8GiB Row of chips LPDDR3 Synchronous 1867 MHz (0.5 ns)
/0/38/1                 memory         8GiB Row of chips LPDDR3 Synchronous 1867 MHz (0.5 ns)
/0/3c                   memory         128KiB L1 cache
/0/3d                   memory         512KiB L2 cache
/0/3e                   memory         4MiB L3 cache
/0/3f                   processor      Intel(R) Core(TM) i7-7500U CPU @ 2.70GHz
/0/100                  bridge         Intel Corporation
/0/100/2                display        Intel Corporation
/0/100/4                generic        Skylake Processor Thermal Subsystem
/0/100/14               bus            Sunrise Point-LP USB 3.0 xHCI Controller
/0/100/14/0    usb1     bus            xHCI Host Controller
/0/100/14/0/3           communication  Bluetooth wireless interface
/0/100/14/0/5           multimedia     Integrated_Webcam_HD
/0/100/14/1    usb2     bus            xHCI Host Controller
/0/100/14.2             generic        Sunrise Point-LP Thermal subsystem
/0/100/15               generic        Sunrise Point-LP Serial IO I2C Controller #0
/0/100/15.1             generic        Sunrise Point-LP Serial IO I2C Controller #1
/0/100/16               communication  Sunrise Point-LP CSME HECI #1
/0/100/1c               bridge         Intel Corporation
/0/100/1c.4             bridge         Sunrise Point-LP PCI Express Root Port #5
/0/100/1c.4/0  wlp58s0  network        Wireless 8260
/0/100/1c.5             bridge         Sunrise Point-LP PCI Express Root Port #6
/0/100/1c.5/0           generic        RTS525A PCI Express Card Reader
/0/100/1d               bridge         Sunrise Point-LP PCI Express Root Port #9
/0/100/1d/0             storage        Lite-On Technology Corporation
/0/100/1f               bridge         Intel Corporation
/0/100/1f.2             memory         Memory controller
/0/100/1f.3             multimedia     Intel Corporation
/0/100/1f.4             bus            Sunrise Point-LP SMBus
/1                      power          DELL TP1GT61
```


[Intel chip]: http://ark.intel.com/products/86068/Intel-Dual-Band-Wireless-AC-8260

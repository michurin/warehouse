# Access to your computer's `sshd` from Android phone:

## Prepare phone

On the phone

- *Settings*
- *About phone*
- *Software information*
- tap seven times on *Build number*
- unlock
- *Settings* (go back)
- *Developer options*
- turn on *USB Debugging*

## Manual

On computer:

    adb devices
    adb reverse tcp:3022 tcp:22

    adb detach
    adb kill-server

On phone:

    ssh -p 3022 user@localhost

Useful commands:

    until adb reverse tcp:3022 tcp:22; do sleep 1; done

    adb wait-for-usb-device

    adb shell cmd notification post -S bigtext -t "'Title x'" 'Tag' "'Multiline text'"

## Automation (udev+systemd)

### udev

Useful commands to find out ids:

    sudo udevadm monitor # useful option: --environment
                         # or
                         # udevadm monitor -p -u -s tty
                         # and
                         # udevadm info -a /dev/ttyACM0
    # connect usb cable
    # KERNEL[000] add      /devices/pci0000:00/0000:00:14.0/usb1/1-3/1-3.2 (usb)
    # ...
    # UDEV  [000] bind     /devices/pci0000:00/0000:00:14.0/usb1/1-3/1-3.2 (usb)
    # UDEV  [000] add      /devices/pci0000:00/0000:00:14.0/usb1/1-3/1-3.2/1-3.2:1.1/tty/ttyACM0 (tty)
    # ...
    # UDEV  [000] remove   /devices/pci0000:00/0000:00:14.0/usb1/1-3/1-3.2 (usb)
    # UDEV  [000] remove   /devices/pci0000:00/0000:00:14.0/usb1/1-3/1-3.2/1-3.2:1.1/tty/ttyACM0 (tty)

    sudo dmesg
    # New USB device found, idVendor=04e8, idProduct=6860, bcdDevice= 4.04

Rules `/etc/udev/rules.d/90-android.rules` (use `90`! It's important to have access to correct ENVs)

    SUBSYSTEMS=="tty", ENV{ID_SERIAL}=="SAMSUNG_SAMSUNG_Android_R58RB140BDR", ACTION=="add", RUN+="/usr/bin/systemctl --no-block start android.service"
    SUBSYSTEMS=="tty", ENV{ID_SERIAL}=="SAMSUNG_SAMSUNG_Android_R58RB140BDR", ACTION=="remove", RUN+="/usr/bin/systemctl --no-block stop android.service"

*Be vary careful:* udev just ignores incorrect rules! Check twice typos, `==`s, `"`s and other minor details.

Reload rules:

    sudo udevadm control --reload-rules

Useful commands for debugging:

    udevadm control --log-priority=debug
    journalctl -f

### systemd

Systemd service:

    # /etc/systemd/system/android.service
    # systemctl daemon-reload

    [Unit]
    Description=Android ADB reverse proxy

    [Service]
    Type=forking
    ExecStart=/usr/bin/adb start-server
    ExecStartPost=/usr/bin/adb wait-for-usb-device
    ExecStartPost=/usr/bin/adb shell cmd notification post TagWayfarer "'Connected: port 3022'"
    ExecStartPost=/usr/bin/adb reverse tcp:3022 tcp:22
    ExecStop=/usr/bin/adb kill-server
    Restart=no

    [Install]
    WantedBy=multi-user.target

*Note:* You cannot start `adb` from rules directly. You need systemd service.

*Note:* There is an option to create device-unit. However, systemd is not completely ready for this.

## Links:

- <https://gist.github.com/Pulimet/5013acf2cd5b28e55036c82c91bd56d8>

<!-- ::: vi: set ft=markdown ::: -->

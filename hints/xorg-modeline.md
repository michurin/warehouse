# Subject

I connected Linux-box to TV with native nonstandard resolution 1366x768.
Tools like `cvt` don't help with that.

I'm sharing how to setup it.

This is just examples worked for me.

# Once from command line

Check available and active modes:

```
xrandr
```

If you do not see your resolution, you have to add it:

```
xrandr --newmode "1366x768" 75.61  1366 1406 1438 1574  768 771 777 800 -hsync -vsync
xrandr --addmode HDMI-1 "1366x768"
```

Now you can check the mod is available and switch to it:

```
xrandr --output HDMI-1 --mode "1366x768"
```

# Setup permanent configuration

Just example of `/etc/X11/xorg.conf.d/10-monitor.conf`

```
Section "Monitor"
    Identifier "TCL LCD TV"
    Modeline "1366x768" 75.61 1366 1406 1438 1574 768 771 777 800 -hsync -vsync
    Option "PreferredMode" "1366x768" # intel drivers do not like this option
EndSection

Section "Screen"
    Identifier "Screen0"
    Monitor "TCL LCD TV"
    DefaultDepth 24
    SubSection "Display"
        Depth 24
        Modes "1366x768" # for intel drivers
    EndSubSection
EndSection
```

In fact, for your driver, you need only one section.

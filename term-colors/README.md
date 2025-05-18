# ECMA-48 Select Graphic Rendition (SGR) Codes

This script shows the effect of some basic SGR codes. Not all terminals support all codes.

Example how to grep terminfo to figure out all codes the terminal supports.

```sh
infocmp xterm-256color | grep -E '\\E\[[0-9;]+m'
```

See, [`man console_codes`](https://man7.org/linux/man-pages/man4/console_codes.4.html) for more details.

Amazing collection of color schemes: [github.com/mbadolato/iTerm2-Color-Schemes](https://github.com/mbadolato/iTerm2-Color-Schemes/).
[Scripts](https://github.com/mbadolato/iTerm2-Color-Schemes/tree/master/dynamic-colors) work almost everywhere, not in iTerm2 only.

For example `./dynamic-colors/iTerm2 Tango Dark.sh`:

```sh
#!/bin/sh
# iTerm2 Tango Dark
printf "\033]4;0;#000000;1;#d81e00;2;#5ea702;3;#cfae00;4;#427ab3;5;#89658e;6;#00a7aa;7;#dbded8;8;#686a66;9;#f54235;10;#99e343;11;#fdeb61;12;#84b0d8;13;#bc94b7;14;#37e6e8;15;#f1f1f0\007"
printf "\033]10;#ffffff;#000000;#ffffff\007"
printf "\033]17;#c1deff\007"
printf "\033]19;#000000\007"
printf "\033]5;0;#ffffff\007"
```

# ECMA-48 Select Graphic Rendition (SGR) Codes

This script shows the effect of some basic SGR codes. Not all terminals support all codes.

Example how to grep terminfo to figure out all codes the terminal supports.

```sh
infocmp xterm-256color | grep -E '\\E\[[0-9;]+m'
```

See, [`man console_codes`](https://man7.org/linux/man-pages/man4/console_codes.4.html) for more details.

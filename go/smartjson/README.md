# smartjson

Draft. Not finished, however, ready for playing.

The goal is to format json to fit terminal window.

```sh
for x in `seq 3`; do perl -e "\$x=q{1}x$x;"'print(qq/{"x":[$x, "t", null, false, true]}/);' | go run ./main.go ; done
```

More details in code.

Further improvement:
[term size](https://cs.opensource.google/go/x/term/+/refs/tags/v0.4.0:term_unix.go;drc=1efcd90d861e239a7719db7012b81621e6f7d297;l=60),
[more](https://pkg.go.dev/golang.org/x/sys/unix#IoctlGetWinsize) and
[more](https://pkg.go.dev/golang.org/x/sys/unix#TIOCGWINSZ)
